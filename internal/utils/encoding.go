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

func ParseInt(str string) int {
	if str == "" {
		return 0
	}
	return int(ParseInt64(str))
}

func ParseInt64(str string) int64 {
	if str == "" {
		return 0
	}
	var i int64
	fmt.Sscanf(str, "%d", &i)
	return i
}
