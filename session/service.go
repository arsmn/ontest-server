package session

import (
	"context"

	t "github.com/arsmn/ontest-server/transport"
)

type ServiceProvider interface {
	SessionService() Service
}

type Service interface {
	GetSession(context.Context, string) (*Session, error)
	IssueSession(context.Context, *t.SigninRequest) (*Session, error)
}
