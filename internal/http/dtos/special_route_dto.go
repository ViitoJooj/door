package dtos

import "time"

type SpecialRouteInput struct {
	Path                string `json:"path"`
	MaxDistinctRequests int    `json:"max_distinct_requests"`
	WindowSeconds       int    `json:"window_seconds"`
	BanSeconds          int    `json:"ban_seconds"`
	Enabled             bool   `json:"enabled"`
}

type SpecialRouteData struct {
	ID                  int       `json:"id"`
	RouteType           string    `json:"route_type"`
	Path                string    `json:"path"`
	MaxDistinctRequests int       `json:"max_distinct_requests"`
	WindowSeconds       int       `json:"window_seconds"`
	BanSeconds          int       `json:"ban_seconds"`
	Enabled             bool      `json:"enabled"`
	CreatedBy           int       `json:"created_by"`
	UpdatedBy           int       `json:"updated_by"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type SpecialRouteOutput struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    SpecialRouteData `json:"data"`
}

type SpecialRouteListOutput struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    []SpecialRouteData `json:"data"`
}
