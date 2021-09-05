package handler

import (
	stderr "errors"
	"fmt"
	"net/http"

	"github.com/arsmn/ontest/module/context"
	"github.com/arsmn/ontest/module/errors"
)

type HTTPError struct {
	// The status code
	//
	// example: 404
	CodeField int `json:"code,omitempty"`

	// The status description
	//
	// example: Not Found
	StatusField string `json:"status,omitempty"`

	// The request ID
	//
	// The request ID is often exposed internally in order to trace
	// errors across service architectures. This is often a UUID.
	//
	// example: d7ef54b1-ec15-46e6-bccb-524b82c035e6
	RIDField string `json:"request,omitempty"`

	// A human-readable reason for the error
	//
	// example: User with ID 1234 does not exist.
	ReasonField string `json:"reason,omitempty"`

	// Debug information
	//
	// This field is often not exposed to protect against leaking
	// sensitive information.
	//
	// example: SQL field "foo" is not a bool.
	DebugField string `json:"debug,omitempty"`

	// Further error details
	DetailsField map[string]interface{} `json:"details,omitempty"`

	// Error message
	//
	// The error's message.
	//
	// example: The resource could not be found
	// required: true
	ErrorField string `json:"message"`

	err error
}

func (e HTTPError) Unwrap() error {
	return e.err
}

func (e *HTTPError) Wrap(err error) {
	e.err = err
}

func (e HTTPError) Status() string {
	return e.StatusField
}

func (e HTTPError) Error() string {
	return e.ErrorField
}

func (e HTTPError) RequestID() string {
	return e.RIDField
}

func (e HTTPError) Reason() string {
	return e.ReasonField
}

func (e HTTPError) Debug() string {
	return e.DebugField
}

func (e HTTPError) Details() map[string]interface{} {
	return e.DetailsField
}

func (e HTTPError) StatusCode() int {
	return e.CodeField
}

func (e HTTPError) WithReason(reason string) *HTTPError {
	e.ReasonField = reason
	return &e
}

func (e HTTPError) WithReasonf(reason string, args ...interface{}) *HTTPError {
	return e.WithReason(fmt.Sprintf(reason, args...))
}

func (e HTTPError) WithError(message string) *HTTPError {
	e.ErrorField = message
	return &e
}

func (e HTTPError) WithErrorf(message string, args ...interface{}) *HTTPError {
	return e.WithError(fmt.Sprintf(message, args...))
}

func (e HTTPError) WithDebugf(debug string, args ...interface{}) *HTTPError {
	return e.WithDebug(fmt.Sprintf(debug, args...))
}

func (e HTTPError) WithDebug(debug string) *HTTPError {
	e.DebugField = debug
	return &e
}

func (e HTTPError) WithDetail(key string, detail interface{}) *HTTPError {
	if e.DetailsField == nil {
		e.DetailsField = map[string]interface{}{}
	}
	e.DetailsField[key] = detail
	return &e
}

func (e HTTPError) WithDetailf(key string, message string, args ...interface{}) *HTTPError {
	if e.DetailsField == nil {
		e.DetailsField = map[string]interface{}{}
	}
	e.DetailsField[key] = fmt.Sprintf(message, args...)
	return &e
}

func ToHTTPError(err error, id string) *HTTPError {
	de := &HTTPError{
		RIDField:     id,
		CodeField:    http.StatusInternalServerError,
		DetailsField: map[string]interface{}{},
		ErrorField:   err.Error(),
	}
	de.Wrap(err)

	if c := errors.ReasonCarrier(nil); stderr.As(err, &c) {
		de.ReasonField = c.Reason()
	}
	if c := errors.RequestIDCarrier(nil); stderr.As(err, &c) && c.RequestID() != "" {
		de.RIDField = c.RequestID()
	}
	if c := errors.DetailsCarrier(nil); stderr.As(err, &c) && c.Details() != nil {
		de.DetailsField = c.Details()
	}
	if c := errors.StatusCarrier(nil); stderr.As(err, &c) && c.Status() != "" {
		de.StatusField = c.Status()
	}
	if c := errors.StatusCodeCarrier(nil); stderr.As(err, &c) && c.StatusCode() != 0 {
		de.CodeField = c.StatusCode()
	}
	if c := errors.DebugCarrier(nil); stderr.As(err, &c) {
		de.DebugField = c.Debug()
	}

	if de.StatusField == "" {
		de.StatusField = http.StatusText(de.StatusCode())
	}

	return de
}

type jsonError struct {
	Error *HTTPError `json:"error"`
}

func enhanceError(r *http.Request, err error) interface{} {
	return &jsonError{Error: ToHTTPError(err, r.Header.Get("X-Request-ID"))}
}

func handleError(ctx *context.Context, err error) error {
	status := http.StatusInternalServerError
	payload := enhanceError(ctx.Request(), err)

	if c := errors.StatusCodeCarrier(nil); stderr.As(err, &c) {
		status = c.StatusCode()
	}

	return ctx.Status(status).JSON(payload)
}
