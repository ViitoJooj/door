package ip

import (
	"strings"

	"github.com/valyala/fasthttp"
)

func GetIP(ctx *fasthttp.RequestCtx) string {
	if ip := strings.TrimSpace(string(ctx.Request.Header.Peek("X-Real-IP"))); ip != "" {
		return ip
	}

	if forwarded := strings.TrimSpace(string(ctx.Request.Header.Peek("X-Forwarded-For"))); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	return ctx.RemoteIP().String()
}
