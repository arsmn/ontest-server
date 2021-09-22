package handler

import (
	"github.com/arsmn/ontest-server/user"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) accountRouter(r chi.Router) {
	r.Post("/change-password", h.clown(h.changePassword, h.withUser))
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
