package handler

import (
	"github.com/shj1081/sso/internal/config"
	"github.com/shj1081/sso/internal/service"
	"github.com/shj1081/sso/internal/storer"
)

type Handler struct {
	cfg   *config.Config
	st    storer.Storer
	OAuth *service.OAuthService
	// JWT   *service.JWTService
	Email *service.EmailService // ✅ 이메일 서비스 추가
}

func NewHandler(cfg *config.Config, st storer.Storer, oauth *service.OAuthService, email *service.EmailService) *Handler {
	return &Handler{
		cfg:   cfg,
		st:    st,
		OAuth: oauth,
		// JWT:   jwt,
		Email: email,
	}
}
