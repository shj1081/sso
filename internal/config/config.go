package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	KakaoClientID    string
	KakaoRedirectURI string
	KakaoTokenURI    string
	KaKaoUserInfoURI string
	JWTSecret        string
	SSOFeSignupURL   string
	DBURL            string
	ServerAddress    string
}

func LoadConfig() (*Config, error) {
	// .env 파일 로드
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	cfg := &Config{
		KakaoClientID:    getEnv("KAKAO_CLIENT_ID"),
		KakaoRedirectURI: getEnv("KAKAO_REDIRECT_URI"),
		KakaoTokenURI:    getEnv("KAKAO_TOKEN_URI"),
		KaKaoUserInfoURI: getEnv("KAKAO_USER_INFO_URI"),
		JWTSecret:        getEnv("JWT_SECRET"),
		SSOFeSignupURL:   getEnv("SSO_FE_SIGNUP_URL"),
		DBURL:            getEnv("DB_URL"),
		ServerAddress:    getEnv("SERVER_ADDRESS"),
	}

	return cfg, nil
}

func getEnv(key string, defaultValue ...string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}
