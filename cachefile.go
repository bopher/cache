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

func (this fCache) err(pattern string, params ...interface{}) error {
	return utils.TaggedError([]string{"FileCache"}, pattern, params...)
}

func (this *fCache) init(prefix string, dir string) {
	this.prefix = prefix
	this.dir = dir
}

func (this fCache) hashPath(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(utils.ConcatStr("-", this.prefix, key)))
	fileName := hex.EncodeToString(hasher.Sum(nil))
	fileName = path.Join(this.dir, fileName)
	return fileName
}

func (this fCache) delete(key string) error {
	if err := os.Remove(this.hashPath(key)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return this.err(err.Error())
	}
	return nil
}

func (this fCache) read(key string) (*record, error) {
	bytes, err := ioutil.ReadFile(this.hashPath(key))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}

	if err != nil {
		return nil, this.err(err.Error())
	}

	rec := record{}
	if err := rec.Deserialize(string(bytes)); err != nil {
		return nil, this.err(err.Error())
	}

	if rec.IsExpired() {
		err := this.delete(key)
		if err != nil {
			err = this.err(err.Error())
		}
		return nil, err
	}

	return &rec, nil
}

func (this fCache) write(key string, record record) error {
	err := utils.CreateDirectory(this.dir)
	if err != nil {
		return this.err(err.Error())
	}

	encoded, err := record.Serialize()
	if err != nil {
		return this.err(err.Error())
	}

	err = ioutil.WriteFile(this.hashPath(key), []byte(encoded), 0644)
	if err != nil {
		return this.err(err.Error())
	}

	return nil
}

func (this fCache) Put(key string, value interface{}, ttl time.Duration) error {
	rec := record{
		TTL:  time.Now().UTC().Add(ttl),
		Data: value,
	}
	return this.write(key, rec)
}

func (this fCache) PutForever(key string, value interface{}) error {
	rec := record{
		TTL:  time.Unix(math.MaxInt64, 0),
		Data: value,
	}
	return this.write(key, rec)
}

func (this fCache) Set(key string, value interface{}) (bool, error) {
	rec, err := this.read(key)
	if err != nil || rec == nil {
		return false, err
	}

	rec.Data = value
	return true, this.write(key, *rec)
}

func (this fCache) Get(key string) (interface{}, error) {
	rec, err := this.read(key)
	if err != nil || rec == nil {
		return nil, err
	}

	return rec.Data, nil
}

func (this fCache) Exists(key string) (bool, error) {
	rec, err := this.read(key)
	return rec == nil, err
}

func (this fCache) Forget(key string) error {
	return this.delete(key)
}

func (this fCache) Pull(key string) (interface{}, error) {
	if v, err := this.Get(key); err != nil {
		return nil, err
	} else {
		return v, this.delete(key)
	}
}

func (this fCache) TTL(key string) (time.Duration, error) {
	rec, err := this.read(key)
	if err != nil || rec == nil {
		return -1, err
	}

	return rec.TTL.UTC().Sub(time.Now().UTC()), nil
}

func (this fCache) Cast(key string) (caster.Caster, error) {
	v, err := this.Get(key)
	return caster.NewCaster(v), err
}

func (this fCache) IncrementFloat(key string, value float64) (bool, error) {
	if c, err := this.Cast(key); err != nil {
		return false, err
	} else {
		if v, err := c.Float64(); err != nil {
			return false, this.err(err.Error())
		} else {
			return this.Set(key, v+value)
		}
	}
}

func (this fCache) Increment(key string, value int64) (bool, error) {
	if c, err := this.Cast(key); err != nil {
		return false, err
	} else {
		if v, err := c.Int64(); err != nil {
			return false, this.err(err.Error())
		} else {
			return this.Set(key, v+value)
		}
	}
}

func (this fCache) DecrementFloat(key string, value float64) (bool, error) {
	if c, err := this.Cast(key); err != nil {
		return false, err
	} else {
		if v, err := c.Float64(); err != nil {
			return false, this.err(err.Error())
		} else {
			return this.Set(key, v-value)
		}
	}
}

func (this fCache) Decrement(key string, value int64) (bool, error) {
	if c, err := this.Cast(key); err != nil {
		return false, err
	} else {
		if v, err := c.Int64(); err != nil {
			return false, this.err(err.Error())
		} else {
			return this.Set(key, v-value)
		}
	}
}
