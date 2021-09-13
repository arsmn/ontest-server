package driver

import (
	"context"
	"sync"

	"github.com/arsmn/ontest/app"
	"github.com/arsmn/ontest/app/service"
	"github.com/arsmn/ontest/module/hash"
	"github.com/arsmn/ontest/module/xlog"
	"github.com/arsmn/ontest/persistence"
	"github.com/arsmn/ontest/persistence/sql"
	"github.com/arsmn/ontest/settings"
)

type RegistryCore struct {
	mtx sync.RWMutex
	l   *xlog.Logger
	c   *settings.Config

	app            app.App
	persister      persistence.Persister
	passwordHasher hash.Hasher
}

func NewRegistryCore() *RegistryCore {
	return &RegistryCore{}
}

func (r *RegistryCore) Init(ctx context.Context) error {
	p, err := sql.NewPersister(r)
	if err != nil {
		return err
	}

	r.persister = p
	r.app = service.NewApp(r)
	r.passwordHasher = hash.NewHasherArgon2(r)

	return nil
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

func (r *RegistryCore) App() app.App {
	return r.app
}

func (r *RegistryCore) Persister() persistence.Persister {
	return r.persister
}

func (r *RegistryCore) Hasher() hash.Hasher {
	return r.passwordHasher
}
