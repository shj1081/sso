package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/shj1081/sso/internal/storer"
)

type SessionService struct {
	st storer.Storer
}

func NewSessionService(st storer.Storer) *SessionService {
	return &SessionService{st: st}
}

func (s *SessionService) CreateSession(ctx context.Context, session *storer.Session) error {
	session.CreatedAt = time.Now()
	session.ExpiresAt = session.CreatedAt.Add(10 * time.Minute)
	return s.st.CreateSession(ctx, session)
}

func (s *SessionService) GetSession(ctx context.Context, sessionID string) (*storer.Session, error) {
	return s.st.GetSession(ctx, sessionID)
}

func (s *SessionService) DeleteSession(ctx context.Context, sessionID string) error {
	return s.st.DeleteSession(ctx, sessionID)
}

func (s *SessionService) CreateUser(ctx context.Context, user *storer.User) (*storer.User, error) {
	return s.st.CreateUser(ctx, user)
}

func (s *SessionService) UpdateUser(ctx context.Context, user *storer.User) (*storer.User, error) {
	return s.st.UpdateUser(ctx, user)
}

func (s *SessionService) GenerateRandomString(n int) string {
	bytes := make([]byte, n/2)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err) // 보안적으로 중요한 함수이므로 실패 시 패닉 처리
	}
	return hex.EncodeToString(bytes)
}
