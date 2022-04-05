package gocache

import "time"

func isExpired[V any](e *entry[V]) bool {
	return e.exp && e.ts.Before(time.Now())
}
