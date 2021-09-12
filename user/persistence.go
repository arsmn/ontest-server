package user

type PersistenceProvider interface {
	UserPersister() Persister
}

type Persister interface {
	FindUser(id int64) (*User, error)
}
