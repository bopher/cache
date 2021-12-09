package cache

import (
	"time"

	"github.com/bopher/utils"
)

type vcDriver struct {
	key   string
	ttl   time.Duration
	cache Cache
}

func (this vcDriver) err(pattern string, params ...interface{}) error {
	return utils.TaggedError([]string{"VerificationCode", this.key}, pattern, params...)
}

func (this vcDriver) notExistsErr() error {
	return utils.TaggedError([]string{"VerificationCode", "NotExists", this.key}, "%s not exists", this.key)
}

func (this *vcDriver) init(key string, ttl time.Duration, cache Cache) error {
	this.key = key
	this.cache = cache

	exists, err := cache.Exists(key)
	if err != nil {
		return this.err(err.Error())
	}

	if !exists {
		return cache.Put(key, "", ttl)
	}

	return nil
}

func (this vcDriver) Set(value string) error {
	exists, err := this.cache.Set(this.key, value)
	if err != nil {
		return this.err(err.Error())
	}

	if !exists {
		return this.notExistsErr()
	}
	return nil
}

func (this vcDriver) Generate() (string, error) {
	if val, err := utils.RandomStringFromCharset(5, "0123456789"); err != nil {
		return "", this.err(err.Error())
	} else {
		return val, this.Set(val)
	}
}

func (this vcDriver) GenerateN(count uint) (string, error) {
	if val, err := utils.RandomStringFromCharset(count, "0123456789"); err != nil {
		return "", err
	} else {
		return val, this.Set(val)
	}
}

func (this vcDriver) Clear() error {
	if err := this.cache.Forget(this.key); err != nil {
		return this.err(err.Error())
	}
	return nil
}

func (this vcDriver) Get() (string, error) {
	caster, err := this.cache.Cast(this.key)
	if err != nil {
		return "", this.err(err.Error())
	}

	if caster.IsNil() {
		return "", this.notExistsErr()
	}

	v, err := caster.String()
	if err != nil {
		err = this.err(err.Error())
	}

	return v, err
}

func (this vcDriver) Exists() (bool, error) {
	exists, err := this.cache.Exists(this.key)
	if err != nil {
		err = this.err(err.Error())
	}
	return exists, err
}
