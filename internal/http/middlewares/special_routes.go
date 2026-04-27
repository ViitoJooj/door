package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/pkg/database"
	"github.com/ViitoJooj/ward/pkg/ip"
	"github.com/valyala/fasthttp"
)

type specialRouteRuntimeState struct {
	DistinctAttempts map[string]time.Time
	BannedUntil      time.Time
}

var (
	specialRouteRulesByPath = map[string]*domain.SpecialRouteRule{}
	specialRouteStates      = map[string]*specialRouteRuntimeState{}
	specialRouteMu          sync.RWMutex
)

func specialRouteKey(ruleID int, clientIP string) string {
	return fmt.Sprintf("%d|%s", ruleID, clientIP)
}

func requestFingerprint(ctx *fasthttp.RequestCtx) string {
	payload := strings.Join([]string{
		string(ctx.Method()),
		string(ctx.Path()),
		string(ctx.URI().QueryString()),
		string(ctx.PostBody()),
	}, "|")
	hash := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(hash[:])
}

func canonicalSpecialPath(path string) string {
	normalized := strings.TrimSpace(path)
	if normalized == "" {
		return normalized
	}
	if len(normalized) > 1 {
		return strings.TrimSuffix(normalized, "/")
	}
	return normalized
}

func LoadSpecialRoutesFromDB() {
	rows, err := database.DB.Query(`
		SELECT id, route_type, path, max_distinct_requests, window_seconds, ban_seconds, enabled, created_by, updated_by, created_at, updated_at
		FROM special_route_rules
		WHERE enabled = 1
	`)
	if err != nil {
		log.Println("error loading special routes from db:", err)
		return
	}
	defer rows.Close()

	newRules := map[string]*domain.SpecialRouteRule{}
	for rows.Next() {
		rule := &domain.SpecialRouteRule{}
		var enabledInt int
		if err := rows.Scan(
			&rule.ID,
			&rule.RouteType,
			&rule.Path,
			&rule.MaxDistinctRequests,
			&rule.WindowSeconds,
			&rule.BanSeconds,
			&enabledInt,
			&rule.CreatedBy,
			&rule.UpdatedBy,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		); err != nil {
			continue
		}
		rule.Enabled = enabledInt == 1
		newRules[canonicalSpecialPath(rule.Path)] = rule
	}

	specialRouteMu.Lock()
	specialRouteRulesByPath = newRules
	specialRouteStates = map[string]*specialRouteRuntimeState{}
	specialRouteMu.Unlock()
}

func applySpecialRouteRule(ctx *fasthttp.RequestCtx, rule *domain.SpecialRouteRule, clientIP string, now time.Time) (blocked bool, retryAfter int) {
	specialRouteMu.Lock()
	defer specialRouteMu.Unlock()

	stateKey := specialRouteKey(rule.ID, clientIP)
	state, exists := specialRouteStates[stateKey]
	if !exists {
		state = &specialRouteRuntimeState{
			DistinctAttempts: map[string]time.Time{},
		}
		specialRouteStates[stateKey] = state
	}

	if state.BannedUntil.After(now) {
		return true, int(time.Until(state.BannedUntil).Seconds())
	}

	windowStart := now.Add(-time.Duration(rule.WindowSeconds) * time.Second)
	for fp, ts := range state.DistinctAttempts {
		if ts.Before(windowStart) {
			delete(state.DistinctAttempts, fp)
		}
	}

	fingerprint := requestFingerprint(ctx)
	state.DistinctAttempts[fingerprint] = now
	if len(state.DistinctAttempts) > rule.MaxDistinctRequests {
		state.BannedUntil = now.Add(time.Duration(rule.BanSeconds) * time.Second)
		state.DistinctAttempts = map[string]time.Time{}
		return true, rule.BanSeconds
	}

	return false, 0
}

func SpecialRoutesMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if !ShouldApplySecurityConfigs(ctx) {
			next(ctx)
			return
		}

		path := canonicalSpecialPath(string(ctx.Path()))

		specialRouteMu.RLock()
		rule, exists := specialRouteRulesByPath[path]
		specialRouteMu.RUnlock()
		if !exists || !rule.Enabled {
			next(ctx)
			return
		}

		clientIP := ip.GetIP(ctx)
		blocked, retryAfter := applySpecialRouteRule(ctx, rule, clientIP, time.Now())
		if blocked {
			if retryAfter > 0 {
				ctx.Response.Header.Set("Retry-After", fmt.Sprintf("%d", retryAfter))
			}
			ctx.Error("too many different attempts.", fasthttp.StatusTooManyRequests)
			return
		}

		next(ctx)
	}
}
