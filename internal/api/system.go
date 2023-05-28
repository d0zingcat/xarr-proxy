package api

import (
	"net/http"
)

var userWrongPassCnt int

func userLogin(w http.ResponseWriter, r *http.Request) {
	if userWrongPassCnt > 10 {
		w.WriteHeader(http.StatusForbidden)
		// get message from http status code
		w.Write([]byte(http.StatusText(http.StatusForbidden)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}

func userInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"id":1, "role":"admin"}`))
}
