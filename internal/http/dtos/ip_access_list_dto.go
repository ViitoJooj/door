package dtos

import "time"

type IPAccessInput struct {
	IP string `json:"ip"`
}

type IPAccessData struct {
	ID        int       `json:"id"`
	IP        string    `json:"ip"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IPAccessOutput struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    IPAccessData `json:"data"`
}

type IPAccessListOutput struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    []IPAccessData `json:"data"`
}
