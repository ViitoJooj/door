package handler

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/ViitoJooj/door/internal/domain"
	"github.com/ViitoJooj/door/internal/http/dtos"
	dto_utils "github.com/ViitoJooj/door/internal/http/dtos/utils"
	"github.com/ViitoJooj/door/internal/services"
	"github.com/ViitoJooj/door/pkg/ip"
	"github.com/valyala/fasthttp"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(ctx *fasthttp.RequestCtx) {
	var input dtos.RegisterInput

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println("invalid json.")
		output := dto_utils.Error{
			Success: false,
			Message: "invalid json.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	user := &domain.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	createdUser, err := h.authService.Register(user)
	if err != nil {
		log.Println("internal error.")
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	output := dtos.RegisterOutput{
		Success: true,
		Message: "User created.",
		Data: dto_utils.UserData{
			Username:   createdUser.Username,
			Email:      createdUser.Email,
			Updated_at: createdUser.Updated_at.String(),
			Created_at: createdUser.Created_at.String(),
		},
	}

	res, _ := json.Marshal(output)

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (c *AuthHandler) Login(ctx *fasthttp.RequestCtx) {
	var input dtos.LoginInput

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println("invalid json.")
		output := dto_utils.Error{
			Success: false,
			Message: "invalid json.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	userIP := ip.GetIP(ctx)
	user, token, err := c.authService.Login(input.Username, input.Email, input.Password, userIP)
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	output := dtos.LoginOutput{
		Success: true,
		Message: "Login successful.",
		Data: dto_utils.UserData{
			ID:         user.ID,
			Username:   user.Username,
			Email:      user.Email,
			Updated_at: user.Updated_at.String(),
			Created_at: user.Created_at.String(),
		},
		Token: token,
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (c *AuthHandler) Token(ctx *fasthttp.RequestCtx) {
	authHeader := string(ctx.Request.Header.Peek("Authorization"))

	if authHeader == "" {
		log.Println("Authorization header not found")
		output := dto_utils.Error{
			Success: false,
			Message: "Authorization header not found",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		log.Println("Invalid Authorization format")
		output := dto_utils.Error{
			Success: false,
			Message: "Invalid Authorization format",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	_, err := c.authService.Token(tokenString)
	if err != nil {
		log.Println("invalid token.")
		output := dto_utils.Error{
			Success: false,
			Message: "invalid token.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	var cookie fasthttp.Cookie
	cookie.SetKey("token")
	cookie.SetValue(tokenString)
	cookie.SetHTTPOnly(true)
	cookie.SetPath("/")
	cookie.SetSecure(false)

	ctx.Response.Header.SetCookie(&cookie)

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBodyString(`{"success":true}`)
}

func (c *AuthHandler) Logout(ctx *fasthttp.RequestCtx) {
	var cookie fasthttp.Cookie
	cookie.SetKey("token")
	cookie.SetValue("")
	cookie.SetExpire(fasthttp.CookieExpireDelete)
	cookie.SetPath("/")

	ctx.Response.Header.SetCookie(&cookie)

	res, _ := json.Marshal(map[string]any{
		"success": true,
		"message": "Logout successful",
	})

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}
