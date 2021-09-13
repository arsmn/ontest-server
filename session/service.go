package session

import (
	"context"

	t "github.com/arsmn/ontest/transport"
)

type ServiceProvider interface {
	SessionService() Service
}

type Service interface {
	IssueSession(context.Context, *t.SigninRequest) (*t.SigninResponse, error)
}
