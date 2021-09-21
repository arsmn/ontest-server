package service

import (
	"github.com/arsmn/ontest-server/app"
	"github.com/arsmn/ontest-server/module/cache"
	"github.com/arsmn/ontest-server/module/hash"
	"github.com/arsmn/ontest-server/module/mail"
	"github.com/arsmn/ontest-server/module/xlog"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/settings"
)

var _ app.App = new(Service)

type (
	serviceDependencies interface {
		xlog.Provider
		settings.Provider
		persistence.Provider
		cache.Provider
		hash.Provider
		mail.Provider
	}
	Service struct {
		dx serviceDependencies
	}
)

func NewApp(dx serviceDependencies) *Service {
	return &Service{dx}
}
