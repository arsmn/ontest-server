package cache

import (
	"net/http"

	"github.com/arsmn/ontest-server/module/errors"
)

var (
	ErrCacheMissing = &errors.Error{
		CodeField:   http.StatusNotFound,
		StatusField: http.StatusText(http.StatusNotFound),
		ErrorField:  "unable to locate the resource",
	}
)
