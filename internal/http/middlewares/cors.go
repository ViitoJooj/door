package middlewares

import (
	"github.com/ViitoJooj/ward/pkg/dotenv"
	"github.com/valyala/fasthttp"
)

func CorsMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		origin := string(ctx.Request.Header.Peek("Origin"))

		if _, ok := dotenv.CorsOriginsMap[origin]; ok {
			ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
			ctx.Response.Header.Set("Vary", "Origin")
			ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
			ctx.Response.Header.Set("Access-Control-Max-Age", "3600")
		}

		if string(ctx.Method()) == "OPTIONS" {
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.Response.Header.Set("Content-Length", "0")
			return
		}

		next(ctx)
	}
}
