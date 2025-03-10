package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/shj1081/sso/internal/storer"
)

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

func (h *Handler) SubmitSignup(w http.ResponseWriter, r *http.Request) {
	sessCookie, err := r.Cookie("sso_session")
	if err != nil {
		http.Error(w, "no session cookie", http.StatusUnauthorized)
		return
	}

	sessionID := sessCookie.Value
	sd, err := h.Session.GetSession(r.Context(), sessionID)
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

	updateUser := &storer.User{
		ID:       sd.UserId,
		Name:     req.Name,
		SkkuMail: req.SkkuMail,
		Phone:    req.Phone,
		UserType: func() string {
			if req.SkkuMail == "" {
				return "external"
			} else {
				return "skkuin"
			}
		}(),
	}
	user, err := h.Session.UpdateUser(r.Context(), updateUser)
	if err != nil {
		http.Error(w, "failed to create user:"+err.Error(), http.StatusInternalServerError)
		return
	}

	jwtToken, _ := h.JWT.CreateJWT(user.ID)
	h.JWT.SetAuthCookie(w, jwtToken)

	_ = h.Session.DeleteSession(r.Context(), sessionID)

	// 원래 서비스로 redirect
	http.Redirect(w, r, sd.OriginalURL, http.StatusFound)
}
