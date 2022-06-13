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
	Data any
}

func (rc record) Serialize() (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(rc)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b.Bytes()), nil
}

func (rc *record) Deserialize(data string) error {
	by, err := hex.DecodeString(data)
	if err != nil {
		return err
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(rc)
	if err != nil {
		return err
	}
	return nil
}

func (rc record) IsExpired() bool {
	return rc.TTL.UTC().Before(time.Now().UTC())
}
