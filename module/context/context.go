package context

import (
	"context"
	"net/http"

	"github.com/arsmn/ontest-server/session"
	"github.com/arsmn/ontest-server/user"
	"github.com/go-chi/chi/v5"
)

type Context struct {
	request  *http.Request
	response http.ResponseWriter

	user *user.User
	sess *session.Session
}

func (ctx *Context) WithUser(u *user.User) *Context {
	ctx.user = u
	return ctx
}

func (ctx *Context) WithSession(s *session.Session) *Context {
	ctx.sess = s
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

func (ctx *Context) Token() string {
	return ctx.sess.Token
}
