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

for creating redis based driver you must pass prefix, redis host, maxIdle, maxActive and db number to constructor function.

```go
import "github.com/bopher/cache"
if rCache := cache.NewRedisCache("myApp", "localhost:6379", 50, 10000, 1); rCache != nil {
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
Put(key string, value interface{}, ttl time.Duration) bool

// Example:
ok := rCache.Put("total-debt", 410203, 100 * time.Hour)
```

### PutForever

Put a new value to cache with infinite ttl.

```go
// Signature:
PutForever(key string, value interface{}) bool

// Example:
ok := rCache.Put("total-users", 5000000)
```

### Set

Change value of cache item

```go
// Signature:
Set(key string, value interface{}) bool

// Example:
ok := rCache.Set("total-users", 5005)
```

### Get

Get item from cache. Get function return cache item as `interface{}`. if you need get cache item with type use helper get functions described later.

```go
// Signature:
Get(key string) interface{}

// Example:
v := rCache.Get("total-users")
```

### Pull

Get item from cache and remove it

```go
// Signature:
Pull(key string) interface{}

// Example:
v := rCache.Pull("total-users")
```

### Exists

Check if item exists in cache

```go
// Signature:
Exists(key string) bool

// Example:
exists := rCache.Exists("total-users");
```

### Forget

Delete item from cache.

```go
// Signature:
Forget(key string) bool

// Example:
deleted := rCache.Forget("total-users")
```

### TTL

Get cache item ttl.

```go
// Signature:
TTL(key string) time.Duration

// Example:
ttl := rCache.TTL("total-users")
```

### Increment

Increment numeric item in cache.

```go
// Signature:
Increment(key string) bool

// Example:
ok := rCache.Increment("total-users")
```

### IncrementBy

Increment numeric item in cache by number.

```go
// Signature:
IncrementBy(key string, value interface{}) bool

// Example:
ok := rCache.IncrementBy("total-users", 10)
```

### Decrement

Decrement numeric item in cache.

```go
// Signature:
Decrement(key string) bool

// Example:
ok := rCache.Decrement("total-users")
```

### DecrementBy

Decrement numeric item in cache by number.

```go
DecrementBy(key string, value interface{}) bool
```

### Get By Type Methods

Helper get methods return fallback value if value not exists in cache.

```go
// Bool parse dependency as boolean
Bool(key string, fallback bool) bool
// Int parse dependency as int
Int(key string, fallback int) int
// Int8 parse dependency as int8
Int8(key string, fallback int8) int8
// Int16 parse dependency as int16
Int16(key string, fallback int16) int16
// Int32 parse dependency as int32
Int32(key string, fallback int32) int32
// Int64 parse dependency as int64
Int64(key string, fallback int64) int64
// UInt parse dependency as uint
UInt(key string, fallback uint) uint
// UInt8 parse dependency as uint8
UInt8(key string, fallback uint8) uint8
// UInt16 parse dependency as uint16
UInt16(key string, fallback uint16) uint16
// UInt32 parse dependency as uint32
UInt32(key string, fallback uint32) uint32
// UInt64 parse dependency as uint64
UInt64(key string, fallback uint64) uint64
// Float32 parse dependency as float64
Float32(key string, fallback float32) float32
// Float64 parse dependency as float64
Float64(key string, fallback float64) float64
// String parse dependency as string
String(key string, fallback string) string
// Bytes parse dependency as bytes array
Bytes(key string, fallback []byte) []byte
```

## Create New Rate Limiter Driver

**Note:** Rate limiter based on cache, For creating rate limiter driver you must pass a cache driver instance to constructor function.

```go
// Signature:
NewRateLimiter(key string, maxAttempts uint32, ttl time.Duration, cache Cache)

// Example:
import "github.com/bopher/cache"
limiter := cache.NewRateLimiter("login-attempts", 3, 60 * time.Second, rCache)
```

## Usage

Rate limiter interface contains following methods:

### Hit

Decrease the allowed times.

```go
// Signature:
Hit()

// Example:
limiter.Hit()
```

### Lock

Lock rate limiter.

```go
// Signature:
Lock()

// Example:
limiter.Lock() // no more attempts left
```

### Reset

Reset rate limiter.

```go
// Signature:
Reset()

// Example:
limiter.Reset()
```

### MustLock

Check if rate limiter must lock access.

```go
// Signature:
MustLock() bool

// Example:
if limiter.MustLock() {
  // Block access
}
```

### TotalAttempts

Get user attempts count.

```go
// Signature:
TotalAttempts() uint32

// Example:
totalAtt := limiter.TotalAttempts() // 3
```

### RetriesLeft

Get user retries left.

```go
// Signature:
RetriesLeft() uint32

// Example:
leftRet := limiter.RetriesLeft() // 2
```

### AvailableIn

Get time until unlock.

```go
// Signature:
AvailableIn() time.Duration

// Example:
availableIn := limiter.AvailableIn()
```

## Create New Verification Code Driver

verification code used for managing verification code sent to user.

**Note:** Verification code based on cache, For creating verification code driver you must pass a cache driver instance to constructor function.

```go
// Signature:
NewVerificationCode(key string, ttl time.Duration, cache Cache)

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
Set(value string)

// Example:
vCode.Set("ABD531")
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
Clear()

// Example:
vCode.Clear()
```

### Get

Get code.

```go
// Signature:
Get() string

// Example:
code := vCode.Get()
```

### Exists

Exists check if code exists in cache and not empty.

```go
// Signature:
Exists() bool

// Example:
exists := vCode.Exists()
```
