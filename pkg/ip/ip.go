package ip

import "github.com/valyala/fasthttp"

func GetIP(ctx *fasthttp.RequestCtx) string {
	if ip := string(ctx.Request.Header.Peek("X-Real-IP")); ip != "" {
		return ip
	}
	if forwarded := string(ctx.Request.Header.Peek("X-Forwarded-For")); forwarded != "" {
		return forwarded
	}

	return ctx.RemoteIP().String()

}
