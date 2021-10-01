package context

import (
	"context"
	"net/http"

	"github.com/arsmn/ontest-server/exam"
	"github.com/arsmn/ontest-server/module/httplib"
	"github.com/arsmn/ontest-server/question"
	"github.com/arsmn/ontest-server/session"
	"github.com/arsmn/ontest-server/user"
	"github.com/go-chi/chi/v5"
)

type Context struct {
	request  *http.Request
	response http.ResponseWriter

	user     *user.User
	sess     *session.Session
	exam     *exam.Exam
	question *question.Question
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

func (ctx *Context) IP() string {
	return httplib.FetchIP(ctx.request)
}

func (ctx *Context) Signed() bool {
	return ctx.user != nil && ctx.sess != nil
}

func (ctx *Context) WithUser(u *user.User) *Context {
	ctx.user = u
	c := context.WithValue(ctx.request.Context(), userKey, u)
	ctx.request = ctx.request.WithContext(c)
	return ctx
}

func (ctx *Context) User() *user.User {
	return ctx.user
}

func (ctx *Context) WithSession(s *session.Session) *Context {
	ctx.sess = s
	c := context.WithValue(ctx.request.Context(), sessionKey, s)
	ctx.request = ctx.request.WithContext(c)
	return ctx
}

func (ctx *Context) Session() *session.Session {
	return ctx.sess
}

func (ctx *Context) WithExam(e *exam.Exam) *Context {
	ctx.exam = e
	c := context.WithValue(ctx.request.Context(), examKey, e)
	ctx.request = ctx.request.WithContext(c)
	return ctx
}

func (ctx *Context) Exam() *exam.Exam {
	return ctx.exam
}

func (ctx *Context) WithQuestion(q *question.Question) *Context {
	ctx.question = q
	c := context.WithValue(ctx.request.Context(), questionKey, q)
	ctx.request = ctx.request.WithContext(c)
	return ctx
}

func (ctx *Context) Question() *question.Question {
	return ctx.question
}
