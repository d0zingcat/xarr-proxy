package api

import (
	"errors"
	"net/http"

	"xarr-proxy/internal/api/req"
	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/helper"
	"xarr-proxy/internal/model"
	"xarr-proxy/internal/services"

	"github.com/go-chi/render"
)

var userWrongPassCnt int

const (
	MAX_LOGIN_FAIL_CNT = 10
)

func userLogin(w http.ResponseWriter, r *http.Request) {
	if userWrongPassCnt > MAX_LOGIN_FAIL_CNT {
		w.WriteHeader(http.StatusForbidden)
		// get message from http status code
		w.Write([]byte(http.StatusText(http.StatusForbidden)))
		return
	}
	req := new(req.SystemUserLoginReq)
	if err := render.Bind(r, req); err != nil {
		userWrongPassCnt++
		render.JSON(w, r, ErrInvalidRequest(err))
		return
	}
	token, err := services.SystemUser.Login(req.Username, req.Password)
	if token == "" || err != nil {
		userWrongPassCnt++
		render.JSON(w, r, ErrInvalidRequest(err))
		return
	}
	render.JSON(w, r, token)
}

func userLogout(w http.ResponseWriter, r *http.Request) {
	token := helper.ExtractToken(r)
	userInfo := r.Context().Value(consts.USER_INFO_CTX_KEY).(model.SystemUser)
	data := services.SystemUser.Logout(userInfo, token)
	render.JSON(w, r, data)
}

func userInfo(w http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(consts.USER_INFO_CTX_KEY)
	if userInfo == nil {
		render.JSON(w, r, ErrInvalidRequest(errors.New("invalid user info")))
		return
	}
	render.JSON(w, r, userInfo)
}
