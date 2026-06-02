package cache

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/goark/gnkf/newline"
)

type CachedResponse struct {
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
}

// Entries field is made public for the sake of simple json serializaton
// It is not intended to be used directly, since no locks will be enforced
type CacheMap struct {
	Entries map[string]*CachedResponse `json:"entries"`
	mutex   sync.RWMutex
}

func newCacheMap() *CacheMap {
	return &CacheMap{
		Entries: map[string]*CachedResponse{},
		mutex:   sync.RWMutex{},
	}
}

func LoadCacheMap(path string) (*CacheMap, error) {	
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return newCacheMap(), nil
		}
		
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if stat.Size() == 0 {
		return newCacheMap(), nil 
	}

	var cacheMap *CacheMap
	if err := json.NewDecoder(file).Decode(&cacheMap); err != nil {
		return nil, err
	}

	return cacheMap, nil
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
	cm.mutex.RLock()
	data, err := json.Marshal(cm)
	if err != nil {
		return err
	}
	cm.mutex.RUnlock()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (cm *CacheMap) Flush(path string) error {
	// to flush the cash is to override the file
	// os.Create does the job here
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	return nil
}
