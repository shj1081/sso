package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shj1081/sso/internal/config"
	"github.com/shj1081/sso/internal/handler"
	"github.com/shj1081/sso/internal/service"
	"github.com/shj1081/sso/internal/storer"
)

type Server struct {
	cfg *config.Config
	st  storer.Storer
	h   *handler.Handler
}

func NewServer(cfg *config.Config, st storer.Storer) *Server {

	// 서비스 계층 생성
	oauthSvc := service.NewOAuthService(cfg, st)
	jwtSvc := service.NewJWTService(cfg)
	emailSvc := service.NewEmailService(cfg, st)

	// 핸들러 생성
	h := handler.NewHandler(cfg, st, oauthSvc, jwtSvc, emailSvc)

	return &Server{cfg: cfg, st: st, h: h}
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// OAuth
	r.Get("/auth/kakao/callback", s.h.KakaoCallback)
	r.Post("/signup", s.h.SubmitSignup)

	// skku 메일 인증
	r.Post("/send-verification", s.h.SendVerification)
	r.Post("/verify-code/{userId}", s.h.VerifyCode)
	r.Post("/verify-code", s.h.VerifyCode)
	r.Post("/verify-code/{userId}", s.h.VerifyCode)

	return r
}
