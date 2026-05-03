package domain

import (
	"errors"
	"strings"
	"time"
)

const (
	SpecialRouteTypeLogin    = "login"
	SpecialRouteTypeRegister = "register"
)

type SpecialRouteRule struct {
	ID                  int
	RouteType           string
	Path                string
	MaxDistinctRequests int
	WindowSeconds       int
	BanSeconds          int
	Enabled             bool
	CreatedBy           int
	UpdatedBy           int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func NormalizeSpecialRouteType(routeType string) (string, error) {
	value := strings.ToLower(strings.TrimSpace(routeType))
	if value == "" {
		return "", errors.New("route_type is required")
	}
	return value, nil
}

func NormalizeSpecialRoutePath(path string) (string, error) {
	normalized := strings.TrimSpace(path)
	if normalized == "" {
		return "", errors.New("path is required")
	}
	if !strings.HasPrefix(normalized, "/") {
		normalized = "/" + normalized
	}
	if len(normalized) > 1 {
		normalized = strings.TrimSuffix(normalized, "/")
	}
	return normalized, nil
}

func NewSpecialRouteRule(routeType string, path string, maxDistinctRequests int, windowSeconds int, banSeconds int, enabled bool) (*SpecialRouteRule, error) {
	normalizedType, err := NormalizeSpecialRouteType(routeType)
	if err != nil {
		return nil, err
	}
	normalizedPath, err := NormalizeSpecialRoutePath(path)
	if err != nil {
		return nil, err
	}

	if maxDistinctRequests <= 0 {
		return nil, errors.New("max_distinct_requests must be greater than 0")
	}
	if windowSeconds <= 0 {
		return nil, errors.New("window_seconds must be greater than 0")
	}
	if banSeconds <= 0 {
		return nil, errors.New("ban_seconds must be greater than 0")
	}

	return &SpecialRouteRule{
		RouteType:           normalizedType,
		Path:                normalizedPath,
		MaxDistinctRequests: maxDistinctRequests,
		WindowSeconds:       windowSeconds,
		BanSeconds:          banSeconds,
		Enabled:             enabled,
	}, nil
}
