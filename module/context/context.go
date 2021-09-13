package context

import (
	"net/http"
)

type Context struct {
	request    *http.Request
	response   http.ResponseWriter
	statusCode int
}

func (ctx *Context) Request() *http.Request {
	return ctx.request
}

func (ctx *Context) Response() http.ResponseWriter {
	return ctx.response
}
