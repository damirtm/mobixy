package main

import (
	"github.com/elazarl/goproxy"
	"net/http"
	mobixy "github.com/dtim1985/mobixy/domain"
)

func getProxy(sessionsStorage mobixy.SessionsStorage, filtersStorage mobixy.FiltersStorage) *goproxy.ProxyHttpServer {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		for _, filter := range filtersStorage.GetAllRequestFilters() {
			r = mobixy.FilterWith(filter, r)
		}
		return r, nil
	})
	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		sessionsStorage.SaveHttp(mobixy.HttpSession{
			Request:  ctx.Req,
			Response: r,
		})
		return r
	})
	return proxy
}
