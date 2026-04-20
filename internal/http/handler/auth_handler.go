package handler

import (
	"encoding/json"
	"log"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/internal/http/dtos"
	dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"
	"github.com/ViitoJooj/ward/internal/services"
	"github.com/ViitoJooj/ward/pkg/ip"
	"github.com/ViitoJooj/ward/pkg/jwtTokens"
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

	userIP := ip.GetIP(ctx)
	user, accessToken, refreshToken, err := h.authService.Login(input.Username, input.Email, input.Password, userIP)
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

	var accessCookie fasthttp.Cookie
	accessCookie.SetKey("access_token")
	accessCookie.SetValue(accessToken)
	accessCookie.SetHTTPOnly(true)
	accessCookie.SetPath("/")
	accessCookie.SetSecure(false)

	var refreshCookie fasthttp.Cookie
	refreshCookie.SetKey("refresh_token")
	refreshCookie.SetValue(refreshToken)
	refreshCookie.SetHTTPOnly(true)
	refreshCookie.SetPath("/ward/api/v1/auth/token")
	refreshCookie.SetSecure(false)

	ctx.Response.Header.SetCookie(&accessCookie)
	ctx.Response.Header.SetCookie(&refreshCookie)

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *AuthHandler) Login(ctx *fasthttp.RequestCtx) {
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
	user, accessToken, refreshToken, err := h.authService.Login(input.Username, input.Email, input.Password, userIP)
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
	}

	var accessCookie fasthttp.Cookie
	accessCookie.SetKey("access_token")
	accessCookie.SetValue(accessToken)
	accessCookie.SetHTTPOnly(true)
	accessCookie.SetPath("/")
	accessCookie.SetSecure(false)

	var refreshCookie fasthttp.Cookie
	refreshCookie.SetKey("refresh_token")
	refreshCookie.SetValue(refreshToken)
	refreshCookie.SetHTTPOnly(true)
	refreshCookie.SetPath("/ward/api/v1/auth/token")
	refreshCookie.SetSecure(false)

	ctx.Response.Header.SetCookie(&accessCookie)
	ctx.Response.Header.SetCookie(&refreshCookie)

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *AuthHandler) Token(ctx *fasthttp.RequestCtx) {
	refreshToken := string(ctx.Request.Header.Cookie("refresh_token"))
	if refreshToken == "" {
		output := dto_utils.Error{
			Success: false,
			Message: "refresh token not found",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	user, err := h.authService.Token(refreshToken, true)
	if err != nil {
		output := dto_utils.Error{
			Success: false,
			Message: "invalid refresh token",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	newAccessToken, err := jwtTokens.GenerateAccessToken(user.ID)
	if err != nil {
		output := dto_utils.Error{
			Success: false,
			Message: "internal error",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	var accessCookie fasthttp.Cookie
	accessCookie.SetKey("access_token")
	accessCookie.SetValue(newAccessToken)
	accessCookie.SetHTTPOnly(true)
	accessCookie.SetPath("/")
	accessCookie.SetSecure(false)

	ctx.Response.Header.SetCookie(&accessCookie)

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBodyString(`{"success":true}`)
}

func (h *AuthHandler) Logout(ctx *fasthttp.RequestCtx) {
	var accessCookie fasthttp.Cookie
	accessCookie.SetKey("access_token")
	accessCookie.SetValue("")
	accessCookie.SetExpire(fasthttp.CookieExpireDelete)
	accessCookie.SetPath("/")

	var refreshCookie fasthttp.Cookie
	refreshCookie.SetKey("refresh_token")
	refreshCookie.SetValue("")
	refreshCookie.SetExpire(fasthttp.CookieExpireDelete)
	refreshCookie.SetPath("/ward/api/v1/token")

	ctx.Response.Header.SetCookie(&accessCookie)
	ctx.Response.Header.SetCookie(&refreshCookie)

	res, _ := json.Marshal(map[string]any{
		"success": true,
		"message": "Logout successful",
	})

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}
