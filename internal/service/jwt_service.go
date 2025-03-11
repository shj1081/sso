package service

// import (
// 	"net/http"
// 	"time"

// 	"github.com/golang-jwt/jwt/v4"
// 	"github.com/shj1081/sso/internal/config"
// )

// type JWTService struct {
// 	cfg *config.Config
// }

// func NewJWTService(cfg *config.Config) *JWTService {
// 	return &JWTService{cfg: cfg}
// }

// type Claims struct {
// 	UserID int64 `json:"user_id"`
// 	jwt.RegisteredClaims
// }

// // Access Token 생성 (유효 기간: 15분)
// func (j *JWTService) CreateAccessToken(userID int64) (string, error) {
// 	claims := &Claims{
// 		UserID: userID,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString([]byte(j.cfg.JWTSecret))
// }

// // Refresh Token 생성 (유효 기간: 7일)
// func (j *JWTService) CreateRefreshToken(userID int64) (string, error) {
// 	claims := &Claims{
// 		UserID: userID,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString([]byte(j.cfg.JWTSecret))
// }

// // jwt cookie 생성
// func (j *JWTService) SetAuthCookies(w http.ResponseWriter, userID int64) {
// 	accessToken, err := j.CreateAccessToken(userID)
// 	if err != nil {
// 		http.Error(w, "failed to create access token", http.StatusInternalServerError)
// 		return
// 	}

// 	refreshToken, err := j.CreateRefreshToken(userID)
// 	if err != nil {
// 		http.Error(w, "failed to create refresh token", http.StatusInternalServerError)
// 		return
// 	}

// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "access_token",
// 		Value:    accessToken,
// 		Expires:  time.Now().Add(15 * time.Minute),
// 		HttpOnly: false,
// 		Secure:   true,
// 	})

// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "refresh_token",
// 		Value:    refreshToken,
// 		Expires:  time.Now().Add(7 * 24 * time.Hour),
// 		HttpOnly: false,
// 		Secure:   true,
// 	})
// }

// 서비스에서 충분히 reissue 가능

// func (j *JWTService) ParseToken(tokenStr string) (*Claims, error) {
// 	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(j.cfg.JWTSecret), nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
// 		return claims, nil
// 	}

// 	return nil, errors.New("invalid token")
// }

// func (j *JWTService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
// 	claims, err := j.ParseToken(refreshToken)
// 	if err != nil {
// 		return "", errors.New("invalid refresh token")
// 	}

// 	// 사용자 검증
// 	user, err := j.st.GetUserByID(ctx, claims.UserID)
// 	if err != nil || user == nil {
// 		return "", errors.New("user not found")
// 	}

// 	// 새로운 Access Token 발급
// 	return j.CreateAccessToken(user.ID)
// }
