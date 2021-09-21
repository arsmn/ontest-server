package driver

import (
	"context"

	"github.com/arsmn/ontest-server/app"
	"github.com/arsmn/ontest-server/module/cache"
	"github.com/arsmn/ontest-server/module/hash"
	"github.com/arsmn/ontest-server/module/mail"
	"github.com/arsmn/ontest-server/module/oauth"
	"github.com/arsmn/ontest-server/module/xlog"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/settings"
)

type Registry interface {
	Init(ctx context.Context) error
	WithLogger(l *xlog.Logger) Registry
	WithConfig(c *settings.Config) Registry

	settings.Provider
	xlog.Provider
	app.Provider
	persistence.Provider
	hash.Provider
	oauth.Provider
	cache.Provider
	mail.Provider
}
