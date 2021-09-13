package driver

import (
	"context"

	"github.com/arsmn/ontest/module/xlog"
	"github.com/arsmn/ontest/settings"
)

func New(ctx context.Context) Registry {
	l := xlog.New()
	c := settings.New(l)

	var r Registry = NewRegistryCore().WithConfig(c).WithLogger(l)

	if err := r.Init(ctx); err != nil {
		l.Fatal("Unable to initialize service registry.", xlog.Err(err))
	}

	return r
}
