package service

import (
	"context"
	stderr "errors"
	"time"

	"github.com/arsmn/ontest-server/app"
	"github.com/arsmn/ontest-server/module/errors"
	"github.com/arsmn/ontest-server/module/generate"
	v "github.com/arsmn/ontest-server/module/validation"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/session"
	t "github.com/arsmn/ontest-server/transport"
	"github.com/arsmn/ontest-server/user"
)

var _ app.App = new(Service)

func (s *Service) GetSession(ctx context.Context, token string) (*session.Session, error) {
	return s.dx.Persister().FindSessionByToken(ctx, token)
}

func (s *Service) DeleteSession(ctx context.Context, token string) error {
	sess, err := s.dx.Persister().FindSessionByToken(ctx, token)
	if err != nil {
		if stderr.Is(err, persistence.ErrNoRows) {
			return errors.ErrUnauthorized
		}
		return err
	}

	if !sess.IsActive() {
		return errors.ErrUnauthorized
	}

	return s.dx.Persister().RemoveSession(ctx, sess.ID)
}

func (s *Service) createSession(ctx context.Context, userID uint64) (*session.Session, error) {
	sess := &session.Session{
		ID:        generate.UID(),
		UserID:    userID,
		Token:     generate.RandomString(35, generate.AlphaNum),
		Active:    true,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(s.dx.Settings().Session().Lifespan),
	}
	return sess, s.dx.Persister().CreateSession(ctx, sess)
}

func (s *Service) IssueSession(ctx context.Context, req *t.SigninRequest) (*session.Session, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	user, err := s.dx.Persister().FindUserByEmail(ctx, req.Email)
	if err != nil {
		if stderr.Is(err, persistence.ErrNoRows) {
			return nil, app.ErrInvalidCredentials
		}
		return nil, err
	}

	if user.Password == "" {
		return nil, app.ErrInvalidCredentials
	}

	if err := s.dx.Hasher().Compare(ctx, []byte(req.Password), []byte(user.Password)); err != nil {
		return nil, app.ErrInvalidCredentials
	}

	return s.createSession(ctx, user.ID)
}

func (s *Service) OAuthIssueSession(ctx context.Context, req *t.OAuthSignRequest) (*session.Session, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	var newuser bool
	u, err := s.dx.Persister().FindUserByEmail(ctx, req.Email)
	if err != nil {
		if stderr.Is(err, persistence.ErrNoRows) {
			newuser = true
		} else {
			return nil, err
		}
	}

	if newuser {
		u = &user.User{
			ID:            generate.UID(),
			Username:      generate.HFUID(),
			Email:         req.Email,
			FirstName:     req.FirstName,
			LastName:      req.LastName,
			IsActive:      true,
			EmailVerified: true,
			Rands:         generate.UserRandCode(),
		}
		if _, err := s.createUser(ctx, u); err != nil {
			return nil, err
		}
	}

	return s.createSession(ctx, u.ID)
}
