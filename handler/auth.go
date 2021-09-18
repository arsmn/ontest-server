package handler

import (
	t "github.com/arsmn/ontest-server/transport"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) authRouter(r chi.Router) {
	r.Post("/signin", h.clown(h.signin))
	r.Post("/signup", h.clown(h.signup))
	r.Get("/whoami", h.clown(h.withUser(h.whoami)))
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
	ctx.SetSecureCookie(s.Cookie, res.Token, int(s.Lifespan.Seconds()), "/", s.Domain)

	return ctx.OK(success)
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

	return ctx.Created("/auth/whoami", success)
}

func (h *Handler) whoami(ctx *Context) error {
	u := ctx.User().CopySanitize("Password")
	return ctx.OK(payload(u))
}
