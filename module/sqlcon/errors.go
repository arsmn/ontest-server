package sqlcon

import (
	"database/sql"
	stdErr "errors"
	"net/http"

	"github.com/arsmn/ontest/module/errors"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgconn"
	"github.com/lib/pq"
)

var (
	ErrUniqueViolation = &errors.DefaultError{
		CodeField:   http.StatusConflict,
		StatusField: http.StatusText(http.StatusConflict),
		ErrorField:  "Unable to insert or update resource because a resource with that value exists already",
	}
	ErrNoRows = &errors.DefaultError{
		CodeField:   http.StatusNotFound,
		StatusField: http.StatusText(http.StatusNotFound),
		ErrorField:  "Unable to locate the resource",
	}
	ErrConcurrentUpdate = &errors.DefaultError{
		CodeField:   http.StatusBadRequest,
		StatusField: http.StatusText(http.StatusBadRequest),
		ErrorField:  "Unable to serialize access due to a concurrent update in another session",
	}
	ErrNoSuchTable = &errors.DefaultError{
		CodeField:   http.StatusInternalServerError,
		StatusField: http.StatusText(http.StatusInternalServerError),
		ErrorField:  "Unable to locate the table",
	}
)

func HandleError(err error) error {
	if err == nil {
		return nil
	}

	if stdErr.Is(err, sql.ErrNoRows) {
		return ErrNoRows
	}

	switch e := err.(type) {
	case interface{ SQLState() string }:
		return handlePostgres(err, e.SQLState())
	case *pq.Error:
		return handlePostgres(err, string(e.Code))
	case *pgconn.PgError:
		return handlePostgres(err, e.Code)
	case *mysql.MySQLError:
		return handleMySQL(err, e.Number)
	}

	return err
}

func handlePostgres(err error, state string) error {
	switch state {
	case "23505":
		return ErrUniqueViolation
	case "40001":
		return ErrConcurrentUpdate
	case "42P01":
		return ErrNoSuchTable
	}
	return err
}

func handleMySQL(err error, number uint16) error {
	switch number {
	case 1062:
		return ErrUniqueViolation
	case 1146:
		return ErrNoSuchTable
	}
	return err
}
