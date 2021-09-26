package handler

import (
	stderr "errors"
	"net/http"

	"github.com/arsmn/ontest-server/module/errors"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/rs/cors"
)

func (h *Handler) cors(hh http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   h.dx.Settings().CORS.AllowedOrigins,
		AllowedMethods:   h.dx.Settings().CORS.AllowedMethods,
		AllowedHeaders:   h.dx.Settings().CORS.AllowedHeaders,
		AllowCredentials: h.dx.Settings().CORS.AllowCredentials,
		Debug:            !h.dx.Settings().IsProduction(),
	}).Handler(hh)
}

func (h *Handler) withUser(fn HandleFunc) HandleFunc {
	return func(ctx *Context) error {
		if err := h.fetchAuth(ctx); err != nil {
			return err
		}
		return fn(ctx)
	}
}

func (h *Handler) withAuth(fn HandleFunc) HandleFunc {
	return func(ctx *Context) error {
		if err := h.fetchAuth(ctx); err != nil {
			return err
		}

		if !ctx.Signed() {
			return errors.ErrUnauthorized
		}

		return fn(ctx)
	}
}

func (h *Handler) fetchAuth(ctx *Context) error {
	s := h.dx.Settings().Session
	c, err := ctx.Cookie(s.Cookie)
	if err != nil {
		return nil
	}

	sess, err := h.dx.App().GetSession(ctx.Request().Context(), c)
	if err != nil {
		if stderr.Is(err, persistence.ErrNoRows) {
			return nil
		}
		return err
	}

	if !sess.IsActive() {
		return nil
	}

	u, err := h.dx.App().GetUser(ctx.Request().Context(), sess.UserID)
	if err != nil {
		if stderr.Is(err, persistence.ErrNoRows) {
			return nil
		}
		return err
	}

	ctx.WithUser(u).WithSession(sess)

	return nil
}
