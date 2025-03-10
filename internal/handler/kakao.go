package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) KakaoCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	tokenResp, err := h.OAuth.GetKakaoAccessToken(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := h.OAuth.GetKakaoUserInfo(tokenResp.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Kakao User Info: %+v", userInfo)
}
