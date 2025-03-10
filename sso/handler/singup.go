package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/shj1081/sso/sso/storer"
)

// 회원가입 처리
func processSignup(h *handler, w http.ResponseWriter, r *http.Request) {
	sessCookie, err := r.Cookie("sso_session")
	if err != nil {
		http.Error(w, "no session cookie", http.StatusUnauthorized)
		return
	}
	sessionID := sessCookie.Value

	sd, err := h.server.GetSession(h.ctx, sessionID)
	if err != nil {
		http.Error(w, "session db error:"+err.Error(), http.StatusInternalServerError)
		return
	}
	if sd == nil {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user, err := h.server.CreateUser(h.ctx, &storer.User{
		KakaoID:  sd.KakaoID,
		Name:     req.Name,
		SkkuMail: req.SkkuMail,
		Phone:    req.Phone,
		UserType: req.UserType,
	})
	if err != nil {
		http.Error(w, "failed to create user:"+err.Error(), http.StatusInternalServerError)
		return
	}

	jwtToken, _ := createJWTForUser(user.ID)
	setAuthCookie(w, jwtToken)

	_ = h.server.DeleteSession(h.ctx, sessionID)
	expired := &http.Cookie{
		Name:    "sso_session",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	}
	http.SetCookie(w, expired)

	resp := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		KakaoID:   user.KakaoID,
		SkkuMail:  user.SkkuMail,
		Phone:     user.Phone,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
