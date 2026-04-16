package dtos

import (
	"time"

	dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"
)

type ApplicationData struct {
	ID         int                `json:"id"`
	Url        string             `json:"url"`
	Country    string             `json:"country"`
	Created_by dto_utils.UserData `json:"created_by"`
	Updated_at time.Time          `json:"updated_at"`
	Created_at time.Time          `json:"created_at"`
}

type ApplicationOutput struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    ApplicationData `json:"data"`
}

type ApplicationInput struct {
	Url     string `json:"url"`
	Country string `json:"country"`
}

// delete dto
