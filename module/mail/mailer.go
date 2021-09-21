package mail

import (
	"context"

	"github.com/arsmn/ontest-server/user"
)

type Mailer interface {
	SendResetPassword(context.Context, *user.User, string)
}

type Provider interface {
	Mailer() Mailer
}
