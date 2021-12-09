package cache

import "time"

// RateLimiter interface for rate limiter
type RateLimiter interface {
	// Hit decrease the allowed times
	Hit() error
	// Lock lock rate limiter
	Lock() error
	// Reset reset rate limiter
	Reset() error
	// Clear remove rate limiter record
	Clear() error
	// MustLock check if rate limiter must lock access
	MustLock() (bool, error)
	// TotalAttempts get user attempts count
	TotalAttempts() (uint32, error)
	// RetriesLeft get user retries left
	RetriesLeft() (uint32, error)
	// AvailableIn get time until unlock
	AvailableIn() (time.Duration, error)
}
