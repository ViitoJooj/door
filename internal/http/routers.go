package http

import (
	"github.com/ViitoJooj/ward/internal/http/handler"
	"github.com/ViitoJooj/ward/internal/http/middlewares"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func RegisterAuthRoutes(r *router.Router, userController *handler.AuthHandler) {
	r.POST("/ward/api/v1/auth/register", userController.Register)
	r.POST("/ward/api/v1/auth/login", userController.Login)
	r.GET("/ward/api/v1/auth/token", userController.Token)
	r.POST("/ward/api/v1/auth/logout", userController.Logout)
}

func RegisterProxyRoutes(r *router.Router, proxyController *handler.ProxyHandler) {
	r.ANY("/{path:*}", proxyController.Proxy)
}

func RegisterRequestLogRoutes(r *router.Router, requestLogController *handler.RequestLogHandler) {
	r.GET("/ward/api/v1/logs", middlewares.UserIdMiddleware(requestLogController.GetAll))
}

func RegisterHealthRoutes(r *router.Router, healthHandler *handler.HealthHandler) {
	r.GET("/ward/api/v1/health", middlewares.UserIdMiddleware(healthHandler.GetOverview))
	r.GET("/ward/api/v1/health/routes", middlewares.UserIdMiddleware(healthHandler.GetRoutes))
}

func RegisterApplicationRouters(r *router.Router, applicationController *handler.ApplicationHandler) {
	r.GET("/ward/api/v1/applications", middlewares.UserIdMiddleware(applicationController.GetAll))
	r.GET("/ward/api/v1/applications/{path:*}", middlewares.UserIdMiddleware(applicationController.GetByID))
	r.POST("/ward/api/v1/applications", middlewares.UserIdMiddleware(applicationController.Create))
	r.DELETE("/ward/api/v1/applications/{path:*}", middlewares.UserIdMiddleware(applicationController.DeleteById))
}

func RegisterEnvRouters(r *router.Router, envController *handler.DotEnvHandler) {
	r.GET("/ward/api/v1/env/", middlewares.UserIdMiddleware(envController.GetAll))
	r.GET("/ward/api/v1/env/{path:*}", middlewares.UserIdMiddleware(envController.GetVar))
	r.PUT("/ward/api/v1/env/{path:*}", middlewares.UserIdMiddleware(envController.ChangeVar))
}

func RegisterCorsOriginsRouters(r *router.Router, corsController *handler.CorsHandler) {
	r.GET("/ward/api/v1/cors/{path:*}", middlewares.UserIdMiddleware(corsController.GetByID))
	r.GET("/ward/api/v1/cors/", middlewares.UserIdMiddleware(corsController.GetAll))
	r.POST("/ward/api/v1/cors/", middlewares.UserIdMiddleware(corsController.Create))
	r.PUT("/ward/api/v1/cors/{path:*}", middlewares.UserIdMiddleware(corsController.Update))
	r.DELETE("/ward/api/v1/cors/{path:*}", middlewares.UserIdMiddleware(corsController.DeleteById))
}

func RegisterUserRouters(r *router.Router, userController *handler.UserHandler) {
	adminOnly := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return middlewares.UserIdMiddleware(middlewares.AdminOnlyMiddleware(next))
	}

	r.PUT("/ward/api/v1/users/me", middlewares.UserIdMiddleware(userController.UpdateMe))

	r.GET("/ward/api/v1/users", adminOnly(userController.GetAll))
	r.POST("/ward/api/v1/users", adminOnly(userController.Create))
	r.GET("/ward/api/v1/users/{path:*}", adminOnly(userController.GetByID))
	r.PUT("/ward/api/v1/users/{path:*}", adminOnly(userController.UpdateByID))
	r.DELETE("/ward/api/v1/users/{path:*}", adminOnly(userController.DeleteByID))
}

func RegisterRateLimitRouters(r *router.Router, rateLimitHandler *handler.RateLimitHandler) {
	adminOnly := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return middlewares.UserIdMiddleware(middlewares.AdminOnlyMiddleware(next))
	}

	r.GET("/ward/api/v1/rate-limit", adminOnly(rateLimitHandler.Get))
	r.PUT("/ward/api/v1/rate-limit", adminOnly(rateLimitHandler.Update))
}

func RegisterIPAccessListRouters(r *router.Router, ipAccessListHandler *handler.IPAccessListHandler) {

	r.GET("/ward/api/v1/ip-whitelist", middlewares.UserIdMiddleware(ipAccessListHandler.GetWhitelist))
	r.POST("/ward/api/v1/ip-whitelist", middlewares.UserIdMiddleware(ipAccessListHandler.CreateWhitelist))
	r.PUT("/ward/api/v1/ip-whitelist/{path:*}", middlewares.UserIdMiddleware(ipAccessListHandler.UpdateWhitelist))
	r.DELETE("/ward/api/v1/ip-whitelist/{path:*}", middlewares.UserIdMiddleware(ipAccessListHandler.DeleteWhitelist))

	r.GET("/ward/api/v1/ip-blacklist", middlewares.UserIdMiddleware(ipAccessListHandler.GetBlacklist))
	r.POST("/ward/api/v1/ip-blacklist", middlewares.UserIdMiddleware(ipAccessListHandler.CreateBlacklist))
	r.PUT("/ward/api/v1/ip-blacklist/{path:*}", middlewares.UserIdMiddleware(ipAccessListHandler.UpdateBlacklist))
	r.DELETE("/ward/api/v1/ip-blacklist/{path:*}", middlewares.UserIdMiddleware(ipAccessListHandler.DeleteBlacklist))
}

func RegisterProtocolSettingsRouters(r *router.Router, protocolHandler *handler.ProtocolSettingsHandler) {
	adminOnly := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return middlewares.UserIdMiddleware(middlewares.AdminOnlyMiddleware(next))
	}

	r.GET("/ward/api/v1/protocol-mode", adminOnly(protocolHandler.Get))
	r.PUT("/ward/api/v1/protocol-mode", adminOnly(protocolHandler.Update))
}

func RegisterRouteRuleRouters(r *router.Router, routeRuleHandler *handler.RouteRuleHandler) {
	adminOnly := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return middlewares.UserIdMiddleware(middlewares.AdminOnlyMiddleware(next))
	}

	r.GET("/ward/api/v1/route-rules", adminOnly(routeRuleHandler.GetAll))
	r.POST("/ward/api/v1/route-rules", adminOnly(routeRuleHandler.Create))
	r.PUT("/ward/api/v1/route-rules/{path:*}", adminOnly(routeRuleHandler.Update))
	r.DELETE("/ward/api/v1/route-rules/{path:*}", adminOnly(routeRuleHandler.Delete))
}

func RegisterSpecialRouteRouters(r *router.Router, specialRouteHandler *handler.SpecialRouteHandler) {
	adminOnly := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return middlewares.UserIdMiddleware(middlewares.AdminOnlyMiddleware(next))
	}

	r.GET("/ward/api/v1/special-routes/login", adminOnly(specialRouteHandler.GetByType))
	r.POST("/ward/api/v1/special-routes/login", adminOnly(specialRouteHandler.Create))
	r.PUT("/ward/api/v1/special-routes/login/{path:*}", adminOnly(specialRouteHandler.Update))
	r.DELETE("/ward/api/v1/special-routes/login/{path:*}", adminOnly(specialRouteHandler.Delete))

	r.GET("/ward/api/v1/special-routes/register", adminOnly(specialRouteHandler.GetByType))
	r.POST("/ward/api/v1/special-routes/register", adminOnly(specialRouteHandler.Create))
	r.PUT("/ward/api/v1/special-routes/register/{path:*}", adminOnly(specialRouteHandler.Update))
	r.DELETE("/ward/api/v1/special-routes/register/{path:*}", adminOnly(specialRouteHandler.Delete))
}
