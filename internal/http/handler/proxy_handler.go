package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"
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

func (h *ProxyHandler) Proxy(ctx *fasthttp.RequestCtx) {
	proxy_path := string(ctx.Path())
	method := string(ctx.Method())
	body_raw := ctx.PostBody()

	clientIP := ip.GetIP(ctx)
	country := ip2location.GetCountry(clientIP)

	application, err := s.applicationService.ApplicationRepo.FindApplicationByCountry(country)
	if err != nil {
		log.Println(err)
		output := dto_utils.Error{
			Success: false,
			Message: "internal error",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	if application == nil {
		application, err = s.applicationService.ApplicationRepo.GetRandomApplication()
		if err != nil {
			log.Println(err)
			output := dto_utils.Error{
				Success: false,
				Message: "internal error",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}
	}

	switch method {
	case "GET":
		resp, err := http.Get(application.Url + proxy_path)
		if err != nil {
			log.Println(err)
			output := dto_utils.Error{
				Success: false,
				Message: "internal error",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetContentType(resp.Header.Get("Content-Type"))
			ctx.SetBody(res)
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		ctx.SetStatusCode(resp.StatusCode)
		ctx.SetContentType(resp.Header.Get("Content-Type"))
		ctx.SetBody(body)
	case "POST":
		resp, err := http.Post(application.Url+proxy_path, "application/json", bytes.NewReader(body_raw))
		if err != nil {
			log.Println(err)
			output := dto_utils.Error{
				Success: false,
				Message: "internal error",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		ctx.SetStatusCode(resp.StatusCode)
		ctx.SetContentType(resp.Header.Get("Content-Type"))
		ctx.SetBody(body)
	case "PUT":
		resp, err := http.NewRequest("PUT", application.Url+proxy_path, bytes.NewReader(body_raw))
		if err != nil {
			log.Println(err)
			output := dto_utils.Error{
				Success: false,
				Message: "internal error",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetContentType(resp.Header.Get("Content-Type"))
			ctx.SetBody(res)
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		ctx.SetStatusCode(resp.Response.StatusCode)
		ctx.SetContentType(resp.Header.Get("Content-Type"))
		ctx.SetBody(body)
	case "DELETE":
		resp, err := http.NewRequest("DELETE", application.Url+proxy_path, bytes.NewReader(body_raw))
		if err != nil {
			log.Println(err)
			output := dto_utils.Error{
				Success: false,
				Message: "internal error",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		ctx.SetStatusCode(resp.Response.StatusCode)
		ctx.SetContentType(resp.Header.Get("Content-Type"))
		ctx.SetBody(body)
	}
	log.Printf("REQ %s %s ip=%s country=%s", method, proxy_path, clientIP, country)
}
