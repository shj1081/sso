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

	userID, redirectURL, err := h.OAuth.AuthenticateKakaoUser(r.Context(), code, originalURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect url이 SSOFeSignupURL가 아닌 경우만 jwt 발급
	if redirectURL != h.cfg.SSOFeSignupURL {
		h.JWT.SetAuthCookies(w, userID)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}
