package cachestorage

import (
	"sync"
)

type CacheStorage[Key comparable, Value any] struct {
	data map[Key]Value
	mu   sync.Mutex
}

func (c *CacheStorage[Key, Value]) NumCached() uint {
	c.mu.Lock()
	numCached := len(c.data)
	c.mu.Unlock()
	return uint(numCached)
}

func (c *CacheStorage[Key, Value]) Clear() uint {
	c.mu.Lock()
	numDeleted := len(c.data)
	c.data = make(map[Key]Value)
	c.mu.Unlock()
	return uint(numDeleted)
}

func NewCacheStorage[Key comparable, Value any]() *CacheStorage[Key, Value] {
	return &CacheStorage[Key, Value]{
		data: make(map[Key]Value),
	}
}

func (c *CacheStorage[Key, Value]) Set(key Key, value Value) {
	c.mu.Lock()
	c.data[key] = value
	c.mu.Unlock()
}

func (c *CacheStorage[Key, Value]) Get(key Key) *Value {
	c.mu.Lock()
	value, ok := c.data[key]
	c.mu.Unlock()
	if !ok {
		return nil
	}
	return &value
}
