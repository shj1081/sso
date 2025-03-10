package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

func setSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     "sso_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(10 * time.Minute),
	}
	http.SetCookie(w, cookie)
}

func generateSecureRandomString(n int) string {
	bytes := make([]byte, n/2)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err) // 보안적으로 중요한 함수이므로 실패 시 패닉 처리
	}
	return hex.EncodeToString(bytes)
}
