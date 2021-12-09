package cache

import (
	"time"

	"github.com/bopher/caster"
)

// Cache interface for cache drivers.
type Cache interface {
	// Put a new value to cache
	Put(key string, value interface{}, ttl time.Duration) error
	// PutForever put value with infinite ttl
	PutForever(key string, value interface{}) error
	// Set Change value of cache item, return false if item not exists
	Set(key string, value interface{}) (bool, error)
	// Get item from cache
	Get(key string) (interface{}, error)
	// Exists check if item exists in cache
	Exists(key string) (bool, error)
	// Forget delete Item from cache
	Forget(key string) error
	// Pull item from cache and remove it
	Pull(key string) (interface{}, error)
	// TTL get cache item ttl. this method returns -1 if item not exists
	TTL(key string) (time.Duration, error)
	// Cast parse cache item as caster
	Cast(key string) (caster.Caster, error)
	// IncrementFloat increment numeric item by float, return false if item not exists
	IncrementFloat(key string, value float64) (bool, error)
	// Increment increment numeric item by int, return false if item not exists
	Increment(key string, value int64) (bool, error)
	// DecrementFloat decrement numeric item by float, return false if item not exists
	DecrementFloat(key string, value float64) (bool, error)
	// Decrement decrement numeric item by int, return false if item not exists
	Decrement(key string, value int64) (bool, error)
}
