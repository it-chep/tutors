package cache

import (
	"time"

	"github.com/hashicorp/golang-lru/arc/v2"
)

type Cache[K comparable, V any] struct {
	arc *arc.ARCCache[K, entry[V]]

	cap int
	ttl time.Duration
}

type entry[V any] struct {
	value   V
	expires time.Time
}

func NewCache[K comparable, V any](cap int, ttl time.Duration) *Cache[K, V] {
	arcInstance, _ := arc.NewARC[K, entry[V]](cap)
	c := &Cache[K, V]{
		arc: arcInstance,
		cap: cap,
		ttl: ttl,
	}

	return c
}

func (c *Cache[K, V]) Put(key K, value V) {
	c.putWithTTL(key, value, c.ttl)
}

func (c *Cache[K, V]) putWithTTL(key K, value V, ttl time.Duration) {
	var expires time.Time
	if ttl > 0 {
		expires = time.Now().Add(ttl)
	}

	c.arc.Add(key, entry[V]{
		expires: expires,
		value:   value,
	})
}

func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	v, ok := c.arc.Get(key)
	if ok {
		if v.expires.IsZero() || time.Now().Before(v.expires) {
			return v.value, true
		}
	}
	var zero V
	return zero, false
}

func (c *Cache[K, V]) Remove(key K) {
	c.arc.Remove(key)
}
