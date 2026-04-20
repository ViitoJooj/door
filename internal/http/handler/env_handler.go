package handler

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/internal/http/dtos"
	dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"
	"github.com/ViitoJooj/ward/internal/services"
	"github.com/valyala/fasthttp"
)

type DotEnvHandler struct {
	dotEnvService *services.DotEnvService
}

func NewDotEnvHandler(service *services.DotEnvService) *DotEnvHandler {
	return &DotEnvHandler{
		dotEnvService: service,
	}
}

func (h *DotEnvHandler) GetAll(ctx *fasthttp.RequestCtx) {
	var data []*domain.Env
	data, err := h.dotEnvService.GetAll()
	if err != nil {
		log.Println("internal error.")
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	res, _ := json.Marshal(data)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *DotEnvHandler) ChangeVar(ctx *fasthttp.RequestCtx) {
	var input dtos.UpgradeEnvInput

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println("invalid json.")
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

	env := domain.Env{
		Id:    input.Id,
		Name:  input.Name,
		Value: input.Value,
	}

	err := h.dotEnvService.ChangeVar(env)
	if err != nil {
		log.Println("internal error.")
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	data, err := h.dotEnvService.GetVar(input.Id)
	if err != nil {
		log.Println("internal error.")
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	output := dtos.EnvOutput{
		Success: true,
		Message: "Changed!",
		Data: dtos.Env{
			Id:    data.Id,
			Name:  data.Name,
			Value: data.Value,
		},
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *DotEnvHandler) GetVar(ctx *fasthttp.RequestCtx) {
	proxy_path := string(ctx.Path())

	application_id_str := strings.ReplaceAll(proxy_path, "/ward/api/v1/applications/", "")
	if application_id_str == "" {
		log.Println("Invalid id.")
		output := dto_utils.Error{
			Success: false,
			Message: "invalid id.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	varId, err := strconv.Atoi(application_id_str)
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

	data, err := h.dotEnvService.GetVar(varId)
	if err != nil {
		log.Println("invalid id.")
		output := dto_utils.Error{
			Success: false,
			Message: "invalid id.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	output := dtos.EnvOutput{
		Success: true,
		Message: "Changed!",
		Data: dtos.Env{
			Id:    data.Id,
			Name:  data.Name,
			Value: data.Value,
		},
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}
