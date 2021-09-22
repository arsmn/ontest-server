package handler

import (
	stderr "errors"

	"github.com/arsmn/ontest-server/module/errors"
	"github.com/arsmn/ontest-server/persistence"
)

func (h *Handler) withUser(fn HandleFunc) HandleFunc {
	return func(ctx *Context) error {
		s := h.dx.Settings().Session
		c, err := ctx.Cookie(s.Cookie)
		if err != nil {
			ctx.RemoveCookie(s.Cookie)
			return errors.ErrUnauthorized
		}

		sess, err := h.dx.App().GetSession(ctx.Request().Context(), c)
		if err != nil {
			if stderr.Is(err, persistence.ErrNoRows) {
				ctx.RemoveCookie(s.Cookie)
				return errors.ErrUnauthorized
			}
			return err
		}

		if !sess.IsActive() {
			return errors.ErrUnauthorized
		}

		u, err := h.dx.App().GetUser(ctx.Request().Context(), sess.UserID)
		if err != nil {
			if stderr.Is(err, persistence.ErrNoRows) {
				return errors.ErrUnauthorized
			}
			return err
		}

		ctx.WithUser(u).WithSession(sess)

		return fn(ctx)
	}
}
