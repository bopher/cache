package cache

import (
	"time"

	"github.com/bopher/utils"
)

type rLimiter struct {
	Key   string
	Max   uint32
	Cache Cache
}

func (this rLimiter) err(pattern string, params ...interface{}) error {
	return utils.TaggedError([]string{"RateLimiter", this.Key}, pattern, params...)
}

func (this *rLimiter) init(key string, maxAttempts uint32, ttl time.Duration, cache Cache) error {
	this.Key = key
	this.Max = maxAttempts
	this.Cache = cache

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
	if i, err := this.Cache.IntE(this.Key); err != nil {
		return this.err(err.Error())
	} else {
		if i > 0 {
			if err := this.Cache.Decrement(this.Key); err != nil {
				return this.err(err.Error())
			}
		}
	}
	return nil
}

func (this rLimiter) Lock() error {
	if exists, err := this.Cache.Exists(this.Key); err != nil {
		return this.err(err.Error())
	} else if exists {
		if err := this.Cache.Set(this.Key, 0); err != nil {
			return this.err(err.Error())
		}
	}
	return nil
}

func (this rLimiter) Reset() error {
	if err := this.Cache.Forget(this.Key); err != nil {
		return this.err(err.Error())
	}
	return nil
}

func (this rLimiter) MustLock() (bool, error) {
	if v, err := this.Cache.IntE(this.Key); err != nil {
		return false, this.err(err.Error())
	} else {
		return v <= 0, nil
	}
}

func (this rLimiter) TotalAttempts() (uint32, error) {
	i, err := this.Cache.IntE(this.Key)
	if err != nil {
		return 0, this.err(err.Error())
	}
	if i < 0 {
		i = 0
	}
	if uint32(i) > this.Max {
		i = int(this.Max)
	}

	return this.Max - uint32(i), nil
}

func (this rLimiter) RetriesLeft() (uint32, error) {
	if v, err := this.Cache.UInt32E(this.Key); err != nil {
		return 0, this.err(err.Error())
	} else {
		return v, nil
	}
}

func (this rLimiter) AvailableIn() (time.Duration, error) {
	if v, err := this.Cache.TTL(this.Key); err != nil {
		return 0, this.err(err.Error())
	} else {
		return v, nil
	}
}
