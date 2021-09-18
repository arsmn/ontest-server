package user

import (
	"context"

	t "github.com/arsmn/ontest-server/transport"
)

type ServiceProvider interface {
	UserService() Service
}

type Service interface {
	GetUser(context.Context, uint64) (*User, error)
	RegisterUser(context.Context, *t.SignupRequest) (*User, error)
}
