package cache

import (
	"net/http"
	"strings"
	"sync"
)

type CachedResponse struct {
	Status  int
	Headers http.Header
	Body    []byte
}

type CacheMap struct {
	entries map[string]CachedResponse
	mutex   sync.RWMutex
}

func NewCacheMap() *CacheMap {
	return &CacheMap{
		entries: map[string]CachedResponse{},
		mutex:   sync.RWMutex{},
	}
}

func (cm *CacheMap) Get(key string) (*CachedResponse, bool) {
	cm.mutex.RLock()
	val, ok := cm.entries[strings.ToLower(key)]
	cm.mutex.RUnlock()

	return &val, ok
}

func (cm *CacheMap) Set(key string, value *CachedResponse) {
	cm.mutex.Lock()
	cm.entries[strings.ToLower(key)] = *value
	cm.mutex.Unlock()
}

func (cm *CacheMap) Flush() {
	cm.mutex.Lock()
	clear(cm.entries)
	cm.mutex.Unlock()
}
