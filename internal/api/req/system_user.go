package req

import (
	"net/http"
)

type (
	SystemUserUpdateReq struct {
		Username string
		Password string
	}
	SystemUserLoginReq struct {
		Username string
		Password string
	}
)

func (*SystemUserUpdateReq) Bind(r *http.Request) error {
	return nil
}

func (*SystemUserLoginReq) Bind(r *http.Request) error {
	return nil
}
