package cert

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5String(data string) string {
	d := []byte(data)
	return MD5(d)
}

//md5加密
func MD5(data []byte) string {
	m := md5.New()
	m.Write(data)
	c := m.Sum(nil)
	return hex.EncodeToString(c)
}
