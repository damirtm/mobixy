package api

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	mobixy "github.com/dtim1985/mobixy/domain"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

type FilterId string

type FiltersStorage interface {
	SaveRequestFilter(id FilterId, filter mobixy.RequestFilter, tags... string) error
	GetRequestFilter(id FilterId) (mobixy.RequestFilter, error)
	DeleteRequestFilter(id FilterId) error
	ListRequestFilter(tags... string) ([]struct{FilterId; mobixy.RequestFilter}, error)
}

var typeHandlers = map[string]FilterOperationsHandler{
	"replaceHost": {
		Get:    GetReplaceHostFilter,
		Delete: DeleteReplaceHostFilter,
		Post:   PostReplaceHostFilter,
		Put:    PutReplaceHostFilter,
	},
}

type FiltersHandler struct {
	filtersStorage FiltersStorage
}

type PostHandler func(filtersStorage FiltersStorage, body []byte) (FilterId, error)
type PutHandler func(filtersStorage FiltersStorage, id FilterId, body []byte) error
type DeleteHandler func(filtersStorage FiltersStorage, id FilterId) error
type GetHandler func(filtersStorage FiltersStorage, id FilterId) ([]byte, error)

type FilterOperationsHandler struct {
	Post PostHandler
	Put PutHandler
	Delete DeleteHandler
	Get GetHandler
}

func (handler FiltersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fltType, ok := vars["type"]
	if !ok {
		http.Error(w, "Wrong request, filter type was not present", 400)
		return
	}
	opHandler, ok := typeHandlers[fltType]
	if !ok {
		http.Error(w, "Wrong request, filter was not found", 400)
		return
	}
	var error error
	var body []byte
	var id FilterId
	if r.Method == "POST" || r.Method == "PUT" {
		defer r.Body.Close()
		body, error = ioutil.ReadAll(r.Body)
		if error != nil {
			http.Error(w, "Wrong request, can't get body", 400)
			return
		}
	}
	if r.Method == "PUT" || r.Method == "DELETE" || r.Method == "GET" {
		strId, ok := vars["id"]
		if !ok {
			http.Error(w, "Wrong request, id does not present", 400)
			return
		}
		id = FilterId(strId)
	}
	switch r.Method {
	case "POST":
		id, error = opHandler.Post(handler.filtersStorage, body)
		if error == nil {
			bytes, _ := json.Marshal(PostOperationResult{Id: id})
			WriteJson(w, bytes)
			return
		}
	case "PUT":
		error = opHandler.Put(handler.filtersStorage, id, body)
		if error == nil {
			return
		}
	case "DELETE":
		error = opHandler.Delete(handler.filtersStorage, id)
		if error == nil {
			return
		}
	case "GET":
		bytes, error := opHandler.Get(handler.filtersStorage, id)
		if error == nil {
			WriteJson(w, bytes)
			return
		}
	default:
		http.Error(w, "Method does not allowed", 400)
	}
	if error != nil {
		HandleOperationError(error, w)
	}
	return
}

func HandleOperationError(err error, w http.ResponseWriter){
	switch err.(type) {
	case StorageError:
		http.Error(w, "Unexpected storage error", 500)
		return
	case InputError:
		http.Error(w, "Wrong operation call error", 400)
		return
	default:
		http.Error(w, "Unexpected error", 500)
		return
	}
}

func WriteJson(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

type PostOperationResult struct {
	Id FilterId
}

type ReplaceHostPost struct {
	FromHost string
	ToHost	string
}

type StorageError error
type InputError error

func PostReplaceHostFilter(storage FiltersStorage, body []byte) (FilterId, error) {
	id := NewFilterId()
	error := PutReplaceHostFilter(storage, id, body)
	return id, error
}

func PutReplaceHostFilter(storage FiltersStorage, id FilterId, body []byte) error {
	var replaceHostPost ReplaceHostPost
	_ = json.Unmarshal(body, &replaceHostPost)
	error := storage.SaveRequestFilter(id, replaceHostPost)
	if error != nil {
		return StorageError(error)
	}
	return nil
}

func GetReplaceHostFilter(storage FiltersStorage, id FilterId) ([]byte, error) {
	filter, error := storage.GetRequestFilter(id)
	if error != nil {
		return nil, StorageError(error)
	}
	result, err := json.Marshal(filter)
	return result, err
}

func DeleteReplaceHostFilter(storage FiltersStorage, id FilterId) error {
	return StorageError(storage.DeleteRequestFilter(id))
}

func (filter ReplaceHostPost) Condition() mobixy.RequestCondition {
	return mobixy.HostEqualCondition{HostToCompare:filter.FromHost}
}

func (filter ReplaceHostPost) Transform() mobixy.RequestTransform {
	return mobixy.HostReplaceTransform{HostToReplace:filter.ToHost}
}

func NewFilterId() FilterId {
	return FilterId(uuid.NewV4().String())
}
