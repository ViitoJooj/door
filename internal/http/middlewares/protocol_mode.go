package middlewares

import (
	"strings"
	"sync"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/valyala/fasthttp"
)

var (
	allowedProtocolMode = domain.ProtocolModeBoth
	configApplyScope    = domain.ConfigScopeAll
	protocolModeMu      sync.RWMutex
)

func UpdateAllowedProtocol(mode string) {
	normalized, err := domain.NormalizeProtocolMode(mode)
	if err != nil {
		normalized = domain.ProtocolModeBoth
	}

	protocolModeMu.Lock()
	allowedProtocolMode = normalized
	protocolModeMu.Unlock()
}

func currentProtocolMode() string {
	protocolModeMu.RLock()
	defer protocolModeMu.RUnlock()
	return allowedProtocolMode
}

func UpdateConfigApplyScope(scope string) {
	normalized, err := domain.NormalizeConfigScope(scope)
	if err != nil {
		normalized = domain.ConfigScopeAll
	}

	protocolModeMu.Lock()
	configApplyScope = normalized
	protocolModeMu.Unlock()
}

func currentConfigApplyScope() string {
	protocolModeMu.RLock()
	defer protocolModeMu.RUnlock()
	return configApplyScope
}

func ShouldApplySecurityConfigs(ctx *fasthttp.RequestCtx) bool {
	scope := currentConfigApplyScope()
	if scope == domain.ConfigScopeAll {
		return true
	}

	path := string(ctx.Path())
	return !strings.HasPrefix(path, "/ward/api/")
}

func requestIsHTTPS(ctx *fasthttp.RequestCtx) bool {
	if ctx.IsTLS() {
		return true
	}

	forwardedProto := strings.ToLower(strings.TrimSpace(string(ctx.Request.Header.Peek("X-Forwarded-Proto"))))
	return forwardedProto == "https"
}

func ProtocolModeMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if !ShouldApplySecurityConfigs(ctx) {
			next(ctx)
			return
		}

		mode := currentProtocolMode()
		isHTTPS := requestIsHTTPS(ctx)

		if mode == domain.ProtocolModeHTTPS && !isHTTPS {
			ctx.Error("HTTPS only.", fasthttp.StatusForbidden)
			return
		}

		if mode == domain.ProtocolModeHTTP && isHTTPS {
			ctx.Error("HTTP only.", fasthttp.StatusForbidden)
			return
		}

		next(ctx)
	}
}
