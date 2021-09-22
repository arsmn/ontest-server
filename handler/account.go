package handler

import (
	stderr "errors"

	"github.com/arsmn/ontest-server/module/errors"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/user"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) accountRouter(r chi.Router) {
	r.Put("/", h.clown(h.updateProfile, h.withUser))
	r.Get("/whoami", h.clown(h.whoami, h.withUser))
	r.Post("/change-password", h.clown(h.changePassword, h.withUser))
	r.Get("/check/{type}/{value}", h.clown(h.check))
	r.Post("/send/verification", h.clown(h.sendVerification, h.withUser))
	r.Post("/verify", h.clown(h.verify, h.withUser))
}

func (h *Handler) updateProfile(ctx *Context) error {
	req := new(user.UpdateProfileRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	req.WithUser(ctx.User())
	err := h.dx.App().UpdateProfile(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(success)
}

func (h *Handler) whoami(ctx *Context) error {
	return ctx.OK(payload(ctx.User()))
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

func (h *Handler) check(ctx *Context) error {
	var err error
	typ := ctx.Param("type")
	value := ctx.Param("value")

	switch typ {
	case "email":
		_, err = h.dx.App().GetUserByUsername(ctx.Request().Context(), value)
	case "username":
		_, err = h.dx.App().GetUserByUsername(ctx.Request().Context(), value)
	default:
		return errors.ErrBadRequest.WithError("invalid type")
	}

	if err != nil {
		if stderr.Is(err, persistence.ErrNoRows) {
			return ctx.OK(success)
		}
		return err
	}

	return errors.ErrConflict.WithError(typ + " is taken")
}

func (h *Handler) sendVerification(ctx *Context) error {
	req := new(user.SendVerificationRequest)
	req.WithUser(ctx.User())
	req.Email = ctx.User().Email

	err := h.dx.App().SendVerification(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(success)
}

func (h *Handler) verify(ctx *Context) error {
	req := new(user.VerificationRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	req.WithUser(ctx.User())

	err := h.dx.App().Verify(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(success)
}
