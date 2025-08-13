package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Cache interface {
	IncrWithExpire(key string, expireSeconds int) (int, error)
	SetWithTTL(key string, value string, ttlSeconds int) error
	Get(key string) (string, error)
	Delete(key string) error
}

type InMemoryCache struct {
	mu   sync.RWMutex
	data map[string]cacheItem
}

type cacheItem struct {
	value      string
	expireTime time.Time
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]cacheItem),
	}
}

func (c *InMemoryCache) IncrWithExpire(key string, expireSeconds int) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.data[key]
	if !exists || time.Now().After(item.expireTime) {
		c.data[key] = cacheItem{value: "1", expireTime: time.Now().Add(time.Duration(expireSeconds) * time.Second)}
		return 1, nil
	}
	val := 0
	fmt.Sscanf(item.value, "%d", &val)
	val++
	c.data[key] = cacheItem{value: fmt.Sprintf("%d", val), expireTime: item.expireTime}
	return val, nil
}

func (c *InMemoryCache) SetWithTTL(key string, value string, ttlSeconds int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheItem{
		value:      value,
		expireTime: time.Now().Add(time.Duration(ttlSeconds) * time.Second),
	}
	return nil
}

func (c *InMemoryCache) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, exists := c.data[key]
	if !exists || time.Now().After(item.expireTime) {
		return "", errors.New("key not found or expired")
	}
	return item.value, nil
}

func (c *InMemoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}
