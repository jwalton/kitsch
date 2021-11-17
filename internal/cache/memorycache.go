package cache

type memoryCache struct {
	cache map[string][]byte
}

// NewMemoryCache creates an in-memory Cache.
func NewMemoryCache() Cache {
	return &memoryCache{
		cache: map[string][]byte{},
	}
}

// Get returns the value for the given key.  If the value is not found,
// returns nil.
func (cache *memoryCache) Get(key string) []byte {
	return cache.cache[key]
}

// Set sets the value for the given key.
func (cache *memoryCache) Set(key string, value []byte) {
	cacheValue := make([]byte, len(value))
	copy(cacheValue, value)
	cache.cache[key] = cacheValue
}

// Delete deletes the value for the given key.
func (cache *memoryCache) Delete(key string) {
	delete(cache.cache, key)
}
