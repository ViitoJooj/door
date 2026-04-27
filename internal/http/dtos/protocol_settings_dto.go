package dtos

import "time"

type ProtocolSettingsInput struct {
	AllowedProtocol string `json:"allowed_protocol"`
	ApplyScope      string `json:"apply_scope"`
}

type ProtocolSettingsData struct {
	ID              int       `json:"id"`
	AllowedProtocol string    `json:"allowed_protocol"`
	ApplyScope      string    `json:"apply_scope"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`
}

type ProtocolSettingsOutput struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    ProtocolSettingsData `json:"data"`
}
