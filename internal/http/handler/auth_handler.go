package handler

import (
	"encoding/json"

	"github.com/ViitoJooj/door/internal/domain"
	"github.com/ViitoJooj/door/internal/http/dtos"
	"github.com/ViitoJooj/door/internal/services"
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
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte(`{"error":"invalid json"}`))
		return
	}

	user := &domain.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	createdUser, err := h.authService.Register(user)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}

	output := dtos.RegisterOutput{
		Success: true,
		Message: "User created.",
		Data: dtos.UserData{
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
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte(`{"error":"invalid json"}`))
		return
	}

	user, token, err := c.authService.Login(input.Username, input.Email, input.Password)
	if err != nil {
		res, _ := json.Marshal(map[string]string{
			"error": err.Error(),
		})
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	output := dtos.LoginOutput{
		Success: true,
		Message: "Login successful.",
		Data: dtos.UserData{
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
	tokenString := string(ctx.Request.Header.Cookie("token"))

	if tokenString == "" {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte(`{"error":"Token not found"}`))
		return
	}

	_, err := c.authService.Token(tokenString)
	if err != nil {
		res, _ := json.Marshal(map[string]string{
			"error": err.Error(),
		})
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	res, _ := json.Marshal(map[string]any{
		"success": true,
		"message": "Token is valid",
	})

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
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
