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
	"github.com/valyala/fasthttp"
)

func main() {
	dotenv.GetEnv()
	database.Conn()

	authRepo := repository.NewSQLiteUserRepository(database.DB)
	log := logger.NewLogger(os.Stdout)
	authService := services.NewAuthService(authRepo, log)
	authHandler := handler.NewAuthHandler(authService)

	r := httpx.SetupRouter(authHandler)
	handlerWithCors := middlewares.CorsMiddleware(r)
	fasthttp.ListenAndServe(":7171", handlerWithCors)
}
