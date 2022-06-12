package cache

import (
	"time"

	"github.com/bopher/utils"
)

type rLimiter struct {
	key   string
	max   uint32
	ttl   time.Duration
	cache Cache
}

func (rl rLimiter) err(pattern string, params ...interface{}) error {
	return utils.TaggedError([]string{"RateLimiter", rl.key}, pattern, params...)
}

func (rl rLimiter) notExistsErr() error {
	return utils.TaggedError([]string{"RateLimiter", "NotExists", rl.key}, "%s not exists", rl.key)
}

func (rl *rLimiter) init(key string, maxAttempts uint32, ttl time.Duration, cache Cache) error {
	rl.key = key
	rl.max = maxAttempts
	rl.ttl = ttl
	rl.cache = cache

	exists, err := cache.Exists(key)
	if err != nil {
		return rl.err(err.Error())
	}

	if !exists {
		return cache.Put(key, maxAttempts, ttl)
	}

	return nil
}

func (rl rLimiter) Hit() error {
	exists, err := rl.cache.Decrement(rl.key, 1)
	if err != nil {
		return rl.err(err.Error())
	}

	if !exists {
		return rl.notExistsErr()
	}
	return nil
}

func (rl rLimiter) Lock() error {
	exists, err := rl.cache.Set(rl.key, 0)
	if err != nil {
		return rl.err(err.Error())
	}

	if !exists {
		return rl.notExistsErr()
	}

	return nil
}

func (rl rLimiter) Reset() error {
	err := rl.cache.Put(rl.key, rl.max, rl.ttl)
	if err != nil {
		return rl.err(err.Error())
	}

	return nil
}

func (rl rLimiter) Clear() error {
	err := rl.cache.Forget(rl.key)
	if err != nil {
		return rl.err(err.Error())
	}

	return nil
}

func (rl rLimiter) MustLock() (bool, error) {
	caster, err := rl.cache.Cast(rl.key)
	if err != nil {
		return true, rl.err(err.Error())
	}

	if caster.IsNil() {
		return false, nil
	}

	v, err := caster.Int()
	if err != nil {
		err = rl.err(err.Error())
	}
	return v <= 0, err
}

func (rl rLimiter) TotalAttempts() (uint32, error) {
	caster, err := rl.cache.Cast(rl.key)
	if err != nil {
		return rl.max, rl.err(err.Error())
	}

	if caster.IsNil() {
		return rl.max, nil
	}

	v, err := caster.Int()
	if err != nil {
		return rl.max, rl.err(err.Error())
	}

	if v > int(rl.max) {
		v = int(rl.max)
	}

	return rl.max - uint32(v), nil
}

func (rl rLimiter) RetriesLeft() (uint32, error) {
	caster, err := rl.cache.Cast(rl.key)
	if err != nil {
		return 0, rl.err(err.Error())
	}

	if caster.IsNil() {
		return 0, nil
	}

	v, err := caster.Int()
	if err != nil {
		err = rl.err(err.Error())
	}
	if v < 0 {
		v = 0
	}
	return uint32(v), err
}

func (rl rLimiter) AvailableIn() (time.Duration, error) {
	if v, err := rl.cache.TTL(rl.key); err != nil {
		return 0, rl.err(err.Error())
	} else {
		return v, nil
	}
}
