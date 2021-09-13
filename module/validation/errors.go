package validation

import (
	"net/http"

	"github.com/arsmn/ontest/module/errors"
)

var (
	ErrValidation = &errors.Error{
		CodeField:   http.StatusBadRequest,
		StatusField: http.StatusText(http.StatusBadRequest),
		ErrorField:  "resource is invalid",
	}
)
