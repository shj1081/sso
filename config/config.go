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
	JWTSecret        string
	SSOFeSignupURL   string
	DBDriver         string
	DBURL            string
	ServerAddress    string
}

// LoadConfig 함수: .env 파일을 로드하고 Config 구조체를 반환
func LoadConfig() *Config {
	// .env 파일 로드
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	return &Config{
		KakaoClientID:    getEnv("KAKAO_CLIENT_ID", ""),
		KakaoRedirectURI: getEnv("KAKAO_REDIRECT_URI", ""),
		KakaoTokenURI:    getEnv("KAKAO_TOKEN_URI", ""),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		SSOFeSignupURL:   getEnv("SSO_FE_SIGNUP_URL", ""),
		DBDriver:         getEnv("DB_DRIVER", "mysql"),
		DBURL:            getEnv("DB_URL", "root:1234@tcp(localhost:3306)/sso?parseTime=true"),
		ServerAddress:    getEnv("SERVER_ADDRESS", ":8080"),
	}
}

// getEnv 함수: 환경 변수를 가져오고, 없으면 기본값을 반환
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
