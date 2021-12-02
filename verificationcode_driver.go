package cache

import (
	"time"

	"github.com/bopher/utils"
)

type vcDriver struct {
	Key   string
	TTL   time.Duration
	Cache Cache
}

func (this vcDriver) err(pattern string, params ...interface{}) error {
	return utils.TaggedError([]string{"VerificationCode", this.Key}, pattern, params...)
}

func (this vcDriver) Set(value string) error {
	if err := this.Cache.Forget(this.Key); err != nil {
		return this.err(err.Error())
	}
	if err := this.Cache.Put(this.Key, value, this.TTL); err != nil {
		return this.err(err.Error())
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
	if err := this.Cache.Forget(this.Key); err != nil {
		return this.err(err.Error())
	}
	return nil
}

func (this vcDriver) Get() (string, error) {
	if v, err := this.Cache.StringE(this.Key); err != nil {
		return "", this.err(err.Error())
	} else {
		return v, nil
	}
}

func (this vcDriver) Exists() (bool, error) {
	if exists, err := this.Cache.Exists(this.Key); err != nil {
		return false, this.err(err.Error())
	} else {
		return exists, nil
	}
}
