package sql

import (
	"context"

	"github.com/arsmn/ontest/persistence"
	"github.com/arsmn/ontest/user"
)

var _ persistence.Persister = new(Persister)

func (p *Persister) FindUser(_ context.Context, id uint64) (*user.User, error) {
	u := new(user.User)
	has, err := p.engine.ID(id).Get(u)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, persistence.ErrNoRows
	}
	return u, nil
}

func (p *Persister) CreateUser(_ context.Context, u *user.User) error {
	_, err := p.engine.InsertOne(u)
	return handleError(err)
}
