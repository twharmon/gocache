package gocache

import (
	"math"
	"time"
)

type EvictionPolicy byte

type Config struct {
	capacityHint       int
	defaultEvictionTtl time.Duration
	maxCapacity        int
	evictionPolicy     EvictionPolicy
}

// NewConfig creates a new config with default settings.
func NewConfig() *Config {
	return &Config{
		capacityHint: 8,
		maxCapacity:  math.MaxInt,
	}
}

// WithCapacityHint sets the capacity hint for the cache.
func (c *Config) WithCapacityHint(size int) *Config {
	c.capacityHint = size
	return c
}

// WithDefaultEvictionPolicy sets the default eviction policy.
func (c *Config) WithDefaultEvictionPolicy(ttl time.Duration, mode *EvictionPolicy) *Config {
	c.defaultEvictionTtl = ttl
	c.evictionPolicy = *mode
	return c
}

// WithMaxCapacity sets the max capacity of the cache. If max
// capacity is reached, items expiring soonest are removed from the
// cache. If no items have a TTL, the oldest items are removed from
// the cache.
func (c *Config) WithMaxCapacity(size int) *Config {
	c.maxCapacity = size
	return c
}

const (
	evictionPolicySetMask           = 0b00000001
	evictionPolicyGetMask           = 0b00000010
	evictionPolicyHasMask           = 0b00000100
	evictionPolicyActivePruningMask = 0b10000000
)

// NewEvictionPolicy creates a new EvictionPolicy. By default, TTL is
// not affected  when Get or Has is called.
func NewEvictionPolicy() *EvictionPolicy {
	mode := EvictionPolicy(evictionPolicySetMask)
	return &mode
}

// UpdateOnGet sets the eviction policy to update TTL when Get is
// called.
func (m *EvictionPolicy) UpdateOnGet() *EvictionPolicy {
	*m |= evictionPolicyGetMask
	return m
}

// UpdateOnHas sets the eviction policy to update TTL when Has is
// called.
func (m *EvictionPolicy) UpdateOnHas() *EvictionPolicy {
	*m |= evictionPolicyHasMask
	return m
}

// ActivePruning sets the eviction policy to regularly scan the cache
// to expire old items.
func (m *EvictionPolicy) ActivePruning() *EvictionPolicy {
	*m |= evictionPolicyActivePruningMask
	return m
}
