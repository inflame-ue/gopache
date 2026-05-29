package cache

import (
	"strings"
	"sync"
)

type CacheMap struct {
	data map[string][]byte
	mutex sync.RWMutex 
}

func NewCacheMap() CacheMap {
	return CacheMap{
		data: map[string][]byte{},
		mutex: sync.RWMutex{},
	}
}

func (cm *CacheMap) Get(key string) ([]byte, bool) {
	cm.mutex.RLock()
	val, ok := cm.data[strings.ToLower(key)]
	cm.mutex.RUnlock()
	
	return val, ok
}

func (cm *CacheMap) Set(key string, value []byte) {
	cm.mutex.Lock()
	cm.data[strings.ToLower(key)] = value
	cm.mutex.Unlock()
}

func (cm *CacheMap) Flush() {
	cm.mutex.Lock()
	clear(cm.data)
	cm.mutex.Unlock()
}
