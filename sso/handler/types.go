package handler

import "time"

type UserRequest struct {
	Name     string `json:"name"`
	KakaoID  string `json:"kakao_id"`
	SkkuMail string `json:"skku_mail"`
	Phone    string `json:"phone"`
	UserType string `json:"usertype"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	KakaoID   string    `json:"kakao_id"`
	SkkuMail  string    `json:"skku_mail"`
	Phone     string    `json:"phone"`
	UserType  string    `json:"usertype"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type KakaoTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type KakaoUserInfoResponse struct {
	ID int64 `json:"id"`
}
