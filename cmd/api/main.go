package main

import (
	httpx "github.com/ViitoJooj/door/internal/http"
	"github.com/ViitoJooj/door/internal/http/handler"
	"github.com/ViitoJooj/door/internal/repository"
	"github.com/ViitoJooj/door/internal/services"
	"github.com/ViitoJooj/door/pkg/database"
	"github.com/ViitoJooj/door/pkg/dotenv"
	"github.com/valyala/fasthttp"
)

func main() {
	dotenv.GetEnv()
	database.Conn()

	authRepo := repository.NewSQLiteUserRepository(database.DB)
	authService := services.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

	r := httpx.SetupRouter(authHandler)
	fasthttp.ListenAndServe(":7171", r)
}
