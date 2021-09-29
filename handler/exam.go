package handler

import (
	"fmt"

	"github.com/arsmn/ontest-server/exam"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) examRouter(r chi.Router) {
	r.Post("/", h.clown(h.createExam, h.withAuth))
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

	return ctx.Created(fmt.Sprintf("/exams/%d", exam.ID), payload(exam))
}
