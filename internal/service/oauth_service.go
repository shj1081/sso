package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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
	AccessToken string `json:"access_token"`
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

func (o *OAuthService) GetKakaoUserInfo(accessToken string) (*KakaoUserInfoResponse, error) {
	req, err := http.NewRequest("GET", o.cfg.KaKaoUserInfoURI, nil)
	if err != nil {
		return nil, err
	}

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

func (o *OAuthService) AuthenticateKakaoUser(ctx context.Context, code, originalURL string) (string, error) {
	tokenResp, err := o.GetKakaoAccessToken(code)
	if err != nil {
		return "", fmt.Errorf("failed to get Kakao access token: %w", err)
	}

	userInfo, err := o.GetKakaoUserInfo(tokenResp.AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed to get Kakao user info: %w", err)
	}

	user, err := o.st.FindByKakaoID(ctx, userInfo.ID)
	if err != nil {
		return "", err
	}

	if user == nil {
		// Temp 유저 생성
		user = &storer.User{
			UserType:   "temp",
			KakaoID:    userInfo.ID,
			VerifyCode: GenerateRandomString(6),
		}

		if _, err := o.st.CreateUser(ctx, user); err != nil {
			return "", err
		}

		// 세션 생성
		session := &storer.Session{
			SessionID:   GenerateRandomString(16),
			UserId:      user.ID,
			VerifyCode:  user.VerifyCode,
			OriginalURL: originalURL,
			CreatedAt:   time.Now(),
			ExpiresAt:   time.Now().Add(10 * time.Minute),
		}

		if err := o.st.CreateSession(ctx, session); err != nil {
			return "", err
		}

		return fmt.Sprintf("%s?session_id=%s", o.cfg.SSOFeSignupURL, session.SessionID), nil
	}

	return originalURL, nil
}

func GenerateRandomString(n int) string {
	bytes := make([]byte, n/2)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err) // 보안적으로 중요한 함수이므로 실패 시 패닉 처리
	}
	return hex.EncodeToString(bytes)
}
