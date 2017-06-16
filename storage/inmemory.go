package storage

import (
	mobixy "github.com/dtim1985/mobixy/domain"
	mobapi "github.com/dtim1985/mobixy/api"
	"sync"
	"fmt"
	"errors"
)

type InMemoryStorage struct {
	channel            chan mobixy.HttpSession
	sessionsData       []mobixy.HttpSession
	requestFiltersMutex *sync.Mutex
	requestFiltersData map[mobapi.FilterId]RequestFilterStoreCell
	requestFiltersTags map[string]map[mobapi.FilterId]bool
}

type RequestFilterStoreCell struct {
	filter mobixy.RequestFilter
	tags map[string]bool
}

func NewInMemoryStorage(chanBufferSize int, initialBufferCapacity int) *InMemoryStorage {
	storage := &InMemoryStorage{
		channel:            make(chan mobixy.HttpSession, chanBufferSize),
		sessionsData:       make([]mobixy.HttpSession, 0, initialBufferCapacity),
		requestFiltersData: make(map[mobapi.FilterId]RequestFilterStoreCell, 0),
		requestFiltersTags: make(map[string]map[mobapi.FilterId]bool, 0),
		requestFiltersMutex: &sync.Mutex{},
	}
	go func() {
		for ss := range storage.channel {
			storage.sessionsData = append(storage.sessionsData, ss)
		}
	}()
	return storage
}

func (storage *InMemoryStorage) SaveHttp(session mobixy.HttpSession) {
	storage.channel <- session
	return
}

func (storage *InMemoryStorage) GetAllHttp() []mobixy.HttpSession {
	len := len(storage.sessionsData)
	dataToReturn := make([]mobixy.HttpSession, len, len)
	copy(dataToReturn, storage.sessionsData)
	return dataToReturn
}

func (storage *InMemoryStorage) SaveRequestFilter(id mobapi.FilterId, filter mobixy.RequestFilter, tags... string) error {
	storage.requestFiltersMutex.Lock()
	defer storage.requestFiltersMutex.Unlock()
	storeCell, ok := storage.requestFiltersData[id]
	if !ok {
		storeCell = RequestFilterStoreCell{filter:nil, tags: make(map[string]bool, 0) }
	}
	newTags := make(map[string]bool, 0)
	for _, tag := range tags {
		delete(storeCell.tags, tag)
		newTags[tag] = true
	}
	for tag, _ := range newTags {
		ids, ok := storage.requestFiltersTags[tag]
		if !ok {
			ids = make(map[mobapi.FilterId]bool)
		}
		ids[id] = true
		storage.requestFiltersTags[tag] = ids
	}
	for tag, _ := range storeCell.tags {
		ids, ok := storage.requestFiltersTags[tag]
		if ok {
			delete(ids, id)
		}
	}
	storeCell.filter = filter
	storeCell.tags = newTags
	storage.requestFiltersData[id] = storeCell
	return nil
}

func (storage *InMemoryStorage) GetRequestFilter(id mobapi.FilterId) (mobixy.RequestFilter, error) {
	filter, ok := storage.requestFiltersData[id]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Filter with id %s was not found", id))
	}
	return filter, nil
}

func (storage *InMemoryStorage) DeleteRequestFilter(id mobapi.FilterId) error {
	storage.requestFiltersMutex.Lock()
	defer storage.requestFiltersMutex.Unlock()
	delete(storage.requestFiltersData, id)
	return nil
}

func (storage *InMemoryStorage) ListRequestFilter(tags... string) ([]struct{mobapi.FilterId; mobixy.RequestFilter}, error) {
	ids := make(map[mobapi.FilterId]bool, 0)
	for _, tag := range tags {
		tagsId, ok := storage.requestFiltersTags[tag]
		if !ok {
			continue
		}
		for id, _ := range tagsId {
			ids[id] = true
		}
	}
	result := make([]struct{mobapi.FilterId; mobixy.RequestFilter}, 0)
	for id, _ := range ids {
		filterStoreCell, ok := storage.requestFiltersData[id]
		if ok {
			result = append(result, struct {
				mobapi.FilterId
				mobixy.RequestFilter
			}{id, filterStoreCell.filter})
		}
	}
	return result, nil
}

func (storage *InMemoryStorage) GetAllRequestFilters() []mobixy.RequestFilter {
	len := len(storage.requestFiltersData)
	dataToReturn := make([]mobixy.RequestFilter, 0, len)
	for _, filter := range storage.requestFiltersData {
		dataToReturn = append(dataToReturn, filter)
	}
	return dataToReturn
}
