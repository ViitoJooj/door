package http

import (
	"github.com/ViitoJooj/ward/internal/http/handler"
	"github.com/ViitoJooj/ward/internal/http/middlewares"
	"github.com/fasthttp/router"
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
