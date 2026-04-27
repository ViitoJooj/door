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

type ProtocolSettingsHandler struct {
	service *services.ProtocolSettingsService
}

func NewProtocolSettingsHandler(service *services.ProtocolSettingsService) *ProtocolSettingsHandler {
	return &ProtocolSettingsHandler{
		service: service,
	}
}

func (h *ProtocolSettingsHandler) Get(ctx *fasthttp.RequestCtx) {
	settings, err := h.service.Get()
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: "internal error."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	output := dtos.ProtocolSettingsOutput{
		Success: true,
		Message: "ok.",
		Data: dtos.ProtocolSettingsData{
			ID:              settings.ID,
			AllowedProtocol: settings.AllowedProtocol,
			ApplyScope:      settings.ApplyScope,
			UpdatedAt:       settings.UpdatedAt,
			CreatedAt:       settings.CreatedAt,
		},
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *ProtocolSettingsHandler) Update(ctx *fasthttp.RequestCtx) {
	var input dtos.ProtocolSettingsInput
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: "invalid json."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	settings, err := h.service.Update(input.AllowedProtocol, input.ApplyScope)
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: err.Error()}
		statusCode := fasthttp.StatusBadGateway
		if err.Error() == "allowed_protocol must be one of: http, https, both" || err.Error() == "apply_scope must be one of: all, external" {
			statusCode = fasthttp.StatusBadRequest
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(statusCode)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	middlewares.UpdateAllowedProtocol(settings.AllowedProtocol)
	middlewares.UpdateConfigApplyScope(settings.ApplyScope)

	output := dtos.ProtocolSettingsOutput{
		Success: true,
		Message: "protocol mode updated.",
		Data: dtos.ProtocolSettingsData{
			ID:              settings.ID,
			AllowedProtocol: settings.AllowedProtocol,
			ApplyScope:      settings.ApplyScope,
			UpdatedAt:       settings.UpdatedAt,
			CreatedAt:       settings.CreatedAt,
		},
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}
