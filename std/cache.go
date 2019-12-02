// Package std represents the interface caching implementation.
package std

import (
	"github.com/golang-mixins/caching"
	"golang.org/x/xerrors"
	"sync"
	"time"
)

// Cache predetermines the consistency of the interfaces caching implementation.
type Cache struct {
	defaultValidity time.Duration
	mutex           *sync.RWMutex
	storage         map[interface{}]caching.Item
}

func (c *Cache) expirationControl(key interface{}, expiration time.Time) {
	defer func() {
		if failure := recover(); failure != nil {
			defer c.delete(key)
		}
	}()

	expiration = expiration.UTC()
	now := time.Now().UTC()

	if expiration.IsZero() || expiration.Before(now) {
		c.delete(key)
		return
	}

	time.Sleep(expiration.Sub(now))

	item, ok := c.load(key)
	if !ok {
		return
	}

	if item.Expiration.UTC().After(time.Now().UTC()) {
		return
	}

	c.delete(key)
}

func (c *Cache) delete(key interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.storage, key)
}

func (c *Cache) load(key interface{}) (item caching.Item, ok bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, ok = c.storage[key]
	return item, ok
}

func (c *Cache) store(key interface{}, item caching.Item) {
	if item.Expiration.IsZero() {
		item.Expiration = time.Now().Add(c.defaultValidity)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.storage[key] = item

	go c.expirationControl(key, item.Expiration)
}

// Get - provides an item by key and a sign of its being in the cache.
func (c *Cache) Get(key interface{}) (item caching.Item, ok bool) {
	return c.load(key)
}

// Set - sets an item by key regardless of whether the item is in the cache.
func (c *Cache) Set(key interface{}, item caching.Item) {
	c.store(key, item)
}

// Add - adds an item to the cache only if the item by key is missing.
func (c *Cache) Add(key interface{}, item caching.Item) {
	_, ok := c.load(key)
	if ok {
		return
	}

	c.store(key, item)
}

// Flush - flush cache.
func (c *Cache) Flush() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.storage = make(map[interface{}]caching.Item)
}

// New - Cache constructor.
func New(defaultValidity time.Duration) (*Cache, error) {
	if defaultValidity == 0 {
		return nil, xerrors.New("argument 'defaultValidity' can't be empty")
	}

	return &Cache{
		defaultValidity,
		new(sync.RWMutex),
		make(map[interface{}]caching.Item),
	}, nil
}
