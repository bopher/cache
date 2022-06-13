package cache

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"math"
	"os"
	"path"
	"time"

	"github.com/bopher/caster"
	"github.com/bopher/utils"
)

type fCache struct {
	prefix string
	dir    string
}

func (rc fCache) err(pattern string, params ...any) error {
	return utils.TaggedError([]string{"FileCache"}, pattern, params...)
}

func (rc *fCache) init(prefix string, dir string) {
	rc.prefix = prefix
	rc.dir = dir
}

func (rc fCache) hashPath(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(utils.ConcatStr("-", rc.prefix, key)))
	fileName := hex.EncodeToString(hasher.Sum(nil))
	fileName = path.Join(rc.dir, fileName)
	return fileName
}

func (rc fCache) delete(key string) error {
	if err := os.Remove(rc.hashPath(key)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return rc.err(err.Error())
	}
	return nil
}

func (rc fCache) read(key string) (*record, error) {
	bytes, err := ioutil.ReadFile(rc.hashPath(key))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}

	if err != nil {
		return nil, rc.err(err.Error())
	}

	rec := record{}
	if err := rec.Deserialize(string(bytes)); err != nil {
		return nil, rc.err(err.Error())
	}

	if rec.IsExpired() {
		err := rc.delete(key)
		if err != nil {
			err = rc.err(err.Error())
		}
		return nil, err
	}

	return &rec, nil
}

func (rc fCache) write(key string, record record) error {
	err := utils.CreateDirectory(rc.dir)
	if err != nil {
		return rc.err(err.Error())
	}

	encoded, err := record.Serialize()
	if err != nil {
		return rc.err(err.Error())
	}

	err = ioutil.WriteFile(rc.hashPath(key), []byte(encoded), 0644)
	if err != nil {
		return rc.err(err.Error())
	}

	return nil
}

func (rc fCache) Put(key string, value any, ttl time.Duration) error {
	rec := record{
		TTL:  time.Now().UTC().Add(ttl),
		Data: value,
	}
	return rc.write(key, rec)
}

func (rc fCache) PutForever(key string, value any) error {
	rec := record{
		TTL:  time.Unix(math.MaxInt64, 0),
		Data: value,
	}
	return rc.write(key, rec)
}

func (rc fCache) Set(key string, value any) (bool, error) {
	rec, err := rc.read(key)
	if err != nil || rec == nil {
		return false, err
	}

	rec.Data = value
	return true, rc.write(key, *rec)
}

func (rc fCache) Get(key string) (any, error) {
	rec, err := rc.read(key)
	if err != nil || rec == nil {
		return nil, err
	}

	return rec.Data, nil
}

func (rc fCache) Exists(key string) (bool, error) {
	rec, err := rc.read(key)
	return rec == nil, err
}

func (rc fCache) Forget(key string) error {
	return rc.delete(key)
}

func (rc fCache) Pull(key string) (any, error) {
	if v, err := rc.Get(key); err != nil {
		return nil, err
	} else {
		return v, rc.delete(key)
	}
}

func (rc fCache) TTL(key string) (time.Duration, error) {
	rec, err := rc.read(key)
	if err != nil || rec == nil {
		return -1, err
	}

	return rec.TTL.UTC().Sub(time.Now().UTC()), nil
}

func (rc fCache) Cast(key string) (caster.Caster, error) {
	v, err := rc.Get(key)
	return caster.NewCaster(v), err
}

func (rc fCache) IncrementFloat(key string, value float64) (bool, error) {
	if c, err := rc.Cast(key); err != nil {
		return false, err
	} else {
		if v, err := c.Float64(); err != nil {
			return false, rc.err(err.Error())
		} else {
			return rc.Set(key, v+value)
		}
	}
}

func (rc fCache) Increment(key string, value int64) (bool, error) {
	if c, err := rc.Cast(key); err != nil {
		return false, err
	} else {
		if v, err := c.Int64(); err != nil {
			return false, rc.err(err.Error())
		} else {
			return rc.Set(key, v+value)
		}
	}
}

func (rc fCache) DecrementFloat(key string, value float64) (bool, error) {
	if c, err := rc.Cast(key); err != nil {
		return false, err
	} else {
		if v, err := c.Float64(); err != nil {
			return false, rc.err(err.Error())
		} else {
			return rc.Set(key, v-value)
		}
	}
}

func (rc fCache) Decrement(key string, value int64) (bool, error) {
	if c, err := rc.Cast(key); err != nil {
		return false, err
	} else {
		if v, err := c.Int64(); err != nil {
			return false, rc.err(err.Error())
		} else {
			return rc.Set(key, v-value)
		}
	}
}
