package gocache

import (
	"math"
)

// Config controls how the cache is configured.
type Config struct {
	capacityHint   int
	maxCapacity    int
	evictionPolicy *EvictionPolicy
}

// NewConfig creates a new config with default settings.
func NewConfig() *Config {
	return &Config{
		capacityHint:   8,
		maxCapacity:    math.MaxInt,
		evictionPolicy: &EvictionPolicy{},
	}
}

// WithCapacityHint sets the capacity hint for the cache.
func (c *Config) WithCapacityHint(size int) *Config {
	c.capacityHint = size
	return c
}

// WithDefaultEvictionPolicy sets the default eviction policy for
// values in the cache. The default policy ttl can be overridden by
// passing a third argument to Set.
func (c *Config) WithDefaultEvictionPolicy(policy *EvictionPolicy) *Config {
	c.evictionPolicy = policy
	return c
}

// WithMaxCapacity sets the max capacity of the cache. If max
// capacity is reached, items expiring soonest are removed from the
// cache. If no items have a ttl, the oldest items are removed from
// the cache. If some items have a ttl and some do not, the items
// with a ttl are removed first.
func (c *Config) WithMaxCapacity(size int) *Config {
	c.maxCapacity = size
	return c
}
