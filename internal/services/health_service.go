package services

import (
	"sort"
	"strings"
	"time"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/internal/repository"
)

const (
	defaultHealthWindowMinutes = 15
	maxHealthWindowMinutes     = 24 * 60
	defaultHealthRouteLimit    = 20
	maxHealthRouteLimit        = 100
	maxHealthLogScan           = 20000
)

type HealthService struct {
	logRepo repository.RequestLogRepository
}

func NewHealthService(logRepo repository.RequestLogRepository) *HealthService {
	return &HealthService{logRepo: logRepo}
}

func (s *HealthService) GetOverview(windowMinutes int) (*domain.HealthOverview, error) {
	window := normalizeWindowMinutes(windowMinutes)
	logs, err := s.logRepo.ListRequestLogsSince(time.Now().Add(-time.Duration(window)*time.Minute), maxHealthLogScan)
	if err != nil {
		return nil, err
	}

	overview := &domain.HealthOverview{
		WindowMinutes: window,
		GeneratedAt:   time.Now().UTC(),
	}

	if len(logs) == 0 {
		overview.Status = domain.HealthStatusUnknown
		return overview, nil
	}

	var (
		totalLatency int64
		latencies    = make([]int64, 0, len(logs))
		uniqueIPs    = map[string]struct{}{}
		uniquePaths  = map[string]struct{}{}
	)

	for _, entry := range logs {
		if ignoreHealthPath(entry.Path) {
			continue
		}
		overview.TotalRequests++
		totalLatency += entry.ResponseTimeMs
		latencies = append(latencies, entry.ResponseTimeMs)

		if entry.StatusCode >= 500 {
			overview.ServerErrors++
		} else if entry.StatusCode >= 400 {
			overview.ClientErrors++
		}

		if entry.IP != "" {
			uniqueIPs[entry.IP] = struct{}{}
		}
		uniquePaths[entry.Path] = struct{}{}
	}

	overview.UniqueIPs = len(uniqueIPs)
	overview.UniquePaths = len(uniquePaths)
	if overview.TotalRequests == 0 {
		overview.Status = domain.HealthStatusUnknown
		return overview, nil
	}
	overview.AverageLatencyMs = float64(totalLatency) / float64(overview.TotalRequests)
	overview.P95LatencyMs = percentile95(latencies)
	overview.ServerErrorRate = percent(overview.ServerErrors, overview.TotalRequests)
	overview.ClientErrorRate = percent(overview.ClientErrors, overview.TotalRequests)
	overview.Availability = 100 - overview.ServerErrorRate
	overview.RequestsPerMinute = float64(overview.TotalRequests) / float64(window)
	overview.Status = evaluateStatus(overview.ServerErrorRate, overview.AverageLatencyMs, overview.Availability)

	return overview, nil
}

func (s *HealthService) GetRouteStats(windowMinutes int, limit int) ([]*domain.HealthRouteStat, error) {
	window := normalizeWindowMinutes(windowMinutes)
	logs, err := s.logRepo.ListRequestLogsSince(time.Now().Add(-time.Duration(window)*time.Minute), maxHealthLogScan)
	if err != nil {
		return nil, err
	}

	type routeAccumulator struct {
		method       string
		path         string
		total        int
		serverErrors int
		clientErrors int
		latencySum   int64
		latencies    []int64
		lastSeen     time.Time
	}

	byRoute := map[string]*routeAccumulator{}

	for _, entry := range logs {
		if ignoreHealthPath(entry.Path) {
			continue
		}
		key := strings.ToUpper(entry.Method) + "|" + entry.Path
		acc, ok := byRoute[key]
		if !ok {
			acc = &routeAccumulator{
				method:    strings.ToUpper(entry.Method),
				path:      entry.Path,
				latencies: make([]int64, 0, 32),
			}
			byRoute[key] = acc
		}

		acc.total++
		acc.latencySum += entry.ResponseTimeMs
		acc.latencies = append(acc.latencies, entry.ResponseTimeMs)
		if entry.StatusCode >= 500 {
			acc.serverErrors++
		} else if entry.StatusCode >= 400 {
			acc.clientErrors++
		}
		if entry.CreatedAt.After(acc.lastSeen) {
			acc.lastSeen = entry.CreatedAt
		}
	}

	stats := make([]*domain.HealthRouteStat, 0, len(byRoute))
	for _, acc := range byRoute {
		avgLatency := 0.0
		if acc.total > 0 {
			avgLatency = float64(acc.latencySum) / float64(acc.total)
		}

		serverErrorRate := percent(acc.serverErrors, acc.total)
		availability := 100 - serverErrorRate
		clientErrorRate := percent(acc.clientErrors, acc.total)

		stats = append(stats, &domain.HealthRouteStat{
			Method:            acc.method,
			Path:              acc.path,
			Status:            evaluateStatus(serverErrorRate, avgLatency, availability),
			WindowMinutes:     window,
			LastSeenAt:        acc.lastSeen,
			RequestCount:      acc.total,
			ServerErrors:      acc.serverErrors,
			ClientErrors:      acc.clientErrors,
			Availability:      availability,
			ServerErrorRate:   serverErrorRate,
			ClientErrorRate:   clientErrorRate,
			AverageLatencyMs:  avgLatency,
			P95LatencyMs:      percentile95(acc.latencies),
			RequestsPerMinute: float64(acc.total) / float64(window),
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		left := stats[i]
		right := stats[j]
		if statusRank(left.Status) != statusRank(right.Status) {
			return statusRank(left.Status) > statusRank(right.Status)
		}
		if left.ServerErrorRate != right.ServerErrorRate {
			return left.ServerErrorRate > right.ServerErrorRate
		}
		if left.AverageLatencyMs != right.AverageLatencyMs {
			return left.AverageLatencyMs > right.AverageLatencyMs
		}
		return left.RequestCount > right.RequestCount
	})

	normalizedLimit := normalizeRouteLimit(limit)
	if len(stats) > normalizedLimit {
		stats = stats[:normalizedLimit]
	}

	return stats, nil
}

func normalizeWindowMinutes(window int) int {
	if window <= 0 {
		return defaultHealthWindowMinutes
	}
	if window > maxHealthWindowMinutes {
		return maxHealthWindowMinutes
	}
	return window
}

func normalizeRouteLimit(limit int) int {
	if limit <= 0 {
		return defaultHealthRouteLimit
	}
	if limit > maxHealthRouteLimit {
		return maxHealthRouteLimit
	}
	return limit
}

func percent(value int, total int) float64 {
	if total == 0 {
		return 0
	}
	return (float64(value) / float64(total)) * 100
}

func percentile95(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}
	sorted := append([]int64(nil), values...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	index := int(float64(len(sorted)-1) * 0.95)
	return sorted[index]
}

func evaluateStatus(serverErrorRate float64, averageLatencyMs float64, availability float64) domain.HealthStatus {
	if serverErrorRate >= 10 || averageLatencyMs >= 1500 || availability < 95 {
		return domain.HealthStatusUnhealthy
	}
	if serverErrorRate >= 3 || averageLatencyMs >= 700 || availability < 98 {
		return domain.HealthStatusDegraded
	}
	return domain.HealthStatusHealthy
}

func statusRank(status domain.HealthStatus) int {
	switch status {
	case domain.HealthStatusUnhealthy:
		return 4
	case domain.HealthStatusDegraded:
		return 3
	case domain.HealthStatusUnknown:
		return 2
	default:
		return 1
	}
}

func ignoreHealthPath(path string) bool {
	return path == "/ward/api/v1/health" ||
		path == "/ward/api/v1/health/routes" ||
		path == "/ward/api/v1/logs"
}
