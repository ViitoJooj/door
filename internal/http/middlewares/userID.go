package middlewares

import (
	"encoding/json"

	dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"
	"github.com/ViitoJooj/ward/pkg/jwtTokens"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
)

func UserIdMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		accessToken := string(ctx.Request.Header.Cookie("access_token"))
		if accessToken == "" {
			output := dto_utils.Error{
				Success: false,
				Message: "access token missing",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}

		token, err := jwtTokens.ValidateAccessToken(accessToken)
		if err != nil {
			output := dto_utils.Error{
				Success: false,
				Message: "invalid access token",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			output := dto_utils.Error{
				Success: false,
				Message: "invalid token claims",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}

		userIdFloat, ok := claims["user_id"].(float64)
		if !ok {
			output := dto_utils.Error{
				Success: false,
				Message: "user_id not found",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}

		ctx.SetUserValue("userId", int(userIdFloat))
		next(ctx)
	}
}
