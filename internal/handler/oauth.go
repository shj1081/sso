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

	// return user id가 -1이면 fe로 redirect 되므로, 그 외의 경우만 JWT 발급
	if userID > 0 {
		h.JWT.SetAuthCookies(w, userID)
	}

	// return user id가 -2이면 sso_session 생성 (redirectURL 뒤에서 16글자)
	if userID == -1 {
		session_id := redirectURL[len(redirectURL)-16:]
		session_cookie := &http.Cookie{
			Name:  "sso_session",
			Value: session_id,
		}

		http.SetCookie(w, session_cookie)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}
