package handler

import (
	"log"
	"strings"

	"github.com/ViitoJooj/ward/internal/services"
	"github.com/ViitoJooj/ward/pkg/ip"
	"github.com/ViitoJooj/ward/pkg/ip2location"
	"github.com/valyala/fasthttp"
)

type ProxyHandler struct {
	proxyHandler *services.ProxyService
}

func NewProxyHandler(proxyHandler *services.ProxyService) *ProxyHandler {
	return &ProxyHandler{
		proxyHandler: proxyHandler,
	}
}

func (h *ProxyHandler) Proxy(ctx *fasthttp.RequestCtx) {
	proxy_path := string(ctx.Path())
	real_path := strings.Replace(proxy_path, "/proxy", "", 1)
	method := string(ctx.Method())

	clientIP := ip.GetIP(ctx)
	country := ip2location.GetCountry(clientIP)

	log.Printf("REQ %s %s ip=%s country=%s", method, real_path, clientIP, country)
}
