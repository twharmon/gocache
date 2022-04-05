package gocache

import (
	"sync"
	"time"
)

var never = time.Now().Add(time.Hour * 24 * 365 * 100)

type Key interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | string
}

type Cache[K Key, V any] struct {
	store  map[K]*entry[V]
	mu     sync.Mutex
	config *Config
}

type entry[V any] struct {
	val V
	ts  time.Time
	exp bool
}

func New[K Key, V any](config ...*Config) *Cache[K, V] {
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
	if c.config.ttlMode&ttlModeActivePruningMask == 0 {
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
		if isExpired(e) {
			delete(c.store, k)
		}
	}
}

func (c *Cache[K, V]) expel() {
	if c.Size() < c.config.maxCapacity {
		return
	}
	var someExpire bool
	var expel K
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
			expel = k
			continue
		}
		if e.ts.Before(ts) {
			ts = e.ts
			expel = k
		}
	}
	delete(c.store, expel)
}

func (c *Cache[K, V]) Set(key K, val V, ttl ...time.Duration) {
	c.expel()
	e := entry[V]{
		val: val,
	}
	e.ts = time.Now()
	if len(ttl) > 0 {
		e.ts = e.ts.Add(ttl[0])
		e.exp = true
	} else if c.config.ttlMode&ttlModeSetMask > 0 {
		e.ts = e.ts.Add(c.config.defaultTTL)
		e.exp = true
	}
	c.mu.Lock()
	c.store[key] = &e
	c.mu.Unlock()
}

func (c *Cache[K, V]) Get(key K) V {
	c.mu.Lock()
	defer c.mu.Unlock()
	e := c.store[key]
	if e == nil || isExpired(e) {
		c.expireUnsafe()
		var v V
		return v
	}
	if c.config.ttlMode&ttlModeGetMask > 0 {
		e.ts = time.Now().Add(c.config.defaultTTL)
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
	if isExpired(e) {
		c.expireUnsafe()
		return false
	}
	if c.config.ttlMode&ttlModeHasMask > 0 {
		e.ts = time.Now().Add(c.config.defaultTTL)
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
