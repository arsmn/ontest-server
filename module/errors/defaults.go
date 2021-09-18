package errors

import "net/http"

var (
	ErrNotFound = Error{
		StatusField: http.StatusText(http.StatusNotFound),
		ErrorField:  "The requested resource could not be found",
		CodeField:   http.StatusNotFound,
	}

	ErrUnauthorized = Error{
		StatusField: http.StatusText(http.StatusUnauthorized),
		ErrorField:  "The request could not be authorized",
		CodeField:   http.StatusUnauthorized,
	}

	ErrForbidden = Error{
		StatusField: http.StatusText(http.StatusForbidden),
		ErrorField:  "The requested action was forbidden",
		CodeField:   http.StatusForbidden,
	}

	ErrInternalServerError = Error{
		StatusField: http.StatusText(http.StatusInternalServerError),
		ErrorField:  "An internal server error occurred, please contact the system administrator",
		CodeField:   http.StatusInternalServerError,
	}

	ErrBadRequest = Error{
		StatusField: http.StatusText(http.StatusBadRequest),
		ErrorField:  "The request was malformed or contained invalid parameters",
		CodeField:   http.StatusBadRequest,
	}

	ErrUnsupportedMediaType = Error{
		StatusField: http.StatusText(http.StatusUnsupportedMediaType),
		ErrorField:  "The request is using an unknown content type",
		CodeField:   http.StatusUnsupportedMediaType,
	}

	ErrConflict = Error{
		StatusField: http.StatusText(http.StatusConflict),
		ErrorField:  "The resource could not be created due to a conflict",
		CodeField:   http.StatusConflict,
	}
)
