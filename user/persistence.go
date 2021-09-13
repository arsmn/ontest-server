package user

import "context"

type PersistenceProvider interface {
	UserPersister() Persister
}

type Persister interface {
	FindUser(context.Context, uint64) (*User, error)
	CreateUser(context.Context, *User) error
}
