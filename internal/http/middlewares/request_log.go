package middlewares

import (
	"time"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/internal/repository"
	"github.com/ViitoJooj/ward/pkg/ip"
	"github.com/valyala/fasthttp"
)

func RequestLoggerMiddleware(next fasthttp.RequestHandler, repo repository.RequestLogRepository) fasthttp.RequestHandler {
	ch := make(chan *domain.RequestLog, 512)

	go func() {
		for log := range ch {
			repo.InsertRequestLog(log)
		}
	}()

	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()

		next(ctx)

		elapsed := time.Since(start).Milliseconds()

		ch <- &domain.RequestLog{
			Method:         string(ctx.Method()),
			Path:           string(ctx.Path()),
			QueryString:    string(ctx.URI().QueryString()),
			StatusCode:     ctx.Response.StatusCode(),
			ResponseTimeMs: elapsed,
			IP:             ip.GetIP(ctx),
			UserAgent:      string(ctx.UserAgent()),
			Referer:        string(ctx.Referer()),
			RequestSize:    len(ctx.Request.Body()),
			ResponseSize:   len(ctx.Response.Body()),
			Internal:       false,
		}
	}
}
