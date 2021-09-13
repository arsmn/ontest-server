package sql

import (
	"database/sql"
	stdErr "errors"

	"github.com/arsmn/ontest/persistence"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgconn"
	"github.com/lib/pq"
)

func handleError(err error) error {
	if err == nil {
		return nil
	}

	if stdErr.Is(err, sql.ErrNoRows) {
		return persistence.ErrNoRows
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
		return persistence.ErrUniqueViolation
	}
	return err
}

func handleMySQL(err error, number uint16) error {
	switch number {
	case 1062:
		return persistence.ErrUniqueViolation
	}
	return err
}
