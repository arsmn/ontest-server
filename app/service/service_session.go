package service

import (
	"context"
	"time"

	"github.com/arsmn/ontest/module/generate"
	v "github.com/arsmn/ontest/module/validation"
	"github.com/arsmn/ontest/session"
	t "github.com/arsmn/ontest/transport"
)

func (s *Service) IssueSession(ctx context.Context, req *t.SigninRequest) (*t.SigninResponse, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	user, err := s.dx.Persister().FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if err := s.dx.Hasher().Compare(ctx, []byte(req.Password), []byte(user.Password)); err != nil {
		return nil, err
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

	return &t.SigninResponse{
		Token: sess.Token,
	}, nil
}
