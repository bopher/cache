package cache

import (
	"context"
	"errors"
	"time"

	"github.com/bopher/caster"
	"github.com/bopher/utils"
	"github.com/go-redis/redis/v8"
)

type rCache struct {
	prefix string
	client *redis.Client
}

func (this rCache) err(pattern string, params ...interface{}) error {
	return utils.TaggedError([]string{"RedisCache"}, pattern, params...)
}

func (this *rCache) init(prefix string, opt redis.Options) {
	this.prefix = prefix
	this.client = redis.NewClient(&opt)
}

func (this rCache) perfixer(key string) string {
	return utils.ConcatStr("-", this.prefix, key)
}

func (this rCache) Put(key string, value interface{}, ttl time.Duration) error {
	if err := this.client.SetEX(
		context.TODO(),
		this.perfixer(key),
		value,
		ttl,
	).Err(); err != nil {
		return this.err(err.Error())
	}
	return nil
}

func (this rCache) PutForever(key string, value interface{}) error {
	if err := this.client.Set(
		context.TODO(),
		this.perfixer(key),
		value,
		0,
	).Err(); err != nil {
		return this.err(err.Error())
	}
	return nil
}

func (this rCache) Set(key string, value interface{}) (bool, error) {
	exists, err := this.Exists(key)
	if err != nil || !exists {
		return false, err
	}

	err = this.client.Set(
		context.TODO(),
		this.perfixer(key),
		value,
		redis.KeepTTL,
	).Err()

	if err != nil {
		err = this.err(err.Error())
	}

	return true, err
}

func (this rCache) Get(key string) (interface{}, error) {
	v, err := this.client.Get(
		context.TODO(),
		this.perfixer(key),
	).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		err = this.err(err.Error())
	}

	return v, err
}

func (this rCache) Exists(key string) (bool, error) {
	if exists, err := this.client.Exists(
		context.TODO(),
		this.perfixer(key),
	).Result(); err != nil {
		return false, this.err(err.Error())
	} else {
		return exists > 0, nil
	}
}

func (this rCache) Forget(key string) error {
	if err := this.client.Del(
		context.TODO(),
		this.perfixer(key),
	).Err(); err != nil && !errors.Is(err, redis.Nil) {
		return this.err(err.Error())
	}
	return nil
}

func (this rCache) Pull(key string) (interface{}, error) {
	if v, err := this.Get(key); err != nil {
		return nil, err
	} else {
		return v, this.Forget(key)
	}
}

func (this rCache) TTL(key string) (time.Duration, error) {
	if ttl, err := this.client.TTL(
		context.TODO(),
		this.perfixer(key),
	).Result(); err != nil {
		return 0, this.err(err.Error())
	} else {
		return ttl, nil
	}
}

func (this rCache) Cast(key string) (caster.Caster, error) {
	v, err := this.Get(key)
	return caster.NewCaster(v), err
}

func (this rCache) IncrementFloat(key string, value float64) (bool, error) {
	exists, err := this.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = this.client.IncrByFloat(
		context.TODO(),
		this.perfixer(key),
		value,
	).Err()
	if err != nil {
		err = this.err(err.Error())
	}
	return true, err
}

func (this rCache) Increment(key string, value int64) (bool, error) {
	exists, err := this.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = this.client.IncrBy(
		context.TODO(),
		this.perfixer(key),
		value,
	).Err()
	if err != nil {
		err = this.err(err.Error())
	}
	return true, err
}

func (this rCache) DecrementFloat(key string, value float64) (bool, error) {
	exists, err := this.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = this.client.IncrByFloat(
		context.TODO(),
		this.perfixer(key),
		-value,
	).Err()
	if err != nil {
		err = this.err(err.Error())
	}
	return true, err
}

func (this rCache) Decrement(key string, value int64) (bool, error) {
	exists, err := this.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = this.client.DecrBy(
		context.TODO(),
		this.perfixer(key),
		value,
	).Err()
	if err != nil {
		err = this.err(err.Error())
	}
	return true, err
}
