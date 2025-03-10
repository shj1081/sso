package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/shj1081/sso/config"
	"github.com/shj1081/sso/sso/storer"
)

var cfg = config.LoadConfig()

// processKakaoCallback handles the code callback from Kakao OAuth
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

	// 1) Exchange the code for an Access Token
	tokenResp, err := getKakaoAccessToken(code)
	if err != nil {
		http.Error(w, "failed to get token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2) Use the Access Token to call the Kakao User Info API
	userInfo, err := getKakaoUserInfo(tokenResp.AccessToken)
	if err != nil {
		http.Error(w, "failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3) The Kakao user’s unique ID is userInfo.ID
	kakaoIDStr := fmt.Sprintf("%d", userInfo.ID)

	// 4) Look up your local user by the Kakao user’s ID
	user, err := h.server.FindUserByKakaoID(h.ctx, kakaoIDStr)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5) If the user is found, generate the JWT & redirect. Otherwise, create a session for onboarding.
	if user != nil {
		jwtToken, _ := createJWTForUser(user.ID)
		setAuthCookie(w, jwtToken)
		http.Redirect(w, r, originalURL, http.StatusFound)
		return
	}

	// Create a session for the new user to complete sign-up
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

// getKakaoAccessToken exchanges an authorization code for a Kakao access token
func getKakaoAccessToken(code string) (*KakaoTokenResponse, error) {

	clientID := cfg.KakaoClientID
	redirectURI := cfg.KakaoRedirectURI
	tokenURI := cfg.KakaoTokenURI

	if clientID == "" || redirectURI == "" || tokenURI == "" {
		return nil, errors.New("missing Kakao OAuth environment variables")
	}

	data := fmt.Sprintf(
		"grant_type=authorization_code&client_id=%s&redirect_uri=%s&code=%s",
		clientID, redirectURI, code,
	)

	req, err := http.NewRequest("POST", tokenURI, bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response from Kakao token endpoint: %s", body)
	}

	var tokenResp KakaoTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

// getKakaoUserInfo uses an access token to fetch a user’s Kakao profile
func getKakaoUserInfo(accessToken string) (*KakaoUserInfoResponse, error) {
	req, err := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		return nil, err
	}

	// Bearer token
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info from Kakao: %s", body)
	}

	var userInfo KakaoUserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
