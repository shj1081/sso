package server

import (
	"context"

	"github.com/shj1081/sso/sso/storer"
)

type Server struct {
	storer *storer.MySQLStorer
}

func NewServer(storer *storer.MySQLStorer) *Server {
	return &Server{
		storer: storer,
	}
}

func (s *Server) FindUserByKakaoID(ctx context.Context, kakaoID string) (*storer.User, error) {
	return s.storer.FindByKakaoID(kakaoID)
}

func (s *Server) CreateUser(ctx context.Context, u *storer.User) (*storer.User, error) {
	return s.storer.CreateUser(ctx, u)
}

func (s *Server) UpdateUser(ctx context.Context, u *storer.User) (*storer.User, error) {
	return s.storer.UpdateUser(ctx, u)
}

func (s *Server) GetUser(ctx context.Context, id int64) (*storer.User, error) {
	return s.storer.GetUserByID(ctx, id)
}

func (s *Server) DeleteUser(ctx context.Context, id int64) error {
	return s.storer.DeleteUser(ctx, id)
}

func (s *Server) CreateSkkuIn(ctx context.Context, si *storer.SkkuIn) (*storer.SkkuIn, error) {
	return s.storer.CreateSkkuIn(ctx, si)
}

func (s *Server) GetSkkuIn(ctx context.Context, id int64) (*storer.SkkuIn, error) {
	return s.storer.GetSkkuInByUserID(ctx, id)
}

func (s *Server) UpdateSkkuIn(ctx context.Context, si *storer.SkkuIn) (*storer.SkkuIn, error) {
	return s.storer.UpdateSkkuIn(ctx, si)
}

func (s *Server) CreateSession(ctx context.Context, session *storer.Session) error {
	return s.storer.CreateSession(ctx, session)
}

func (s *Server) GetSession(ctx context.Context, sessionID string) (*storer.Session, error) {

	return s.storer.GetSession(ctx, sessionID)
}

func (s *Server) DeleteSession(ctx context.Context, sessionID string) error {
	return s.storer.DeleteSession(ctx, sessionID)
}
