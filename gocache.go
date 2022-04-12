package gocache

import (
	"sync"
	"time"
)

var never = time.Now().Add(time.Hour * 24 * 365 * 100)

type Cache[K comparable, V any] struct {
	store  map[K]*entry[V]
	mu     sync.Mutex
	config *Config
}

func New[K comparable, V any](config ...*Config) *Cache[K, V] {
	c := &Cache[K, V]{}
	if len(config) > 0 {
		c.config = config[0]
	} else {
		c.config = NewConfig()
	}
	c.store = make(map[K]*entry[V], c.config.capacityHint)
	go c.watch()
	return c
}

func (c *Cache[K, V]) watch() {
	if c.config.evictionPolicy&evictionPolicyActivePruningMask == 0 {
		return
	}
	for {
		time.Sleep(time.Minute)
		c.mu.Lock()
		c.expireUnsafe()
		c.mu.Unlock()
	}
}

func (c *Cache[K, V]) expireUnsafe() {
	for k, e := range c.store {
		if e.isExpired() {
			delete(c.store, k)
		}
	}
}

func (c *Cache[K, V]) evict() {
	if c.Size() < c.config.maxCapacity {
		return
	}
	var someExpire bool
	var evict K
	ts := never
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, e := range c.store {
		if someExpire && !e.exp {
			continue
		}
		if !someExpire && e.exp {
			someExpire = true
			ts = e.ts
			evict = k
			continue
		}
		if e.ts.Before(ts) {
			ts = e.ts
			evict = k
		}
	}
	delete(c.store, evict)
}

func (c *Cache[K, V]) Set(key K, val V, ttl ...time.Duration) {
	c.evict()
	e := newEntry(val)
	if len(ttl) > 0 {
		e.ts = e.ts.Add(ttl[0])
		e.exp = true
	} else if c.config.evictionPolicy&evictionPolicySetMask > 0 {
		e.ts = e.ts.Add(c.config.defaultEvictionTtl)
		e.exp = true
	}
	c.mu.Lock()
	c.store[key] = e
	c.mu.Unlock()
}

func (c *Cache[K, V]) Get(key K) V {
	c.mu.Lock()
	defer c.mu.Unlock()
	e := c.store[key]
	if e == nil || e.isExpired() {
		c.expireUnsafe()
		var v V
		return v
	}
	if c.config.evictionPolicy&evictionPolicyGetMask > 0 {
		e.ts = time.Now().Add(c.config.defaultEvictionTtl)
	}
	return e.val
}

func (c *Cache[K, V]) Has(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.store[key]
	if !ok {
		return false
	}
	if e.isExpired() {
		c.expireUnsafe()
		return false
	}
	if c.config.evictionPolicy&evictionPolicyHasMask > 0 {
		e.ts = time.Now().Add(c.config.defaultEvictionTtl)
	}
	return true
}

func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

func (c *Cache[K, V]) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.expireUnsafe()
	return len(c.store)
}

func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[K]*entry[V], c.config.capacityHint)
}

func (c *Cache[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.expireUnsafe()
	keys := make([]K, 0, len(c.store))
	for k := range c.store {
		keys = append(keys, k)
	}
	return keys
}
