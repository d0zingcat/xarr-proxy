package auth

import (
	"net/http"
	"time"

	"xarr-proxy/internal/config"

	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func Init(cfg *config.Config) {
	tokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
}

func Issue(cfg *config.Config, data map[string]any) (string, string, error) {
	jwtauth.SetExpiryIn(data, time.Duration(cfg.TokenTTL)*time.Second)
	_, tokenString, err := tokenAuth.Encode(data)
	return tokenString, "Bearer", err
}

func GetVerifier() *jwtauth.JWTAuth {
	return tokenAuth
}

func SignJWT(cfg *config.Config, id int, username, role string, validStatus int) (string, error) {
	token, _, err := Issue(cfg, map[string]any{
		"user_id":      id,
		"username":     username,
		"role":         role,
		"valid_status": validStatus,
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func GetUserInfo(r *http.Request) (int, string, string, int) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	userId := int(claims["user_id"].(float64))
	username := claims["username"].(string)
	role := claims["role"].(string)
	validStatus := int(claims["valid_status"].(float64))
	return userId, username, role, validStatus
}
