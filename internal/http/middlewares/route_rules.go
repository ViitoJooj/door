package middlewares

import (
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/ViitoJooj/ward/pkg/database"
	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"
)

type routeRuleEntry struct {
	path              string
	method            string
	rateLimitEnabled  bool
	rateLimitRPS      float64
	rateLimitBurst    int
	targetURL         string
	geoRoutingEnabled bool
}

// RouteRuleInfo holds the publicly readable fields of a matched route rule.
type RouteRuleInfo struct {
	TargetURL         string
	GeoRoutingEnabled bool
	RateLimitEnabled  bool
}

var (
	routeRules   []routeRuleEntry
	routeRulesMu sync.RWMutex

	routeRateLimiters   = map[string]*IPRateLimiter{}
	routeRateLimitersMu sync.Mutex
)

func LoadRouteRulesFromDB() {
	rows, err := database.DB.Query(`
		SELECT path, method, rate_limit_enabled, rate_limit_rps, rate_limit_burst, target_url, geo_routing_enabled
		FROM route_rules WHERE enabled = 1
	`)
	if err != nil {
		log.Println("error loading route rules:", err)
		return
	}
	defer rows.Close()

	newRules := make([]routeRuleEntry, 0)
	for rows.Next() {
		var entry routeRuleEntry
		var rle, geo int
		if err := rows.Scan(&entry.path, &entry.method, &rle, &entry.rateLimitRPS, &entry.rateLimitBurst, &entry.targetURL, &geo); err != nil {
			continue
		}
		entry.rateLimitEnabled = rle == 1
		entry.geoRoutingEnabled = geo == 1
		newRules = append(newRules, entry)
	}

	routeRulesMu.Lock()
	routeRules = newRules
	routeRulesMu.Unlock()

	routeRateLimitersMu.Lock()
	routeRateLimiters = map[string]*IPRateLimiter{}
	for _, r := range newRules {
		if r.rateLimitEnabled {
			key := r.method + "|" + r.path
			routeRateLimiters[key] = NewIPRateLimiter(rate.Limit(r.rateLimitRPS), r.rateLimitBurst)
		}
	}
	routeRateLimitersMu.Unlock()
}

// FindRouteRule returns public info for the matching rule for a request path + method.
func FindRouteRule(path, method string) *RouteRuleInfo {
	routeRulesMu.RLock()
	defer routeRulesMu.RUnlock()
	for i := range routeRules {
		r := &routeRules[i]
		if r.path != path {
			continue
		}
		if r.method == "" || r.method == method {
			return &RouteRuleInfo{
				TargetURL:         r.targetURL,
				GeoRoutingEnabled: r.geoRoutingEnabled,
				RateLimitEnabled:  r.rateLimitEnabled,
			}
		}
	}
	return nil
}

func findRouteRuleEntry(path, method string) *routeRuleEntry {
	routeRulesMu.RLock()
	defer routeRulesMu.RUnlock()
	for i := range routeRules {
		r := &routeRules[i]
		if r.path != path {
			continue
		}
		if r.method == "" || r.method == method {
			return r
		}
	}
	return nil
}

// RouteRuleRateLimitMiddleware applies per-route rate limiting when a rule defines it.
func RouteRuleRateLimitMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if !ShouldApplySecurityConfigs(ctx) {
			next(ctx)
			return
		}

		path := string(ctx.Path())
		method := string(ctx.Method())
		rule := findRouteRuleEntry(path, method)
		if rule == nil || !rule.rateLimitEnabled {
			next(ctx)
			return
		}

		key := rule.method + "|" + rule.path
		routeRateLimitersMu.Lock()
		limiter := routeRateLimiters[key]
		routeRateLimitersMu.Unlock()

		if limiter == nil {
			next(ctx)
			return
		}

		clientIP, _, err := net.SplitHostPort(ctx.RemoteAddr().String())
		if err != nil {
			next(ctx)
			return
		}

		now := time.Now()
		blocked, blockDuration := limiter.IsBlocked(clientIP, now)
		if blocked {
			ctx.Response.Header.Set("Retry-After", strconv.FormatInt(int64(blockDuration.Seconds()), 10))
			ctx.Error("Rate limit.", fasthttp.StatusTooManyRequests)
			return
		}

		l := limiter.GetLimiter(clientIP)
		if !l.Allow() {
			blockDuration = limiter.RegisterAbuse(clientIP, now)
			if blockDuration > 0 {
				ctx.Response.Header.Set("Retry-After", strconv.FormatInt(int64(blockDuration.Seconds()), 10))
			}
			ctx.Error("Rate limit.", fasthttp.StatusTooManyRequests)
			return
		}

		next(ctx)
	}
}
