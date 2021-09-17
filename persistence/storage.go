package persistence

import (
	"context"

	"github.com/arsmn/ontest-server/session"
	"github.com/arsmn/ontest-server/user"
)

type Provider interface {
	Persister() Persister
}

type Persister interface {
	user.Persister
	session.Persister

	Close(context.Context) error
}
