package cache_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bopher/cache"
)

func fileCache() cache.Cache {
	return cache.NewFileCache("mine", "./caches")
}

func TestFileCachePut(t *testing.T) {
	err := fileCache().Put("name", "kim", time.Second)
	if err != nil {
		t.Fatal(err)
	}

	v, err := fileCache().Get("name")
	if err != nil {
		t.Fatal(err)
	}

	if v != "kim" {
		t.Fatalf("failed put %s", v)
	}
}

func TestFileCacheSet(t *testing.T) {
	exists, err := fileCache().Set("non-exists", "Bla")
	if err != nil {
		t.Fatal(err)
	}

	if exists {
		t.Fatalf(`failed exists check!`)
	}

	err = fileCache().Put("name", "John", time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	exists, err = fileCache().Set("name", "Kate")
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatalf(`failed exists check!`)
	}

	v, err := fileCache().Get("name")
	if err != nil {
		t.Fatal(err)
	}

	if v != "Kate" {
		t.Fatalf(`Want "Kate" get %s`, v)
	}
}

func TestFileCacheForget(t *testing.T) {
	err := fileCache().Put("name", "kim", time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	err = fileCache().Forget("name")
	if err != nil {
		t.Fatal(err)
	}

	v, err := fileCache().Get("name")
	if err != nil {
		t.Fatal(err)
	}

	if v != nil {
		t.Fatal("failed forget!")
	}
}

func TestFileCachePull(t *testing.T) {
	err := fileCache().Put("name", "kim", time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	v, err := fileCache().Pull("name")
	if err != nil {
		t.Fatal(err)
	}

	if v == nil {
		t.Fatal("failed pull get!")
	}

	v, err = fileCache().Get("name")
	if err != nil {
		t.Fatal(err)
	}

	if v != nil {
		t.Fatal("failed pull forget!")
	}
}

func TestFileCacheTTL(t *testing.T) {
	err := fileCache().Put("name", "kim", time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	ttl, err := fileCache().TTL("name")
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

func TestFileCacheIncDecFloat(t *testing.T) {
	err := fileCache().Put("float-val", 10.1, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	exists, err := fileCache().IncrementFloat("float-val", 0.3)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("item not exists!")
	}

	v, err := fileCache().Get("float-val")
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprint(v) != "10.4" {
		t.Fatal("failed increment")
	}

	exists, err = fileCache().DecrementFloat("float-val", 0.5)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("item not exists!")
	}

	v, err = fileCache().Get("float-val")
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprint(v) != "9.9" {
		t.Fatal("failed decrement")
	}
}

func TestFileCacheIncDec(t *testing.T) {
	err := fileCache().Put("int-val", 3, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	exists, err := fileCache().Increment("int-val", 6)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("item not exists!")
	}

	v, err := fileCache().Get("int-val")
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprint(v) != "9" {
		t.Fatalf("failed increment")
	}

	exists, err = fileCache().Decrement("int-val", 2)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("item not exists!")
	}

	v, err = fileCache().Get("int-val")
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprint(v) != "7" {
		t.Fatal("failed decrement")
	}
}

func TestCleanup(t *testing.T) {
	err := os.RemoveAll("./caches")
	if err != nil {
		t.Fatal(err)
	}
}
