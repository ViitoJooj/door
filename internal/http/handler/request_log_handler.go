package handler

import (
	"encoding/json"
	"log"

	"github.com/ViitoJooj/door/internal/http/dtos"
	dto_utils "github.com/ViitoJooj/door/internal/http/dtos/utils"
	"github.com/ViitoJooj/door/internal/services"
	"github.com/valyala/fasthttp"
)

type RequestLogHandler struct {
	logService *services.RequestLogService
}

func NewRequestLogHandler(logService *services.RequestLogService) *RequestLogHandler {
	return &RequestLogHandler{
		logService: logService,
	}
}

func (h *RequestLogHandler) GetAll(ctx *fasthttp.RequestCtx) {
	logs, err := h.logService.GetAll()
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	data := make([]dtos.RequestLogData, 0, len(logs))
	for _, entry := range logs {
		data = append(data, dtos.RequestLogData{
			ID:             entry.ID,
			Method:         entry.Method,
			Path:           entry.Path,
			QueryString:    entry.QueryString,
			StatusCode:     entry.StatusCode,
			ResponseTimeMs: entry.ResponseTimeMs,
			IP:             entry.IP,
			Country:        entry.Country,
			UserAgent:      entry.UserAgent,
			Referer:        entry.Referer,
			RequestSize:    entry.RequestSize,
			ResponseSize:   entry.ResponseSize,
			Internal:       entry.Internal,
			CreatedAt:      entry.CreatedAt,
		})
	}

	output := dtos.RequestLogListOutput{
		Success: true,
		Message: "ok.",
		Data:    data,
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}
