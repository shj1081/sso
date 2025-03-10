package storer

import "time"

type User struct {
	ID         int64  `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	KakaoID    int64  `json:"kakao_id," db:"kakao_id"`
	SkkuMail   string `json:"skku_mail," db:"skku_mail"`
	Phone      string `json:"phone" db:"phone"`
	UserType   string `json:"usertype" db:"usertype"`       // ENUM('external', 'skkuin')
	VerifyCode string `json:"verify_code" db:"verify_code"` // skku mail 인증 코드

	// skkuin info
	Department string `json:"department" db:"department"`
	StudentID  string `json:"student_id" db:"student_id"`
	SkkuinType string `json:"skkuin_type" db:"skkuin_type"` // ENUM('student', 'professor', 'staff')

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Session struct {
	SessionID   string    `json:"session_id" db:"session_id"`
	UserId      int64     `json:"user_id" db:"user_id"`
	VerifyCode  string    `json:"verify_code" db:"verify_code"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
}
