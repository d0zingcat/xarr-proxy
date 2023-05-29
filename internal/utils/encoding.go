package utils

import (
	"crypto/md5"
	"fmt"
)

// encoding string to md5
func MD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return fmt.Sprintf("%x", m.Sum(nil))
}
