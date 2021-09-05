package persistence

import "github.com/arsmn/ontest/user"

type Provider interface {
	Persister() Persister
}

type Persister interface {
	user.Persister
}
