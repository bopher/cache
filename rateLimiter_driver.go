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

func (this rLimiter) err(pattern string, params ...interface{}) error {
	return utils.TaggedError([]string{"RateLimiter", this.key}, pattern, params...)
}

func (this rLimiter) notExistsErr() error {
	return utils.TaggedError([]string{"RateLimiter", "NotExists", this.key}, "%s not exists", this.key)
}

func (this *rLimiter) init(key string, maxAttempts uint32, ttl time.Duration, cache Cache) error {
	this.key = key
	this.max = maxAttempts
	this.ttl = ttl
	this.cache = cache

	exists, err := cache.Exists(key)
	if err != nil {
		return this.err(err.Error())
	}

	if !exists {
		return cache.Put(key, maxAttempts, ttl)
	}

	return nil
}

func (this rLimiter) Hit() error {
	exists, err := this.cache.Decrement(this.key, 1)
	if err != nil {
		return this.err(err.Error())
	}

	if !exists {
		return this.notExistsErr()
	}
	return nil
}

func (this rLimiter) Lock() error {
	exists, err := this.cache.Set(this.key, 0)
	if err != nil {
		return this.err(err.Error())
	}

	if !exists {
		return this.notExistsErr()
	}

	return nil
}

func (this rLimiter) Reset() error {
	err := this.cache.Put(this.key, this.max, this.ttl)
	if err != nil {
		return this.err(err.Error())
	}

	return nil
}

func (this rLimiter) Clear() error {
	err := this.cache.Forget(this.key)
	if err != nil {
		return this.err(err.Error())
	}

	return nil
}

func (this rLimiter) MustLock() (bool, error) {
	caster, err := this.cache.Cast(this.key)
	if err != nil {
		return true, this.err(err.Error())
	}

	if caster.IsNil() {
		return false, nil
	}

	v, err := caster.Int()
	if err != nil {
		err = this.err(err.Error())
	}
	return v <= 0, err
}

func (this rLimiter) TotalAttempts() (uint32, error) {
	caster, err := this.cache.Cast(this.key)
	if err != nil {
		return this.max, this.err(err.Error())
	}

	if caster.IsNil() {
		return this.max, nil
	}

	v, err := caster.Int()
	if err != nil {
		return this.max, this.err(err.Error())
	}

	if v > int(this.max) {
		v = int(this.max)
	}

	return this.max - uint32(v), nil
}

func (this rLimiter) RetriesLeft() (uint32, error) {
	caster, err := this.cache.Cast(this.key)
	if err != nil {
		return 0, this.err(err.Error())
	}

	if caster.IsNil() {
		return 0, nil
	}

	v, err := caster.Int()
	if err != nil {
		err = this.err(err.Error())
	}
	if v < 0 {
		v = 0
	}
	return uint32(v), err
}

func (this rLimiter) AvailableIn() (time.Duration, error) {
	if v, err := this.cache.TTL(this.key); err != nil {
		return 0, this.err(err.Error())
	} else {
		return v, nil
	}
}
