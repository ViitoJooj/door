package domain

import "time"

type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusUnknown   HealthStatus = "unknown"
)

type HealthOverview struct {
	Status            HealthStatus
	WindowMinutes     int
	GeneratedAt       time.Time
	TotalRequests     int
	ServerErrors      int
	ClientErrors      int
	Availability      float64
	ServerErrorRate   float64
	ClientErrorRate   float64
	AverageLatencyMs  float64
	P95LatencyMs      int64
	RequestsPerMinute float64
	UniqueIPs         int
	UniquePaths       int
}

type HealthRouteStat struct {
	Method            string
	Path              string
	Status            HealthStatus
	WindowMinutes     int
	LastSeenAt        time.Time
	RequestCount      int
	ServerErrors      int
	ClientErrors      int
	Availability      float64
	ServerErrorRate   float64
	ClientErrorRate   float64
	AverageLatencyMs  float64
	P95LatencyMs      int64
	RequestsPerMinute float64
}
