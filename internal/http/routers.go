package http

import (
	"time"

	"github.com/ViitoJooj/door/internal/http/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userController *handler.AuthHandler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	users := r.Group("/api/v1/auth")
	{
		users.POST("/register", userController.Register)
		users.POST("/login", userController.Login)
		users.GET("/token", userController.Token)
		users.POST("/logout", userController.Logout)

	}

	return r
}
