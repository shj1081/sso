package handler

import (
	"fmt"
	"net/http"
	"time"
)

func createJWTForUser(userID int64) (string, error) {
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create JWT: %v", err)
	}
	return tokenString, nil
}

// 인증 쿠키 설정
func setAuthCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)
}
