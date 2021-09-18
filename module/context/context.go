package context

import (
	"net/http"

	"github.com/arsmn/ontest-server/user"
)

type Context struct {
	request  *http.Request
	response http.ResponseWriter

	user *user.User
}

func (ctx *Context) WithUser(u *user.User) *Context {
	ctx.user = u
	return ctx
}

func (ctx *Context) Request() *http.Request {
	return ctx.request
}

func (ctx *Context) Response() http.ResponseWriter {
	return ctx.response
}

func (ctx *Context) User() *user.User {
	return ctx.user
}
