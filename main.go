package main

import (
	"encoding/json"
	"github.com/elazarl/goproxy"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
)

func main() {
	storage := NewInMemoryStorage(100, 100)
	proxy := getProxy(storage)
	api := getApi(storage)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Fatal(http.ListenAndServe(":8889", proxy))
	}()
	go func() {
		defer wg.Done()
		log.Fatal(http.ListenAndServe(":8000", api))
	}()
	wg.Wait()
}

func getProxy(storage SessionsStorage) *goproxy.ProxyHttpServer {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return r, nil
	})
	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		storage.SaveHttp(HttpSession{
			request:  ctx.Req,
			response: r,
		})
		return r
	})
	return proxy
}

func getApi(sessionsStorage SessionsStorage) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/http/sessions", GetSessionsHandler{sessionsStorage: sessionsStorage}).Methods("GET")
	return r
}

type GetSessionsHandler struct {
	sessionsStorage SessionsStorage
}

func (handler GetSessionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	allHttp := handler.sessionsStorage.GetAllHttp()
	list := make([]HttpSessionListItem, len(allHttp))
	for i := 0; i < len(allHttp); i++ {
		request := allHttp[i].request
		list[i] = HttpSessionListItem{
			From: request.RemoteAddr,
			To:   request.Host,
		}
	}
	json, _ := json.Marshal(list)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

type HttpSessionListItem struct {
	From string
	To   string
}

type HttpSession struct {
	request  *http.Request
	response *http.Response
	stats    Stats
}

type Stats struct {
}

type SessionsStorage interface {
	SaveHttp(session HttpSession)
	GetAllHttp() []HttpSession
}

type InMemoryStorage struct {
	channel chan HttpSession
	data    []HttpSession
}

func NewInMemoryStorage(chanBufferSize int, initialBufferCapacity int) *InMemoryStorage {
	storage := &InMemoryStorage{
		channel: make(chan HttpSession, chanBufferSize),
		data:    make([]HttpSession, 0, initialBufferCapacity),
	}
	go func() {
		for ss := range storage.channel {
			storage.data = append(storage.data, ss)
		}
	}()
	return storage
}

func (storage *InMemoryStorage) SaveHttp(session HttpSession) {
	storage.channel <- session
	return
}

func (storage *InMemoryStorage) GetAllHttp() []HttpSession {
	len := len(storage.data)
	dataToReturn := make([]HttpSession, len, len)
	copy(dataToReturn, storage.data)
	return dataToReturn
}
