package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	httpx "github.com/ViitoJooj/ward/internal/http"
	"github.com/ViitoJooj/ward/internal/http/handler"
	"github.com/ViitoJooj/ward/internal/http/middlewares"
	"github.com/ViitoJooj/ward/internal/repository"
	"github.com/ViitoJooj/ward/internal/services"
	"github.com/ViitoJooj/ward/pkg/database"
	"github.com/ViitoJooj/ward/pkg/initializer"
	"github.com/ViitoJooj/ward/pkg/ip2location"
	"github.com/ViitoJooj/ward/pkg/logger"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func main() {
	initializer.Init_project()
	database.Conn()
	middlewares.LoadCorsFromDB()
	middlewares.LoadIPAccessListsFromDB()
	ip2location.Open()
	router := router.New()

	logger := logger.NewLogger(os.Stdout)
	envRepo, authRepo, applicationRepo, logRepo, corsRepo, rateLimitRepo, ipAccessListRepo, protocolSettingsRepo, specialRouteRepo := repository.NewSQLiteRepository(database.DB)

	//Cors router
	corsService := services.NewCorsService(corsRepo, authRepo)
	corsHandler := handler.NewCorsHandler(corsService)

	//Env
	envService := services.NewDotEnvService(envRepo)
	envHandler := handler.NewDotEnvHandler(envService)

	//Auth
	authService := services.NewAuthService(authRepo, logger)
	authHandler := handler.NewAuthHandler(authService)

	//Users
	userService := services.NewUserService(authRepo, logger)
	userHandler := handler.NewUserHandler(userService)

	//Application
	applicationService := services.NewApplicationService(applicationRepo, authRepo)
	applicationHandler := handler.NewApplicationHandler(applicationService)

	//Proxy
	proxyService := services.NewProxyService()
	proxyHandler := handler.NewProxyHandler(proxyService, applicationService)

	//RequestLog
	requestLogService := services.NewRequestLogService(logRepo)
	requestLogHandler := handler.NewRequestLogHandler(requestLogService)

	//Health
	healthService := services.NewHealthService(logRepo)
	healthHandler := handler.NewHealthHandler(healthService)

	//RateLimit
	rateLimitService := services.NewRateLimitService(rateLimitRepo)
	rateLimitHandler := handler.NewRateLimitHandler(rateLimitService)

	//IP Access List
	ipAccessListService := services.NewIPAccessListService(ipAccessListRepo, authRepo)
	ipAccessListHandler := handler.NewIPAccessListHandler(ipAccessListService)

	//Protocol settings
	protocolSettingsService := services.NewProtocolSettingsService(protocolSettingsRepo)
	protocolSettingsHandler := handler.NewProtocolSettingsHandler(protocolSettingsService)

	//Special route rules
	specialRouteService := services.NewSpecialRouteService(specialRouteRepo, authRepo)
	specialRouteHandler := handler.NewSpecialRouteHandler(specialRouteService)

	//Routers
	httpx.RegisterEnvRouters(router, envHandler)
	httpx.RegisterAuthRoutes(router, authHandler)
	httpx.RegisterApplicationRouters(router, applicationHandler)
	httpx.RegisterRequestLogRoutes(router, requestLogHandler)
	httpx.RegisterHealthRoutes(router, healthHandler)
	httpx.RegisterCorsOriginsRouters(router, corsHandler)
	httpx.RegisterUserRouters(router, userHandler)
	httpx.RegisterRateLimitRouters(router, rateLimitHandler)
	httpx.RegisterIPAccessListRouters(router, ipAccessListHandler)
	httpx.RegisterProtocolSettingsRouters(router, protocolSettingsHandler)
	httpx.RegisterSpecialRouteRouters(router, specialRouteHandler)
	httpx.RegisterProxyRoutes(router, proxyHandler)

	rateLimitSettings, err := rateLimitService.Get()
	if err != nil {
		log.Println(err)
	} else {
		middlewares.UpdateRateLimitConfig(rateLimitSettings.RequestsPerSecond, rateLimitSettings.Burst, rateLimitSettings.Progressive)
	}

	protocolSettings, err := protocolSettingsService.Get()
	if err != nil {
		log.Println(err)
	} else {
		middlewares.UpdateAllowedProtocol(protocolSettings.AllowedProtocol)
		middlewares.UpdateConfigApplyScope(protocolSettings.ApplyScope)
	}

	middlewares.LoadSpecialRoutesFromDB()

	//Middelwares
	handlerWithLog := middlewares.RequestLoggerMiddleware(router.Handler, logRepo)
	handlerWithCors := middlewares.CorsMiddleware(handlerWithLog)
	handlerWithProtocolMode := middlewares.ProtocolModeMiddleware(handlerWithCors)
	handlerWithSpecialRoutes := middlewares.SpecialRoutesMiddleware(handlerWithProtocolMode)
	handlerWithRateLimit := middlewares.RateLimitMiddleware(handlerWithSpecialRoutes)

	port := initializer.EnsureAppPort(database.DB)
	os.Setenv("APP_PORT", strconv.Itoa(port))
	address := fmt.Sprintf(":%d", port)

	log.Printf("Ward running on port: %d", port)
	fasthttp.ListenAndServe(address, handlerWithRateLimit)
}
