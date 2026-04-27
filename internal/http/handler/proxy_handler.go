package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"
	"github.com/ViitoJooj/ward/internal/security"
	"github.com/ViitoJooj/ward/internal/services"
	"github.com/ViitoJooj/ward/pkg/ip"
	"github.com/ViitoJooj/ward/pkg/ip2location"
	"github.com/valyala/fasthttp"
)

type ProxyHandler struct {
	proxyService       *services.ProxyService
	applicationService *services.ApplicationService
}

func NewProxyHandler(proxyHandler *services.ProxyService, applicationService *services.ApplicationService) *ProxyHandler {
	return &ProxyHandler{
		proxyService:       proxyHandler,
		applicationService: applicationService,
	}
}

func writeProxyError(ctx *fasthttp.RequestCtx, statusCode int, message string) {
	output := dto_utils.Error{
		Success: false,
		Message: message,
	}
	res, _ := json.Marshal(output)
	ctx.SetStatusCode(statusCode)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func hasSQLCommandInput(ctx *fasthttp.RequestCtx) bool {
	body := strings.TrimSpace(string(ctx.PostBody()))
	if body != "" && security.ContainsSQLCommand(body) {
		return true
	}

	query := strings.TrimSpace(string(ctx.URI().QueryString()))
	if query != "" && security.ContainsSQLCommand(query) {
		return true
	}

	return false
}

func hasSQLCommandOutput(body []byte) bool {
	value := strings.TrimSpace(string(body))
	if value == "" {
		return false
	}
	return security.ContainsSQLCommand(value)
}

func (h *ProxyHandler) Proxy(ctx *fasthttp.RequestCtx) {
	proxyPath := string(ctx.Path())
	method := string(ctx.Method())
	bodyRaw := ctx.PostBody()

	if hasSQLCommandInput(ctx) {
		writeProxyError(ctx, fasthttp.StatusBadRequest, "invalid data.")
		return
	}

	clientIP := ip.GetIP(ctx)
	country := ip2location.GetCountry(clientIP)

	application, err := h.applicationService.ApplicationRepo.FindApplicationByCountry(country)
	if err != nil {
		log.Println(err)
		writeProxyError(ctx, fasthttp.StatusInternalServerError, "internal error")
		return
	}

	if application == nil {
		application, err = h.applicationService.ApplicationRepo.GetRandomApplication()
		if err != nil {
			log.Println(err)
			writeProxyError(ctx, fasthttp.StatusInternalServerError, "internal error")
			return
		}
	}

	targetURL := application.Url + proxyPath
	if query := string(ctx.URI().QueryString()); query != "" {
		targetURL += "?" + query
	}

	req, err := http.NewRequest(method, targetURL, bytes.NewReader(bodyRaw))
	if err != nil {
		log.Println(err)
		writeProxyError(ctx, fasthttp.StatusInternalServerError, "internal error")
		return
	}
	req.Header.Set("Content-Type", string(ctx.Request.Header.ContentType()))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		writeProxyError(ctx, fasthttp.StatusInternalServerError, "internal error")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		writeProxyError(ctx, fasthttp.StatusInternalServerError, "internal error")
		return
	}

	if hasSQLCommandOutput(body) {
		writeProxyError(ctx, fasthttp.StatusInternalServerError, "internal error")
		return
	}

	ctx.SetStatusCode(resp.StatusCode)
	ctx.SetContentType(resp.Header.Get("Content-Type"))
	ctx.SetBody(body)
	log.Printf("REQ %s %s ip=%s country=%s", method, proxyPath, clientIP, country)
}
