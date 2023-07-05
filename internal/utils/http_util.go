package utils

import (
	"net/http"
	"net/http/httputil"
	"time"
)

func GetClient(timeout *int) *http.Client {
	c := &http.Client{}
	if timeout == nil {
		c.Timeout = 30 * time.Second
	} else {
		c.Timeout = time.Duration(*timeout) * time.Second
	}
	return c
}

func DumpResponse(resp *http.Response) string {
	bs, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return ""
	}
	return string(bs)
}

func DumpRequest(req http.Request) string {
	bs, err := httputil.DumpRequest(&req, true)
	if err != nil {
		return ""
	}
	return string(bs)
}
