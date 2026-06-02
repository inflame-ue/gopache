package cache

type Cache interface {
	Get(key string) (*CachedResponse, bool)
	Set(key string, value *CachedResponse)
	Save(path string) error
	Load(path string) error
	Flush()
}
