package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/shj1081/sso/internal/config"
)

type OAuthService struct {
	cfg *config.Config
}

func NewOAuthService(cfg *config.Config) *OAuthService {
	return &OAuthService{cfg: cfg}
}

type KakaoTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

type KakaoUserInfoResponse struct {
	ID int64 `json:"id"`
}

func (o *OAuthService) GetKakaoAccessToken(code string) (*KakaoTokenResponse, error) {
	data := fmt.Sprintf(
		"grant_type=authorization_code&client_id=%s&redirect_uri=%s&code=%s",
		o.cfg.KakaoClientID, o.cfg.KakaoRedirectURI, code,
	)

	req, err := http.NewRequest("POST", o.cfg.KakaoTokenURI, bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResp KakaoTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// getKakaoUserInfo uses an access token to fetch a user’s Kakao profile
func (o *OAuthService) GetKakaoUserInfo(accessToken string) (*KakaoUserInfoResponse, error) {
	req, err := http.NewRequest("GET", o.cfg.KaKaoUserInfoURI, nil)
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
