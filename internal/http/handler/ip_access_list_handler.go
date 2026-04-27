package handler

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/ViitoJooj/ward/internal/http/dtos"
	dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"
	"github.com/ViitoJooj/ward/internal/http/middlewares"
	"github.com/ViitoJooj/ward/internal/services"
	"github.com/valyala/fasthttp"
)

type IPAccessListHandler struct {
	service *services.IPAccessListService
}

func NewIPAccessListHandler(service *services.IPAccessListService) *IPAccessListHandler {
	return &IPAccessListHandler{
		service: service,
	}
}

func ipAccessErrorStatus(err error) int {
	if err == nil {
		return fasthttp.StatusBadGateway
	}
	if err.Error() == "invalid ip" || err.Error() == "ip is required" {
		return fasthttp.StatusBadRequest
	}
	if err.Error() == "ip already exists" {
		return fasthttp.StatusBadRequest
	}
	if err.Error() == "ip not found" {
		return fasthttp.StatusNotFound
	}
	return fasthttp.StatusBadGateway
}

func parseIDFromPath(path string, prefix string) (int, error) {
	idStr := strings.TrimSpace(strings.ReplaceAll(path, prefix, ""))
	return strconv.Atoi(idStr)
}

func extractUserID(ctx *fasthttp.RequestCtx) (int, error) {
	userIdRaw := ctx.UserValue("userId")
	if userIdRaw == nil {
		return 0, errors.New("userId not found")
	}

	userID, ok := userIdRaw.(int)
	if !ok {
		return 0, errors.New("invalid userId")
	}
	return userID, nil
}

func (h *IPAccessListHandler) GetWhitelist(ctx *fasthttp.RequestCtx) {
	entries, err := h.service.GetWhitelist()
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: "internal error."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	data := make([]dtos.IPAccessData, 0, len(entries))
	for _, entry := range entries {
		data = append(data, dtos.IPAccessData{
			ID:        entry.ID,
			IP:        entry.IP,
			CreatedBy: entry.CreatedBy,
			UpdatedBy: entry.UpdatedBy,
			CreatedAt: entry.CreatedAt,
			UpdatedAt: entry.UpdatedAt,
		})
	}

	output := dtos.IPAccessListOutput{Success: true, Message: "ok.", Data: data}
	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *IPAccessListHandler) GetBlacklist(ctx *fasthttp.RequestCtx) {
	entries, err := h.service.GetBlacklist()
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: "internal error."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	data := make([]dtos.IPAccessData, 0, len(entries))
	for _, entry := range entries {
		data = append(data, dtos.IPAccessData{
			ID:        entry.ID,
			IP:        entry.IP,
			CreatedBy: entry.CreatedBy,
			UpdatedBy: entry.UpdatedBy,
			CreatedAt: entry.CreatedAt,
			UpdatedAt: entry.UpdatedAt,
		})
	}

	output := dtos.IPAccessListOutput{Success: true, Message: "ok.", Data: data}
	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *IPAccessListHandler) CreateWhitelist(ctx *fasthttp.RequestCtx) {
	userID, err := extractUserID(ctx)
	if err != nil {
		output := dto_utils.Error{Success: false, Message: "internal error."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	var input dtos.IPAccessInput
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: "invalid json."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	entry, err := h.service.CreateWhitelist(input.IP, userID)
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: err.Error()}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(ipAccessErrorStatus(err))
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}
	middlewares.LoadIPAccessListsFromDB()

	output := dtos.IPAccessOutput{
		Success: true,
		Message: "whitelist ip created.",
		Data: dtos.IPAccessData{
			ID:        entry.ID,
			IP:        entry.IP,
			CreatedBy: entry.CreatedBy,
			UpdatedBy: entry.UpdatedBy,
			CreatedAt: entry.CreatedAt,
			UpdatedAt: entry.UpdatedAt,
		},
	}
	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *IPAccessListHandler) CreateBlacklist(ctx *fasthttp.RequestCtx) {
	userID, err := extractUserID(ctx)
	if err != nil {
		output := dto_utils.Error{Success: false, Message: "internal error."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	var input dtos.IPAccessInput
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: "invalid json."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	entry, err := h.service.CreateBlacklist(input.IP, userID)
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: err.Error()}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(ipAccessErrorStatus(err))
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}
	middlewares.LoadIPAccessListsFromDB()

	output := dtos.IPAccessOutput{
		Success: true,
		Message: "blacklist ip created.",
		Data: dtos.IPAccessData{
			ID:        entry.ID,
			IP:        entry.IP,
			CreatedBy: entry.CreatedBy,
			UpdatedBy: entry.UpdatedBy,
			CreatedAt: entry.CreatedAt,
			UpdatedAt: entry.UpdatedAt,
		},
	}
	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *IPAccessListHandler) UpdateWhitelist(ctx *fasthttp.RequestCtx) {
	userID, err := extractUserID(ctx)
	if err != nil {
		output := dto_utils.Error{Success: false, Message: "internal error."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	id, err := parseIDFromPath(string(ctx.Path()), "/ward/api/v1/ip-whitelist/")
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: "invalid id."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	var input dtos.IPAccessInput
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: "invalid json."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	entry, err := h.service.UpdateWhitelist(id, input.IP, userID)
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: err.Error()}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(ipAccessErrorStatus(err))
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}
	middlewares.LoadIPAccessListsFromDB()

	output := dtos.IPAccessOutput{
		Success: true,
		Message: "whitelist ip updated.",
		Data: dtos.IPAccessData{
			ID:        entry.ID,
			IP:        entry.IP,
			CreatedBy: entry.CreatedBy,
			UpdatedBy: entry.UpdatedBy,
			CreatedAt: entry.CreatedAt,
			UpdatedAt: entry.UpdatedAt,
		},
	}
	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *IPAccessListHandler) UpdateBlacklist(ctx *fasthttp.RequestCtx) {
	userID, err := extractUserID(ctx)
	if err != nil {
		output := dto_utils.Error{Success: false, Message: "internal error."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	id, err := parseIDFromPath(string(ctx.Path()), "/ward/api/v1/ip-blacklist/")
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: "invalid id."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	var input dtos.IPAccessInput
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: "invalid json."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	entry, err := h.service.UpdateBlacklist(id, input.IP, userID)
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: err.Error()}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(ipAccessErrorStatus(err))
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}
	middlewares.LoadIPAccessListsFromDB()

	output := dtos.IPAccessOutput{
		Success: true,
		Message: "blacklist ip updated.",
		Data: dtos.IPAccessData{
			ID:        entry.ID,
			IP:        entry.IP,
			CreatedBy: entry.CreatedBy,
			UpdatedBy: entry.UpdatedBy,
			CreatedAt: entry.CreatedAt,
			UpdatedAt: entry.UpdatedAt,
		},
	}
	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *IPAccessListHandler) DeleteWhitelist(ctx *fasthttp.RequestCtx) {
	id, err := parseIDFromPath(string(ctx.Path()), "/ward/api/v1/ip-whitelist/")
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: "invalid id."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	entry, err := h.service.DeleteWhitelist(id)
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: err.Error()}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(ipAccessErrorStatus(err))
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}
	middlewares.LoadIPAccessListsFromDB()

	output := dtos.IPAccessOutput{
		Success: true,
		Message: "whitelist ip deleted.",
		Data: dtos.IPAccessData{
			ID:        entry.ID,
			IP:        entry.IP,
			CreatedBy: entry.CreatedBy,
			UpdatedBy: entry.UpdatedBy,
			CreatedAt: entry.CreatedAt,
			UpdatedAt: entry.UpdatedAt,
		},
	}
	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *IPAccessListHandler) DeleteBlacklist(ctx *fasthttp.RequestCtx) {
	id, err := parseIDFromPath(string(ctx.Path()), "/ward/api/v1/ip-blacklist/")
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: "invalid id."}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	entry, err := h.service.DeleteBlacklist(id)
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{Success: false, Message: err.Error()}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(ipAccessErrorStatus(err))
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}
	middlewares.LoadIPAccessListsFromDB()

	output := dtos.IPAccessOutput{
		Success: true,
		Message: "blacklist ip deleted.",
		Data: dtos.IPAccessData{
			ID:        entry.ID,
			IP:        entry.IP,
			CreatedBy: entry.CreatedBy,
			UpdatedBy: entry.UpdatedBy,
			CreatedAt: entry.CreatedAt,
			UpdatedAt: entry.UpdatedAt,
		},
	}
	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}
