package handler

import (
	"net/http"

	t "github.com/arsmn/ontest/transport"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) authRouter(r chi.Router) {
	r.Post("/signin", nil)
	r.Post("/signup", h.clown(h.signup))
}

func (h *Handler) signup(ctx *Context) error {
	req := new(t.SignupRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	_, err := h.dx.App().Signup(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusCreated)
}
