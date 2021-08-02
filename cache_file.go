package cache

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/bopher/utils"
)

type cacheRecord struct {
	TTL  time.Time
	Data interface{}
}

// Serialize data
func (r *cacheRecord) Serialize() (string, bool) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(*r)
	if err != nil {
		return "", false
	}
	return hex.EncodeToString(b.Bytes()), true
}

// Deserialize data
func (r *cacheRecord) Deserialize(data string) bool {
	by, err := hex.DecodeString(data)
	if err != nil {
		return false
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(r)
	if err != nil {
		return false
	}
	return true
}

// IsExpired check if record is expired
func (r *cacheRecord) IsExpired() bool {
	return r.TTL.UTC().Before(time.Now().UTC())
}

// ParseAsInt64 parse data as int64
func (r *cacheRecord) ParseAsInt64() (int64, bool) {
	switch v := r.Data.(type) {
	case int8:
		return int64(v), true
	case uint8:
		return int64(v), true
	case int16:
		return int64(v), true
	case uint16:
		return int64(v), true
	case int32:
		return int64(v), true
	case uint32:
		return int64(v), true
	case int:
		return int64(v), true
	case uint:
		return int64(v), true
	case int64:
		return int64(v), true
	case uint64:
		return int64(v), true
	case float32:
		return int64(v), true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}

// ParseAsUint64 parse data as uint64
func (r *cacheRecord) ParseAsUint64() (uint64, bool) {
	switch v := r.Data.(type) {
	case int8:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case int16:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case int32:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case int:
		return uint64(v), true
	case uint:
		return uint64(v), true
	case int64:
		return uint64(v), true
	case uint64:
		return uint64(v), true
	case float32:
		return uint64(v), true
	case float64:
		return uint64(v), true
	default:
		return 0, false
	}
}

// ParseAsFloat64 parse data as float64
func (r *cacheRecord) ParseAsFloat64() (float64, bool) {
	switch v := r.Data.(type) {
	case int8:
		return float64(v), true
	case uint8:
		return float64(v), true
	case int16:
		return float64(v), true
	case uint16:
		return float64(v), true
	case int32:
		return float64(v), true
	case uint32:
		return float64(v), true
	case int:
		return float64(v), true
	case uint:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return float64(v), true
	default:
		return 0, false
	}
}

type fileCache struct {
	prefix string
	dir    string
}

func (c *fileCache) init(prefix string, dir string) {
	c.prefix = prefix
	c.dir = dir
}

func (c *fileCache) pathResolver(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(c.prefix + "-" + key))
	fileName := hex.EncodeToString(hasher.Sum(nil))
	fileName = path.Join(c.dir, fileName)
	return fileName
}

func (c *fileCache) read(key string) (*cacheRecord, bool) {
	bytes, err := ioutil.ReadFile(c.pathResolver(key))
	if err != nil {
		return nil, false
	}
	rec := cacheRecord{}
	if !rec.Deserialize(string(bytes)) {
		return nil, false
	}

	if rec.IsExpired() {
		c.delete(key)
		return nil, false
	}

	return &rec, true
}

func (c *fileCache) write(key string, record cacheRecord) bool {
	utils.CreateDirectory(c.dir)
	encoded, ok := record.Serialize()
	if ok {
		err := ioutil.WriteFile(c.pathResolver(key), []byte(encoded), 0644)
		if err == nil {
			return true
		}
	}
	return false
}

func (c *fileCache) delete(key string) bool {
	return os.Remove(c.pathResolver(key)) == nil
}

// Put a new value to cache
func (c *fileCache) Put(key string, value interface{}, ttl time.Duration) bool {
	record := cacheRecord{
		TTL:  time.Now().UTC().Add(ttl),
		Data: value,
	}
	return c.write(key, record)
}

// PutForever put value with infinite ttl
func (c *fileCache) PutForever(key string, value interface{}) bool {
	record := cacheRecord{
		TTL:  time.Now().UTC().AddDate(100, 0, 0),
		Data: value,
	}
	return c.write(key, record)
}

// Set Change value of cache item
func (c *fileCache) Set(key string, value interface{}) bool {
	rec, exists := c.read(key)
	if exists {
		rec.Data = value
		return c.write(key, *rec)
	}
	return false
}

// Get item from cache
func (c *fileCache) Get(key string) interface{} {
	rec, exists := c.read(key)
	if exists {
		return rec.Data
	}
	return nil
}

// Pull item from cache and remove it
func (c *fileCache) Pull(key string) interface{} {
	defer c.delete(key)
	rec, exists := c.read(key)
	if exists {
		return rec.Data
	}
	return nil
}

// Check if item exists in cache
func (c *fileCache) Exists(key string) bool {
	_, exists := c.read(key)
	return exists
}

// Forget item from cache (delete item)
func (c *fileCache) Forget(key string) bool {
	return c.delete(key)
}

// TTL get cache item ttl
func (c *fileCache) TTL(key string) time.Duration {
	rec, exists := c.read(key)
	if exists {
		return time.Now().UTC().Sub(rec.TTL.UTC())
	}
	return 0
}

// Bool parse dependency as boolean
func (c *fileCache) Bool(key string, fallback bool) bool {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.Data.(bool); ok {
			return res
		}
	}
	return fallback
}

// Int parse dependency as int
func (c *fileCache) Int(key string, fallback int) int {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsInt64(); ok {
			return int(res)
		}
	}
	return fallback
}

// Int8 parse dependency as int8
func (c *fileCache) Int8(key string, fallback int8) int8 {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsInt64(); ok {
			return int8(res)
		}
	}
	return fallback
}

// Int16 parse dependency as int16
func (c *fileCache) Int16(key string, fallback int16) int16 {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsInt64(); ok {
			return int16(res)
		}
	}
	return fallback
}

// Int32 parse dependency as int32
func (c *fileCache) Int32(key string, fallback int32) int32 {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsInt64(); ok {
			return int32(res)
		}
	}
	return fallback
}

// Int64 parse dependency as int64
func (c *fileCache) Int64(key string, fallback int64) int64 {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsInt64(); ok {
			return res
		}
	}
	return fallback
}

// UInt parse dependency as uint
func (c *fileCache) UInt(key string, fallback uint) uint {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsUint64(); ok {
			return uint(res)
		}
	}
	return fallback
}

// UInt8 parse dependency as uint8
func (c *fileCache) UInt8(key string, fallback uint8) uint8 {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsUint64(); ok {
			return uint8(res)
		}
	}
	return fallback
}

// UInt16 parse dependency as uint16
func (c *fileCache) UInt16(key string, fallback uint16) uint16 {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsUint64(); ok {
			return uint16(res)
		}
	}
	return fallback
}

// UInt32 parse dependency as uint32
func (c *fileCache) UInt32(key string, fallback uint32) uint32 {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsUint64(); ok {
			return uint32(res)
		}
	}
	return fallback
}

// UInt64 parse dependency as uint64
func (c *fileCache) UInt64(key string, fallback uint64) uint64 {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsUint64(); ok {
			return uint64(res)
		}
	}
	return fallback
}

// Float32 parse dependency as float64
func (c *fileCache) Float32(key string, fallback float32) float32 {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsFloat64(); ok {
			return float32(res)
		}
	}
	return fallback
}

// Float64 parse dependency as float64
func (c *fileCache) Float64(key string, fallback float64) float64 {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsFloat64(); ok {
			return res
		}
	}
	return fallback
}

// String parse dependency as string
func (c *fileCache) String(key string, fallback string) string {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.Data.(string); ok {
			return res
		}
	}
	return fallback
}

// Bytes parse dependency as bytes array
func (c *fileCache) Bytes(key string, fallback []byte) []byte {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.Data.([]byte); ok {
			return res
		}
	}
	return fallback
}

// Increment numeric item in cache
func (c *fileCache) Increment(key string) bool {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsFloat64(); ok {
			return c.Set(key, res+1)
		}
	}
	return false
}

// IncrementBy numeric item in cache by number
func (c *fileCache) IncrementBy(key string, value interface{}) bool {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsFloat64(); ok {
			temp := cacheRecord{
				Data: value,
			}
			if val, ok := temp.ParseAsFloat64(); ok {
				return c.Set(key, res+val)
			}
		}
	}
	return false
}

// Decrement numeric item in cache
func (c *fileCache) Decrement(key string) bool {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsFloat64(); ok {
			return c.Set(key, res-1)
		}
	}
	return false
}

// DecrementBy numeric item in cache by number
func (c *fileCache) DecrementBy(key string, value interface{}) bool {
	rec, exists := c.read(key)
	if exists {
		if res, ok := rec.ParseAsFloat64(); ok {
			temp := cacheRecord{
				Data: value,
			}
			if val, ok := temp.ParseAsFloat64(); ok {
				return c.Set(key, res-val)
			}
		}
	}
	return false
}
