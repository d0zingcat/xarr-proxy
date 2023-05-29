package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"xarr-proxy/internal/auth"
	"xarr-proxy/internal/cache"
	"xarr-proxy/internal/config"
	"xarr-proxy/internal/db"
	"xarr-proxy/internal/model"

	"golang.org/x/crypto/bcrypt"
)

var SystemUser = &systemUserService{}

type systemUserService struct{}

func (*systemUserService) Login(username, password string) (string, error) {
	password = strings.ToLower(password)
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user := new(db.SystemUser)
	if err := db.Get().First(&user, "username = ?", username, pass).Error; err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}
	if user.ValidStatus == 0 {
		return "", errors.New("user invalidated")
	}
	token, err := auth.SignJWT(config.Get(), user.Id, user.Username, user.Role, user.ValidStatus)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (*systemUserService) UserInfo() {
}

func (*systemUserService) Logout(userInfo model.SystemUser, token string) bool {
	cache.Get().Set(token, true, time.Second*time.Duration(config.Get().TokenBlockTTL))
	fmt.Println(cache.Get().Items())
	return true
}
