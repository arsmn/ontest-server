package user

type PersistenceProvider interface {
	UserPersister() Persister
}

type Persister interface {
	Find(id int64) (*User, error)
}
