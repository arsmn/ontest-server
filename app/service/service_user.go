package service

import (
	"context"

	"github.com/arsmn/ontest/app"
	v "github.com/arsmn/ontest/module/validation"
	t "github.com/arsmn/ontest/transport"
	"github.com/arsmn/ontest/user"
)

var _ app.App = new(Service)

func (s *Service) Signup(ctx context.Context, req *t.SignupRequest) (*t.SignupResponse, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	return nil, s.dx.Persister().CreateUser(ctx, &user.User{})
}
