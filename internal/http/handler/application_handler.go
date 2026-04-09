package handler

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/ViitoJooj/door/internal/domain"
	"github.com/ViitoJooj/door/internal/http/dtos"
	"github.com/ViitoJooj/door/internal/services"
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
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte(`{"error":"invalid json"}`))
		return
	}

	userIdRaw := ctx.UserValue("userId")
	if userIdRaw == nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	userId, ok := userIdRaw.(int64)
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	createApplication, user, err := h.applicationService.Create(input.Url, input.Country, userId)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		return
	}

	output := dtos.ApplicationOutput{
		Success: true,
		Message: "Application has created with successfull.",
		Data: dtos.ApplicationData{
			ID:      createApplication.ID,
			Url:     createApplication.Url,
			Country: createApplication.Country,
			Created_by: dtos.UserData{
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
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *ApplicationHandler) GetAll(ctx *fasthttp.RequestCtx) {
	var data []*domain.Application
	data, err := h.applicationService.GetAll()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		return
	}

	res, _ := json.Marshal(data)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *ApplicationHandler) GetByID(ctx *fasthttp.RequestCtx) {
	proxy_path := string(ctx.Path())

	application_id_str := strings.ReplaceAll(proxy_path, "/api/v1/applications/", "")
	if application_id_str == "" {
		log.Println("invalid path")
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		return
	}

	applicationId, err := strconv.ParseInt(application_id_str, 10, 64)
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		return
	}

	data, err := h.applicationService.GetByID(applicationId)
	if err != nil {
		log.Println("invalid path")
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		return
	}

	res, _ := json.Marshal(data)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *ApplicationHandler) DeleteById(ctx *fasthttp.RequestCtx) {
	proxy_path := string(ctx.Path())
	application_id_str := strings.ReplaceAll(proxy_path, "/api/v1/applications/", "")

	applicationId, err := strconv.ParseInt(application_id_str, 10, 64)
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		return
	}

	application, user, err := h.applicationService.DeleteById(applicationId)
	if err != nil {
		log.Println("invalid path")
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		return
	}

	output := dtos.ApplicationOutput{
		Success: true,
		Message: "application deleted.",
		Data: dtos.ApplicationData{
			ID:      application.ID,
			Url:     application.Url,
			Country: application.Country,
			Created_by: dtos.UserData{
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
