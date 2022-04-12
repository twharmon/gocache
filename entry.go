package gocache

import "time"

type entry[V any] struct {
	val V
	ts  time.Time
	exp bool
}

func newEntry[V any](v V) *entry[V] {
	return &entry[V]{
		val: v,
		ts:  time.Now(),
	}
}

func (e *entry[V]) isExpired(reference time.Time) bool {
	return e.exp && e.ts.Before(reference)
}
