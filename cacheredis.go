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

func (rc rCache) err(pattern string, params ...any) error {
	return utils.TaggedError([]string{"RedisCache"}, pattern, params...)
}

func (rc *rCache) init(prefix string, opt redis.Options) {
	rc.prefix = prefix
	rc.client = redis.NewClient(&opt)
}

func (rc rCache) perfixer(key string) string {
	return utils.ConcatStr("-", rc.prefix, key)
}

func (rc rCache) Put(key string, value any, ttl time.Duration) error {
	if err := rc.client.SetEX(
		context.TODO(),
		rc.perfixer(key),
		value,
		ttl,
	).Err(); err != nil {
		return rc.err(err.Error())
	}
	return nil
}

func (rc rCache) PutForever(key string, value any) error {
	if err := rc.client.Set(
		context.TODO(),
		rc.perfixer(key),
		value,
		0,
	).Err(); err != nil {
		return rc.err(err.Error())
	}
	return nil
}

func (rc rCache) Set(key string, value any) (bool, error) {
	exists, err := rc.Exists(key)
	if err != nil || !exists {
		return false, err
	}

	err = rc.client.Set(
		context.TODO(),
		rc.perfixer(key),
		value,
		redis.KeepTTL,
	).Err()

	if err != nil {
		err = rc.err(err.Error())
	}

	return true, err
}

func (rc rCache) Get(key string) (any, error) {
	v, err := rc.client.Get(
		context.TODO(),
		rc.perfixer(key),
	).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		err = rc.err(err.Error())
	}

	return v, err
}

func (rc rCache) Exists(key string) (bool, error) {
	if exists, err := rc.client.Exists(
		context.TODO(),
		rc.perfixer(key),
	).Result(); err != nil {
		return false, rc.err(err.Error())
	} else {
		return exists > 0, nil
	}
}

func (rc rCache) Forget(key string) error {
	if err := rc.client.Del(
		context.TODO(),
		rc.perfixer(key),
	).Err(); err != nil && !errors.Is(err, redis.Nil) {
		return rc.err(err.Error())
	}
	return nil
}

func (rc rCache) Pull(key string) (any, error) {
	if v, err := rc.Get(key); err != nil {
		return nil, err
	} else {
		return v, rc.Forget(key)
	}
}

func (rc rCache) TTL(key string) (time.Duration, error) {
	if ttl, err := rc.client.TTL(
		context.TODO(),
		rc.perfixer(key),
	).Result(); err != nil {
		return 0, rc.err(err.Error())
	} else {
		return ttl, nil
	}
}

func (rc rCache) Cast(key string) (caster.Caster, error) {
	v, err := rc.Get(key)
	return caster.NewCaster(v), err
}

func (rc rCache) IncrementFloat(key string, value float64) (bool, error) {
	exists, err := rc.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = rc.client.IncrByFloat(
		context.TODO(),
		rc.perfixer(key),
		value,
	).Err()
	if err != nil {
		err = rc.err(err.Error())
	}
	return true, err
}

func (rc rCache) Increment(key string, value int64) (bool, error) {
	exists, err := rc.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = rc.client.IncrBy(
		context.TODO(),
		rc.perfixer(key),
		value,
	).Err()
	if err != nil {
		err = rc.err(err.Error())
	}
	return true, err
}

func (rc rCache) DecrementFloat(key string, value float64) (bool, error) {
	exists, err := rc.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = rc.client.IncrByFloat(
		context.TODO(),
		rc.perfixer(key),
		-value,
	).Err()
	if err != nil {
		err = rc.err(err.Error())
	}
	return true, err
}

func (rc rCache) Decrement(key string, value int64) (bool, error) {
	exists, err := rc.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = rc.client.DecrBy(
		context.TODO(),
		rc.perfixer(key),
		value,
	).Err()
	if err != nil {
		err = rc.err(err.Error())
	}
	return true, err
}
