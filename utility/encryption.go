package utility

import (
	"crypto/md5"
	"encoding/hex"
)

func EncryptMd5(str string) string {
	m := md5.New()
	_, err := m.Write([]byte(str))
	if err != nil {
		panic(err)
	}
	bytes := m.Sum(nil)
	return hex.EncodeToString(bytes)
}
