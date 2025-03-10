package storer

import (
	"context"
)

// Storer는 MySQLStorer가 구현해야 할 메서드를 정의한 인터페이스입니다.
type Storer interface {
	Close() error

	// User 관련
	FindByKakaoID(ctx context.Context, kakaoID int64) (*User, error)
	CreateUser(ctx context.Context, u *User) (*User, error)
	UpdateUser(ctx context.Context, u *User) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetVerifyCodeByID(ctx context.Context, id int64) (string, error)
	DeleteUser(ctx context.Context, id int64) error

	// Session 관련
	CreateSession(ctx context.Context, s *Session) error
	GetSession(ctx context.Context, sessionID string) (*Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
}
