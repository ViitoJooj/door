package handler

import (
	"net/http"

	"github.com/ViitoJooj/door/internal/domain"
	"github.com/ViitoJooj/door/internal/http/dtos"
	"github.com/ViitoJooj/door/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (c *AuthHandler) Register(ctx *gin.Context) {
	var input dtos.RegisterInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := &domain.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	createdUser, err := c.authService.Register(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
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

	ctx.JSON(http.StatusCreated, output)
}

func (c *AuthHandler) Login(ctx *gin.Context) {
	var input dtos.LoginInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, token, err := c.authService.Login(input.Username, input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
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

	ctx.JSON(http.StatusOK, output)
}

func (c *AuthHandler) Token(ctx *gin.Context) {
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token not found",
		})
		return
	}

	_, err = c.authService.Token(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token is valid",
	})
}

func (c *AuthHandler) Logout(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}
