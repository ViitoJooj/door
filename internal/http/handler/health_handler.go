package handler

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/ViitoJooj/ward/internal/http/dtos"
	dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"
	"github.com/ViitoJooj/ward/internal/services"
	"github.com/valyala/fasthttp"
)

type HealthHandler struct {
	healthService *services.HealthService
}

func NewHealthHandler(healthService *services.HealthService) *HealthHandler {
	return &HealthHandler{healthService: healthService}
}

func (h *HealthHandler) GetOverview(ctx *fasthttp.RequestCtx) {
	windowMinutes, err := parsePositiveIntArg(ctx, "window_minutes")
	if err != nil {
		writeBadRequest(ctx, "invalid window_minutes.")
		return
	}

	overview, svcErr := h.healthService.GetOverview(windowMinutes)
	if svcErr != nil {
		log.Println(svcErr)
		writeInternalError(ctx)
		return
	}

	output := dtos.HealthOverviewOutput{
		Success: true,
		Message: "ok.",
		Data: dtos.HealthOverviewData{
			Status:            string(overview.Status),
			WindowMinutes:     overview.WindowMinutes,
			GeneratedAt:       overview.GeneratedAt,
			TotalRequests:     overview.TotalRequests,
			ServerErrors:      overview.ServerErrors,
			ClientErrors:      overview.ClientErrors,
			Availability:      overview.Availability,
			ServerErrorRate:   overview.ServerErrorRate,
			ClientErrorRate:   overview.ClientErrorRate,
			AverageLatencyMs:  overview.AverageLatencyMs,
			P95LatencyMs:      overview.P95LatencyMs,
			RequestsPerMinute: overview.RequestsPerMinute,
			UniqueIPs:         overview.UniqueIPs,
			UniquePaths:       overview.UniquePaths,
		},
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *HealthHandler) GetRoutes(ctx *fasthttp.RequestCtx) {
	windowMinutes, err := parsePositiveIntArg(ctx, "window_minutes")
	if err != nil {
		writeBadRequest(ctx, "invalid window_minutes.")
		return
	}

	limit, err := parsePositiveIntArg(ctx, "limit")
	if err != nil {
		writeBadRequest(ctx, "invalid limit.")
		return
	}

	stats, svcErr := h.healthService.GetRouteStats(windowMinutes, limit)
	if svcErr != nil {
		log.Println(svcErr)
		writeInternalError(ctx)
		return
	}

	data := make([]dtos.HealthRouteData, 0, len(stats))
	for _, item := range stats {
		data = append(data, dtos.HealthRouteData{
			Method:            item.Method,
			Path:              item.Path,
			Status:            string(item.Status),
			WindowMinutes:     item.WindowMinutes,
			LastSeenAt:        item.LastSeenAt,
			RequestCount:      item.RequestCount,
			ServerErrors:      item.ServerErrors,
			ClientErrors:      item.ClientErrors,
			Availability:      item.Availability,
			ServerErrorRate:   item.ServerErrorRate,
			ClientErrorRate:   item.ClientErrorRate,
			AverageLatencyMs:  item.AverageLatencyMs,
			P95LatencyMs:      item.P95LatencyMs,
			RequestsPerMinute: item.RequestsPerMinute,
		})
	}

	output := dtos.HealthRoutesOutput{
		Success: true,
		Message: "ok.",
		Data:    data,
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func parsePositiveIntArg(ctx *fasthttp.RequestCtx, key string) (int, error) {
	raw := string(ctx.QueryArgs().Peek(key))
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}
	if value <= 0 {
		return 0, errors.New("invalid value")
	}
	return value, nil
}

func writeBadRequest(ctx *fasthttp.RequestCtx, message string) {
	output := dto_utils.Error{
		Success: false,
		Message: message,
	}
	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func writeInternalError(ctx *fasthttp.RequestCtx) {
	output := dto_utils.Error{
		Success: false,
		Message: "internal error.",
	}
	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusBadGateway)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}
