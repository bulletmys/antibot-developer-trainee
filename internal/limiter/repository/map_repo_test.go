package repository

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestMapRepo_AddToBlackList(t *testing.T) {
	testKey := "192.168.0.0"
	testKey2 := "192.169.0.0"

	testKey2Time := time.Now()

	mapRepo := MapRepo{
		blackList:    map[string]time.Time{},
		requests:     map[string][]time.Time{testKey: {time.Now()}, testKey2: {testKey2Time}},
		BlackListTTL: time.Second * 2,
		muBL:         &sync.RWMutex{},
		muR:          &sync.RWMutex{},
	}

	AfterFuncR := map[string][]time.Time{testKey2: {testKey2Time}}

	mapRepo.AddToBlackList(testKey)

	_, ok := mapRepo.blackList[testKey]
	assert.True(t, ok)

	assert.True(t, reflect.DeepEqual(mapRepo.requests, AfterFuncR))
}

func TestMapRepo_CheckBlackList(t *testing.T) {
	testKey := "192.168.0.0"

	type testCases struct {
		BlackList map[string]time.Time
		Result    bool
	}

	cases := []testCases{
		{BlackList: map[string]time.Time{testKey: time.Unix(1, 0)}, Result: true},
		{BlackList: map[string]time.Time{}, Result: true},
		{BlackList: map[string]time.Time{testKey: time.Now().Add(time.Hour)}, Result: false},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("CheckBlackList_Test %d", i), func(t *testing.T) {
			mapRepo := MapRepo{
				blackList: c.BlackList,
				muBL:      &sync.RWMutex{},
			}
			res := mapRepo.CheckBlackList(testKey)
			assert.Equal(t, c.Result, res)
		})
	}
}

func TestMapRepo_CountAndAddRequest(t *testing.T) {
	testKey := "192.168.0.0"
	requestTTl := time.Hour
	requestLimit := 2

	type testCases struct {
		Requests map[string][]time.Time
		Result   bool
	}

	cases := []testCases{
		{Requests: map[string][]time.Time{}, Result: true},
		{Requests: map[string][]time.Time{testKey: {}}, Result: true},
		{Requests: map[string][]time.Time{testKey: {time.Unix(1, 0)}}, Result: true},
		{Requests: map[string][]time.Time{testKey: {time.Now().Add(time.Hour), time.Now().Add(time.Hour)}}, Result: false},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("CountAndAddRequest_Test %d", i), func(t *testing.T) {
			mapRepo := MapRepo{
				requests:      c.Requests,
				RequestTTL:    requestTTl,
				RequestsLimit: requestLimit,
				muR:           &sync.RWMutex{},
			}
			res := mapRepo.CountAndAddRequest(testKey)
			assert.Equal(t, c.Result, res)
		})
	}
}

func TestMapRepo_ResetLimitByPrefix(t *testing.T) {
	prefix := "192.168"

	testKey := "192.168.0.1"
	testKey2 := "192.168.0.0"
	testKey3 := "195.175.0.0"

	type testCases struct {
		Prefix            string
		Requests          map[string][]time.Time
		RequestsExpected  map[string][]time.Time
		BlackList         map[string]time.Time
		BlackListExpected map[string]time.Time
	}

	cases := []testCases{
		{
			Prefix:            prefix,
			Requests:          map[string][]time.Time{},
			RequestsExpected:  map[string][]time.Time{},
			BlackList:         map[string]time.Time{},
			BlackListExpected: map[string]time.Time{},
		},
		{
			Prefix:            prefix,
			Requests:          map[string][]time.Time{testKey: {}},
			RequestsExpected:  map[string][]time.Time{},
			BlackList:         map[string]time.Time{},
			BlackListExpected: map[string]time.Time{},
		},
		{
			Prefix:            prefix,
			Requests:          map[string][]time.Time{testKey: {}, testKey2: {}},
			RequestsExpected:  map[string][]time.Time{},
			BlackList:         map[string]time.Time{testKey: {}, testKey3: {}},
			BlackListExpected: map[string]time.Time{testKey3: {}},
		},
		{
			Prefix:            prefix,
			Requests:          map[string][]time.Time{testKey3: {}},
			RequestsExpected:  map[string][]time.Time{testKey3: {}},
			BlackList:         map[string]time.Time{testKey: {}, testKey3: {}},
			BlackListExpected: map[string]time.Time{testKey3: {}},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("ResetLimitByPrefix_Test %d", i), func(t *testing.T) {
			mapRepo := &MapRepo{
				requests:  c.Requests,
				blackList: c.BlackList,
				muR:       &sync.RWMutex{},
				muBL:      &sync.RWMutex{},
			}
			mapRepo.ResetLimitByPrefix(c.Prefix)

			//нужно для того, чтобы подождать конец выполнения горутины
			time.Sleep(time.Millisecond * 50)

			assert.True(t, reflect.DeepEqual(mapRepo.blackList, c.BlackListExpected))
			assert.True(t, reflect.DeepEqual(mapRepo.requests, c.RequestsExpected))
		})
	}
}
