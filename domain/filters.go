package domain

import (
	"net/http"
	"net/url"
	"log"
)

type FiltersStorage interface {
	GetAllRequestFilters() []RequestFilter
}

func FilterWith(filter RequestFilter, req *http.Request) *http.Request {
	switch v := filter.(type) {
	case ConditionRequestTransformFilter:
		if v.Condition().Suited(req) {
			return v.Transform().Apply(req)
		}
		return req
	case SimpleFilter:
		return v.Filter(req)
	default:
		return req
	}
}

type RequestFilter interface {}

type SimpleFilter interface {
	RequestFilter
	Filter(req *http.Request) *http.Request
}

type ConditionRequestTransformFilter interface {
	RequestFilter
	Condition() RequestCondition
	Transform() RequestTransform
}

type RequestCondition interface {
	Suited(req *http.Request) bool
}

type RequestTransform interface {
	Apply(req *http.Request) *http.Request
}

type HostEqualCondition struct {
	HostToCompare string
}

func (condition HostEqualCondition) Suited(req *http.Request) bool {
	remoteUrl, err := url.Parse(req.RemoteAddr)
	if err != nil {
		log.Fatalf("Error while parsing Request RemoteAddr to url, err: %s", err)
		return false
	}
	return remoteUrl.Host == condition.HostToCompare
}

type HostReplaceTransform struct {
	HostToReplace string
}

func (transform HostReplaceTransform) Apply(req *http.Request) *http.Request {
	req.Host = transform.HostToReplace
	return req
}
