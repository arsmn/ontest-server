package sql

import (
	"context"

	"github.com/arsmn/ontest/persistence"
	"github.com/arsmn/ontest/user"
	"xorm.io/xorm"
)

var _ persistence.Persister = new(Persister)

func (p *Persister) findUser(_ context.Context, x *xorm.Session) (*user.User, error) {
	u := new(user.User)
	has, err := x.Get(u)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, persistence.ErrNoRows
	}
	return u, nil
}

func (p *Persister) FindUser(ctx context.Context, id uint64) (*user.User, error) {
	return p.findUser(ctx, p.engine.ID(id))
}

func (p *Persister) FindUserByEmail(ctx context.Context, email string) (*user.User, error) {
	return p.findUser(ctx, p.engine.Where("email = ?", email))
}

func (p *Persister) CreateUser(_ context.Context, u *user.User) error {
	_, err := p.engine.InsertOne(u)
	return handleError(err)
}
