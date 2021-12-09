package cache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bopher/cache"
	"github.com/go-redis/redis/v8"
)

func redisCache() cache.Cache {
	return cache.NewRedisCache("test", redis.Options{Addr: "localhost:6379"})
}

func TestRedisCachePut(t *testing.T) {
	err := redisCache().Put("name", "kim", time.Second)
	if err != nil {
		t.Fatal(err)
	}

	v, err := redisCache().Get("name")
	if err != nil {
		t.Fatal(err)
	}

	if v != "kim" {
		t.Fatalf("failed put %s", v)
	}
}

func TestRedisCacheSet(t *testing.T) {
	exists, err := redisCache().Set("non-exists", "Bla")
	if err != nil {
		t.Fatal(err)
	}

	if exists {
		t.Fatalf(`failed exists check!`)
	}

	err = redisCache().Put("name", "John", time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	exists, err = redisCache().Set("name", "Kate")
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatalf(`failed exists check!`)
	}

	v, err := redisCache().Get("name")
	if err != nil {
		t.Fatal(err)
	}

	if v != "Kate" {
		t.Fatalf(`Want "Kate" get %s`, v)
	}
}

func TestRedisCacheForget(t *testing.T) {
	err := redisCache().Put("name", "kim", time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	err = redisCache().Forget("name")
	if err != nil {
		t.Fatal(err)
	}

	v, err := redisCache().Get("name")
	if err != nil {
		t.Fatal(err)
	}

	if v != nil {
		t.Fatal("failed forget!")
	}
}

func TestRedisCachePull(t *testing.T) {
	err := redisCache().Put("name", "kim", time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	v, err := redisCache().Pull("name")
	if err != nil {
		t.Fatal(err)
	}

	if v == nil {
		t.Fatal("failed pull get!")
	}

	v, err = redisCache().Get("name")
	if err != nil {
		t.Fatal(err)
	}

	if v != nil {
		t.Fatal("failed pull forget!")
	}
}

func TestRedisCacheTTL(t *testing.T) {
	err := redisCache().Put("name", "kim", time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	ttl, err := redisCache().TTL("name")
	if err != nil {
		t.Fatal(err)
	}

	if ttl < 59*time.Second {
		t.Fail()
	}

	ttl, err = redisCache().TTL("non-exists")
	if err != nil {
		t.Fatal(err)
	}

	if ttl > 0 {
		t.Fatal("failed non exists ttl")
	}
}

func TestRedisCacheIncDecFloat(t *testing.T) {
	err := redisCache().Put("float-val", 10.1, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	exists, err := redisCache().IncrementFloat("float-val", 0.3)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("item not exists!")
	}

	v, err := redisCache().Get("float-val")
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprint(v) != "10.4" {
		t.Fatalf("failed increment, %v", v)
	}

	exists, err = redisCache().DecrementFloat("float-val", 0.5)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("item not exists!")
	}

	v, err = redisCache().Get("float-val")
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprint(v) != "9.9" {
		t.Fatal("failed decrement")
	}
}

func TestRedisCacheIncDec(t *testing.T) {
	err := redisCache().Put("int-val", 3, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	exists, err := redisCache().Increment("int-val", 6)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("item not exists!")
	}

	v, err := redisCache().Get("int-val")
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprint(v) != "9" {
		t.Fatalf("failed increment")
	}

	exists, err = redisCache().Decrement("int-val", 2)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("item not exists!")
	}

	v, err = redisCache().Get("int-val")
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprint(v) != "7" {
		t.Fatal("failed decrement")
	}
}
