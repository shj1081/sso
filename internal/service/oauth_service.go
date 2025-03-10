package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/shj1081/sso/internal/config"
	"github.com/shj1081/sso/internal/storer"
)

type OAuthService struct {
	cfg *config.Config
	st  storer.Storer
}

func NewOAuthService(cfg *config.Config, st storer.Storer) *OAuthService {
	return &OAuthService{cfg: cfg, st: st}
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

func (o *OAuthService) FindByKakaoID(ctx context.Context, kakaoID int64) (*storer.User, error) {
	user, err := o.st.FindByKakaoID(ctx, kakaoID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (o *OAuthService) redirectFrontend(w http.ResponseWriter, r *http.Request, user *storer.User) {
	http.Redirect(w, r, fmt.Sprintf("%s?user_id=%d", o.cfg.SSOFeSignupURL, user.ID), http.StatusFound)
	return
}
