package handler

import (
	stderr "errors"

	"github.com/arsmn/ontest-server/module/errors"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/user"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) accountRouter(r chi.Router) {
	r.Get("/check-username/{username}", h.clown(h.checkUsername))
	r.Post("/change-password", h.clown(h.changePassword, h.withUser))
}

func (h *Handler) checkUsername(ctx *Context) error {
	_, err := h.dx.App().GetUserByUsername(ctx.Request().Context(), ctx.Param("username"))
	if err != nil {
		if stderr.Is(err, persistence.ErrNoRows) {
			return ctx.OK(success)
		}
		return err
	}

	return errors.ErrConflict.WithError("Username is taken")
}

func (h *Handler) changePassword(ctx *Context) error {
	req := new(user.ChangePasswordRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	req.WithUser(ctx.User())
	err := h.dx.App().ChangePassword(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(success)
}
