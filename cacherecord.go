package cache

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"time"
)

// cache record used for working with file cache

type record struct {
	TTL  time.Time
	Data interface{}
}

func (this record) Serialize() (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(this)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b.Bytes()), nil
}

func (this *record) Deserialize(data string) error {
	by, err := hex.DecodeString(data)
	if err != nil {
		return err
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(this)
	if err != nil {
		return err
	}
	return nil
}

func (this record) IsExpired() bool {
	return this.TTL.UTC().Before(time.Now().UTC())
}
