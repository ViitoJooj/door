package middlewares

import (
	"encoding/json"
	"log"
	"strings"

	dto_utils "github.com/ViitoJooj/door/internal/http/dtos/utils"
	"github.com/ViitoJooj/door/pkg/jwtTokens"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
)

func UserIdMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		tokenString := strings.TrimPrefix(string(ctx.Request.Header.Peek("Authorization")), "Bearer ")
		if tokenString == "" {
			log.Println("No token.")
			output := dto_utils.Error{
				Success: false,
				Message: "invalid token.",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}

		token, err := jwtTokens.ValidateToken(tokenString)
		if err != nil {
			log.Println(err)
			output := dto_utils.Error{
				Success: false,
				Message: "internal error.",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Invalid token claims.")
			output := dto_utils.Error{
				Success: false,
				Message: "invalid token claims.",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}

		userIdFloat, ok := claims["user_id"].(float64)
		if !ok {
			log.Println("user_id not found or invalid type")
			output := dto_utils.Error{
				Success: false,
				Message: "user_id not found or invalid type",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}

		userId := int(userIdFloat)

		ctx.SetUserValue("userId", userId)
		next(ctx)
	}
}
