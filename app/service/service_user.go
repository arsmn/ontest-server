package service

import (
	"context"

	"github.com/arsmn/ontest/app"
	"github.com/arsmn/ontest/module/generate"
	v "github.com/arsmn/ontest/module/validation"
	t "github.com/arsmn/ontest/transport"
	"github.com/arsmn/ontest/user"
)

var _ app.App = new(Service)

func (s *Service) Signup(ctx context.Context, req *t.SignupRequest) (*t.SignupResponse, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	pswd, err := s.dx.Hasher().Generate(ctx, []byte(req.Password))
	if err != nil {
		return nil, err
	}

	u := &user.User{
		ID:       generate.UID(),
		Username: generate.HFUID(),
		Email:    req.Email,
		IsActive: true,
		Password: string(pswd),
	}

	return nil, s.dx.Persister().CreateUser(ctx, u)
}
