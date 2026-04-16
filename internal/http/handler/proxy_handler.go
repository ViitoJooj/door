package handler

import (
	"log"
	"strings"

	"github.com/ViitoJooj/ward/internal/services"
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

func (s *ProxyHandler) Proxy(ctx *fasthttp.RequestCtx) {
	proxy_path := string(ctx.Path())
	real_path := strings.Replace(proxy_path, "/proxy", "", 1)
	method := string(ctx.Method())

	log.Printf("REQ %s %s", method, real_path)
}
