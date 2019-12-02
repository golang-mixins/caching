// Package caching presents interface (and its implementation sets) of a caching with item storage.
package caching

import "time"

// Item - typifies cache item.
type Item struct {
	// Value - item value.
	Value interface{}
	// Expiration - item expiration.
	Expiration time.Time
}

// IsExpired - predicate characterizing the expired item in the repository at a given time or not.
func (e *Item) IsExpired() bool {
	return time.Now().UTC().After(e.Expiration.UTC())
}

// Cache - provides a caching interface.
type Cache interface {
	// Get - provides an item by key and a sign of its being in the cache.
	Get(key interface{}) (item Item, ok bool)
	// Set - sets an item by key regardless of whether the item is in the cache.
	Set(key interface{}, item Item)
	// Add - adds an item to the cache only if the item by key is missing.
	Add(key interface{}, item Item)
	// Flush - flush cache.
	Flush()
}
