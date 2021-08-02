package cache

import "time"

// RateLimiter interface for rate limiter
type RateLimiter interface {
	// Hit decrease the allowed times
	Hit()
	// Lock lock rate limiter
	Lock()
	// Reset reset rate limiter
	Reset()
	// MustLock check if rate limiter must lock access
	MustLock() bool
	// TotalAttempts get user attempts count
	TotalAttempts() uint32
	// RetriesLeft get user retries left
	RetriesLeft() uint32
	// AvailableIn get time until unlock
	AvailableIn() time.Duration
}
