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

func (h *Handler) signup(c *ctx) error {
	req := new(t.SignupRequest)
	if err := c.BindJson(req); err != nil {
		return err
	}

	_, err := h.dx.App().Signup(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusCreated)
}
