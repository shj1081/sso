package storer

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MySQLStorer struct {
	db *sqlx.DB
}

func NewMySQLStorer(db *sqlx.DB) *MySQLStorer {
	return &MySQLStorer{db: db}
}

func (ms *MySQLStorer) Close() error {
	return ms.db.Close()
}

// FindByKakaoID는 카카오 ID로 사용자를 조회합니다.
func (ms *MySQLStorer) FindByKakaoID(ctx context.Context, kakaoID string) (*User, error) {
	var u User
	err := ms.db.GetContext(ctx, &u, "SELECT * FROM users WHERE kakao_id=?", kakaoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting user by kakao id: %v", err)
	}
	return &u, nil
}

// CreateUser는 새 사용자를 생성합니다.
func (ms *MySQLStorer) CreateUser(ctx context.Context, u *User) (*User, error) {
	now := time.Now()
	res, err := ms.db.ExecContext(ctx,
		`INSERT INTO users (name, kakao_id, skku_mail, phone, usertype, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?, ?, ?)`,
		u.Name, u.KakaoID, u.SkkuMail, u.Phone, u.UserType, now, now)
	if err != nil {
		return nil, fmt.Errorf("error inserting user: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %v", err)
	}

	u.ID = id
	u.CreatedAt = now
	u.UpdatedAt = now
	return u, nil
}

// UpdateUser는 기존 사용자를 수정합니다.
func (ms *MySQLStorer) UpdateUser(ctx context.Context, u *User) (*User, error) {
	now := time.Now()
	_, err := ms.db.ExecContext(ctx,
		`UPDATE users
         SET name=?, kakao_id=?, skku_mail=?, phone=?, usertype=?, updated_at=?
         WHERE id=?`,
		u.Name, u.KakaoID, u.SkkuMail, u.Phone, u.UserType, now, u.ID)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}
	u.UpdatedAt = now
	return u, nil
}

// GetUserByID는 ID로 사용자를 조회합니다.
func (ms *MySQLStorer) GetUserByID(ctx context.Context, id int64) (*User, error) {
	var u User
	err := ms.db.GetContext(ctx, &u, "SELECT * FROM users WHERE id=?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting user by id: %v", err)
	}
	return &u, nil
}

// DeleteUser는 사용자 정보를 삭제합니다.
func (ms *MySQLStorer) DeleteUser(ctx context.Context, id int64) error {
	_, err := ms.db.ExecContext(ctx, "DELETE FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}

// CreateSkkuIn, GetSkkuInByUserID, UpdateSkkuIn도 동일한 패턴으로 작성
func (ms *MySQLStorer) CreateSkkuIn(ctx context.Context, si *SkkuIn) (*SkkuIn, error) {
	now := time.Now()
	res, err := ms.db.ExecContext(ctx,
		`INSERT INTO skkuin (skkuin_type, department, student_id, user_id, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?, ?)`,
		si.SkkuinType, si.Department, si.StudentID, si.UserID, now, now)
	if err != nil {
		return nil, fmt.Errorf("error inserting skkuin: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %v", err)
	}
	si.ID = id
	si.CreatedAt = now
	si.UpdatedAt = now
	return si, nil
}

func (ms *MySQLStorer) GetSkkuInByUserID(ctx context.Context, userId int64) (*SkkuIn, error) {
	var si SkkuIn
	err := ms.db.GetContext(ctx, &si, "SELECT * FROM skkuin WHERE user_id=?", userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting skkuin by user id: %v", err)
	}
	return &si, nil
}

func (ms *MySQLStorer) UpdateSkkuIn(ctx context.Context, si *SkkuIn) (*SkkuIn, error) {
	now := time.Now()
	_, err := ms.db.ExecContext(ctx,
		`UPDATE skkuin
         SET skkuin_type=?, department=?, student_id=?, user_id=?, updated_at=?
         WHERE id=?`,
		si.SkkuinType, si.Department, si.StudentID, si.UserID, now, si.ID)
	if err != nil {
		return nil, fmt.Errorf("error updating skkuin: %v", err)
	}
	si.UpdatedAt = now
	return si, nil
}

// 세션 관련
func (ms *MySQLStorer) CreateSession(ctx context.Context, s *Session) error {
	_, err := ms.db.ExecContext(ctx,
		`INSERT INTO sessions (session_id, kakao_id, original_url, created_at, expires_at)
         VALUES (?, ?, ?, ?, ?)`,
		s.SessionID, s.KakaoID, s.OriginalURL, s.CreatedAt, s.ExpiresAt)
	if err != nil {
		return fmt.Errorf("error inserting session: %v", err)
	}
	return nil
}

func (ms *MySQLStorer) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	var sd Session
	err := ms.db.GetContext(ctx, &sd,
		`SELECT session_id, kakao_id, original_url, created_at, expires_at
         FROM sessions WHERE session_id=?`,
		sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting session: %v", err)
	}
	return &sd, nil
}

func (ms *MySQLStorer) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := ms.db.ExecContext(ctx, "DELETE FROM sessions WHERE session_id=?", sessionID)
	if err != nil {
		return fmt.Errorf("error deleting session: %v", err)
	}
	return nil
}
