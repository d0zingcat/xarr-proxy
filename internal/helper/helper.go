package helper

import (
	"net/http"
	"strings"
)

func ExtractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token == "" {
		return ""
	}
	return strings.Replace(token, "Bearer ", "", 1)
}
