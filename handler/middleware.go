package handler

import (
	stderr "errors"
	"net/http"
	"strconv"

	"github.com/arsmn/ontest-server/module/context"
	"github.com/arsmn/ontest-server/module/errors"
	"github.com/arsmn/ontest-server/module/httplib"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/rs/cors"
)

func (h *Handler) cors(hh http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   h.dx.Settings().CORS.AllowedOrigins,
		AllowedMethods:   h.dx.Settings().CORS.AllowedMethods,
		AllowedHeaders:   h.dx.Settings().CORS.AllowedHeaders,
		AllowCredentials: h.dx.Settings().CORS.AllowCredentials,
		Debug:            !h.dx.Settings().IsProduction(),
	}).Handler(hh)
}

func (h *Handler) httpValues(hh http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.HTTPValuesContext(r.Context(), &context.HttpValues{
			IP:        httplib.FetchIP(r),
			UserAgent: r.UserAgent(),
		}))
		hh.ServeHTTP(rw, r)
	})
}

func (h *Handler) withUser(fn HandleFunc) HandleFunc {
	return func(ctx *Context) error {
		if err := h.fetchAuth(ctx); err != nil {
			return err
		}
		return fn(ctx)
	}
}

func (h *Handler) withAuth(fn HandleFunc) HandleFunc {
	return func(ctx *Context) error {
		if err := h.fetchAuth(ctx); err != nil {
			return err
		}

		if !ctx.Signed() {
			return errors.ErrUnauthorized
		}

		return fn(ctx)
	}
}

func (h *Handler) fetchAuth(ctx *Context) error {
	s := h.dx.Settings().Session
	c, err := ctx.Cookie(s.Cookie)
	if err != nil {
		return nil
	}

	sess, err := h.dx.App().GetSession(ctx.Request().Context(), c)
	if err != nil {
		if stderr.Is(err, persistence.ErrNoRows) {
			return nil
		}
		return err
	}

	if !sess.IsActive() {
		return nil
	}

	u, err := h.dx.App().GetUser(ctx.Request().Context(), sess.UserID)
	if err != nil {
		if stderr.Is(err, persistence.ErrNoRows) {
			return nil
		}
		return err
	}

	ctx.WithUser(u).WithSession(sess)

	return nil
}

func (h *Handler) withExam(fn HandleFunc) HandleFunc {
	return func(ctx *Context) error {
		id := ctx.Param("id")
		if len(id) == 0 {
			return fn(ctx)
		}

		eid, err := strconv.ParseUint(id, 10, 0)
		if err != nil {
			return err
		}

		exam, err := h.dx.App().GetExam(ctx.Context(), eid)
		if err != nil {
			if stderr.Is(err, persistence.ErrNoRows) {
				return nil
			}
			return err
		}

		ctx.WithExam(exam)

		return fn(ctx)
	}
}

func (h *Handler) withExamOwner(fn HandleFunc) HandleFunc {
	return func(ctx *Context) error {
		if !ctx.Signed() {
			return errors.ErrUnauthorized
		}

		if ctx.Exam() == nil {
			return errors.ErrForbidden
		}

		if ctx.Exam().Examiner != ctx.User().ID {
			return errors.ErrForbidden
		}

		return fn(ctx)
	}
}

func (h *Handler) withQuestion(fn HandleFunc) HandleFunc {
	return func(ctx *Context) error {
		id := ctx.Param("qid")
		if len(id) == 0 {
			return fn(ctx)
		}

		qid, err := strconv.ParseUint(id, 10, 0)
		if err != nil {
			return err
		}

		question, err := h.dx.App().GetQuestion(ctx.Context(), qid)
		if err != nil {
			if stderr.Is(err, persistence.ErrNoRows) {
				return nil
			}
			return err
		}

		ctx.WithQuestion(question)

		return fn(ctx)
	}
}

func (h *Handler) withQuestionOwner(fn HandleFunc) HandleFunc {
	return func(ctx *Context) error {
		if !ctx.Signed() {
			return errors.ErrUnauthorized
		}

		if ctx.Question() == nil {
			return errors.ErrForbidden
		}

		if ctx.Question().Examiner != ctx.User().ID {
			return errors.ErrForbidden
		}

		return fn(ctx)
	}
}
