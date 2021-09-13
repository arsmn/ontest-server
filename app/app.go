package app

import (
	"github.com/arsmn/ontest/session"
	"github.com/arsmn/ontest/user"
)

type Provider interface {
	App() App
}

type App interface {
	user.Service
	session.Service
}
