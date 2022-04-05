package gocache

import (
	"math"
	"time"
)

type TTLMode byte

type Config struct {
	capacityHint int
	defaultTTL   time.Duration
	maxCapacity  int
	ttlMode      TTLMode
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

// WithDefaultTTL sets the default TTL settings.
func (c *Config) WithDefaultTTL(ttl time.Duration, mode *TTLMode) *Config {
	c.defaultTTL = ttl
	c.ttlMode = *mode
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
	ttlModeSetMask           = 0b00000001
	ttlModeGetMask           = 0b00000010
	ttlModeHasMask           = 0b00000100
	ttlModeActivePruningMask = 0b10000000
)

// NewTTLMode creates a new TTL mode. By default, TTL is not affected
// when Get or Has is called.
func NewTTLMode() *TTLMode {
	mode := TTLMode(ttlModeSetMask)
	return &mode
}

// UpdateOnGet sets the mode to update TTL when Get is called.
func (m *TTLMode) UpdateOnGet() *TTLMode {
	*m |= ttlModeGetMask
	return m
}

// UpdateOnHas sets the mode to update TTL when Has is called.
func (m *TTLMode) UpdateOnHas() *TTLMode {
	*m |= ttlModeHasMask
	return m
}

// ActivePruning sets the mode to regularly scan the cache to expire
// old items.
func (m *TTLMode) ActivePruning() *TTLMode {
	*m |= ttlModeActivePruningMask
	return m
}
