package cache

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"time"

	"github.com/bopher/utils"
)

type record struct {
	TTL  time.Time
	Data interface{}
}

func (this record) Serialize() (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(this)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b.Bytes()), nil
}

func (this *record) Deserialize(data string) error {
	by, err := hex.DecodeString(data)
	if err != nil {
		return err
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(this)
	if err != nil {
		return err
	}
	return nil
}

func (this record) IsExpired() bool {
	return this.TTL.UTC().Before(time.Now().UTC())
}

type fCache struct {
	prefix string
	dir    string
}

func (this *fCache) init(prefix string, dir string) {
	this.prefix = prefix
	this.dir = dir
}

func (this fCache) err(pattern string, params ...interface{}) error {
	return utils.TaggedError([]string{"FileCache"}, pattern, params...)
}

func (this fCache) pathResolver(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(utils.ConcatStr("-", this.prefix, key)))
	fileName := hex.EncodeToString(hasher.Sum(nil))
	fileName = path.Join(this.dir, fileName)
	return fileName
}

func (this fCache) delete(key string) error {
	if err := os.Remove(this.pathResolver(key)); err != nil {
		return this.err(err.Error())
	}
	return nil
}

func (this fCache) read(key string) (*record, error) {
	bytes, err := ioutil.ReadFile(this.pathResolver(key))
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

	err = ioutil.WriteFile(this.pathResolver(key), []byte(encoded), 0644)
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

func (this fCache) Set(key string, value interface{}) error {
	rec, err := this.read(key)
	if err != nil {
		return err
	}

	rec.Data = value
	return this.write(key, *rec)
}

func (this fCache) Get(key string) (interface{}, error) {
	rec, err := this.read(key)
	if err != nil {
		return nil, err
	}

	return rec.Data, nil
}

func (this fCache) Exists(key string) (bool, error) {
	_, err := this.read(key)
	if err != nil {
		return false, err
	}

	return true, nil
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
	if err != nil {
		return 0, err
	}

	return time.Now().UTC().Sub(rec.TTL.UTC()), nil
}

func (this fCache) BoolE(key string) (bool, error) {
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

func (this fCache) Bool(key string, fallback bool) bool {
	if v, err := this.BoolE(key); err == nil {
		return v
	}
	return fallback
}

func (this fCache) IntE(key string) (int, error) {
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

func (this fCache) Int(key string, fallback int) int {
	if v, err := this.IntE(key); err == nil {
		return v
	}
	return fallback
}

func (this fCache) Int8E(key string) (int8, error) {
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

func (this fCache) Int8(key string, fallback int8) int8 {
	if v, err := this.Int8E(key); err == nil {
		return v
	}
	return fallback
}

func (this fCache) Int16E(key string) (int16, error) {
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

func (this fCache) Int16(key string, fallback int16) int16 {
	if v, err := this.Int16E(key); err == nil {
		return v
	}
	return fallback
}

func (this fCache) Int32E(key string) (int32, error) {
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

func (this fCache) Int32(key string, fallback int32) int32 {
	if v, err := this.Int32E(key); err == nil {
		return v
	}
	return fallback
}

func (this fCache) Int64E(key string) (int64, error) {
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

func (this fCache) Int64(key string, fallback int64) int64 {
	if v, err := this.Int64E(key); err == nil {
		return v
	}
	return fallback
}

func (this fCache) UIntE(key string) (uint, error) {
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

func (this fCache) UInt(key string, fallback uint) uint {
	if v, err := this.UIntE(key); err == nil {
		return v
	}
	return fallback
}

func (this fCache) UInt8E(key string) (uint8, error) {
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

func (this fCache) UInt8(key string, fallback uint8) uint8 {
	if v, err := this.UInt8E(key); err == nil {
		return v
	}
	return fallback
}

func (this fCache) UInt16E(key string) (uint16, error) {
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

func (this fCache) UInt16(key string, fallback uint16) uint16 {
	if v, err := this.UInt16E(key); err == nil {
		return v
	}
	return fallback
}

func (this fCache) UInt32E(key string) (uint32, error) {
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

func (this fCache) UInt32(key string, fallback uint32) uint32 {
	if v, err := this.UInt32E(key); err == nil {
		return v
	}
	return fallback
}

func (this fCache) UInt64E(key string) (uint64, error) {
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

func (this fCache) UInt64(key string, fallback uint64) uint64 {
	if v, err := this.UInt64E(key); err == nil {
		return v
	}
	return fallback
}

func (this fCache) Float64E(key string) (float64, error) {
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

func (this fCache) Float64(key string, fallback float64) float64 {
	if v, err := this.Float64E(key); err == nil {
		return v
	}
	return fallback
}

func (this fCache) StringE(key string) (string, error) {
	if vs, err := this.Get(key); err != nil {
		return "", err
	} else {
		if vs != nil {
			return fmt.Sprint(vs), nil
		}
	}
	return "", this.err("cant get %s as string!", key)
}

func (this fCache) String(key string, fallback string) string {
	if v, err := this.StringE(key); err == nil {
		return v
	}
	return fallback
}

func (this fCache) IncrementBy(key string, value interface{}) error {
	var err error = nil
	switch value.(type) {
	case float32, float64:
		org, err := this.Float64E(key)
		if err == nil {
			v, err := utils.CastFloat64E(value)
			if err == nil {
				return this.Set(key, org+v)
			} else {
				err = this.err(err.Error())
			}
		}
	default:
		org, err := this.Int64E(key)
		if err == nil {
			v, err := utils.CastInt64E(value)
			if err == nil {
				return this.Set(key, org+v)
			} else {
				err = this.err(err.Error())
			}
		}
	}
	return err
}

func (this fCache) Increment(key string) error {
	return this.IncrementBy(key, 1)
}

func (this fCache) DecrementBy(key string, value interface{}) error {
	var err error = nil
	switch value.(type) {
	case float32, float64:
		org, err := this.Float64E(key)
		if err == nil {
			v, err := utils.CastFloat64E(value)
			if err == nil {
				return this.Set(key, org-v)
			} else {
				err = this.err(err.Error())
			}
		}
	default:
		org, err := this.Int64E(key)
		if err == nil {
			v, err := utils.CastInt64E(value)
			if err == nil {
				return this.Set(key, org-v)
			} else {
				err = this.err(err.Error())
			}
		}
	}
	return err
}

func (this fCache) Decrement(key string) error {
	return this.DecrementBy(key, 1)
}
