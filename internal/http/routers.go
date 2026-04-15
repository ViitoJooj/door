package http

import (
	"github.com/ViitoJooj/door/internal/http/handler"
	"github.com/ViitoJooj/door/internal/http/middlewares"
	"github.com/fasthttp/router"
)

func RegisterAuthRoutes(r *router.Router, userController *handler.AuthHandler) {
	r.POST("/api/v1/auth/register", userController.Register)
	r.POST("/api/v1/auth/login", userController.Login)
	r.GET("/api/v1/auth/token", userController.Token)
	r.POST("/api/v1/auth/logout", userController.Logout)
}

func RegisterProxyRoutes(r *router.Router, proxyController *handler.ProxyHandler) {
	r.ANY("/proxy/{path:*}", proxyController.Proxy)
}

func RegisterRequestLogRoutes(r *router.Router, requestLogController *handler.RequestLogHandler) {
	r.GET("/api/v1/logs", middlewares.UserIdMiddleware(requestLogController.GetAll))
}

func RegisterApplicationRouters(r *router.Router, applicationController *handler.ApplicationHandler) {
	r.GET("/api/v1/applications", middlewares.UserIdMiddleware(applicationController.GetAll))
	r.GET("/api/v1/applications/{path:*}", middlewares.UserIdMiddleware(applicationController.GetByID))
	r.POST("/api/v1/applications", middlewares.UserIdMiddleware(applicationController.Create))
	r.DELETE("/api/v1/applications/{path:*}", middlewares.UserIdMiddleware(applicationController.DeleteById))
}
