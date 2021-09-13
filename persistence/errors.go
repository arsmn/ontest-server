package persistence

import (
	"net/http"

	"github.com/arsmn/ontest/module/errors"
)

var (
	ErrNoRows = &errors.Error{
		CodeField:   http.StatusNotFound,
		StatusField: http.StatusText(http.StatusNotFound),
		ErrorField:  "Unable to locate the resource",
	}
	ErrUniqueViolation = &errors.Error{
		CodeField:   http.StatusConflict,
		StatusField: http.StatusText(http.StatusConflict),
		ErrorField:  "Unable to insert or update resource because a resource with that value exists already",
	}
)
