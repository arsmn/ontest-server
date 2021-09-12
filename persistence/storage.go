package persistence

import (
	"context"

	"github.com/arsmn/ontest/user"
)

type Provider interface {
	Persister() Persister
}

type Persister interface {
	user.Persister

	Close(context.Context) error
}
