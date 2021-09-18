package service

import (
	"context"
	"errors"
	"time"

	"github.com/arsmn/ontest-server/app"
	"github.com/arsmn/ontest-server/module/generate"
	v "github.com/arsmn/ontest-server/module/validation"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/session"
	t "github.com/arsmn/ontest-server/transport"
)

var _ app.App = new(Service)

func (s *Service) GetSession(ctx context.Context, token string) (*session.Session, error) {
	return s.dx.Persister().FindSessionByToken(ctx, token)
}

func (s *Service) IssueSession(ctx context.Context, req *t.SigninRequest) (*session.Session, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	user, err := s.dx.Persister().FindUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, persistence.ErrNoRows) {
			return nil, app.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := s.dx.Hasher().Compare(ctx, []byte(req.Password), []byte(user.Password)); err != nil {
		return nil, app.ErrInvalidCredentials
	}

	sess := &session.Session{
		ID:        generate.UID(),
		UserID:    user.ID,
		Token:     generate.RandomString(35, generate.AlphaNum),
		Active:    true,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(s.dx.Settings().Session().Lifespan),
	}

	if err := s.dx.Persister().CreateSession(ctx, sess); err != nil {
		return nil, err
	}

	return sess, nil
}
