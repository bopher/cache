# Cache

Cache manager with default file and redis driver (rate limiter and verification code manager included).

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

A new value to cache.

```go
// Signature:
Put(key string, value interface{}, ttl time.Duration) error

// Example:
err := rCache.Put("total-debt", 410203, 100 * time.Hour)
```

### PutForever

Put value with infinite ttl.

```go
// Signature:
PutForever(key string, value interface{}) error

// Example:
err := rCache.PutForever("base-discount", 10)
```

### Set

Change value of cache item and return false if item not exists (this. methods keep cache ttl).

**Cation:** set value on non exists item will generate error. please check if item exists before set or use put method instead!

```go
// Signature:
Set(key string, value interface{}) (bool, error)

// Example:
ok, err := rCache.Set("base-discount", 15)
```

### Get

Get item from cache.

```go
// Signature:
Get(key string) (interface{}, error)

// Example:
v, err := rCache.Get("total-users")
```

### Exists

Check if item exists in cache.

```go
// Signature:
Exists(key string) (bool, error)

// Example:
exists, err := rCache.Exists("total-users");
```

### Forget

Delete Item from cache.

```go
// Signature:
Forget(key string) error

// Example:
err := rCache.Forget("total-users")
```

### Pull

Item from cache and then remove it.

```go
// Signature:
Pull(key string) (interface{}, error)

// Example:
v, err := rCache.Pull("total-users")
```

### TTL

Get cache item ttl. This method returns -1 if item not exists.

```go
// Signature:
TTL(key string) (time.Duration, error)

// Example:
ttl, err := rCache.TTL("total-users")
```

### Cast

Parse cache item as caster.

```go
// Signature:
Cast(key string) (caster.Caster, error)

// Example:
c, err := rCache.Cast("total-users")
v, err := c.Int32()
```

### IncrementFloat

Increment numeric item by float, return false if item not exists

```go
// Signature:
IncrementFloat(key string, value float64) (bool, error)

// Example:
err := rCache.IncrementFloat("some-float", 0.01)
```

### Increment

Increment numeric item by int, return false if item not exists

```go
// Signature:
Increment(key string, value int64) (bool, error)

// Example:
err := rCache.Increment("some-number", 10)
```

### DecrementFloat

Decrement numeric item by float, return false if item not exists

```go
// Signature:
DecrementFloat(key string, value float64) (bool, error)

// Example:
err := rCache.DecrementFloat("some-float", 0.29)
```

### Decrement

Decrement numeric item by int, return false if item not exists

```go
// Signature:
Decrement(key string, value int64) (bool, error)

// Example:
err := rCache.Decrement("total-try", 1)
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

### Usage

Rate limiter interface contains following methods:

#### Hit

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

#### Reset

Reset rate limiter (clear total attempts).

```go
// Signature:
Reset() error

// Example:
err := limiter.Reset()
```

#### Clear

Remove rate limiter record. call any method after clear with generate `"NotExists"` error!

```go
// Signature:
Clear() error

// Example:
err := limiter.Clear()
```

#### MustLock

Check if rate limiter must lock access.

```go
// Signature:
MustLock() (bool, error)

// Example:
if locked, _:= limiter.MustLock(), locked {
  // Block access
}
```

#### TotalAttempts

Get user attempts count.

```go
// Signature:
TotalAttempts() (uint32, error)

// Example:
totalAtt, err := limiter.TotalAttempts() // 3
```

#### RetriesLeft

Get user retries left.

```go
// Signature:
RetriesLeft() (uint32, error)

// Example:
leftRet, err := limiter.RetriesLeft() // 2
```

#### AvailableIn

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
NewVerificationCode(key string, ttl time.Duration, cache Cache) (VerificationCode, error)

// Example:
import "github.com/bopher/cache"
vCode, err := cache.NewVerificationCode("phone-verification", 5 * time.Minute, rCache)
```

### Usage

Verification code interface contains following methods:

#### Set

Set code. You can set code directly or use generator methods.

```go
// Signature:
Set(value string) error

// Example:
err := vCode.Set("ABD531")
```

#### Generate

Generate a random numeric code with 5 character length and set as code.

```go
// Signature:
Generate() (string, error)

// Example:
code, err := vCode.Generate()
```

#### GenerateN

Generate a random numeric code with special character length and set as code.

```go
// Signature:
GenerateN(count uint) (string, error)

// Example:
code, err := vCode.GenerateN(6)
```

#### Clear

Clear code from cache.

```go
// Signature:
Clear() error

// Example:
err := vCode.Clear()
```

#### Get

Get code.

```go
// Signature:
Get() (string, error)

// Example:
code, err := vCode.Get()
```

#### Exists

Exists check if code exists in cache and not empty.

```go
// Signature:
Exists() (bool, error)

// Example:
exists, err := vCode.Exists()
```

#### TTL

Get token ttl.

```go
// Signature:
TTL() (time.Duration, error)

// Example:
ttl, err := vCode.TTl()
```
