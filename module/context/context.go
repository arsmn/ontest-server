package context

import (
	"context"
	"net/http"

	"github.com/arsmn/ontest-server/exam"
	"github.com/arsmn/ontest-server/session"
	"github.com/arsmn/ontest-server/user"
	"github.com/go-chi/chi/v5"
)

type Context struct {
	request  *http.Request
	response http.ResponseWriter

	user *user.User
	sess *session.Session

	exam *exam.Exam
}

func (ctx *Context) WithUser(u *user.User) *Context {
	ctx.user = u
	return ctx
}

func (ctx *Context) WithSession(s *session.Session) *Context {
	ctx.sess = s
	return ctx
}

func (ctx *Context) WithExam(e *exam.Exam) *Context {
	ctx.exam = e
	return ctx
}

func (ctx *Context) Request() *http.Request {
	return ctx.request
}

func (ctx *Context) Response() http.ResponseWriter {
	return ctx.response
}

func (ctx *Context) Context() context.Context {
	return ctx.request.Context()
}

func (ctx *Context) Param(p string) string {
	return chi.URLParam(ctx.request, p)
}

func (ctx *Context) User() *user.User {
	return ctx.user
}

func (ctx *Context) Session() *session.Session {
	return ctx.sess
}

func (ctx *Context) Exam() *exam.Exam {
	return ctx.exam
}

func (ctx *Context) Token() string {
	return ctx.sess.Token
}

func (ctx *Context) Signed() bool {
	return ctx.user != nil && ctx.sess != nil
}

func (ctx *Context) IP() string {
	forwarded := ctx.request.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return ctx.request.RemoteAddr
}
