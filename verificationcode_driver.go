package cache

import (
	"time"

	"github.com/bopher/utils"
)

type vcDriver struct {
	key   string
	cache Cache
}

func (vc vcDriver) err(pattern string, params ...any) error {
	return utils.TaggedError([]string{"VerificationCode", vc.key}, pattern, params...)
}

func (vc vcDriver) notExistsErr() error {
	return utils.TaggedError([]string{"VerificationCode", "NotExists", vc.key}, "%s not exists", vc.key)
}

func (vc *vcDriver) init(key string, ttl time.Duration, cache Cache) error {
	vc.key = key
	vc.cache = cache

	exists, err := cache.Exists(key)
	if err != nil {
		return vc.err(err.Error())
	}

	if !exists {
		return cache.Put(key, "", ttl)
	}

	return nil
}

func (vc vcDriver) Set(value string) error {
	exists, err := vc.cache.Set(vc.key, value)
	if err != nil {
		return vc.err(err.Error())
	}

	if !exists {
		return vc.notExistsErr()
	}
	return nil
}

func (vc vcDriver) Generate() (string, error) {
	if val, err := utils.RandomStringFromCharset(5, "0123456789"); err != nil {
		return "", vc.err(err.Error())
	} else {
		return val, vc.Set(val)
	}
}

func (vc vcDriver) GenerateN(count uint) (string, error) {
	if val, err := utils.RandomStringFromCharset(count, "0123456789"); err != nil {
		return "", err
	} else {
		return val, vc.Set(val)
	}
}

func (vc vcDriver) Clear() error {
	if err := vc.cache.Forget(vc.key); err != nil {
		return vc.err(err.Error())
	}
	return nil
}

func (vc vcDriver) Get() (string, error) {
	caster, err := vc.cache.Cast(vc.key)
	if err != nil {
		return "", vc.err(err.Error())
	}

	if caster.IsNil() {
		return "", vc.notExistsErr()
	}

	v, err := caster.String()
	if err != nil {
		err = vc.err(err.Error())
	}

	return v, err
}

func (vc vcDriver) Exists() (bool, error) {
	exists, err := vc.cache.Exists(vc.key)
	if err != nil {
		err = vc.err(err.Error())
	}
	return exists, err
}

func (vc vcDriver) TTL() (time.Duration, error) {
	if v, err := vc.cache.TTL(vc.key); err != nil {
		return 0, vc.err(err.Error())
	} else {
		return v, nil
	}
}
