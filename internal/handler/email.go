package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type VerifyRequest struct {
	Email string `json:"email"`
}

type VerifyCodeRequest struct {
	Code string `json:"code"`
}

// 이메일 인증 코드 요청 API
func (h *Handler) SendVerification(w http.ResponseWriter, r *http.Request) {
	sessCookie, err := r.Cookie("sso_session")
	if err != nil {
		http.Error(w, "no session cookie", http.StatusUnauthorized)
		return
	}

	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	err = h.Email.SendVerificationEmailBySession(r.Context(), sessCookie.Value, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Verification email sent"})
}

func (h *Handler) SendVerificationById(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user id is required", http.StatusBadRequest)
		return
	}

	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	err = h.Email.SendVerificationEmailByUserID(r.Context(), userIDInt, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Verification email sent"})
}

func (h *Handler) VerifyCode(w http.ResponseWriter, r *http.Request) {
	sessCookie, err := r.Cookie("sso_session")
	if err != nil {
		http.Error(w, "no session cookie", http.StatusUnauthorized)
		return
	}

	var req VerifyCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = h.Email.VerifyCodeBySession(r.Context(), sessCookie.Value, req.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Verification code verified"})
}

func (h *Handler) VerifyCodeById(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user id is required", http.StatusBadRequest)
		return
	}

	userIDInt, _ := strconv.ParseInt(userID, 10, 64)

	var req VerifyCodeRequest
	json.NewDecoder(r.Body).Decode(&req)

	err := h.Email.VerifyCodeByUserID(r.Context(), userIDInt, req.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Verification code verified"})
}
