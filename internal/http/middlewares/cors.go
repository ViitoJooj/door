package middlewares

import (
	"log"
	"strings"
	"sync"

	"github.com/ViitoJooj/ward/pkg/database"
	"github.com/valyala/fasthttp"
)

var (
	corsOriginsMap = map[string]bool{}
	corsMutex      sync.RWMutex
)

func LoadCorsFromDB() {
	rows, err := database.DB.Query(`SELECT origin FROM cors`)
	if err != nil {
		log.Println("error loading cors from db:", err)
		return
	}
	defer rows.Close()

	temp := map[string]bool{}

	for rows.Next() {
		var origin string
		if err := rows.Scan(&origin); err != nil {
			continue
		}

		origin = strings.TrimSpace(origin)
		if origin != "" {
			temp[origin] = true
		}
	}

	corsMutex.Lock()
	corsOriginsMap = temp
	corsMutex.Unlock()

	log.Println("CORS loaded:", len(corsOriginsMap))
}

func CorsMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		origin := string(ctx.Request.Header.Peek("Origin"))

		corsMutex.RLock()
		allowed := corsOriginsMap[origin]
		corsMutex.RUnlock()

		if allowed {
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
