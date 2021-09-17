package service

import (
	"github.com/arsmn/ontest-server/app"
	"github.com/arsmn/ontest-server/module/hash"
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
		hash.HashProvider
	}
	Service struct {
		dx serviceDependencies
	}
)

func NewApp(dx serviceDependencies) *Service {
	return &Service{dx}
}
