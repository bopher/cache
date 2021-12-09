package cache_test

import (
	"testing"
	"time"

	"github.com/bopher/cache"
)

func TestSetAndGet(t *testing.T) {
	vCode, err := cache.NewVerificationCode("test", time.Minute, redisCache())
	if err != nil {
		t.Fatal(err)
	}

	err = vCode.Set("ABCDE")
	if err != nil {
		t.Fatal(err)
	}

	code, err := vCode.Get()
	if err != nil {
		t.Fatal(err)
	}

	if code != "ABCDE" {
		t.Fail()
	}
}

func TestGenerateAndGenerateN(t *testing.T) {
	vCode, err := cache.NewVerificationCode("test", time.Minute, redisCache())
	if err != nil {
		t.Fatal(err)
	}

	code, err := vCode.Generate()
	if err != nil {
		t.Fatal(err)
	}

	if len(code) != 5 {
		t.Fatal()
	}

	code, err = vCode.GenerateN(10)
	if err != nil {
		t.Fatal(err)
	}

	if len(code) != 10 {
		t.Fatal()
	}
}

func TestClearAndExists(t *testing.T) {
	vCode, err := cache.NewVerificationCode("test", time.Minute, redisCache())
	if err != nil {
		t.Fatal(err)
	}

	err = vCode.Clear()
	if err != nil {
		t.Fatal(err)
	}

	exists, err := vCode.Exists()
	if err != nil {
		t.Fatal(err)
	}

	if exists {
		t.Fail()
	}
}
