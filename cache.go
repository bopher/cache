package cache

import "time"

// Cache interface for cache drivers.
type Cache interface {
	// Put a new value to cache
	Put(key string, value interface{}, ttl time.Duration) error
	// PutForever put value with infinite ttl
	PutForever(key string, value interface{}) error
	// Set Change value of cache item
	Set(key string, value interface{}) error
	// Get item from cache
	Get(key string) (interface{}, error)
	// Check if item exists in cache
	Exists(key string) (bool, error)
	// Forget item from cache (delete item)
	Forget(key string) error
	// Pull item from cache and remove it
	Pull(key string) (interface{}, error)
	// TTL get cache item ttl
	TTL(key string) (time.Duration, error)
	// BoolE parse item as boolean or return error on fail
	BoolE(key string) (bool, error)
	// Bool parse item as boolean or return fallback
	Bool(key string, fallback bool) bool
	// IntE parse item as int or return error on fail
	IntE(key string) (int, error)
	// Int parse item as int or return fallback
	Int(key string, fallback int) int
	// Int8E parse item as int8 or return error on fail
	Int8E(key string) (int8, error)
	// Int8 parse item as int8 or return fallback
	Int8(key string, fallback int8) int8
	// Int16E parse item as int16 or return error on fail
	Int16E(key string) (int16, error)
	// Int16 parse item as int16  or return fallback
	Int16(key string, fallback int16) int16
	// Int32E parse item as int32 or return error on fail
	Int32E(key string) (int32, error)
	// Int32 parse item as int32 or return fallback
	Int32(key string, fallback int32) int32
	// Int64E parse item as int64 or return error on fail
	Int64E(key string) (int64, error)
	// Int64 parse item as int64 or return fallback
	Int64(key string, fallback int64) int64
	// UIntE parse item as uint or return error on fail
	UIntE(key string) (uint, error)
	// UInt parse item as uint or return fallback
	UInt(key string, fallback uint) uint
	// UInt8E parse item as uint8 or return error on fail
	UInt8E(key string) (uint8, error)
	// UInt8 parse item as uint8 or return fallback
	UInt8(key string, fallback uint8) uint8
	// UInt16E parse item as uint16 or return error on fail
	UInt16E(key string) (uint16, error)
	// UInt16 parse item as uint16 or return fallback
	UInt16(key string, fallback uint16) uint16
	// UInt32E parse item as uint32 or return error on fail
	UInt32E(key string) (uint32, error)
	// UInt32 parse item as uint32 or return fallback
	UInt32(key string, fallback uint32) uint32
	// UInt64E parse item as uint64 or return error on fail
	UInt64E(key string) (uint64, error)
	// UInt64 parse item as uint64 or return fallback
	UInt64(key string, fallback uint64) uint64
	// Float64E parse item as float64 or return error on fail
	Float64E(key string) (float64, error)
	// Float64 parse item as float64 or return fallback
	Float64(key string, fallback float64) float64
	// StringE parse item as string or return error on fail
	StringE(key string) (string, error)
	// String parse item as string or return fallback
	String(key string, fallback string) string
	// IncrementBy numeric item in cache by number
	IncrementBy(key string, value interface{}) error
	// Increment numeric item in cache
	Increment(key string) error
	// DecrementBy numeric item in cache by number
	DecrementBy(key string, value interface{}) error
	// Decrement numeric item in cache
	Decrement(key string) error
}
