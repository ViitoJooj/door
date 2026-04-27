package handler

import (
	"encoding/json"
	"log"

	"github.com/ViitoJooj/ward/internal/http/dtos"
	dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"
	"github.com/ViitoJooj/ward/internal/http/middlewares"
	"github.com/ViitoJooj/ward/internal/services"
	"github.com/valyala/fasthttp"
)

type RateLimitHandler struct {
	rateLimitService *services.RateLimitService
}

func NewRateLimitHandler(rateLimitService *services.RateLimitService) *RateLimitHandler {
	return &RateLimitHandler{
		rateLimitService: rateLimitService,
	}
}

func (h *RateLimitHandler) Get(ctx *fasthttp.RequestCtx) {
	settings, err := h.rateLimitService.Get()
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

	output := dtos.RateLimitOutput{
		Success: true,
		Message: "ok.",
		Data: dtos.RateLimitData{
			ID:                settings.ID,
			RequestsPerSecond: settings.RequestsPerSecond,
			Burst:             settings.Burst,
			Progressive:       settings.Progressive,
			UpdatedAt:         settings.UpdatedAt,
			CreatedAt:         settings.CreatedAt,
		},
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *RateLimitHandler) Update(ctx *fasthttp.RequestCtx) {
	var input dtos.RateLimitInput
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println(err)
		output := dto_utils.Error{
			Success: false,
			Message: "invalid json.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	settings, err := h.rateLimitService.Update(input.RequestsPerSecond, input.Burst, input.Progressive)
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{
			Success: false,
			Message: err.Error(),
		}
		statusCode := fasthttp.StatusBadGateway
		if err.Error() == "requests_per_second must be greater than 0" || err.Error() == "burst must be greater than 0" {
			statusCode = fasthttp.StatusBadRequest
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(statusCode)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	middlewares.UpdateRateLimitConfig(settings.RequestsPerSecond, settings.Burst, settings.Progressive)

	output := dtos.RateLimitOutput{
		Success: true,
		Message: "rate limit updated.",
		Data: dtos.RateLimitData{
			ID:                settings.ID,
			RequestsPerSecond: settings.RequestsPerSecond,
			Burst:             settings.Burst,
			Progressive:       settings.Progressive,
			UpdatedAt:         settings.UpdatedAt,
			CreatedAt:         settings.CreatedAt,
		},
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}
