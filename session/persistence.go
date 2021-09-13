package session

import "context"

type PersistenceProvider interface {
	SessionPersister() Persister
}

type Persister interface {
	FindSession(context.Context, uint64) (*Session, error)
	FindSessionByToken(context.Context, string) (*Session, error)
	CreateSession(context.Context, *Session) error
}
