package handler

import (
	"bytes"
	"fmt"
	"image"
	"io"

	"github.com/arsmn/ontest-server/exam"
	"github.com/arsmn/ontest-server/module/avatar"
	"github.com/arsmn/ontest-server/module/storage"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) examRouter(r chi.Router) {
	r.Get("/{id}", h.clown(h.getExam, h.withAuth, h.withExam, h.withOwner))
	r.Post("/", h.clown(h.createExam, h.withAuth))
	r.Put("/{id}", h.clown(h.updateExam, h.withAuth, h.withExam, h.withOwner))
	r.Post("/{id}/cover", h.clown(h.setCover, h.withAuth, h.withExam, h.withOwner))
	r.Delete("/{id}/cover", h.clown(h.deleteCover, h.withAuth, h.withExam, h.withOwner))
}

func (h *Handler) getExam(ctx *Context) error {
	return ctx.OK(payload(ctx.Exam().Map()))
}

func (h *Handler) createExam(ctx *Context) error {
	req := new(exam.CreateExamRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	req.WithUser(ctx.User())
	exam, err := h.dx.App().CreateExam(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.Created(fmt.Sprintf("/exams/%d", exam.ID), payload(exam.Map()))
}

func (h *Handler) updateExam(ctx *Context) error {
	req := new(exam.UpdateExamRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	req.WithUser(ctx.User())
	req.WithExam(ctx.Exam())
	err := h.dx.App().UpdateExam(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(success)
}

func (h *Handler) setCover(ctx *Context) error {
	var img image.Image

	if err := ctx.Request().ParseMultipartForm(10 << 20); err != nil {
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

	img, err = avatar.PrepareAvatar(buf.Bytes(), 3000, 3000, 290)
	if err != nil {
		return err
	}

	_, err = storage.WriteImage(ctx.Exam().Fs(), "cover", img)
	if err != nil {
		return err
	}

	return ctx.OK(payload(ctx.Exam().Map()))
}

func (h *Handler) deleteCover(ctx *Context) error {
	err := ctx.Exam().Fs().Remove("cover")
	if err != nil {
		return err
	}

	return ctx.OK(success)
}
