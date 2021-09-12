package sql

import (
	"context"

	"github.com/arsmn/ontest/module/xlog"
	"github.com/arsmn/ontest/persistence"
	"github.com/arsmn/ontest/settings"
	"github.com/jmoiron/sqlx"
)

var _ persistence.Persister = new(Persister)

type (
	persisterDependencies interface {
		xlog.Provider
		settings.Provider
	}
	Persister struct {
		db *sqlx.DB
		r  persisterDependencies
	}
)

func NewPersister(r persisterDependencies, db *sqlx.DB) *Persister {
	return &Persister{db, r}
}

func (p *Persister) Close(ctx context.Context) error {
	return p.db.Close()
}
