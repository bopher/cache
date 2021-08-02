package cache

import (
	"time"
)

// NewRedisCache create a new redis cache manager instance
func NewRedisCache(prefix string, host string, maxIdle int, maxActive int, db uint8) Cache {
	rc := new(redisCache)
	rc.init(prefix, host, maxIdle, maxActive, db)
	return rc
}

// NewFileCache create a new file cache manager instance
func NewFileCache(prefix string, dir string) Cache {
	fc := new(fileCache)
	fc.init(prefix, dir)
	return fc
}

// NewRateLimiter create a new rate limiter
func NewRateLimiter(key string, maxAttempts uint32, ttl time.Duration, cache Cache) RateLimiter {
	limiter := new(rateLimiterDriver)
	limiter.init(key, maxAttempts, ttl, cache)
	return limiter
}

// NewVerificationCode create a new verification code manager instance
func NewVerificationCode(key string, ttl time.Duration, cache Cache) VerificationCode {
	vc := new(vcDriver)
	vc.Key = key
	vc.TTL = ttl
	vc.Cache = cache
	return vc
}
