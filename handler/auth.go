package handler

import (
	"github.com/arsmn/ontest-server/session"
	"github.com/arsmn/ontest-server/user"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) authRouter(r chi.Router) {
	r.Post("/signin", h.clown(h.signin))
	r.Post("/signup", h.clown(h.signup))
	r.Post("/signout", h.clown(h.signout, h.withAuth))
	r.Post("/send/reset-password", h.clown(h.sendResetPassword))
	r.Post("/reset-password", h.clown(h.resetPassword))
}

func (h *Handler) signin(ctx *Context) error {
	req := new(session.SigninRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	req.IP = ctx.IP()
	req.UserAgent = ctx.Request().UserAgent()

	res, err := h.dx.App().IssueSession(ctx.Context(), req)
	if err != nil {
		return err
	}

	age := 0
	s := h.dx.Settings().Session
	if req.Remember {
		age = int(s.Lifespan.Seconds())
	}

	ctx.SetSecureCookie(s.Cookie, res.Token, age, "/", h.dx.Settings().Serve.Domain)

	return ctx.OK(success)
}

func (h *Handler) signup(ctx *Context) error {
	req := new(user.SignupRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	_, err := h.dx.App().RegisterUser(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.Created("/auth/whoami", success)
}

func (h *Handler) signout(ctx *Context) error {
	req := new(session.DeleteSessionByTokenRequest)

	req.Token = ctx.Token()
	req.WithUser(ctx.User()).WithToken(ctx.Token())
	err := h.dx.App().DeleteSessionByToken(ctx.Context(), req)
	if err != nil {
		return err
	}

	ctx.RemoveCookie(h.dx.Settings().Session.Cookie)

	return ctx.OK(success)
}

func (h *Handler) sendResetPassword(ctx *Context) error {
	req := new(user.SendResetPasswordRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	err := h.dx.App().SendResetPassword(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(success)
}

func (h *Handler) resetPassword(ctx *Context) error {
	req := new(user.ResetPasswordRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	err := h.dx.App().ResetPassword(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(success)
}
