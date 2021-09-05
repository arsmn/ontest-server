package driver

import (
	"github.com/arsmn/ontest/app"
	"github.com/arsmn/ontest/persistence"
	"github.com/arsmn/ontest/settings"
)

type Registry interface {
	app.Provider

	persistence.Provider

	settings.Provider
}
