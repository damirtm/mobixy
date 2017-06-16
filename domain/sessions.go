package domain

import "net/http"

type HttpSession struct {
	Request  *http.Request
	Response *http.Response
	Stats    Stats
}

type Stats struct {
}

type SessionsStorage interface {
	SaveHttp(session HttpSession)
	GetAllHttp() []HttpSession
}
