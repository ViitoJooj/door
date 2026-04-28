package dtos

import "time"

type HealthOverviewData struct {
	Status            string    `json:"status"`
	WindowMinutes     int       `json:"window_minutes"`
	GeneratedAt       time.Time `json:"generated_at"`
	TotalRequests     int       `json:"total_requests"`
	ServerErrors      int       `json:"server_errors"`
	ClientErrors      int       `json:"client_errors"`
	Availability      float64   `json:"availability"`
	ServerErrorRate   float64   `json:"server_error_rate"`
	ClientErrorRate   float64   `json:"client_error_rate"`
	AverageLatencyMs  float64   `json:"average_latency_ms"`
	P95LatencyMs      int64     `json:"p95_latency_ms"`
	RequestsPerMinute float64   `json:"requests_per_minute"`
	UniqueIPs         int       `json:"unique_ips"`
	UniquePaths       int       `json:"unique_paths"`
}

type HealthRouteData struct {
	Method            string    `json:"method"`
	Path              string    `json:"path"`
	Status            string    `json:"status"`
	WindowMinutes     int       `json:"window_minutes"`
	LastSeenAt        time.Time `json:"last_seen_at"`
	RequestCount      int       `json:"request_count"`
	ServerErrors      int       `json:"server_errors"`
	ClientErrors      int       `json:"client_errors"`
	Availability      float64   `json:"availability"`
	ServerErrorRate   float64   `json:"server_error_rate"`
	ClientErrorRate   float64   `json:"client_error_rate"`
	AverageLatencyMs  float64   `json:"average_latency_ms"`
	P95LatencyMs      int64     `json:"p95_latency_ms"`
	RequestsPerMinute float64   `json:"requests_per_minute"`
}

type HealthOverviewOutput struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    HealthOverviewData `json:"data"`
}

type HealthRoutesOutput struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    []HealthRouteData `json:"data"`
}
