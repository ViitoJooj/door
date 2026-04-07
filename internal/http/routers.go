package http

import (
	"github.com/ViitoJooj/door/internal/http/handler"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func SetupRouter(userController *handler.AuthHandler) fasthttp.RequestHandler {
	r := router.New()

	r.POST("/api/v1/auth/register", userController.Register)
	r.POST("/api/v1/auth/login", userController.Login)
	r.GET("/api/v1/auth/token", userController.Token)
	r.POST("/api/v1/auth/logout", userController.Logout)

	return r.Handler
}
