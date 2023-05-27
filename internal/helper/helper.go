package helper

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func GetUserInfo(r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	userID := claims["user_id"]
	_ = userID
}
