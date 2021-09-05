package driver

import (
	"sync"

	"github.com/arsmn/ontest/app"
	"github.com/arsmn/ontest/persistence"
	"github.com/arsmn/ontest/settings"
)

type RegistryCore struct {
	rwl sync.RWMutex

	app app.App

	persister persistence.Persister

	settings settings.Settings
}
