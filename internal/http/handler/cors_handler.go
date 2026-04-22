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

type CorsHandler struct {
	corsService *services.CorsService
}

func NewCorsHandler(corsService *services.CorsService) *CorsHandler {
	return &CorsHandler{
		corsService: corsService,
	}
}

func (h *CorsHandler) Create(ctx *fasthttp.RequestCtx) {
	var input dtos.CorsInput

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
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

	userIdRaw := ctx.UserValue("userId")
	if userIdRaw == nil {
		log.Println("userIdRaw not exists.")
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	userId, ok := userIdRaw.(int)
	if !ok {
		log.Println("userId Not valid.")
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	createCors, user, err := h.corsService.Create(input.Origin, userId)
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

	output := dtos.CorsOutput{
		Success: true,
		Message: "Cors created successfully.",
		Data: dtos.CorsData{
			ID:     createCors.Id,
			Origin: createCors.Origin,
			Created_by: dto_utils.UserData{
				ID:         user.ID,
				Username:   user.Username,
				Email:      user.Email,
				Updated_at: user.Updated_at.String(),
				Created_at: user.Created_at.String(),
			},
		},
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *CorsHandler) GetAll(ctx *fasthttp.RequestCtx) {
	var data []*domain.Cors

	data, err := h.corsService.GetAll()
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

	res, _ := json.Marshal(data)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *CorsHandler) GetByID(ctx *fasthttp.RequestCtx) {
	proxy_path := string(ctx.Path())

	cors_id_str := strings.ReplaceAll(proxy_path, "/ward/api/v1/cors/", "")
	if cors_id_str == "" {
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

	corsId, err := strconv.Atoi(cors_id_str)
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

	data, err := h.corsService.GetByID(corsId)
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

	res, _ := json.Marshal(data)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *CorsHandler) Update(ctx *fasthttp.RequestCtx) {
	proxy_path := string(ctx.Path())
	cors_id_str := strings.ReplaceAll(proxy_path, "/ward/api/v1/cors/", "")

	corsId, err := strconv.Atoi(cors_id_str)
	if err != nil {
		log.Println(err)
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

	var input dtos.CorsInput
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
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

	userIdRaw := ctx.UserValue("userId")
	if userIdRaw == nil {
		log.Println("userIdRaw not exists.")
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	userId, ok := userIdRaw.(int)
	if !ok {
		log.Println("userId Not valid.")
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	updatedCors, user, err := h.corsService.Update(corsId, input.Origin, userId)
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

	output := dtos.CorsOutput{
		Success: true,
		Message: "cors updated.",
		Data: dtos.CorsData{
			ID:     updatedCors.Id,
			Origin: updatedCors.Origin,
			Created_by: dto_utils.UserData{
				ID:         user.ID,
				Username:   user.Username,
				Email:      user.Email,
				Updated_at: user.Updated_at.String(),
				Created_at: user.Created_at.String(),
			},
		},
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *CorsHandler) DeleteById(ctx *fasthttp.RequestCtx) {
	proxy_path := string(ctx.Path())
	cors_id_str := strings.ReplaceAll(proxy_path, "/ward/api/v1/cors/", "")

	corsId, err := strconv.Atoi(cors_id_str)
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

	userIdRaw := ctx.UserValue("userId")
	if userIdRaw == nil {
		log.Println("userIdRaw not exists.")
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	userId, ok := userIdRaw.(int)
	if !ok {
		log.Println("userId Not valid.")
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	cors, user, err := h.corsService.DeleteByID(corsId, userId)
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{
			Success: false,
			Message: "invalid id.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	output := dtos.CorsOutput{
		Success: true,
		Message: "cors deleted.",
		Data: dtos.CorsData{
			ID:     cors.Id,
			Origin: cors.Origin,
			Created_by: dto_utils.UserData{
				ID:         user.ID,
				Username:   user.Username,
				Email:      user.Email,
				Updated_at: user.Updated_at.String(),
				Created_at: user.Created_at.String(),
			},
		},
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}
