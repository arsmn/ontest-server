package sql

import (
	"github.com/arsmn/ontest/persistence"
	"github.com/arsmn/ontest/user"
)

var _ persistence.Persister = new(Persister)

func (p *Persister) FindUser(id int64) (*user.User, error) {
	panic("not implemented")
}
