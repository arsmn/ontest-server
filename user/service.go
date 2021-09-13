package user

import (
	"context"

	t "github.com/arsmn/ontest/transport"
)

type ServiceProvider interface {
	UserService() Service
}

type Service interface {
	Signup(context.Context, *t.SignupRequest) (*t.SignupResponse, error)
}
