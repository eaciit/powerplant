package helpers

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"time"
)

func TrunctTime(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}

func Decode(bytesData []byte, result interface{}) error {
	buf := bytes.NewBuffer(bytesData)
	dec := gob.NewDecoder(buf)
	e := dec.Decode(result)
	return e
}

func Encode(obj interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	gw := gob.NewEncoder(buf)
	err := gw.Encode(obj)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
