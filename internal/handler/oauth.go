package handler

import (
	"net/http"
)

func (h *Handler) KakaoCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	originalURL := r.URL.Query().Get("state")
	if code == "" || originalURL == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	redirectURL, err := h.OAuth.AuthenticateKakaoUser(r.Context(), code, originalURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}
