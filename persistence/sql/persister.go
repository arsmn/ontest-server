package sql

import (
	"context"
	"time"

	"github.com/arsmn/ontest/module/sqlcon"
	"github.com/arsmn/ontest/module/xlog"
	"github.com/arsmn/ontest/persistence"
	"github.com/arsmn/ontest/settings"
	"github.com/cenkalti/backoff"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
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

func NewPersister(r persisterDependencies) (*Persister, error) {
	p := new(Persister)

	bc := backoff.NewExponentialBackOff()
	bc.MaxElapsedTime = time.Minute * 5
	bc.Reset()

	return p, backoff.Retry(
		func() error {
			maxOpenConns, maxIdleConns, connMaxLifetime, cleanedDSN := sqlcon.ParseConnectionOptions(r.Logger(), r.Settings().SQL().DSN)
			r.Logger().
				Debug("Connecting to SQL Database",
					xlog.Int("maxOpenConns", maxOpenConns),
					xlog.Int("maxIdleConns", maxIdleConns),
					xlog.Duration("connMaxLifetime", connMaxLifetime))

			db, err := sqlx.Connect(r.Settings().SQL().Driver, sqlcon.FinalizeDSN(r.Logger(), cleanedDSN))
			if err != nil {
				r.Logger().Warn("Unable to connect to database, retrying.", xlog.Err(err))
				return err
			}

			if err := db.Ping(); err != nil {
				r.Logger().Warn("Unable to ping database, retrying.", xlog.Err(err))
				return err
			}

			db.SetConnMaxLifetime(connMaxLifetime)
			db.SetMaxOpenConns(maxOpenConns)
			db.SetMaxIdleConns(maxIdleConns)

			p.r = r
			p.db = db

			return nil
		}, bc)
}

func (p *Persister) Close(ctx context.Context) error {
	return p.db.Close()
}
