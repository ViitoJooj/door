package dtos

import "time"

// Utils

type ApplicationData struct {
	ID         int64     `json:"id"`
	Url        string    `json:"url"`
	Country    string    `json:"country"`
	Created_by UserData  `json:"created_by"`
	Updated_at time.Time `json:"updated_at"`
	Created_at time.Time `json:"created_at"`
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
