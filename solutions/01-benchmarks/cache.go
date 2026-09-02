package benchmarks

import "sync"

type Cache interface {
	Get(string) (string, bool)
	Set(string, string)
}

type RWCache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewRWCache() *RWCache {
	return &RWCache{data: make(map[string]string)}
}

func (c *RWCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.data[key]
	return v, ok
}

func (c *RWCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

type SyncMapCache struct {
	data sync.Map
}

func NewSyncMapCache() *SyncMapCache {
	return &SyncMapCache{}
}

func (c *SyncMapCache) Get(key string) (string, bool) {
	v, ok := c.data.Load(key)
	if !ok {
		return "", false
	}
	return v.(string), true
}

func (c *SyncMapCache) Set(key, value string) {
	c.data.Store(key, value)
}
