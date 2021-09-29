package app

import (
	"github.com/arsmn/ontest-server/exam"
	"github.com/arsmn/ontest-server/session"
	"github.com/arsmn/ontest-server/user"
)

type Provider interface {
	App() App
}

type App interface {
	user.Service
	session.Service
	exam.Service
}
