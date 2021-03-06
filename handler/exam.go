package handler

import (
	"bytes"
	"fmt"
	"image"
	"io"

	"github.com/arsmn/ontest-server/exam"
	"github.com/arsmn/ontest-server/module/avatar"
	"github.com/arsmn/ontest-server/module/storage"
	"github.com/arsmn/ontest-server/question"
	"github.com/arsmn/ontest-server/shared"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) examRouter(r chi.Router) {
	r.Get("/{id}", h.clown(h.getExam, h.withAuth, h.withExam))
	r.Get("/{id}/stats", h.clown(h.getExamStats, h.withAuth, h.withExam, h.withExamOwner))
	r.Post("/{id}/publish", h.clown(h.publishExam, h.withAuth, h.withExam, h.withExamOwner))
	r.Post("/", h.clown(h.createExam, h.withAuth))
	r.Put("/{id}", h.clown(h.updateExam, h.withAuth, h.withExam, h.withExamOwner))
	r.Post("/{id}/cover", h.clown(h.setCover, h.withAuth, h.withExam, h.withExamOwner))
	r.Delete("/{id}/cover", h.clown(h.deleteCover, h.withAuth, h.withExam, h.withExamOwner))
	r.Get("/{id}/questions", h.clown(h.getQuestions, h.withAuth, h.withExam, h.withExamOwner))
	r.Post("/{id}/question", h.clown(h.createQuestion, h.withAuth, h.withExam, h.withExamOwner))
	r.Put("/{id}/question/{qid}", h.clown(h.updateQuestion, h.withAuth, h.withExam, h.withExamOwner, h.withQuestion, h.withQuestionOwner))
	r.Delete("/{id}/question/{qid}", h.clown(h.deleteQuestion, h.withAuth, h.withExam, h.withExamOwner, h.withQuestion, h.withQuestionOwner))
	r.Get("/{id}/question/{qid}/options", h.clown(h.getOptions, h.withAuth, h.withExam, h.withExamOwner, h.withQuestion, h.withQuestionOwner))
	r.Get("/search", h.clown(h.searchExam, h.withAuth))
	r.Post("/{id}/participate", h.clown(h.participateExam, h.withAuth, h.withExam))
}

func (h *Handler) getExam(ctx *Context) error {
	return ctx.OK(payload(ctx.Exam().Map()))
}

func (h *Handler) getExamStats(ctx *Context) error {
	res, err := h.dx.App().GetExamStats(ctx.Context(), ctx.Exam().ID)
	if err != nil {
		return err
	}

	return ctx.OK(payload(res))
}

func (h *Handler) createExam(ctx *Context) error {
	req := new(exam.CreateExamRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

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

func (h *Handler) getQuestions(ctx *Context) error {
	req := new(question.GetQuestionListRequest)
	req.ExamID = ctx.Exam().ID
	req.PaginatedRequest = shared.NewPaginatedRequest(ctx.Request())

	res, err := h.dx.App().GetQuestionList(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(paginatedPayload(res.Questions, res.PaginatedResponse))
}

func (h *Handler) createQuestion(ctx *Context) error {
	req := new(question.CreateQuestionRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	res, err := h.dx.App().CreateQuestion(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.Created("", payload(res))
}

func (h *Handler) updateQuestion(ctx *Context) error {
	req := new(question.CreateQuestionRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	err := h.dx.App().UpdateQuestion(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(success)
}

func (h *Handler) deleteQuestion(ctx *Context) error {
	err := h.dx.App().DeleteQuestion(ctx.Context(), ctx.Question().ID)
	if err != nil {
		return err
	}

	return ctx.OK(success)
}

func (h *Handler) getOptions(ctx *Context) error {
	res, err := h.dx.App().GetQuestionOptions(ctx.Context(), ctx.Question().ID)
	if err != nil {
		return err
	}

	return ctx.OK(payload(res))
}

func (h *Handler) searchExam(ctx *Context) error {
	req := new(exam.SearchExamRequest)
	req.PaginatedRequest = shared.NewPaginatedRequest(ctx.Request())
	res, err := h.dx.App().SearchExam(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.OK(paginatedPayload(res.Exams, res.PaginatedResponse))
}

func (h *Handler) publishExam(ctx *Context) error {
	err := h.dx.App().PublishExam(ctx.Context(), ctx.Exam())
	if err != nil {
		return err
	}

	return ctx.OK(success)
}

func (h *Handler) participateExam(ctx *Context) error {
	res, err := h.dx.App().Participate(ctx.Context(), ctx.User(), ctx.Exam())
	if err != nil {
		return err
	}

	return ctx.OK(payload(res))
}
