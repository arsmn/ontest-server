package service

import (
	"github.com/arsmn/ontest/app"
	"github.com/arsmn/ontest/settings"
)

var _ app.App = new(Service)

type (
	serviceDependencies interface {
		settings.Provider
	}
	Service struct {
		r serviceDependencies
	}
)

func NewManager(r serviceDependencies) *Service {
	return &Service{r}
}
