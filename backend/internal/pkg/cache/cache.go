// Package cache provides a minimal in-process TTL cache used as a stand-in for
// Redis. It is safe for concurrent use via sync.Map. For multi-instance
// deployments this should be replaced by a shared Redis instance.
package cache

import (
	"sync"
	"time"
)

type entry struct {
	value     interface{}
	expiresAt time.Time
}

// TTLCache is a string-keyed cache with per-key expiry. Expired entries are
// evicted lazily on access.
type TTLCache struct {
	m   sync.Map
	ttl time.Duration
}

func NewTTL(ttl time.Duration) *TTLCache {
	return &TTLCache{ttl: ttl}
}

func (c *TTLCache) Get(key string) (interface{}, bool) {
	v, ok := c.m.Load(key)
	if !ok {
		return nil, false
	}
	e := v.(entry)
	if time.Now().After(e.expiresAt) {
		c.m.Delete(key)
		return nil, false
	}
	return e.value, true
}

func (c *TTLCache) Set(key string, value interface{}) {
	c.m.Store(key, entry{value: value, expiresAt: time.Now().Add(c.ttl)})
}

func (c *TTLCache) Delete(key string) {
	c.m.Delete(key)
}
