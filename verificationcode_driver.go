package cache

import (
	"time"

	"github.com/bopher/utils"
)

// vcDriver verification code manager
type vcDriver struct {
	Key   string
	TTL   time.Duration
	Cache Cache
}

// Set set code
func (vc *vcDriver) Set(value string) {
	vc.Cache.Forget(vc.Key)
	vc.Cache.Put(vc.Key, value, vc.TTL)
}

// Generate generate a random numeric code with 5 character length
func (vc *vcDriver) Generate() (string, error) {
	val, err := utils.RandomStringFromCharset(5, "0123456789")
	if err != nil {
		return "", err
	}
	vc.Set(val)
	return val, nil
}

// GenerateN generate a random numeric code with special character length
func (vc *vcDriver) GenerateN(count uint) (string, error) {
	val, err := utils.RandomStringFromCharset(count, "0123456789")
	if err != nil {
		return "", err
	}
	vc.Set(val)
	return val, nil
}

// Clear clear code
func (vc *vcDriver) Clear() {
	vc.Cache.Forget(vc.Key)
}

// Get get code
func (vc *vcDriver) Get() string {
	return vc.Cache.String(vc.Key, "")
}

// Exists check if code exists
func (vc *vcDriver) Exists() bool {
	return vc.Cache.Exists(vc.Key)
}
