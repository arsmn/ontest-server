package handler

import (
	"github.com/arsmn/ontest-server/module/context"
	"github.com/arsmn/ontest-server/module/errors"
)

type jsonError struct {
	Error *errors.Error `json:"error"`
}

func handleError(ctx *context.Context, err error) error {
	reqId := ctx.Request().Header.Get("X-Request-ID")
	defErr := errors.ToError(err, reqId)
	payload := &jsonError{Error: defErr}
	return ctx.Json(defErr.CodeField, payload)
}
