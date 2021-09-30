package context

import (
	"context"

	"github.com/arsmn/ontest-server/exam"
	"github.com/arsmn/ontest-server/session"
	"github.com/arsmn/ontest-server/user"
)

type contextKey struct {
	name string
}

var userKey = &contextKey{"ctx-user"}

func User(ctx context.Context) *user.User {
	u, _ := ctx.Value(userKey).(*user.User)
	return u
}

var sessionKey = &contextKey{"ctx-session"}

func Session(ctx context.Context) *session.Session {
	s, _ := ctx.Value(sessionKey).(*session.Session)
	return s
}

var examKey = &contextKey{"ctx-exam"}

func Exam(ctx context.Context) *exam.Exam {
	e, _ := ctx.Value(examKey).(*exam.Exam)
	return e
}
