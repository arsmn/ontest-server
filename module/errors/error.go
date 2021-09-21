package errors

import (
	stderr "errors"
	"fmt"
	"log"
	"net/http"
)

type Error struct {
	CodeField    int                    `json:"code,omitempty"`
	StatusField  string                 `json:"status,omitempty"`
	RIDField     string                 `json:"request,omitempty"`
	ReasonField  string                 `json:"reason,omitempty"`
	DebugField   string                 `json:"debug,omitempty"`
	DetailsField map[string]interface{} `json:"details,omitempty"`
	ErrorField   string                 `json:"message"`

	err error
}

func (e Error) Unwrap() error {
	return e.err
}

func (e *Error) Wrap(err error) {
	e.err = err
}

func (e Error) Status() string {
	return e.StatusField
}

func (e Error) Error() string {
	return e.ErrorField
}

func (e Error) RequestID() string {
	return e.RIDField
}

func (e Error) Reason() string {
	return e.ReasonField
}

func (e Error) Debug() string {
	return e.DebugField
}

func (e Error) Details() map[string]interface{} {
	return e.DetailsField
}

func (e Error) StatusCode() int {
	return e.CodeField
}

func (e Error) WithReason(reason string) *Error {
	e.ReasonField = reason
	return &e
}

func (e Error) WithReasonf(reason string, args ...interface{}) *Error {
	return e.WithReason(fmt.Sprintf(reason, args...))
}

func (e Error) WithError(message string) *Error {
	e.ErrorField = message
	return &e
}

func (e Error) WithErrorf(message string, args ...interface{}) *Error {
	return e.WithError(fmt.Sprintf(message, args...))
}

func (e Error) WithDebugf(debug string, args ...interface{}) *Error {
	return e.WithDebug(fmt.Sprintf(debug, args...))
}

func (e Error) WithDebug(debug string) *Error {
	e.DebugField = debug
	return &e
}

func (e Error) WithDetail(key string, detail interface{}) *Error {
	if e.DetailsField == nil {
		e.DetailsField = map[string]interface{}{}
	}
	e.DetailsField[key] = detail
	return &e
}

func (e Error) WithDetailf(key string, message string, args ...interface{}) *Error {
	if e.DetailsField == nil {
		e.DetailsField = map[string]interface{}{}
	}
	e.DetailsField[key] = fmt.Sprintf(message, args...)
	return &e
}

func ToError(err error, id string) *Error {
	de := &Error{
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
		log.Println(de.ErrorField)
		de.ErrorField = "an error occured"
	}

	return de
}
