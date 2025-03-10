package handler

import (
	"github.com/shj1081/sso/internal/config"
	"github.com/shj1081/sso/internal/service"
)

type Handler struct {
	cfg     *config.Config
	OAuth   *service.OAuthService
	JWT     *service.JWTService
	Session *service.SessionService
}

func NewHandler(cfg *config.Config, oauth *service.OAuthService, jwt *service.JWTService, sess *service.SessionService) *Handler {
	return &Handler{
		cfg:     cfg,
		OAuth:   oauth,
		JWT:     jwt,
		Session: sess,
	}
}
