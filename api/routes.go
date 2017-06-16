package api

import (
	"github.com/gorilla/mux"
	mobixy "github.com/dtim1985/mobixy/domain"
)

func GetApi(sessionsStorage mobixy.SessionsStorage, filtersStorage FiltersStorage) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/sessions", GetSessionsHandler{sessionsStorage: sessionsStorage}).Methods("GET")
	r.Handle("/filters/{type}/{id}", FiltersHandler{filtersStorage: filtersStorage}).Methods("POST", "GET", "DELETE", "PUT")
	return r
}