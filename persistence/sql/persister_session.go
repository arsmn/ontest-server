package sql

import (
	"context"

	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/session"
	"xorm.io/xorm"
)

var _ persistence.Persister = new(Persister)

func (p *Persister) findSession(_ context.Context, x *xorm.Session) (*session.Session, error) {
	s := new(session.Session)
	has, err := x.Get(s)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, persistence.ErrNoRows
	}
	return s, nil
}

func (p *Persister) FindSession(ctx context.Context, id uint64) (*session.Session, error) {
	return p.findSession(ctx, p.engine.ID(id))
}

func (p *Persister) FindSessionByToken(ctx context.Context, token string) (*session.Session, error) {
	return p.findSession(ctx, p.engine.Where("token = ?", token))
}

func (p *Persister) CreateSession(_ context.Context, s *session.Session) error {
	_, err := p.engine.InsertOne(s)
	return handleError(err)
}

func (p *Persister) RemoveSession(_ context.Context, id uint64) error {
	_, err := p.engine.ID(id).Delete(new(session.Session))
	return err
}
