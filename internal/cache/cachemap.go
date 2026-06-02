package cache

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
)

type CachedResponse struct {
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
}

type CacheMap struct {
	Entries map[string]*CachedResponse `json:"entries"`
	mutex   sync.RWMutex
}

func NewCacheMap() *CacheMap {
	return &CacheMap{
		Entries: map[string]*CachedResponse{},
		mutex:   sync.RWMutex{},
	}
}

func (cm *CacheMap) Get(key string) (*CachedResponse, bool) {
	cm.mutex.RLock()
	val, ok := cm.Entries[strings.ToLower(key)]
	cm.mutex.RUnlock()

	return val, ok
}

func (cm *CacheMap) Set(key string, value *CachedResponse) {
	cm.mutex.Lock()
	cm.Entries[strings.ToLower(key)] = value
	cm.mutex.Unlock()
}

func (cm *CacheMap) Length() int {
	cm.mutex.RLock()
	length := len(cm.Entries)
	cm.mutex.RUnlock()
	return length
}

func (cm *CacheMap) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// this omits the key for now, testing
	data, err := json.Marshal(cm)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (cm *CacheMap) Load(path string) error {
	return nil
}

func (cm *CacheMap) Flush() {
	cm.mutex.Lock()
	clear(cm.Entries)
	cm.mutex.Unlock()
}
