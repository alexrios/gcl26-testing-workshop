package synctesting

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]entry
}

func New() *Cache {
	return &Cache{data: make(map[string]entry)}
}

func (c *Cache) SetWithTTL(key, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = entry{value: value, expiresAt: time.Now().Add(ttl)}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	e, ok := c.data[key]
	if !ok || !time.Now().Before(e.expiresAt) {
		return "", false
	}
	return e.value, true
}

// Size conta o estado físico. Diferentemente de Get, não esconde entradas que
// expiraram mas ainda não foram removidas pelo janitor.
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

func (c *Cache) StartJanitor(ctx context.Context, interval time.Duration) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.removeExpired()
			case <-ctx.Done():
				return
			}
		}
	}()
	return done
}

func (c *Cache) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, e := range c.data {
		if !now.Before(e.expiresAt) {
			delete(c.data, key)
		}
	}
}

func NewHandler(c *Cache) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /cache/{key}", func(w http.ResponseWriter, r *http.Request) {
		ttl, err := time.ParseDuration(r.URL.Query().Get("ttl"))
		if err != nil || ttl <= 0 {
			http.Error(w, "invalid ttl", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
		if err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		c.SetWithTTL(r.PathValue("key"), string(body), ttl)
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET /cache/{key}", func(w http.ResponseWriter, r *http.Request) {
		value, ok := c.Get(r.PathValue("key"))
		if !ok {
			http.NotFound(w, r)
			return
		}
		_, _ = io.Copy(w, strings.NewReader(value))
	})
	return mux
}
