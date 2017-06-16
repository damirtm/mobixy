package storage

import (
	"testing"
	mobapi "github.com/dtim1985/mobixy/api"
	mobixy "github.com/dtim1985/mobixy/domain"
	"reflect"
)

func TestInMemoryStorage_ListRequestFilter_ShouldReturnFiltersByTags(t *testing.T){
	storage := NewInMemoryStorage(100, 100)
	storage.SaveRequestFilter(mobapi.FilterId("1"), TestFilter{"str1"}, "tag1")
	storage.SaveRequestFilter(mobapi.FilterId("2"), TestFilter{"str2"}, "tag1", "tag2")
	storage.SaveRequestFilter(mobapi.FilterId("3"), TestFilter{"str3"}, "tag2")
	storage.SaveRequestFilter(mobapi.FilterId("4"), TestFilter{"str4"}, "tag3")

	ListRequestFilterAndCompare(t, "tag1", storage,
		[]struct{mobapi.FilterId; mobixy.RequestFilter} {
			{mobapi.FilterId("1"), TestFilter{"str1"}},
			{mobapi.FilterId("2"), TestFilter{"str2"}},
		}...)

	ListRequestFilterAndCompare(t, "tag2", storage,
		[]struct{mobapi.FilterId; mobixy.RequestFilter} {
			{mobapi.FilterId("2"), TestFilter{"str2"}},
			{mobapi.FilterId("3"), TestFilter{"str3"}},
		}...)

	ListRequestFilterAndCompare(t, "tag3", storage,
		struct{mobapi.FilterId; mobixy.RequestFilter} {
			mobapi.FilterId("4"), TestFilter{"str4"},
		})
}

func TestInMemoryStorage_SaveRequestFilter_ShouldRewriteTags(t *testing.T) {
	storage := NewInMemoryStorage(100, 100)
	filterId := mobapi.FilterId("1")
	filter := TestFilter{"str1"}
	filterTuple := struct{mobapi.FilterId; mobixy.RequestFilter}{filterId, filter}
	storage.SaveRequestFilter(filterId, filter, "tag1")
	storage.SaveRequestFilter(filterId, filter, "tag1", "tag2")

	ListRequestFilterAndCompare(t, "tag1", storage, filterTuple)
	ListRequestFilterAndCompare(t, "tag2", storage, filterTuple)

	storage.SaveRequestFilter(filterId, filter, "tag2", "tag3")

	ListRequestFilterAndCompare(t, "tag1", storage)
	ListRequestFilterAndCompare(t, "tag2", storage, filterTuple)
	ListRequestFilterAndCompare(t, "tag3", storage, filterTuple)
}

func ListRequestFilterAndCompare(t *testing.T, tag string, storage *InMemoryStorage, expectedFilters... struct {
	mobapi.FilterId;
	mobixy.RequestFilter
}) {
	actualFilters, _ := storage.ListRequestFilter(tag)
	equivalent(actualFilters, expectedFilters)
	if !equivalent(actualFilters, expectedFilters) {
		t.Errorf("For %s expected %v, but got %v", tag, expectedFilters, actualFilters)
	}
}


type TestFilter struct {
	SomeString string
}

func equivalent(one interface{}, two interface{}) bool {
	valueOne := reflect.ValueOf(one)
	valueTwo := reflect.ValueOf(two)
	kindOne := valueOne.Kind()
	kindTwo := valueTwo.Kind()
	switch  {
	case kindOne == reflect.Slice || kindOne == reflect.Array:
		if kindTwo != reflect.Array && kindTwo != reflect.Slice {
			return false
		}
		lenOne := valueOne.Len()
		lenTwo := valueTwo.Len()
		if lenOne != lenTwo {
			return false
		}
		for i := 0; i < valueOne.Len(); i++ {
			notExist := true
			for j := 0; j < valueTwo.Len(); j++ {
				v1 := valueOne.Index(i).Interface()
				v2 := valueTwo.Index(j).Interface()
				if v1 == v2 {
					notExist = false
				}
			}
			if notExist {
				return false
			}
		}
		return true
	}
	return one == two
}