package caching

import (
	"strings"
	"sync"

	"github.com/aisalamdag23/promptme-cli/internal/infrastructure/config"
)

type (
	InMemoryCache struct {
		Data map[string]string
		mu   sync.RWMutex
	}

	Cache interface {
		Get(key string) (string, bool)
		Set(key, value string)
	}
)

// NewInMemory init new in memory caching
func NewInMemory(cfg *config.Config) Cache {
	return &InMemoryCache{
		Data: make(map[string]string),
	}
}

// Get retrieve if key is existing
func (c *InMemoryCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, exists := c.Data[strings.ToLower(strings.TrimSpace(key))]; exists {
		return value, true
	}

	return "", false
}

// Set sets key value
func (c *InMemoryCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Data[strings.ToLower(strings.TrimSpace(key))] = strings.TrimSpace(value)
}
