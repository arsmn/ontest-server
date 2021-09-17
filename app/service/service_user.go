package service

import (
	"context"

	"github.com/arsmn/ontest-server/app"
	"github.com/arsmn/ontest-server/module/generate"
	v "github.com/arsmn/ontest-server/module/validation"
	t "github.com/arsmn/ontest-server/transport"
	"github.com/arsmn/ontest-server/user"
)

var _ app.App = new(Service)

func (s *Service) RegisterUser(ctx context.Context, req *t.SignupRequest) (*t.SignupResponse, error) {
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
