package api

import (
	"net/http"
	"encoding/json"
	mobixy "github.com/dtim1985/mobixy/domain"
)

type GetSessionsHandler struct {
	sessionsStorage mobixy.SessionsStorage
}

func (handler GetSessionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	allHttp := handler.sessionsStorage.GetAllHttp()
	list := make([]HttpSessionListItem, len(allHttp))
	for i := 0; i < len(allHttp); i++ {
		request := allHttp[i].Request
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
