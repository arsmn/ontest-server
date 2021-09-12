package driver

import (
	"context"

	"github.com/arsmn/ontest/app"
	"github.com/arsmn/ontest/module/xlog"
	"github.com/arsmn/ontest/persistence"
	"github.com/arsmn/ontest/settings"
)

type Registry interface {
	Init(ctx context.Context) error
	WithLogger(l *xlog.Logger) Registry
	WithConfig(c *settings.Config) Registry

	settings.Provider
	xlog.Provider
	app.Provider
	persistence.Provider
	settings.Provider
}
