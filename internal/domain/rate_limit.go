package domain

import (
	"errors"
	"time"
)

type RateLimitSettings struct {
	ID                int
	RequestsPerSecond float64
	Burst             int
	Progressive       bool
	UpdatedAt         time.Time
	CreatedAt         time.Time
}

func NewRateLimitSettings(requestsPerSecond float64, burst int, progressive bool) (*RateLimitSettings, error) {
	if requestsPerSecond <= 0 {
		return nil, errors.New("requests_per_second must be greater than 0")
	}
	if burst <= 0 {
		return nil, errors.New("burst must be greater than 0")
	}

	return &RateLimitSettings{
		ID:                1,
		RequestsPerSecond: requestsPerSecond,
		Burst:             burst,
		Progressive:       progressive,
	}, nil
}
