package handler

import (
	"strconv"

	"github.com/arsmn/ontest-server/exam"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) resultRouter(r chi.Router) {
	r.Get("/{id}", h.clown(h.getResult, h.withAuth))
	r.Post("/{aid}/submit", h.clown(h.submitAnswer, h.withAuth))
}

func (h *Handler) getResult(ctx *Context) error {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 0)

	res, err := h.dx.App().GetResult(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.OK(payload(res))
}

func (h *Handler) submitAnswer(ctx *Context) error {
	req := new(exam.SubmitAnswerRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	id, _ := strconv.ParseUint(ctx.Param("aid"), 10, 0)
	req.ID = id

	err := h.dx.App().SubmitAnswer(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(success)
}
