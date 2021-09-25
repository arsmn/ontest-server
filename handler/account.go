package handler

import (
	"bytes"
	stderr "errors"
	"image"
	"io"
	"strconv"

	"github.com/arsmn/ontest-server/module/avatar"
	"github.com/arsmn/ontest-server/module/errors"
	"github.com/arsmn/ontest-server/module/storage"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/user"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) accountRouter(r chi.Router) {
	r.Put("/", h.clown(h.updateProfile, h.withUser))
	r.Get("/whoami", h.clown(h.whoami, h.withUser))
	r.Post("/change-password", h.clown(h.changePassword, h.withUser))
	r.Post("/set-password", h.clown(h.setPassword, h.withUser))
	r.Get("/check/{type}/{value}", h.clown(h.check))
	r.Post("/send/verification", h.clown(h.sendVerification, h.withUser))
	r.Post("/verify", h.clown(h.verify, h.withUser))
	r.Post("/avatar", h.clown(h.setAvatar, h.withUser))
	r.Delete("/avatar", h.clown(h.deleteAvatar, h.withUser))
	r.Post("/set-preference", h.clown(h.setPreference, h.withUser))
}

func (h *Handler) updateProfile(ctx *Context) error {
	req := new(user.UpdateProfileRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	req.WithUser(ctx.User())
	err := h.dx.App().UpdateProfile(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(payload(ctx.User().Map()))
}

func (h *Handler) whoami(ctx *Context) error {
	return ctx.OK(payload(ctx.User().Map()))
}

func (h *Handler) changePassword(ctx *Context) error {
	req := new(user.ChangePasswordRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	req.WithUser(ctx.User()).WithToken(ctx.Token())
	err := h.dx.App().ChangePassword(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(success)
}

func (h *Handler) setPassword(ctx *Context) error {
	req := new(user.SetPasswordRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	req.WithUser(ctx.User()).WithToken(ctx.Token())
	err := h.dx.App().SetPassword(ctx.Context(), req)
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
		_, err = h.dx.App().GetUserByEmail(ctx.Context(), value)
	case "username":
		_, err = h.dx.App().GetUserByUsername(ctx.Context(), value)
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

	err := h.dx.App().SendVerification(ctx.Context(), req)
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

	err := h.dx.App().Verify(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(success)
}

func (h *Handler) setAvatar(ctx *Context) (err error) {
	var img image.Image
	gen := ctx.Request().URL.Query().Get("gen")

	if gen == "true" {
		seed := ctx.User().Email + strconv.FormatUint(ctx.User().ID, 10)
		img, err = avatar.GenerateRandom([]byte(seed))
		if err != nil {
			return err
		}
	} else {
		if err = ctx.Request().ParseMultipartForm(10 << 20); err != nil {
			return err
		}

		file, _, err := ctx.Request().FormFile("file")
		if err != nil {
			return err
		}
		defer file.Close()

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, file); err != nil {
			return err
		}

		img, err = avatar.PrepareAvatar(buf.Bytes(), 2000, 2000, 290)
		if err != nil {
			return err
		}
	}

	_, err = storage.WriteImage(ctx.User().Fs(), "avatar", img)
	if err != nil {
		return err
	}

	return ctx.OK(payload(ctx.User().Map()))
}

func (h *Handler) deleteAvatar(ctx *Context) error {
	err := ctx.User().Fs().Remove("avatar")
	if err != nil {
		return err
	}

	return ctx.OK(payload(ctx.User().Map()))
}

func (h *Handler) setPreference(ctx *Context) error {
	req := new(user.SetPreferenceRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	req.WithUser(ctx.User())
	err := h.dx.App().SetPreference(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(success)
}
