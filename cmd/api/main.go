package main

import (
	"os"

	httpx "github.com/ViitoJooj/door/internal/http"
	"github.com/ViitoJooj/door/internal/http/handler"
	"github.com/ViitoJooj/door/internal/http/middlewares"
	"github.com/ViitoJooj/door/internal/repository"
	"github.com/ViitoJooj/door/internal/services"
	"github.com/ViitoJooj/door/pkg/database"
	"github.com/ViitoJooj/door/pkg/dotenv"
	"github.com/ViitoJooj/door/pkg/logger"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func main() {
	dotenv.GetEnv()
	database.Conn()

	router := router.New()
	log := logger.NewLogger(os.Stdout)
	authRepo, applicationRepo, logRepo := repository.NewSQLiteRepository(database.DB)

	//Auth
	authService := services.NewAuthService(authRepo, log)
	authHandler := handler.NewAuthHandler(authService)

	//Application
	applicationService := services.NewApplicationService(applicationRepo, authRepo)
	applicationHandler := handler.NewApplicationHandler(applicationService)

	//Proxy
	proxyService := services.NewProxyService()
	proxyHandler := handler.NewProxyHandler(proxyService)

	//Routers
	httpx.RegisterAuthRoutes(router, authHandler)
	httpx.RegisterApplicationRouters(router, applicationHandler)
	httpx.RegisterProxyRoutes(router, proxyHandler)

	//Middelwares
	handlerWithLog := middlewares.RequestLoggerMiddleware(router.Handler, logRepo)
	handlerWithCors := middlewares.CorsMiddleware(handlerWithLog)

	fasthttp.ListenAndServe(":7171", handlerWithCors)

}
