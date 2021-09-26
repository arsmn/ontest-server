package app

import (
	"net/http"

	"github.com/arsmn/ontest-server/module/errors"
)

var (
	ErrInvalidCredentials = &errors.Error{
		CodeField:   http.StatusBadRequest,
		StatusField: http.StatusText(http.StatusBadRequest),
		ErrorField:  "The provided credentials are invalid",
	}
)
