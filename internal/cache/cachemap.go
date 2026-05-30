package cache

import (
	"strings"
	"sync"
)

type CacheMap struct {
	body  map[string][]byte
	mutex sync.RWMutex
}

func NewCacheMap() CacheMap {
	return CacheMap{
		body:  map[string][]byte{},
		mutex: sync.RWMutex{},
	}
}

func (cm *CacheMap) Get(key string) ([]byte, bool) {
	cm.mutex.RLock()
	val, ok := cm.body[strings.ToLower(key)]
	cm.mutex.RUnlock()

	return val, ok
}

func (cm *CacheMap) Set(key string, value []byte) {
	cm.mutex.Lock()
	cm.body[strings.ToLower(key)] = value
	cm.mutex.Unlock()
}

func (cm *CacheMap) Flush() {
	cm.mutex.Lock()
	clear(cm.body)
	cm.mutex.Unlock()
}
