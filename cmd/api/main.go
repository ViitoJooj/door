package main

import (
	"log"
	"os"

	httpx "github.com/ViitoJooj/ward/internal/http"
	"github.com/ViitoJooj/ward/internal/http/handler"
	"github.com/ViitoJooj/ward/internal/http/middlewares"
	"github.com/ViitoJooj/ward/internal/repository"
	"github.com/ViitoJooj/ward/internal/services"
	"github.com/ViitoJooj/ward/pkg/database"
	initproject "github.com/ViitoJooj/ward/pkg/init_project"
	"github.com/ViitoJooj/ward/pkg/ip2location"
	"github.com/ViitoJooj/ward/pkg/logger"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func main() {
	initproject.Init_project()
	middlewares.LoadCorsFromDB()
	database.Conn()
	ip2location.Open()
	router := router.New()

	logger := logger.NewLogger(os.Stdout)
	envRepo, authRepo, applicationRepo, logRepo, corsRepo := repository.NewSQLiteRepository(database.DB)

	//Cors router
	corsService := services.NewCorsService(corsRepo, authRepo)
	corsHandler := handler.NewCorsHandler(corsService)

	//Env
	envService := services.NewDotEnvService(envRepo)
	envHandler := handler.NewDotEnvHandler(envService)

	//Auth
	authService := services.NewAuthService(authRepo, logger)
	authHandler := handler.NewAuthHandler(authService)

	//Application
	applicationService := services.NewApplicationService(applicationRepo, authRepo)
	applicationHandler := handler.NewApplicationHandler(applicationService)

	//Proxy
	proxyService := services.NewProxyService()
	proxyHandler := handler.NewProxyHandler(proxyService, applicationService)

	//RequestLog
	requestLogService := services.NewRequestLogService(logRepo)
	requestLogHandler := handler.NewRequestLogHandler(requestLogService)

	//Routers
	httpx.RegisterEnvRouters(router, envHandler)
	httpx.RegisterAuthRoutes(router, authHandler)
	httpx.RegisterApplicationRouters(router, applicationHandler)
	httpx.RegisterRequestLogRoutes(router, requestLogHandler)
	httpx.RegisterProxyRoutes(router, proxyHandler)
	httpx.RegisterCorsOriginsRouters(router, corsHandler)

	//Middelwares
	handlerWithLog := middlewares.RequestLoggerMiddleware(router.Handler, logRepo)
	handlerWithCors := middlewares.CorsMiddleware(handlerWithLog)

	log.Println("Ward running on port: 7171")
	fasthttp.ListenAndServe(":7171", handlerWithCors)
}
