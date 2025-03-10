package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

var r *chi.Mux

// RegisterRoutes: SSO 라우팅 설정
func RegisterRoutes(h *handler) *chi.Mux {
	r := chi.NewRouter()

	// 카카오 콜백
	r.Get("/auth/kakao/callback", h.kakaoCallback)

	// 회원가입 API
	r.Post("/signup", h.submitSignup)

	return r
}

func StartServer(addr string) error {
	return http.ListenAndServe(addr, r)
}
