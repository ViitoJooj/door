package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"
	"github.com/ViitoJooj/ward/internal/http/middlewares"
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

var hopByHopHeaders = map[string]bool{
	"Connection":          true,
	"Keep-Alive":          true,
	"Proxy-Authenticate":  true,
	"Proxy-Authorization": true,
	"Te":                  true,
	"Trailers":            true,
	"Transfer-Encoding":   true,
	"Upgrade":             true,
}

func isHopByHop(header string) bool {
	return hopByHopHeaders[http.CanonicalHeaderKey(header)]
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
	start := time.Now()
	proxyPath := string(ctx.Path())
	method := string(ctx.Method())
	bodyRaw := ctx.PostBody()

	if hasSQLCommandInput(ctx) {
		writeProxyError(ctx, fasthttp.StatusBadRequest, "invalid data.")
		return
	}

	clientIP := ip.GetIP(ctx)
	country := ip2location.GetCountry(clientIP)

	// Check for a route rule with a target URL override.
	var baseURL string
	routeRule := middlewares.FindRouteRule(proxyPath, method)
	if routeRule != nil && routeRule.TargetURL != "" {
		baseURL = routeRule.TargetURL
	} else {
		useGeo := routeRule == nil || routeRule.GeoRoutingEnabled
		if useGeo {
			app, err := h.applicationService.ApplicationRepo.FindApplicationByCountry(country)
			if err != nil {
				log.Println(err)
				writeProxyError(ctx, fasthttp.StatusInternalServerError, "internal error")
				return
			}
			if app != nil {
				baseURL = app.Url
			}
		}
		if baseURL == "" {
			app, err := h.applicationService.ApplicationRepo.GetRandomApplication()
			if err != nil {
				log.Println(err)
				writeProxyError(ctx, fasthttp.StatusInternalServerError, "internal error")
				return
			}
			if app == nil {
				log.Printf("PROXY %s %s → no application registered", method, proxyPath)
				writeProxyError(ctx, fasthttp.StatusBadGateway, "no application registered")
				return
			}
			baseURL = app.Url
		}
	}

	targetURL := baseURL + proxyPath
	if query := string(ctx.URI().QueryString()); query != "" {
		targetURL += "?" + query
	}

	req, err := http.NewRequest(method, targetURL, bytes.NewReader(bodyRaw))
	if err != nil {
		log.Println(err)
		writeProxyError(ctx, fasthttp.StatusInternalServerError, "internal error")
		return
	}

	ctx.Request.Header.VisitAll(func(key, value []byte) {
		k := string(key)
		if !isHopByHop(k) {
			req.Header.Set(k, string(value))
		}
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("PROXY %s %s → %s | error: %v", method, proxyPath, targetURL, err)
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
	for key, values := range resp.Header {
		if !isHopByHop(key) {
			ctx.Response.Header.Set(key, values[0])
		}
	}
	ctx.SetBody(body)
	log.Printf("PROXY %s %s → %s | %d | %dms | ip=%s country=%s",
		method, proxyPath, targetURL, resp.StatusCode, time.Since(start).Milliseconds(), clientIP, country)
}
