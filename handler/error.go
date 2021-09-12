package handler

import (
	stderr "errors"
	"net/http"

	"github.com/arsmn/ontest/module/context"
	"github.com/arsmn/ontest/module/errors"
)

type jsonError struct {
	Error *errors.DefaultError `json:"error"`
}

func enhanceError(r *http.Request, err error) interface{} {
	return &jsonError{Error: errors.ToDefaultError(err, r.Header.Get("X-Request-ID"))}
}

func handleError(ctx *context.Context, err error) error {
	status := http.StatusInternalServerError
	payload := enhanceError(ctx.Request(), err)

	if c := errors.StatusCodeCarrier(nil); stderr.As(err, &c) {
		status = c.StatusCode()
	}

	return ctx.Status(status).JSON(payload)
}
