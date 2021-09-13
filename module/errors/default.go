package errors

import (
	stderr "errors"
	"fmt"
	"net/http"
)

type DefaultError struct {
	CodeField    int                    `json:"code,omitempty"`
	StatusField  string                 `json:"status,omitempty"`
	RIDField     string                 `json:"request,omitempty"`
	ReasonField  string                 `json:"reason,omitempty"`
	DebugField   string                 `json:"debug,omitempty"`
	DetailsField map[string]interface{} `json:"details,omitempty"`
	ErrorField   string                 `json:"message"`

	err error
}

func (e DefaultError) Unwrap() error {
	return e.err
}

func (e *DefaultError) Wrap(err error) {
	e.err = err
}

func (e DefaultError) Status() string {
	return e.StatusField
}

func (e DefaultError) Error() string {
	return e.ErrorField
}

func (e DefaultError) RequestID() string {
	return e.RIDField
}

func (e DefaultError) Reason() string {
	return e.ReasonField
}

func (e DefaultError) Debug() string {
	return e.DebugField
}

func (e DefaultError) Details() map[string]interface{} {
	return e.DetailsField
}

func (e DefaultError) StatusCode() int {
	return e.CodeField
}

func (e DefaultError) WithReason(reason string) *DefaultError {
	e.ReasonField = reason
	return &e
}

func (e DefaultError) WithReasonf(reason string, args ...interface{}) *DefaultError {
	return e.WithReason(fmt.Sprintf(reason, args...))
}

func (e DefaultError) WithError(message string) *DefaultError {
	e.ErrorField = message
	return &e
}

func (e DefaultError) WithErrorf(message string, args ...interface{}) *DefaultError {
	return e.WithError(fmt.Sprintf(message, args...))
}

func (e DefaultError) WithDebugf(debug string, args ...interface{}) *DefaultError {
	return e.WithDebug(fmt.Sprintf(debug, args...))
}

func (e DefaultError) WithDebug(debug string) *DefaultError {
	e.DebugField = debug
	return &e
}

func (e DefaultError) WithDetail(key string, detail interface{}) *DefaultError {
	if e.DetailsField == nil {
		e.DetailsField = map[string]interface{}{}
	}
	e.DetailsField[key] = detail
	return &e
}

func (e DefaultError) WithDetailf(key string, message string, args ...interface{}) *DefaultError {
	if e.DetailsField == nil {
		e.DetailsField = map[string]interface{}{}
	}
	e.DetailsField[key] = fmt.Sprintf(message, args...)
	return &e
}

func ToDefaultError(err error, id string) *DefaultError {
	de := &DefaultError{
		RIDField:     id,
		CodeField:    http.StatusInternalServerError,
		DetailsField: map[string]interface{}{},
		ErrorField:   err.Error(),
	}
	de.Wrap(err)

	if c := ReasonCarrier(nil); stderr.As(err, &c) {
		de.ReasonField = c.Reason()
	}
	if c := RequestIDCarrier(nil); stderr.As(err, &c) && c.RequestID() != "" {
		de.RIDField = c.RequestID()
	}
	if c := DetailsCarrier(nil); stderr.As(err, &c) && c.Details() != nil {
		de.DetailsField = c.Details()
	}
	if c := StatusCarrier(nil); stderr.As(err, &c) && c.Status() != "" {
		de.StatusField = c.Status()
	}
	if c := StatusCodeCarrier(nil); stderr.As(err, &c) && c.StatusCode() != 0 {
		de.CodeField = c.StatusCode()
	}
	if c := DebugCarrier(nil); stderr.As(err, &c) {
		de.DebugField = c.Debug()
	}

	if de.StatusField == "" {
		de.StatusField = http.StatusText(de.StatusCode())
	}

	if de.StatusCode() == http.StatusInternalServerError {
		de.ErrorField = "an error occured"
	}

	return de
}
