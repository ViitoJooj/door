package dtos

import "time"

type RouteRuleInput struct {
	Path              string  `json:"path"`
	Method            string  `json:"method"`
	RateLimitEnabled  bool    `json:"rate_limit_enabled"`
	RateLimitRPS      float64 `json:"rate_limit_rps"`
	RateLimitBurst    int     `json:"rate_limit_burst"`
	TargetURL         string  `json:"target_url"`
	GeoRoutingEnabled bool    `json:"geo_routing_enabled"`
	Enabled           bool    `json:"enabled"`
}

type RouteRuleData struct {
	ID                int       `json:"id"`
	Path              string    `json:"path"`
	Method            string    `json:"method"`
	RateLimitEnabled  bool      `json:"rate_limit_enabled"`
	RateLimitRPS      float64   `json:"rate_limit_rps"`
	RateLimitBurst    int       `json:"rate_limit_burst"`
	TargetURL         string    `json:"target_url"`
	GeoRoutingEnabled bool      `json:"geo_routing_enabled"`
	Enabled           bool      `json:"enabled"`
	CreatedBy         int       `json:"created_by"`
	UpdatedBy         int       `json:"updated_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type RouteRuleOutput struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    RouteRuleData `json:"data"`
}

type RouteRuleListOutput struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    []RouteRuleData `json:"data"`
}
