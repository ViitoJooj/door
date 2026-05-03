package domain

import (
	"errors"
	"strings"
	"time"
)

type RouteRule struct {
	ID                int
	Path              string
	Method            string
	RateLimitEnabled  bool
	RateLimitRPS      float64
	RateLimitBurst    int
	TargetURL         string
	GeoRoutingEnabled bool
	Enabled           bool
	CreatedBy         int
	UpdatedBy         int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewRouteRule(path, method string, rateLimitEnabled bool, rateLimitRPS float64, rateLimitBurst int, targetURL string, geoRoutingEnabled bool, enabled bool) (*RouteRule, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, errors.New("path is required")
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if rateLimitEnabled {
		if rateLimitRPS <= 0 {
			return nil, errors.New("rate_limit_rps must be greater than 0")
		}
		if rateLimitBurst <= 0 {
			return nil, errors.New("rate_limit_burst must be greater than 0")
		}
	}

	return &RouteRule{
		Path:              path,
		Method:            strings.ToUpper(strings.TrimSpace(method)),
		RateLimitEnabled:  rateLimitEnabled,
		RateLimitRPS:      rateLimitRPS,
		RateLimitBurst:    rateLimitBurst,
		TargetURL:         strings.TrimSpace(targetURL),
		GeoRoutingEnabled: geoRoutingEnabled,
		Enabled:           enabled,
	}, nil
}
