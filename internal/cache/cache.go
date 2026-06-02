package cache

type Cache interface {
	Get(key string) (*CachedResponse, bool)
	Set(key string, value *CachedResponse)
	Save(path string) error
	Flush(path string) error
}
