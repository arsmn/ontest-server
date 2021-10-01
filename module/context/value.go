package context

import (
	"context"

	"github.com/arsmn/ontest-server/exam"
	"github.com/arsmn/ontest-server/question"
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

var questionKey = &contextKey{"ctx-exam"}

func Question(ctx context.Context) *question.Question {
	e, _ := ctx.Value(questionKey).(*question.Question)
	return e
}

type HttpValues struct {
	IP        string
	UserAgent string
}

var httpValuesKey = &contextKey{"ctx-http-values"}

func HTTPValuesContext(ctx context.Context, hv *HttpValues) context.Context {
	return context.WithValue(ctx, httpValuesKey, hv)
}

func HTTPValues(ctx context.Context) *HttpValues {
	hv, _ := ctx.Value(httpValuesKey).(*HttpValues)
	return hv
}
