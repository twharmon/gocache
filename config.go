package gocache

import (
	"math"
)

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

// WithDefaultEvictionPolicy sets the default eviction policy.
func (c *Config) WithDefaultEvictionPolicy(policy *EvictionPolicy) *Config {
	c.evictionPolicy = policy
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
