package app

import (
	"net/http"

	"github.com/arsmn/ontest/module/errors"
)

var (
	ErrInvalidCredentials = &errors.Error{
		CodeField:   http.StatusBadRequest,
		StatusField: http.StatusText(http.StatusBadRequest),
		ErrorField:  "the provided credentials are invalid",
	}
)
