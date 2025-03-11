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

	// sso 도메인과 서비스의 도메인이 달라 쿠키를 사용할 수 없기에 jwt는 서비스 벡엔드에서

	// // return user id가 -1이면 fe로 redirect 되므로, 그 외의 경우만 JWT 발급
	// if userID > 0 {
	// 	h.JWT.SetAuthCookies(w, userID)
	// }

	// return user id가 -1이면 sso_session 생성 (redirectURL 뒤에서 16글자)
	if userID == -1 {
		session_id := redirectURL[len(redirectURL)-16:]

		cookie := &http.Cookie{
			Name:   "session_id",
			Value:  session_id,
			Path:   "/",
			Domain: "localhost",
			Secure: false, // HTTPS가 아니라면 false
		}
		http.SetCookie(w, cookie)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}
