package dtos

import dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"

type CorsInput struct {
	Origin string `json:"origin"`
}

type CorsData struct {
	ID         int                `json:"id"`
	Origin     string             `json:"origin"`
	Created_by dto_utils.UserData `json:"created_by"`
}

type CorsOutput struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    CorsData `json:"data"`
}
