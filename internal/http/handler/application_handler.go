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

type ApplicationHandler struct {
	applicationService *services.ApplicationService
}

func NewApplicationHandler(applicationService *services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		applicationService: applicationService,
	}
}

func (h *ApplicationHandler) Create(ctx *fasthttp.RequestCtx) {
	var input dtos.ApplicationInput

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

	createApplication, user, err := h.applicationService.Create(input.Url, input.Country, userId)
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

	output := dtos.ApplicationOutput{
		Success: true,
		Message: "Application has created with successfull.",
		Data: dtos.ApplicationData{
			ID:      createApplication.ID,
			Url:     createApplication.Url,
			Country: createApplication.Country,
			Created_by: dto_utils.UserData{
				ID:         user.ID,
				Username:   user.Username,
				Email:      user.Email,
				Updated_at: user.Updated_at.String(),
				Created_at: user.Created_at.String(),
			},
			Updated_at: createApplication.Updated_at,
			Created_at: createApplication.Created_at,
		},
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *ApplicationHandler) GetAll(ctx *fasthttp.RequestCtx) {
	var data []*domain.Application
	data, err := h.applicationService.GetAll()
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

func (h *ApplicationHandler) GetByID(ctx *fasthttp.RequestCtx) {
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

	applicationId, err := strconv.Atoi(application_id_str)
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

	data, err := h.applicationService.GetByID(applicationId)
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

func (h *ApplicationHandler) DeleteById(ctx *fasthttp.RequestCtx) {
	proxy_path := string(ctx.Path())
	application_id_str := strings.ReplaceAll(proxy_path, "/ward/api/v1/applications/", "")

	applicationId, err := strconv.Atoi(application_id_str)
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

	application, user, err := h.applicationService.DeleteById(applicationId)
	if err != nil {
		log.Println("no application")
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

	output := dtos.ApplicationOutput{
		Success: true,
		Message: "application deleted.",
		Data: dtos.ApplicationData{
			ID:      application.ID,
			Url:     application.Url,
			Country: application.Country,
			Created_by: dto_utils.UserData{
				ID:         user.ID,
				Username:   user.Username,
				Email:      user.Email,
				Updated_at: user.Updated_at.String(),
				Created_at: user.Created_at.String(),
			},
			Updated_at: application.Updated_at,
			Created_at: application.Created_at,
		},
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}
