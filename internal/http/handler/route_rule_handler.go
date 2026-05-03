package handler

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/ViitoJooj/ward/internal/http/dtos"
	dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"
	"github.com/ViitoJooj/ward/internal/http/middlewares"
	"github.com/ViitoJooj/ward/internal/services"
	"github.com/valyala/fasthttp"
)

type RouteRuleHandler struct {
	service *services.RouteRuleService
}

func NewRouteRuleHandler(service *services.RouteRuleService) *RouteRuleHandler {
	return &RouteRuleHandler{service: service}
}

func routeRuleStatus(err error) int {
	if err == nil {
		return fasthttp.StatusBadGateway
	}
	switch err.Error() {
	case "route rule not found":
		return fasthttp.StatusNotFound
	case "path is required", "rate_limit_rps must be greater than 0", "rate_limit_burst must be greater than 0":
		return fasthttp.StatusBadRequest
	}
	return fasthttp.StatusBadGateway
}

func (h *RouteRuleHandler) GetAll(ctx *fasthttp.RequestCtx) {
	rules, err := h.service.List()
	if err != nil {
		log.Println(err)
		res, _ := json.Marshal(dto_utils.Error{Success: false, Message: "internal error."})
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	data := make([]dtos.RouteRuleData, 0, len(rules))
	for _, r := range rules {
		data = append(data, dtos.RouteRuleData{
			ID: r.ID, Path: r.Path, Method: r.Method,
			RateLimitEnabled: r.RateLimitEnabled, RateLimitRPS: r.RateLimitRPS, RateLimitBurst: r.RateLimitBurst,
			TargetURL: r.TargetURL, GeoRoutingEnabled: r.GeoRoutingEnabled, Enabled: r.Enabled,
			CreatedBy: r.CreatedBy, UpdatedBy: r.UpdatedBy, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
		})
	}

	res, _ := json.Marshal(dtos.RouteRuleListOutput{Success: true, Message: "ok.", Data: data})
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *RouteRuleHandler) Create(ctx *fasthttp.RequestCtx) {
	userID, err := extractUserID(ctx)
	if err != nil {
		res, _ := json.Marshal(dto_utils.Error{Success: false, Message: "internal error."})
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	var input dtos.RouteRuleInput
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		res, _ := json.Marshal(dto_utils.Error{Success: false, Message: "invalid json."})
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	rule, err := h.service.Create(input.Path, input.Method, input.RateLimitEnabled, input.RateLimitRPS, input.RateLimitBurst, input.TargetURL, input.GeoRoutingEnabled, input.Enabled, userID)
	if err != nil {
		log.Println(err)
		res, _ := json.Marshal(dto_utils.Error{Success: false, Message: err.Error()})
		ctx.SetStatusCode(routeRuleStatus(err))
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}
	middlewares.LoadRouteRulesFromDB()

	out := dtos.RouteRuleOutput{Success: true, Message: "route rule created.", Data: dtos.RouteRuleData{
		ID: rule.ID, Path: rule.Path, Method: rule.Method,
		RateLimitEnabled: rule.RateLimitEnabled, RateLimitRPS: rule.RateLimitRPS, RateLimitBurst: rule.RateLimitBurst,
		TargetURL: rule.TargetURL, GeoRoutingEnabled: rule.GeoRoutingEnabled, Enabled: rule.Enabled,
		CreatedBy: rule.CreatedBy, UpdatedBy: rule.UpdatedBy, CreatedAt: rule.CreatedAt, UpdatedAt: rule.UpdatedAt,
	}}
	res, _ := json.Marshal(out)
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *RouteRuleHandler) Update(ctx *fasthttp.RequestCtx) {
	userID, err := extractUserID(ctx)
	if err != nil {
		res, _ := json.Marshal(dto_utils.Error{Success: false, Message: "internal error."})
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	idStr := strings.TrimPrefix(string(ctx.Path()), "/ward/api/v1/route-rules/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res, _ := json.Marshal(dto_utils.Error{Success: false, Message: "invalid id."})
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	var input dtos.RouteRuleInput
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		res, _ := json.Marshal(dto_utils.Error{Success: false, Message: "invalid json."})
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	rule, err := h.service.Update(id, input.Path, input.Method, input.RateLimitEnabled, input.RateLimitRPS, input.RateLimitBurst, input.TargetURL, input.GeoRoutingEnabled, input.Enabled, userID)
	if err != nil {
		log.Println(err)
		res, _ := json.Marshal(dto_utils.Error{Success: false, Message: err.Error()})
		ctx.SetStatusCode(routeRuleStatus(err))
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}
	middlewares.LoadRouteRulesFromDB()

	out := dtos.RouteRuleOutput{Success: true, Message: "route rule updated.", Data: dtos.RouteRuleData{
		ID: rule.ID, Path: rule.Path, Method: rule.Method,
		RateLimitEnabled: rule.RateLimitEnabled, RateLimitRPS: rule.RateLimitRPS, RateLimitBurst: rule.RateLimitBurst,
		TargetURL: rule.TargetURL, GeoRoutingEnabled: rule.GeoRoutingEnabled, Enabled: rule.Enabled,
		CreatedBy: rule.CreatedBy, UpdatedBy: rule.UpdatedBy, CreatedAt: rule.CreatedAt, UpdatedAt: rule.UpdatedAt,
	}}
	res, _ := json.Marshal(out)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *RouteRuleHandler) Delete(ctx *fasthttp.RequestCtx) {
	idStr := strings.TrimPrefix(string(ctx.Path()), "/ward/api/v1/route-rules/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res, _ := json.Marshal(dto_utils.Error{Success: false, Message: "invalid id."})
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	rule, err := h.service.Delete(id)
	if err != nil {
		log.Println(err)
		res, _ := json.Marshal(dto_utils.Error{Success: false, Message: err.Error()})
		ctx.SetStatusCode(routeRuleStatus(err))
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}
	middlewares.LoadRouteRulesFromDB()

	out := dtos.RouteRuleOutput{Success: true, Message: "route rule deleted.", Data: dtos.RouteRuleData{
		ID: rule.ID, Path: rule.Path, Method: rule.Method,
		RateLimitEnabled: rule.RateLimitEnabled, RateLimitRPS: rule.RateLimitRPS, RateLimitBurst: rule.RateLimitBurst,
		TargetURL: rule.TargetURL, GeoRoutingEnabled: rule.GeoRoutingEnabled, Enabled: rule.Enabled,
		CreatedBy: rule.CreatedBy, UpdatedBy: rule.UpdatedBy, CreatedAt: rule.CreatedAt, UpdatedAt: rule.UpdatedAt,
	}}
	res, _ := json.Marshal(out)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}
