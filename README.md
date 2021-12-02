# Cache

Cache manager with default file and redis driver. this package contains rate limiter and verification code driver.

## Create New Cache Driver

Cache library contains two different driver by default.

**NOTE:** You can extend your driver by implementing `Cache` interface.

### Create File Based Driver

for creating file based driver you must pass file name prefix and cache directory to constructor function.

```go
import "github.com/bopher/cache"
if fCache := cache.NewFileCache("myApp", "./caches"); fCache != nil {
  // Cache driver created
} else {
  panic("failed to build cache driver")
}
```

### Create Redis Based Driver

for creating redis based driver you must pass prefix, and redis options to constructor function.

```go
import "github.com/bopher/cache"
if rCache := cache.NewRedisCache("myApp", redis.Options{
  Addr: "localhost:6379",
}); rCache != nil {
  // Cache driver created
} else {
  panic("failed to build cache driver")
}
```

## Usage

Cache interface contains following methods:

### Put

Put a new value to cache.

```go
// Signature:
Put(key string, value interface{}, ttl time.Duration) error

// Example:
err := rCache.Put("total-debt", 410203, 100 * time.Hour)
```

### PutForever

Put a new value to cache with infinite ttl.

```go
// Signature:
PutForever(key string, value interface{}) error

// Example:
err := rCache.PutForever("total-users", 5000000)
```

### Set

Change value of cache item (keep ttl).

```go
// Signature:
Set(key string, value interface{}) error

// Example:
err := rCache.Set("total-users", 5005)
```

### Get

Get item from cache. Get function return cache item as `interface{}`. if you need get cache item with type use helper get functions described later.

```go
// Signature:
Get(key string) (interface{}, error)

// Example:
v, err := rCache.Get("total-users")
```

### Exists

Check if item exists in cache

```go
// Signature:
Exists(key string) (bool, error)

// Example:
exists, err := rCache.Exists("total-users");
```

### Forget

Delete item from cache.

```go
// Signature:
Forget(key string) error

// Example:
err := rCache.Forget("total-users")
```

### Pull

Get item from cache and remove it

```go
// Signature:
Pull(key string) (interface{}, error)

// Example:
v, err := rCache.Pull("total-users")
```

### TTL

Get cache item ttl.

```go
// Signature:
TTL(key string) (time.Duration, error)

// Example:
ttl, err := rCache.TTL("total-users")
```

### IncrementBy

Increment numeric item in cache by number.

```go
// Signature:
IncrementBy(key string, value interface{}) error

// Example:
err := rCache.IncrementBy("total-users", 10)
```

### Increment

Increment numeric item in cache.

```go
// Signature:
Increment(key string) error

// Example:
err := rCache.Increment("total-users")
```

### DecrementBy

Decrement numeric item in cache by number.

```go
// Signature:
DecrementBy(key string, value interface{}) error

// Example:
err := rCache.DecrementBy("total-users", 10)
```

### Decrement

Decrement numeric item in cache.

```go
// Signature:
Decrement(key string) error

// Example:
err := rCache.Decrement("total-users")
```

### Getters

Getters function allow you to cast cache item directly as type. Getters item return error when item not exists or type cast failed!

```go
// BoolE parse item as boolean or return error on fail
BoolE(key string) (bool, error)

// IntE parse item as int or return error on fail
IntE(key string) (int, error)

// Int8E parse item as int8 or return error on fail
Int8E(key string) (int8, error)

// Int16E parse item as int16 or return error on fail
Int16E(key string) (int16, error)

// Int32E parse item as int32 or return error on fail
Int32E(key string) (int32, error)

// Int64E parse item as int64 or return error on fail
Int64E(key string) (int64, error)

// UIntE parse item as uint or return error on fail
UIntE(key string) (uint, error)

// UInt8E parse item as uint8 or return error on fail
UInt8E(key string) (uint8, error)

// UInt16E parse item as uint16 or return error on fail
UInt16E(key string) (uint16, error)

// UInt32E parse item as uint32 or return error on fail
UInt32E(key string) (uint32, error)

// UInt64E parse item as uint64 or return error on fail
UInt64E(key string) (uint64, error)

// Float64E parse item as float64 or return error on fail
Float64E(key string) (float64, error)

// StringE parse item as string or return error on fail
StringE(key string) (string, error)
```

### Error Safe Getters

You can use safe getters to cast cache item and pass fallback value in case of item casting failed!

```go
// Bool parse item as boolean or return fallback
Bool(key string, fallback bool) bool

// Int parse item as int or return fallback
Int(key string, fallback int) int

// Int8 parse item as int8 or return fallback
Int8(key string, fallback int8) int8

// Int16 parse item as int16  or return fallback
Int16(key string, fallback int16) int16

// Int32 parse item as int32 or return fallback
Int32(key string, fallback int32) int32

// Int64 parse item as int64 or return fallback
Int64(key string, fallback int64) int64

// UInt parse item as uint or return fallback
UInt(key string, fallback uint) uint

// UInt8 parse item as uint8 or return fallback
UInt8(key string, fallback uint8) uint8

// UInt16 parse item as uint16 or return fallback
UInt16(key string, fallback uint16) uint16

// UInt32 parse item as uint32 or return fallback
UInt32(key string, fallback uint32) uint32

// UInt64 parse item as uint64 or return fallback
UInt64(key string, fallback uint64) uint64

// Float64 parse item as float64 or return fallback
Float64(key string, fallback float64) float64

// String parse item as string or return fallback
String(key string, fallback string) string
```

## Create New Rate Limiter Driver

**Note:** Rate limiter based on cache, For creating rate limiter driver you must pass a cache driver instance to constructor function.

```go
// Signature:
NewRateLimiter(key string, maxAttempts uint32, ttl time.Duration, cache Cache) (RateLimiter, error)

// Example: allow 3 attempts every 60 seconds
import "github.com/bopher/cache"
limiter, err := cache.NewRateLimiter("login-attempts", 3, 60 * time.Second, rCache)
```

## Usage

Rate limiter interface contains following methods:

### Hit

Decrease the allowed times.

```go
// Signature:
Hit() error

// Example:
err := limiter.Hit()
```

### Lock

Lock rate limiter.

```go
// Signature:
Lock() error

// Example:
err := limiter.Lock() // no more attempts left
```

### Reset

Reset rate limiter.

```go
// Signature:
Reset() error

// Example:
err := limiter.Reset()
```

### MustLock

Check if rate limiter must lock access.

```go
// Signature:
MustLock() (bool, error)

// Example:
if locked, _:= limiter.MustLock(), locked {
  // Block access
}
```

### TotalAttempts

Get user attempts count.

```go
// Signature:
TotalAttempts() (uint32, error)

// Example:
totalAtt, err := limiter.TotalAttempts() // 3
```

### RetriesLeft

Get user retries left.

```go
// Signature:
RetriesLeft() (uint32, error)

// Example:
leftRet, err := limiter.RetriesLeft() // 2
```

### AvailableIn

Get time until unlock.

```go
// Signature:
AvailableIn() (time.Duration, error)

// Example:
availableIn, err := limiter.AvailableIn()
```

## Create New Verification Code Driver

verification code used for managing verification code sent to user.

**Note:** Verification code based on cache, For creating verification code driver you must pass a cache driver instance to constructor function.

```go
// Signature:
NewVerificationCode(key string, ttl time.Duration, cache Cache) VerificationCode

// Example:
import "github.com/bopher/cache"
vCode := cache.NewVerificationCode("phone-verification", 5 * time.Minute, rCache)
```

## Usage

Verification code interface contains following methods:

### Set

Set code. You can set code directly or use generator methods.

```go
// Signature:
Set(value string) error

// Example:
err := vCode.Set("ABD531")
```

### Generate

Generate a random numeric code with 5 character length and set as code.

```go
// Signature:
Generate() (string, error)

// Example:
code, err := vCode.Generate()
```

### GenerateN

Generate a random numeric code with special character length and set as code.

```go
// Signature:
GenerateN(count uint) (string, error)

// Example:
code, err := vCode.GenerateN(6)
```

### Clear

Clear code from cache.

```go
// Signature:
Clear() error

// Example:
err := vCode.Clear()
```

### Get

Get code.

```go
// Signature:
Get() (string, error)

// Example:
code, err := vCode.Get()
```

### Exists

Exists check if code exists in cache and not empty.

```go
// Signature:
Exists() (bool, error)

// Example:
exists, err := vCode.Exists()
```
