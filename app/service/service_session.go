package service

import (
	"context"
	stderr "errors"

	"github.com/arsmn/ontest-server/app"
	"github.com/arsmn/ontest-server/module/errors"
	v "github.com/arsmn/ontest-server/module/validation"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/session"
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

func (s *Service) IssueSession(ctx context.Context, req *session.SigninRequest) (*session.Session, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	find := func(ctx context.Context, typ, value string) (*user.User, error) {
		switch typ {
		case "email":
			return s.dx.Persister().FindUserByEmail(ctx, value)
		case "username":
			return s.dx.Persister().FindUserByUsername(ctx, value)
		default:
			return nil, nil
		}
	}

	var user *user.User
	for _, t := range []string{"email", "username"} {
		u, err := find(ctx, t, req.Identifier)
		if err != nil {
			if stderr.Is(err, persistence.ErrNoRows) {
				continue
			}
			return nil, err
		}
		user = u
		break
	}

	if user == nil {
		return nil, app.ErrInvalidCredentials
	}

	if !user.PasswordSet() {
		return nil, app.ErrInvalidCredentials
	}

	if err := s.dx.Hasher().Compare(ctx, []byte(req.Password), []byte(user.Password)); err != nil {
		return nil, app.ErrInvalidCredentials
	}

	sess := session.NewActiveSession(user.ID, s.dx.Settings().Session.Lifespan)

	if err := s.dx.Persister().CreateSession(ctx, sess); err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *Service) OAuthIssueSession(ctx context.Context, req *session.OAuthSignRequest) (*session.Session, error) {
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
		u = user.NewActiveUser(req.FirstName, req.LastName, req.Email)
		u.EmailVerified = true
		if _, err := s.createUser(ctx, u); err != nil {
			return nil, err
		}
	}

	sess := session.NewActiveSession(u.ID, s.dx.Settings().Session.Lifespan)

	if err := s.dx.Persister().CreateSession(ctx, sess); err != nil {
		return nil, err
	}

	return sess, nil
}
