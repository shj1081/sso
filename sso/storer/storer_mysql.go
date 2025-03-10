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
	return &MySQLStorer{
		db: db,
	}
}

func (ms *MySQLStorer) Close() error {
	return ms.db.Close()
}

func (ms *MySQLStorer) FindByKakaoID(kakaoID string) (*User, error) {
	u := User{}
	err := ms.db.Get(&u, "SELECT * FROM users WHERE kakao_id=?", kakaoID)
	if err != nil {
		return nil, fmt.Errorf("error getting user by kakao id: %v", err)
	}

	return &u, nil
}

func (ms *MySQLStorer) CreateUser(ctx context.Context, u *User) (*User, error) {
	res, err := ms.db.NamedExecContext(ctx,
		`
		INSERT INTO users (name, kakao_id, skku_mail, phone, usertype) 
		VALUES (:name, :kakao_id, :skku_mail, :phone, :usertype)
		`,
		u)

	if err != nil {
		return nil, fmt.Errorf("error inserting user: %v", err)
	}

	userId, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %v", err)
	}

	u.ID = userId
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return u, nil
}

func (ms *MySQLStorer) UpdateUser(ctx context.Context, u *User) (*User, error) {
	_, err := ms.db.NamedExecContext(ctx,
		`
		UPDATE users 
		SET name=:name, kakao_id=:kakao_id, skku_mail=:skku_mail, phone=:phone, usertype=:usertype
		WHERE id=:id
		`,
		u)

	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}

	return u, nil
}

func (ms *MySQLStorer) GetUserByID(ctx context.Context, id int64) (*User, error) {
	u := User{}
	err := ms.db.GetContext(ctx, &u, "SELECT * FROM users WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting user by id: %v", err)
	}

	return &u, nil
}

func (ms *MySQLStorer) DeleteUser(ctx context.Context, id int64) error {
	_, err := ms.db.ExecContext(ctx, "DELETE FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}

func (ms *MySQLStorer) CreateSkkuIn(ctx context.Context, si *SkkuIn) (*SkkuIn, error) {
	res, err := ms.db.NamedExecContext(ctx,
		`
		INSERT INTO skkuin (skkuin_type, department, student_id, user_id)
		VALUES (:skkuin_type, :department, :student_id, :user_id)
		`,
		si)

	if err != nil {
		return nil, fmt.Errorf("error inserting skkuin: %v", err)
	}

	siId, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %v", err)
	}

	si.ID = siId
	return si, nil
}

func (ms *MySQLStorer) UpdateSkkuIn(ctx context.Context, si *SkkuIn) (*SkkuIn, error) {
	_, err := ms.db.NamedExecContext(ctx,
		`
		UPDATE skkuin
		SET skkuin_type=:skkuin_type, department=:department, student_id=:student_id, user_id=:user_id
		WHERE id=:id
		`,
		si)

	if err != nil {
		return nil, fmt.Errorf("error updating skkuin: %v", err)
	}

	return si, nil
}

func (ms *MySQLStorer) GetSkkuInByUserID(ctx context.Context, userId int64) (*SkkuIn, error) {
	si := SkkuIn{}
	err := ms.db.GetContext(ctx, &si, "SELECT * FROM skkuin WHERE user_id=?", userId)
	if err != nil {
		return nil, fmt.Errorf("error getting skkuin by user id: %v", err)
	}

	return &si, nil
}

func (ms *MySQLStorer) CreateSession(ctx context.Context, s *Session) error {

	_, err := ms.db.NamedExecContext(ctx,
		`
		INSERT INTO sessions (session_id, kakao_id, original_url, created_at) 
		VALUES (:session_id, :kakao_id, :original_url, :created_at)
		`,
		s,
	)
	if err != nil {
		return fmt.Errorf("error inserting session: %v", err)
	}
	return nil
}

func (ms *MySQLStorer) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	var sd Session
	err := ms.db.GetContext(ctx, &sd,
		`
		SELECT session_id, kakao_id, original_url, created_at 
		FROM sessions 
		WHERE session_id = ?
		`, sessionID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting session by session_id: %v", err)
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
