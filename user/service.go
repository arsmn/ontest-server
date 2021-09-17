package user

import (
	"context"

	t "github.com/arsmn/ontest-server/transport"
)

type ServiceProvider interface {
	UserService() Service
}

type Service interface {
	RegisterUser(context.Context, *t.SignupRequest) (*t.SignupResponse, error)
}
