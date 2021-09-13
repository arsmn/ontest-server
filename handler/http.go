package handler

import (
	"net/http"

	"github.com/arsmn/ontest/app"
	"github.com/arsmn/ontest/module/context"
	"github.com/arsmn/ontest/settings"

	"github.com/go-chi/chi/v5"
)

type (
	handleFunc func(ctx *context.Context) error

	dependencies interface {
		app.Provider
		settings.Provider
	}

	Handler struct {
		dx      dependencies
		handler http.Handler
		count   uint32
	}
)

func New(dx dependencies) *Handler {
	api := &Handler{}

	api.dx = dx

	root := chi.NewRouter()
	root.Route("/auth", api.authRouter)

	api.handler = root

	return api
}

func (api *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	api.handler.ServeHTTP(rw, r)
}

func (a *Handler) String() string {
	return ""
}

func (a *Handler) clown(fn handleFunc) http.HandlerFunc {
	a.count++

	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := context.Acquire(rw, r)
		if err := fn(ctx); err != nil {
			if catch := handleError(ctx, err); catch != nil {
				_ = ctx.SendStatus(http.StatusInternalServerError)
			}
		}
	})

	return handler
}
