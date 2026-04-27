package handler

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/internal/http/dtos"
	dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"
	"github.com/ViitoJooj/ward/internal/http/middlewares"
	"github.com/ViitoJooj/ward/internal/services"
	"github.com/valyala/fasthttp"
)

type SpecialRouteHandler struct {
	service *services.SpecialRouteService
}

func NewSpecialRouteHandler(service *services.SpecialRouteService) *SpecialRouteHandler {
	return &SpecialRouteHandler{
		service: service,
	}
}

func specialRouteStatus(err error) int {
	if err == nil {
		return fasthttp.StatusBadGateway
	}
	msg := err.Error()
	if msg == "route_type must be one of: login, register" ||
		msg == "path is required" ||
		msg == "max_distinct_requests must be greater than 0" ||
		msg == "window_seconds must be greater than 0" ||
		msg == "ban_seconds must be greater than 0" ||
		msg == "path already exists for this route_type" {
		return fasthttp.StatusBadRequest
	}
	if msg == "special route not found" {
		return fasthttp.StatusNotFound
	}
	return fasthttp.StatusBadGateway
}

func routeTypeFromPath(path string) (string, error) {
	switch {
	case strings.HasPrefix(path, "/ward/api/v1/special-routes/login"):
		return domain.SpecialRouteTypeLogin, nil
	case strings.HasPrefix(path, "/ward/api/v1/special-routes/register"):
		return domain.SpecialRouteTypeRegister, nil
	default:
		return "", errors.New("invalid route_type")
	}
}

func parseSpecialRouteID(path string, routeType string) (int, error) {
	prefix := "/ward/api/v1/special-routes/" + routeType + "/"
	idStr := strings.TrimSpace(strings.ReplaceAll(path, prefix, ""))
	return strconv.Atoi(idStr)
}

func mapSpecialRouteData(rule *domain.SpecialRouteRule) dtos.SpecialRouteData {
	return dtos.SpecialRouteData{
		ID:                  rule.ID,
		RouteType:           rule.RouteType,
		Path:                rule.Path,
		MaxDistinctRequests: rule.MaxDistinctRequests,
		WindowSeconds:       rule.WindowSeconds,
		BanSeconds:          rule.BanSeconds,
		Enabled:             rule.Enabled,
		CreatedBy:           rule.CreatedBy,
		UpdatedBy:           rule.UpdatedBy,
		CreatedAt:           rule.CreatedAt,
		UpdatedAt:           rule.UpdatedAt,
	}
}

func (h *SpecialRouteHandler) GetByType(ctx *fasthttp.RequestCtx) {
	routeType, err := routeTypeFromPath(string(ctx.Path()))
	if err != nil {
		output := dto_utils.Error{Success: false, Message: "invalid route_type."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	rules, err := h.service.GetByType(routeType)
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: err.Error()}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(specialRouteStatus(err))
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	data := make([]dtos.SpecialRouteData, 0, len(rules))
	for _, rule := range rules {
		data = append(data, mapSpecialRouteData(rule))
	}

	output := dtos.SpecialRouteListOutput{Success: true, Message: "ok.", Data: data}
	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *SpecialRouteHandler) Create(ctx *fasthttp.RequestCtx) {
	routeType, err := routeTypeFromPath(string(ctx.Path()))
	if err != nil {
		output := dto_utils.Error{Success: false, Message: "invalid route_type."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	userID, err := extractUserID(ctx)
	if err != nil {
		output := dto_utils.Error{Success: false, Message: "internal error."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	var input dtos.SpecialRouteInput
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		output := dto_utils.Error{Success: false, Message: "invalid json."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	rule, err := h.service.Create(routeType, input.Path, input.MaxDistinctRequests, input.WindowSeconds, input.BanSeconds, input.Enabled, userID)
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: err.Error()}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(specialRouteStatus(err))
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}
	middlewares.LoadSpecialRoutesFromDB()

	output := dtos.SpecialRouteOutput{
		Success: true,
		Message: "special route created.",
		Data:    mapSpecialRouteData(rule),
	}
	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *SpecialRouteHandler) Update(ctx *fasthttp.RequestCtx) {
	routeType, err := routeTypeFromPath(string(ctx.Path()))
	if err != nil {
		output := dto_utils.Error{Success: false, Message: "invalid route_type."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	userID, err := extractUserID(ctx)
	if err != nil {
		output := dto_utils.Error{Success: false, Message: "internal error."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	id, err := parseSpecialRouteID(string(ctx.Path()), routeType)
	if err != nil {
		output := dto_utils.Error{Success: false, Message: "invalid id."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	var input dtos.SpecialRouteInput
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		output := dto_utils.Error{Success: false, Message: "invalid json."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	rule, err := h.service.Update(routeType, id, input.Path, input.MaxDistinctRequests, input.WindowSeconds, input.BanSeconds, input.Enabled, userID)
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: err.Error()}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(specialRouteStatus(err))
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}
	middlewares.LoadSpecialRoutesFromDB()

	output := dtos.SpecialRouteOutput{
		Success: true,
		Message: "special route updated.",
		Data:    mapSpecialRouteData(rule),
	}
	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *SpecialRouteHandler) Delete(ctx *fasthttp.RequestCtx) {
	routeType, err := routeTypeFromPath(string(ctx.Path()))
	if err != nil {
		output := dto_utils.Error{Success: false, Message: "invalid route_type."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	id, err := parseSpecialRouteID(string(ctx.Path()), routeType)
	if err != nil {
		output := dto_utils.Error{Success: false, Message: "invalid id."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	rule, err := h.service.Delete(routeType, id)
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: err.Error()}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(specialRouteStatus(err))
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}
	middlewares.LoadSpecialRoutesFromDB()

	output := dtos.SpecialRouteOutput{
		Success: true,
		Message: "special route deleted.",
		Data:    mapSpecialRouteData(rule),
	}
	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}
