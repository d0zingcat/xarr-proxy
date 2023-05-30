package utils

import (
	"net/http"
	"time"
)

func GetClient(timeout *int) *http.Client {
	c := &http.Client{}
	if timeout == nil {
		c.Timeout = 5 * time.Second
	} else {
		c.Timeout = time.Duration(*timeout) * time.Second
	}
	return c
}
