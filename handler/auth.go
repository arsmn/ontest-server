package handler

import (
	t "github.com/arsmn/ontest-server/transport"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) authRouter(r chi.Router) {
	r.Post("/signin", h.clown(h.signin))
	r.Post("/signup", h.clown(h.signup))
	r.Post("/signout", h.clown(h.withUser(h.signout)))
	r.Get("/whoami", h.clown(h.withUser(h.whoami)))
	r.Post("/forgot-password", h.clown(h.forgotPassword))
	r.Post("/reset-password", h.clown(h.resettPassword))
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
	ctx.SetSecureCookie(s.Cookie, res.Token, int(s.Lifespan.Seconds()), "/", h.dx.Settings().Domain())

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

func (h *Handler) signout(ctx *Context) error {
	c := h.dx.Settings().Session().Cookie
	token, err := ctx.Cookie(c)
	if err != nil {
		return err
	}

	err = h.dx.App().DeleteSession(ctx.Request().Context(), token)
	if err != nil {
		return err
	}

	ctx.RemoveCookie(h.dx.Settings().Session().Cookie)
	return ctx.OK(success)
}

func (h *Handler) whoami(ctx *Context) error {
	u := ctx.User().CopySanitize()
	return ctx.OK(payload(u))
}

func (h *Handler) forgotPassword(ctx *Context) error {
	req := new(t.ForgotPasswordRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	err := h.dx.App().ForgotPassword(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(success)
}

func (h *Handler) resettPassword(ctx *Context) error {
	req := new(t.ResetPasswordRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	err := h.dx.App().ResetPassword(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(success)
}
