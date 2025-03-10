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

func (ms *MySQLStorer) FindByKakaoID(ctx context.Context, kakaoID int64) (*User, error) {
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

func (ms *MySQLStorer) CreateUser(ctx context.Context, u *User) (*User, error) {
	now := time.Now()
	res, err := ms.db.ExecContext(ctx,
		`INSERT INTO users (name, kakao_id, skku_mail, phone, usertype, verify_code, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?, ?, ?)`,
		u.Name, u.KakaoID, u.SkkuMail, u.Phone, u.UserType, u.VerifyCode, now, now)
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

func (ms *MySQLStorer) UpdateUser(ctx context.Context, u *User) (*User, error) {
	now := time.Now()
	_, err := ms.db.ExecContext(ctx,
		`UPDATE users
         SET name=?, skku_mail=?, phone=?, usertype=?, updated_at=?
         WHERE id=?`,
		u.Name, u.KakaoID, u.SkkuMail, u.Phone, u.UserType, now, u.ID)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}
	u.UpdatedAt = now
	return u, nil
}

func (ms *MySQLStorer) UpdateUserByKakaoID(ctx context.Context, u *User) (*User, error) {
	now := time.Now()
	_, err := ms.db.ExecContext(ctx,
		`UPDATE users
		 SET name=?, skku_mail=?, phone=?, usertype=?, updated_at=?
		 WHERE kakao_id=?`,
		u.Name, u.SkkuMail, u.Phone, u.UserType, now, u.KakaoID)
	if err != nil {
		return nil, fmt.Errorf("error updating user by kakao id: %v", err)
	}
	u.UpdatedAt = now
	return u, nil
}

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

func (ms *MySQLStorer) GetVerifyCodeByID(ctx context.Context, id int64) (string, error) {
	var verifyCode string
	err := ms.db.GetContext(ctx, &verifyCode, "SELECT verify_code FROM users WHERE id=?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("error getting verify code by id: %v", err)
	}
	return verifyCode, nil
}

func (ms *MySQLStorer) DeleteUser(ctx context.Context, id int64) error {
	_, err := ms.db.ExecContext(ctx, "DELETE FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}

// 세션 관련
func (ms *MySQLStorer) CreateSession(ctx context.Context, s *Session) error {
	_, err := ms.db.ExecContext(ctx,
		`INSERT INTO sessions (session_id, user_id, verify_code, original_url, created_at, expires_at)
         VALUES (?, ?, ?, ?, ?)`,
		s.SessionID, s.UserId, s.VerifyCode, s.OriginalURL, s.CreatedAt, s.ExpiresAt)
	if err != nil {
		return fmt.Errorf("error inserting session: %v", err)
	}
	return nil
}

func (ms *MySQLStorer) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	var sd Session
	err := ms.db.GetContext(ctx, &sd,
		`SELECT session_id, user_id, verify_code, original_url, created_at, expires_at
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
