package middlewares

import (
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ViitoJooj/ward/pkg/database"
	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"
)

const (
	defaultPenaltyBaseDuration = time.Second
	maxPenaltyStrikes          = 20
)

type progressivePenalty struct {
	blockedUntil time.Time
	strikes      int
}

type IPRateLimiter struct {
	ips               map[string]*rate.Limiter
	penalties         map[string]*progressivePenalty
	mu                sync.Mutex
	r                 rate.Limit
	b                 int
	progressive       bool
	penaltyBaseWindow time.Duration
}

var (
	whitelistIPs = map[string]bool{}
	blacklistIPs = map[string]bool{}
	accessMu     sync.RWMutex
)

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips:               make(map[string]*rate.Limiter),
		penalties:         make(map[string]*progressivePenalty),
		r:                 r,
		b:                 b,
		penaltyBaseWindow: defaultPenaltyBaseDuration,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

func (i *IPRateLimiter) SetConfig(requestsPerSecond rate.Limit, burst int, progressive bool) {
	if requestsPerSecond <= 0 {
		requestsPerSecond = 1
	}
	if burst <= 0 {
		burst = 1
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	i.r = requestsPerSecond
	i.b = burst
	i.progressive = progressive
	i.ips = make(map[string]*rate.Limiter)
	i.penalties = make(map[string]*progressivePenalty)
}

func (i *IPRateLimiter) penaltyDuration(strikes int) time.Duration {
	if strikes < 1 {
		strikes = 1
	}
	if strikes > maxPenaltyStrikes {
		strikes = maxPenaltyStrikes
	}
	return i.penaltyBaseWindow * time.Duration(1<<(strikes-1))
}

func (i *IPRateLimiter) RegisterAbuse(ip string, now time.Time) time.Duration {
	i.mu.Lock()
	defer i.mu.Unlock()

	if !i.progressive {
		return 0
	}

	penalty, exists := i.penalties[ip]
	if !exists {
		penalty = &progressivePenalty{}
		i.penalties[ip] = penalty
	}

	penalty.strikes++
	duration := i.penaltyDuration(penalty.strikes)
	penalty.blockedUntil = now.Add(duration)
	return duration
}

func (i *IPRateLimiter) IsBlocked(ip string, now time.Time) (bool, time.Duration) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if !i.progressive {
		return false, 0
	}

	penalty, exists := i.penalties[ip]
	if !exists || now.After(penalty.blockedUntil) {
		return false, 0
	}

	penalty.strikes++
	duration := i.penaltyDuration(penalty.strikes)
	penalty.blockedUntil = now.Add(duration)
	return true, duration
}

func (i *IPRateLimiter) ClearPenalty(ip string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	delete(i.penalties, ip)
}

var ipLimiter = NewIPRateLimiter(1, 5)

func UpdateRateLimitConfig(requestsPerSecond float64, burst int, progressive bool) {
	ipLimiter.SetConfig(rate.Limit(requestsPerSecond), burst, progressive)
}

func LoadIPAccessListsFromDB() {
	whitelistRows, err := database.DB.Query(`SELECT ip FROM ip_whitelist`)
	if err != nil {
		log.Println("error loading ip whitelist from db:", err)
		return
	}
	defer whitelistRows.Close()

	blacklistRows, err := database.DB.Query(`SELECT ip FROM ip_blacklist`)
	if err != nil {
		log.Println("error loading ip blacklist from db:", err)
		return
	}
	defer blacklistRows.Close()

	newWhitelist := map[string]bool{}
	newBlacklist := map[string]bool{}

	for whitelistRows.Next() {
		var ip string
		if err := whitelistRows.Scan(&ip); err != nil {
			continue
		}
		ip = strings.TrimSpace(ip)
		if ip != "" {
			newWhitelist[ip] = true
		}
	}

	for blacklistRows.Next() {
		var ip string
		if err := blacklistRows.Scan(&ip); err != nil {
			continue
		}
		ip = strings.TrimSpace(ip)
		if ip != "" {
			newBlacklist[ip] = true
		}
	}

	accessMu.Lock()
	whitelistIPs = newWhitelist
	blacklistIPs = newBlacklist
	accessMu.Unlock()
}

func IsWhitelisted(ip string) bool {
	accessMu.RLock()
	defer accessMu.RUnlock()
	return whitelistIPs[ip]
}

func IsBlacklisted(ip string) bool {
	accessMu.RLock()
	defer accessMu.RUnlock()
	return blacklistIPs[ip]
}

func RateLimitMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if !ShouldApplySecurityConfigs(ctx) {
			next(ctx)
			return
		}

		ip, _, err := net.SplitHostPort(string(ctx.RemoteAddr().String()))
		if err != nil {
			ctx.Error("Erro on processing IP", fasthttp.StatusInternalServerError)
			return
		}

		if IsBlacklisted(ip) {
			ctx.Error("IP blocked.", fasthttp.StatusForbidden)
			return
		}

		if IsWhitelisted(ip) {
			ipLimiter.ClearPenalty(ip)
			next(ctx)
			return
		}

		now := time.Now()
		blocked, blockDuration := ipLimiter.IsBlocked(ip, now)
		if blocked {
			ctx.Response.Header.Set("Retry-After", strconv.FormatInt(int64(blockDuration.Seconds()), 10))
			ctx.Error("Rate limit.", fasthttp.StatusTooManyRequests)
			return
		}

		limiter := ipLimiter.GetLimiter(ip)
		if !limiter.Allow() {
			blockDuration = ipLimiter.RegisterAbuse(ip, now)
			if blockDuration > 0 {
				ctx.Response.Header.Set("Retry-After", strconv.FormatInt(int64(blockDuration.Seconds()), 10))
			}
			ctx.Error("Rate limit.", fasthttp.StatusTooManyRequests)
			return
		}

		ipLimiter.ClearPenalty(ip)
		next(ctx)
	}
}
