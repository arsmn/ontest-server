package sql

import (
	"context"

	"github.com/arsmn/ontest/persistence"
	"github.com/arsmn/ontest/user"
)

var _ persistence.Persister = new(Persister)

func (p *Persister) FindUser(ctx context.Context, id int64) (*user.User, error) {
	panic("not implemented")
}

func (p *Persister) CreateUser(ctx context.Context, u *user.User) error {
	panic("not implemented")
}
