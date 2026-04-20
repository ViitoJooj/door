package main

import (
	"fmt"
	"os"

	httpx "github.com/ViitoJooj/ward/internal/http"
	"github.com/ViitoJooj/ward/internal/http/handler"
	"github.com/ViitoJooj/ward/internal/http/middlewares"
	"github.com/ViitoJooj/ward/internal/repository"
	"github.com/ViitoJooj/ward/internal/services"
	"github.com/ViitoJooj/ward/pkg/database"
	"github.com/ViitoJooj/ward/pkg/dotenv"
	"github.com/ViitoJooj/ward/pkg/ip2location"
	"github.com/ViitoJooj/ward/pkg/logger"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func main() {
	fmt.Println("Starting Ward . . .")
	dotenv.GetEnv()
	database.Conn()
	ip2location.Open()
	router := router.New()

	log := logger.NewLogger(os.Stdout)
	envRepo, authRepo, applicationRepo, logRepo := repository.NewSQLiteRepository(database.DB)

	//Env
	envService := services.NewDotEnvService(envRepo)
	envHandler := handler.NewDotEnvHandler(envService)

	//Auth
	authService := services.NewAuthService(authRepo, log)
	authHandler := handler.NewAuthHandler(authService)

	//Application
	applicationService := services.NewApplicationService(applicationRepo, authRepo)
	applicationHandler := handler.NewApplicationHandler(applicationService)

	//Proxy
	proxyService := services.NewProxyService()
	proxyHandler := handler.NewProxyHandler(proxyService)

	//RequestLog
	requestLogService := services.NewRequestLogService(logRepo)
	requestLogHandler := handler.NewRequestLogHandler(requestLogService)

	//Routers
	httpx.RegisterEnvRouters(router, envHandler)
	httpx.RegisterAuthRoutes(router, authHandler)
	httpx.RegisterApplicationRouters(router, applicationHandler)
	httpx.RegisterRequestLogRoutes(router, requestLogHandler)
	httpx.RegisterProxyRoutes(router, proxyHandler)

	//Middelwares
	handlerWithLog := middlewares.RequestLoggerMiddleware(router.Handler, logRepo)
	handlerWithCors := middlewares.CorsMiddleware(handlerWithLog)

	fasthttp.ListenAndServe(":7171", handlerWithCors)

}
