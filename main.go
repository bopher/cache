package cache

import (
	"time"

	"github.com/go-redis/redis/v8"
)

// NewRedisCache create a new redis cache manager instance
func NewRedisCache(prefix string, opt redis.Options) Cache {
	rc := new(rCache)
	rc.init(prefix, opt)
	return rc
}

// NewFileCache create a new file cache manager instance
func NewFileCache(prefix string, dir string) Cache {
	fc := new(fCache)
	fc.init(prefix, dir)
	return fc
}

// NewRateLimiter create a new rate limiter
func NewRateLimiter(key string, maxAttempts uint32, ttl time.Duration, cache Cache) (RateLimiter, error) {
	limiter := new(rLimiter)
	if err := limiter.init(key, maxAttempts, ttl, cache); err != nil {
		return nil, err
	} else {
		return limiter, nil
	}
}

// NewVerificationCode create a new verification code manager instance
func NewVerificationCode(key string, ttl time.Duration, cache Cache) VerificationCode {
	vc := new(vcDriver)
	vc.Key = key
	vc.TTL = ttl
	vc.Cache = cache
	return vc
}
