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

	// session에서 user id 가져오기
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

	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// 이메일 전송
	err = h.Email.SendVerificationEmail(r.Context(), req.Email, sd.VerifyCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Verification email sent"})
}

func (h *Handler) SendVerificationById(w http.ResponseWriter, r *http.Request) {

	// user id로 verify code 가져오기
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
	verifyCode, err := h.Email.GetVerifyCodeByID(r.Context(), userIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	// 이메일 전송
	err = h.Email.SendVerificationEmail(r.Context(), req.Email, verifyCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Verification email sent"})
}

// 이메일 인증 코드 확인 API
func (h *Handler) VerifyCode(w http.ResponseWriter, r *http.Request) {

	// session에서 user id 가져오기
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

	var req VerifyCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	err = h.Email.VerifyCode(r.Context(), req.Code, sd.VerifyCode, sd.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Verification code verified"})
}

func (h *Handler) VerifyCodeById(w http.ResponseWriter, r *http.Request) {

	// user id로 verify code 가져오기
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

	verifyCode, err := h.Email.GetVerifyCodeByID(r.Context(), userIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req VerifyCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	err = h.Email.VerifyCode(r.Context(), req.Code, verifyCode, userIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Verification code verified"})
}
