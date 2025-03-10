package storer

import "time"

type User struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	KakaoID   string    `json:"kakao_id," db:"kakao_id"`
	SkkuMail  string    `json:"skku_mail," db:"skku_mail"`
	Phone     string    `json:"phone" db:"phone"`
	UserType  string    `json:"usertype" db:"usertype"` // ENUM('external', 'skkuin')
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type SkkuIn struct {
	ID         int64     `json:"id" db:"id"`
	SkkuinType string    `json:"skkuin_type" db:"skkuin_type"` // ENUM('student', 'professor', 'staff')
	Department string    `json:"department" db:"department"`
	StudentID  string    `json:"student_id," db:"student_id"`
	UserID     int64     `json:"user_id," db:"user_id"` // FK to users.id, can be NULL
	CreatedAt  time.Time `json:"created_at" db:"created_at "`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Session struct {
	SessionID   string    `json:"session_id" db:"session_id"`
	KakaoID     string    `json:"kakao_id" db:"kakao_id"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
}
