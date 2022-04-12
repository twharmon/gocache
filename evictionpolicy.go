package gocache

import "time"

type EvictionPolicyFlags byte

type EvictionPolicy struct {
	flags EvictionPolicyFlags
	ttl   time.Duration
}

const (
	evictionPolicySetMask           EvictionPolicyFlags = 0b00000001
	evictionPolicyGetMask                               = 0b00000010
	evictionPolicyHasMask                               = 0b00000100
	evictionPolicyActivePruningMask                     = 0b10000000
)

// NewEvictionPolicy creates a new EvictionPolicy. By default, TTL is
// not affected  when Get or Has is called.
func NewEvictionPolicy(ttl time.Duration) *EvictionPolicy {
	return &EvictionPolicy{
		flags: evictionPolicySetMask,
		ttl:   ttl,
	}
}

// UpdateOnGet sets the eviction policy to update TTL when Get is
// called.
func (m *EvictionPolicy) UpdateOnGet() *EvictionPolicy {
	m.flags |= evictionPolicyGetMask
	return m
}

// UpdateOnHas sets the eviction policy to update TTL when Has is
// called.
func (m *EvictionPolicy) UpdateOnHas() *EvictionPolicy {
	m.flags |= evictionPolicyHasMask
	return m
}

// ActivePruning sets the eviction policy to regularly scan the cache
// to expire old items.
func (m *EvictionPolicy) ActivePruning() *EvictionPolicy {
	m.flags |= evictionPolicyActivePruningMask
	return m
}
