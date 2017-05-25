package main

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func main() {
	var endpointMap EndpointMap
	server := &http.Server{
		Addr: ":8080",
	}
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		endpointMap.Register(host(r.RemoteAddr), host(r.Host))
		return r, nil
	})
	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		r.Header.Set("X-GoProxy", "yxorPoG-X")
		return r
	})
	log.Fatal(http.ListenAndServe(":8888", proxy))
}

type appHandler struct {

}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, *http.Request) {
	w.
}


func host(remoteAddr string) Endpoint {
	return Endpoint{[]EndpointAddress{EndpointAddress{remoteAddr, "host"}}}
}

type Endpoint struct {
	Addresses []EndpointAddress
}

type EndpointAddress struct {
	Address  string
	AddrType string
}

type WireRegistry interface {
	Register(requestEndpoint Endpoint, responseEndpoint Endpoint)
	RequestEndpoints() []Endpoint
}

func (m *EndpointMap) Register(requestEndpoint Endpoint, responseEndpoint Endpoint) {
	*m = append(*m, struct {
		requestEndpoint  Endpoint
		responseEndpoint Endpoint
	}{requestEndpoint, responseEndpoint})
}

type EndpointMap []struct {
	requestEndpoint  Endpoint
	responseEndpoint Endpoint
}
