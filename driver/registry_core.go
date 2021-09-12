package driver

import (
	"context"
	"sync"
	"time"

	"github.com/arsmn/ontest/app"
	"github.com/arsmn/ontest/module/sqlcon"
	"github.com/arsmn/ontest/module/xlog"
	"github.com/arsmn/ontest/persistence"
	"github.com/arsmn/ontest/persistence/sql"
	"github.com/arsmn/ontest/settings"
	"github.com/cenkalti/backoff"
	"github.com/jmoiron/sqlx"
)

type RegistryCore struct {
	mtx sync.RWMutex
	l   *xlog.Logger
	c   *settings.Config

	app       app.App
	persister persistence.Persister
}

func NewRegistryCore() *RegistryCore {
	return &RegistryCore{}
}

func (r *RegistryCore) Init(ctx context.Context) error {
	bc := backoff.NewExponentialBackOff()
	bc.MaxElapsedTime = time.Minute * 5
	bc.Reset()

	return backoff.Retry(
		func() error {
			maxOpenConns, maxIdleConns, connMaxLifetime, cleanedDSN := sqlcon.ParseConnectionOptions(r.l, r.Settings().SQL().DSN)
			r.Logger().
				Debug("Connecting to SQL Database",
					xlog.Int("maxOpenConns", maxOpenConns),
					xlog.Int("maxIdleConns", maxIdleConns),
					xlog.Duration("connMaxLifetime", connMaxLifetime))

			db, err := sqlx.ConnectContext(ctx, r.c.SQL().Driver, sqlcon.FinalizeDSN(r.l, cleanedDSN))
			if err != nil {
				r.Logger().Warn("Unable to connect to database, retrying.", xlog.Err(err))
				return err
			}

			if err := db.PingContext(ctx); err != nil {
				r.Logger().Warn("Unable to ping database, retrying.", xlog.Err(err))
				return err
			}

			db.SetConnMaxLifetime(connMaxLifetime)
			db.SetMaxOpenConns(maxOpenConns)
			db.SetMaxIdleConns(maxIdleConns)

			r.persister = sql.NewPersister(r, db)
			return nil
		}, bc)
}

func (r *RegistryCore) WithLogger(l *xlog.Logger) Registry {
	r.l = l
	return r
}

func (r *RegistryCore) WithConfig(c *settings.Config) Registry {
	r.c = c
	return r
}

func (r *RegistryCore) Logger() *xlog.Logger {
	return r.l
}

func (r *RegistryCore) Settings() *settings.Config {
	return r.c
}
