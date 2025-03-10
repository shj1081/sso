package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/shj1081/sso/sso/storer"
)

// 카카오 콜백 핸들러
func processKakaoCallback(h *handler, w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}
	originalURL := r.URL.Query().Get("state")
	if originalURL == "" {
		http.Error(w, "missing state(original url)", http.StatusBadRequest)
		return
	}

	token, err := getKakaoAccessToken(code)
	if err != nil {
		http.Error(w, "failed to get token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	kakaoIDStr := token
	user, err := h.server.FindUserByKakaoID(h.ctx, kakaoIDStr)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if user != nil {
		jwtToken, _ := createJWTForUser(user.ID)
		setAuthCookie(w, jwtToken)
		http.Redirect(w, r, originalURL, http.StatusFound)
		return
	}

	session := &storer.Session{
		SessionID:   generateSecureRandomString(32),
		KakaoID:     kakaoIDStr,
		OriginalURL: originalURL,
	}

	err = h.server.CreateSession(h.ctx, session)
	if err != nil {
		http.Error(w, "failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, session.SessionID)
	http.Redirect(w, r, os.Getenv("SSO_FE_SIGNUP_URL"), http.StatusFound)
}

// 카카오 액세스 토큰 요청
func getKakaoAccessToken(code string) (string, error) {
	clientID := os.Getenv("KAKAO_CLIENT_ID")
	redirectURI := os.Getenv("KAKAO_REDIRECT_URI")
	tokenURI := os.Getenv("KAKAO_TOKEN_URI")

	if clientID == "" || redirectURI == "" || tokenURI == "" {
		return "", errors.New("missing Kakao OAuth environment variables")
	}

	data := fmt.Sprintf("grant_type=authorization_code&client_id=%s&redirect_uri=%s&code=%s",
		clientID, redirectURI, code,
	)

	req, err := http.NewRequest("POST", tokenURI, bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected response from Kakao: %s", body)
	}

	var tokenResp KakaoTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}
