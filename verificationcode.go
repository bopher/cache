package cache

import "time"

// VerificationCode interface for verification code
type VerificationCode interface {
	// Set set code
	Set(value string) error
	// Generate generate a random numeric code with 5 character length
	Generate() (string, error)
	// GenerateN generate a random numeric code with special character length
	GenerateN(count uint) (string, error)
	// Clear clear code
	Clear() error
	// Get get code
	Get() (string, error)
	// Exists check if code exists
	Exists() (bool, error)
	// TTL get ttl
	TTL() (time.Duration, error)
}
