package handler

import (
	"context"
	"net/http"

	"github.com/shj1081/sso/config"
	"github.com/shj1081/sso/sso/server"
)

type handler struct {
	ctx    context.Context
	server *server.Server
	config *config.Config
}

func NewHandler(s *server.Server) *handler {
	return &handler{
		ctx:    context.Background(),
		server: s,
		config: config.LoadConfig(),
	}
}

// 카카오 OAuth 콜백 핸들러
func (h *handler) kakaoCallback(w http.ResponseWriter, r *http.Request) {
	processKakaoCallback(h, w, r) // Kakao 핸들러 분리
}

// 회원가입 핸들러
func (h *handler) submitSignup(w http.ResponseWriter, r *http.Request) {
	processSignup(h, w, r) // 회원가입 핸들러 분리
}
