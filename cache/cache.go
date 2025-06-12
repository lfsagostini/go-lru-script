package cache

import (
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)


type cacheEntry[V any] struct {
	value     V
	expiresAt time.Time
}


type Cache[K comparable, V any] struct {
	mu  sync.RWMutex
	lru *lru.Cache[K, *cacheEntry[V]] 
	ttl time.Duration
}


func New[K comparable, V any](capacity int, ttl time.Duration) (*Cache[K, V], error) {
	
	lruCache, err := lru.New[K, *cacheEntry[V]](capacity)
	if err != nil {
		return nil, err
	}

	return &Cache[K, V]{
		lru: lruCache,
		ttl: ttl,
	}, nil
}


func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := &cacheEntry[V]{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
	
	c.lru.Add(key, entry)
}


func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	entry, found := c.lru.Get(key)
	c.mu.RUnlock()
	if !found {
		var zero V
		return zero, false
	}

	if time.Now().After(entry.expiresAt) {
		c.mu.Lock()
		c.lru.Remove(key)
		c.mu.Unlock()
		var zero V
		return zero, false
	}

	return entry.value, true
}


func (c *Cache[K, V]) Len() int {
	return c.lru.Len()
}