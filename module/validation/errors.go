package validation

import (
	"net/http"

	"github.com/arsmn/ontest-server/module/errors"
)

var (
	ErrValidation = &errors.Error{
		CodeField:   http.StatusBadRequest,
		StatusField: http.StatusText(http.StatusBadRequest),
		ErrorField:  "resource is invalid",
	}
)
