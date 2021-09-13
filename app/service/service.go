package service

import (
	"github.com/arsmn/ontest/app"
	"github.com/arsmn/ontest/module/xlog"
	"github.com/arsmn/ontest/persistence"
	"github.com/arsmn/ontest/settings"
)

var _ app.App = new(Service)

type (
	serviceDependencies interface {
		xlog.Provider
		settings.Provider
		persistence.Provider
	}
	Service struct {
		dx serviceDependencies
	}
)

func NewApp(dx serviceDependencies) *Service {
	return &Service{dx}
}
