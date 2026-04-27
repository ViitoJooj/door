package dtos

import "time"

type RateLimitInput struct {
	RequestsPerSecond float64 `json:"requests_per_second"`
	Burst             int     `json:"burst"`
	Progressive       bool    `json:"progressive_rate_limit"`
}

type RateLimitData struct {
	ID                int       `json:"id"`
	RequestsPerSecond float64   `json:"requests_per_second"`
	Burst             int       `json:"burst"`
	Progressive       bool      `json:"progressive_rate_limit"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type RateLimitOutput struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    RateLimitData `json:"data"`
}
