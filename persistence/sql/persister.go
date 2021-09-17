package sql

import (
	"context"
	"time"

	"github.com/arsmn/ontest-server/module/sqlcon"
	"github.com/arsmn/ontest-server/module/xlog"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/session"
	"github.com/arsmn/ontest-server/settings"
	"github.com/arsmn/ontest-server/user"
	"github.com/cenkalti/backoff"
	"xorm.io/xorm"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var _ persistence.Persister = new(Persister)

var tables = []interface{}{
	new(user.User),
	new(session.Session),
}

type (
	persisterDependencies interface {
		xlog.Provider
		settings.Provider
	}
	Persister struct {
		engine *xorm.Engine
		dx     persisterDependencies
	}
)

func NewPersister(dx persisterDependencies) (*Persister, error) {
	p := new(Persister)

	bc := backoff.NewExponentialBackOff()
	bc.MaxElapsedTime = time.Minute * 5
	bc.Reset()

	return p, backoff.Retry(
		func() error {
			maxOpenConns, maxIdleConns, connMaxLifetime, cleanedDSN := sqlcon.ParseConnectionOptions(dx.Logger(), dx.Settings().SQL().DSN)
			dx.Logger().
				Debug("Connecting to SQL Database",
					xlog.Int("maxOpenConns", maxOpenConns),
					xlog.Int("maxIdleConns", maxIdleConns),
					xlog.Duration("connMaxLifetime", connMaxLifetime))

			engine, err := xorm.NewEngine(dx.Settings().SQL().Driver, sqlcon.FinalizeDSN(dx.Logger(), cleanedDSN))
			if err != nil {
				dx.Logger().Warn("Unable to connect to database, retrying.", xlog.Err(err))
				return err
			}

			if err := engine.Ping(); err != nil {
				dx.Logger().Warn("Unable to ping database, retrying.", xlog.Err(err))
				return err
			}

			if err := engine.Sync2(tables...); err != nil {
				dx.Logger().Warn("Unable to sync database, retrying.", xlog.Err(err))
				return err
			}

			engine.SetConnMaxLifetime(connMaxLifetime)
			engine.SetMaxOpenConns(maxOpenConns)
			engine.SetMaxIdleConns(maxIdleConns)
			engine.ShowSQL(!dx.Settings().IsProd())

			p.dx = dx
			p.engine = engine

			return nil
		}, bc)
}

func (p *Persister) Close(ctx context.Context) error {
	return p.engine.Close()
}
