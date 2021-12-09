package cache_test

import (
	"testing"
	"time"

	"github.com/bopher/cache"
)

func TestHitTotalAndAttempts(t *testing.T) {
	limiter, err := cache.NewRateLimiter("test-hit", 5, time.Minute, redisCache())
	if err != nil {
		t.Fatal(err)
	}

	err = limiter.Reset()
	if err != nil {
		t.Fatal(err)
	}

	err = limiter.Hit()
	if err != nil {
		t.Fatal(err)
	}

	total, err := limiter.TotalAttempts()
	if err != nil {
		t.Fatal(err)
	}

	if total != 1 {
		t.Fail()
	}
}

func TestLockAndMustLock(t *testing.T) {
	limiter, err := cache.NewRateLimiter("test-lock", 5, time.Minute, redisCache())
	if err != nil {
		t.Fatal(err)
	}

	err = limiter.Reset()
	if err != nil {
		t.Fatal(err)
	}

	err = limiter.Lock()
	if err != nil {
		t.Fatal(err)
	}

	mustLock, err := limiter.MustLock()
	if err != nil {
		t.Fatal(err)
	}

	if !mustLock {
		t.Fail()
	}
}

func TestResetAndRemains(t *testing.T) {
	limiter, err := cache.NewRateLimiter("test-reset", 5, time.Minute, redisCache())
	if err != nil {
		t.Fatal(err)
	}

	err = limiter.Hit()
	if err != nil {
		t.Fatal(err)
	}

	err = limiter.Reset()
	if err != nil {
		t.Fatal(err)
	}

	remains, err := limiter.RetriesLeft()
	if err != nil {
		t.Fatal(err)
	}

	if remains != 5 {
		t.Fail()
	}
}

func TestAvailable(t *testing.T) {
	limiter, err := cache.NewRateLimiter("test-available", 5, time.Minute, redisCache())
	if err != nil {
		t.Fatal(err)
	}

	err = limiter.Reset()
	if err != nil {
		t.Fatal(err)
	}

	err = limiter.Lock()
	if err != nil {
		t.Fatal(err)
	}

	ttl, err := limiter.AvailableIn()
	if err != nil {
		t.Fatal(err)
	}

	if ttl < 58*time.Second {
		t.Fail()
	}
}
