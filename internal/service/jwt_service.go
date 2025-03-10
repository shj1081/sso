package service

import (
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/shj1081/sso/internal/config"
)

type JWTService struct {
	tokenAuth *jwtauth.JWTAuth
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{
		tokenAuth: jwtauth.New("HS256", []byte(cfg.JWTSecret), nil),
	}
}

func (j *JWTService) CreateJWT(userID int64) (string, error) {
	claims := map[string]interface{}{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	_, token, err := j.tokenAuth.Encode(claims)
	return token, err
}

func (j *JWTService) SetAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}
