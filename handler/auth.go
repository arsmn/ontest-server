package handler

import (
	"net/http"

	t "github.com/arsmn/ontest/transport"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) authRouter(r chi.Router) {
	r.Post("/signin", h.clown(h.signin))
	r.Post("/signup", h.clown(h.signup))
}

func (h *Handler) signin(ctx *Context) error {
	req := new(t.SigninRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	res, err := h.dx.App().IssueSession(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	s := h.dx.Settings().Session()
	ctx.SetSecureCookie(s.Cookie, res.Token, int(s.Lifespan.Seconds()))

	return ctx.SendStatus(http.StatusOK)
}

func (h *Handler) signup(ctx *Context) error {
	req := new(t.SignupRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	_, err := h.dx.App().RegisterUser(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusCreated)
}
