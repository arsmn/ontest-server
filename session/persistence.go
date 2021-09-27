package session

import "context"

type PersistenceProvider interface {
	SessionPersister() Persister
}

type Persister interface {
	FindSession(context.Context, uint64) (*Session, error)
	FindSessionByToken(context.Context, string) (*Session, error)
	FindUserSessions(context.Context, uint64) ([]*Session, error)
	CreateSession(context.Context, *Session) error
	RemoveSession(context.Context, uint64) error
	RemoveUserSessions(context.Context, uint64, ...string) error
}
