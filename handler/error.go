package handler

import (
	"github.com/arsmn/ontest/module/context"
	"github.com/arsmn/ontest/module/errors"
)

type jsonError struct {
	Error *errors.DefaultError `json:"error"`
}

func handleError(ctx *context.Context, err error) error {
	reqId := ctx.Request().Header.Get("X-Request-ID")
	defErr := errors.ToDefaultError(err, reqId)
	payload := &jsonError{Error: defErr}
	return ctx.Json(defErr.CodeField, payload)
}
