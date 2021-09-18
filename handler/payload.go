package handler

import (
	"github.com/arsmn/ontest-server/module/context"
	"github.com/arsmn/ontest-server/module/errors"
)

type jsonPayload struct {
	Success bool          `json:"success"`
	Data    interface{}   `json:"data,omitempty"`
	Error   *errors.Error `json:"error,omitempty"`
}

var success = &jsonPayload{Success: true}

func handleError(ctx *context.Context, err error) error {
	reqId := ctx.Request().Header.Get("X-Request-ID")
	defErr := errors.ToError(err, reqId)
	payload := &jsonPayload{Success: false, Error: defErr}
	return ctx.Json(defErr.CodeField, payload)
}

func payload(d interface{}) *jsonPayload {
	return &jsonPayload{
		Success: true,
		Data:    d,
	}
}
