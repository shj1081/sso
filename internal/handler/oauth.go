package handler

import (
	"net/http"

	"github.com/shj1081/sso/internal/storer"
)

func (h *Handler) KakaoCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	original_url := r.URL.Query().Get("state")
	if code == "" || original_url == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	tokenResp, err := h.OAuth.GetKakaoAccessToken(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInfoResp, err := h.OAuth.GetKakaoUserInfo(tokenResp.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 카카오 ID로 사용자 조회
	user, err := h.OAuth.FindByKakaoID(r.Context(), userInfoResp.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 사용자가 없으면 temp유저 및 session 생성 후 fe 로 redirect
	if user == nil {

		user := &storer.User{
			UserType:   "temp",
			KakaoID:    userInfoResp.ID,
			VerifyCode: h.Session.GenerateRandomString(6),
		} // temp user
		if _, err := h.Session.CreateUser(r.Context(), user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session := &storer.Session{
			SessionID:   h.Session.GenerateRandomString(16),
			UserId:      user.ID,
			VerifyCode:  user.VerifyCode,
			OriginalURL: original_url,
		}

		if err := h.Session.CreateSession(r.Context(), session); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "sso_session",
			Value: session.SessionID,
		})
		http.Redirect(w, r, h.cfg.SSOFeSignupURL, http.StatusFound)
		return
	} else {
		// 사용자가 있으면 jwt 발급 후 redirect
		jwtToken, _ := h.JWT.CreateJWT(user.ID)
		h.JWT.SetAuthCookie(w, jwtToken)
		http.Redirect(w, r, original_url, http.StatusFound)
		return
	}

}
