package cache

import (
	"context"
	"time"

	"github.com/bopher/utils"
	"github.com/go-redis/redis/v8"
)

type rCache struct {
	prefix string
	client *redis.Client
}

func (this *rCache) init(prefix string, opt redis.Options) {
	this.prefix = prefix
	this.client = redis.NewClient(&opt)
}

func (this rCache) err(pattern string, params ...interface{}) error {
	return utils.TaggedError([]string{"RedisCache"}, pattern, params...)
}

func (this rCache) px(key string) string {
	return utils.ConcatStr("-", this.prefix, key)
}

func (this rCache) Put(key string, value interface{}, ttl time.Duration) error {
	if err := this.client.SetEX(
		context.TODO(),
		this.px(key),
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
		this.px(key),
		value,
		0,
	).Err(); err != nil {
		return this.err(err.Error())
	}
	return nil
}

func (this rCache) Set(key string, value interface{}) error {
	if err := this.client.Set(
		context.TODO(),
		this.px(key),
		value,
		redis.KeepTTL,
	).Err(); err != nil {
		return this.err(err.Error())
	}
	return nil
}

func (this rCache) Exists(key string) (bool, error) {
	if exists, err := this.client.Exists(
		context.TODO(),
		this.px(key),
	).Result(); err != nil {
		return false, this.err(err.Error())
	} else {
		return exists > 0, nil
	}
}

func (this rCache) Forget(key string) error {
	if err := this.client.Del(
		context.TODO(),
		this.px(key),
	).Err(); err != nil {
		return this.err(err.Error())
	}
	return nil
}

func (this rCache) Get(key string) (interface{}, error) {
	if v, err := this.client.Get(
		context.TODO(),
		this.px(key),
	).Result(); err != nil {
		return nil, this.err(err.Error())
	} else {
		return v, nil
	}
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
		this.px(key),
	).Result(); err != nil {
		return 0, this.err(err.Error())
	} else {
		return ttl, nil
	}
}

func (this rCache) BoolE(key string) (bool, error) {
	if vs, err := this.Get(key); err != nil {
		return false, err
	} else {
		if v, err := utils.CastBoolE(vs); err != nil {
			return false, this.err(err.Error())
		} else {
			return v, nil
		}
	}
}

func (this rCache) Bool(key string, fallback bool) bool {
	if v, err := this.BoolE(key); err == nil {
		return v
	}
	return fallback
}

func (this rCache) IntE(key string) (int, error) {
	if vs, err := this.Get(key); err != nil {
		return 0, err
	} else {
		if v, err := utils.CastIntE(vs); err != nil {
			return 0, this.err(err.Error())
		} else {
			return v, nil
		}
	}
}

func (this rCache) Int(key string, fallback int) int {
	if v, err := this.IntE(key); err == nil {
		return v
	}
	return fallback
}

func (this rCache) Int8E(key string) (int8, error) {
	if vs, err := this.Get(key); err != nil {
		return 0, err
	} else {
		if v, err := utils.CastInt8E(vs); err != nil {
			return 0, this.err(err.Error())
		} else {
			return v, nil
		}
	}
}

func (this rCache) Int8(key string, fallback int8) int8 {
	if v, err := this.Int8E(key); err == nil {
		return v
	}
	return fallback
}

func (this rCache) Int16E(key string) (int16, error) {
	if vs, err := this.Get(key); err != nil {
		return 0, err
	} else {
		if v, err := utils.CastInt16E(vs); err != nil {
			return 0, this.err(err.Error())
		} else {
			return v, nil
		}
	}
}

func (this rCache) Int16(key string, fallback int16) int16 {
	if v, err := this.Int16E(key); err == nil {
		return v
	}
	return fallback
}

func (this rCache) Int32E(key string) (int32, error) {
	if vs, err := this.Get(key); err != nil {
		return 0, err
	} else {
		if v, err := utils.CastInt32E(vs); err != nil {
			return 0, this.err(err.Error())
		} else {
			return v, nil
		}
	}
}

func (this rCache) Int32(key string, fallback int32) int32 {
	if v, err := this.Int32E(key); err == nil {
		return v
	}
	return fallback
}

func (this rCache) Int64E(key string) (int64, error) {
	if vs, err := this.Get(key); err != nil {
		return 0, err
	} else {
		if v, err := utils.CastInt64E(vs); err != nil {
			return 0, this.err(err.Error())
		} else {
			return v, nil
		}
	}
}

func (this rCache) Int64(key string, fallback int64) int64 {
	if v, err := this.Int64E(key); err == nil {
		return v
	}
	return fallback
}

func (this rCache) UIntE(key string) (uint, error) {
	if vs, err := this.Get(key); err != nil {
		return 0, err
	} else {
		if v, err := utils.CastUIntE(vs); err != nil {
			return 0, this.err(err.Error())
		} else {
			return v, nil
		}
	}
}

func (this rCache) UInt(key string, fallback uint) uint {
	if v, err := this.UIntE(key); err == nil {
		return v
	}
	return fallback
}

func (this rCache) UInt8E(key string) (uint8, error) {
	if vs, err := this.Get(key); err != nil {
		return 0, err
	} else {
		if v, err := utils.CastUInt8E(vs); err != nil {
			return 0, this.err(err.Error())
		} else {
			return v, nil
		}
	}
}

func (this rCache) UInt8(key string, fallback uint8) uint8 {
	if v, err := this.UInt8E(key); err == nil {
		return v
	}
	return fallback
}

func (this rCache) UInt16E(key string) (uint16, error) {
	if vs, err := this.Get(key); err != nil {
		return 0, err
	} else {
		if v, err := utils.CastUInt16E(vs); err != nil {
			return 0, this.err(err.Error())
		} else {
			return v, nil
		}
	}
}

func (this rCache) UInt16(key string, fallback uint16) uint16 {
	if v, err := this.UInt16E(key); err == nil {
		return v
	}
	return fallback
}

func (this rCache) UInt32E(key string) (uint32, error) {
	if vs, err := this.Get(key); err != nil {
		return 0, err
	} else {
		if v, err := utils.CastUInt32E(vs); err != nil {
			return 0, this.err(err.Error())
		} else {
			return v, nil
		}
	}
}

func (this rCache) UInt32(key string, fallback uint32) uint32 {
	if v, err := this.UInt32E(key); err == nil {
		return v
	}
	return fallback
}

func (this rCache) UInt64E(key string) (uint64, error) {
	if vs, err := this.Get(key); err != nil {
		return 0, err
	} else {
		if v, err := utils.CastUInt64E(vs); err != nil {
			return 0, this.err(err.Error())
		} else {
			return v, nil
		}
	}
}

func (this rCache) UInt64(key string, fallback uint64) uint64 {
	if v, err := this.UInt64E(key); err == nil {
		return v
	}
	return fallback
}

func (this rCache) Float64E(key string) (float64, error) {
	if vs, err := this.Get(key); err != nil {
		return 0, err
	} else {
		if v, err := utils.CastFloat64E(vs); err != nil {
			return 0, this.err(err.Error())
		} else {
			return v, nil
		}
	}
}

func (this rCache) Float64(key string, fallback float64) float64 {
	if v, err := this.Float64E(key); err == nil {
		return v
	}
	return fallback
}

func (this rCache) StringE(key string) (string, error) {
	if v, err := this.client.Get(
		context.TODO(),
		this.px(key),
	).Result(); err != nil {
		return "", this.err(err.Error())
	} else {
		return v, nil
	}
}

func (this rCache) String(key string, fallback string) string {
	if v, err := this.StringE(key); err == nil {
		return v
	}
	return fallback
}

func (this rCache) IncrementBy(key string, value interface{}) error {
	var err error = nil
	switch value.(type) {
	case float32, float64:
		v, err := utils.CastFloat64E(value)
		if err == nil {
			err = this.client.IncrByFloat(
				context.TODO(),
				this.px(key),
				v,
			).Err()
		}
	default:
		v, err := utils.CastInt64E(value)
		if err == nil {
			err = this.client.IncrBy(
				context.TODO(),
				this.px(key),
				v,
			).Err()
		}
	}
	if err != nil {
		err = this.err(err.Error())
	}
	return err
}

func (this rCache) Increment(key string) error {
	return this.IncrementBy(key, 1)
}

func (this rCache) DecrementBy(key string, value interface{}) error {
	var err error = nil
	switch value.(type) {
	case float32, float64:
		v, err := utils.CastFloat64E(value)
		if err == nil {
			err = this.client.IncrByFloat(
				context.TODO(),
				this.px(key),
				-v,
			).Err()
		}
	default:
		v, err := utils.CastInt64E(value)
		if err == nil {
			err = this.client.DecrBy(
				context.TODO(),
				this.px(key),
				v,
			).Err()
		}
	}
	if err != nil {
		err = this.err(err.Error())
	}
	return err
}

func (this rCache) Decrement(key string) error {
	return this.DecrementBy(key, 1)
}
