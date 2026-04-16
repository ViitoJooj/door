package dtos

import dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"

// Register
type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterOutput struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    dto_utils.UserData `json:"data"`
}

// Login
type LoginInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    dto_utils.UserData `json:"data"`
	Token   string             `json:"token"`
}
