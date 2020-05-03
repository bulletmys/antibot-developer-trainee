package repository

import (
	"strings"
	"sync"
	"time"
)

type MapRepo struct {
	blackList     map[string]time.Time
	requests      map[string][]time.Time
	BlackListTTL  time.Duration
	RequestTTL    time.Duration
	RequestsLimit int
	muBL          *sync.RWMutex
	muR           *sync.RWMutex
}

func NewMapRepo(requestsLimit int, blackListTTL, requestTTL time.Duration) *MapRepo {
	repo := &MapRepo{
		blackList:     make(map[string]time.Time),
		requests:      make(map[string][]time.Time),
		BlackListTTL:  blackListTTL,
		RequestTTL:    requestTTL,
		RequestsLimit: requestsLimit,
		muBL:          &sync.RWMutex{},
		muR:           &sync.RWMutex{},
	}

	go repo.cleanupRequests()
	go repo.cleanupBlackList()

	return repo
}

func (r *MapRepo) cleanupRequests() {
	for {
		time.Sleep(r.RequestTTL)

		currTime := time.Now()
		r.muR.Lock()
		for ip, v := range r.requests {
			reqStampsLen := len(v)
			if reqStampsLen > 0 && currTime.After(v[reqStampsLen-1]) {
				delete(r.requests, ip)
			}
		}
		r.muR.Unlock()
	}
}

func (r *MapRepo) cleanupBlackList() {
	for {
		time.Sleep(r.BlackListTTL)

		currTime := time.Now()
		r.muBL.Lock()
		for ip, v := range r.blackList {
			if currTime.After(v) {
				delete(r.blackList, ip)
			}
		}
		r.muBL.Unlock()
	}
}

func (r *MapRepo) CheckBlackList(key string) bool {
	r.muBL.Lock()
	timeStamp, ok := r.blackList[key]
	if !ok || timeStamp.Before(time.Now()) {
		r.muBL.Unlock()
		return true
	}

	r.muBL.Unlock()
	return false
}

// CountAndAddRequest проверяет лимит запросов, чистит старые записи и добавляет пришедший запрос
func (r *MapRepo) CountAndAddRequest(key string) (bool, error) {
	r.muR.Lock()
	defer r.muR.Unlock()

	result, ok := r.requests[key]

	if !ok {
		r.requests[key] = []time.Time{time.Now().Add(r.RequestTTL)}
		return true, nil
	}

	currTime := time.Now()

	counter := 0
	for _, elem := range result {
		if currTime.After(elem) {
			counter++
		}
	}

	if len(result)-counter >= r.RequestsLimit {
		return false, nil
	}

	r.requests[key] = append(result[counter:], currTime.Add(r.RequestTTL))

	return true, nil
}

func (r *MapRepo) AddToBlackList(key string) {
	r.muBL.Lock()
	r.blackList[key] = time.Now().Add(r.BlackListTTL)
	delete(r.requests, key)
	r.muBL.Unlock()
}

func (r *MapRepo) ResetLimitByPrefix(prefix string) {
	go func() {
		r.muR.Lock()
		for key, _ := range r.requests {
			if strings.HasPrefix(key, prefix) {
				delete(r.requests, key)
			}
		}
		r.muR.Unlock()
	}()
	r.muBL.Lock()
	for key, _ := range r.blackList {
		if strings.HasPrefix(key, prefix) {
			delete(r.blackList, key)
		}
	}
	r.muBL.Unlock()
}

func (r * MapRepo) GetBlackListTTL() time.Duration {
	return r.BlackListTTL
}
