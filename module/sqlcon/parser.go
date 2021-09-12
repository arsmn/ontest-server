package sqlcon

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/arsmn/ontest/module/xlog"
)

func ParseConnectionOptions(l *xlog.Logger, dsn string) (maxConns int, maxIdleConns int, maxConnLifetime time.Duration, cleanedDSN string) {
	maxConns = maxParallelism() * 2
	maxIdleConns = maxParallelism()
	maxConnLifetime = time.Duration(0)
	cleanedDSN = dsn

	parts := strings.Split(dsn, "?")
	if len(parts) != 2 {
		l.
			Debug("No SQL connection options have been defined, falling back to default connection options.",
				xlog.Int("sql_max_connections", maxConns),
				xlog.Int("sql_max_idle_connections", maxIdleConns),
				xlog.Duration("sql_max_connection_lifetime", maxConnLifetime))
		return
	}

	query, err := url.ParseQuery(parts[1])
	if err != nil {
		l.
			Warn("Unable to parse SQL DSN query, falling back to default connection options.",
				xlog.Int("sql_max_connections", maxConns),
				xlog.Int("sql_max_idle_connections", maxIdleConns),
				xlog.Duration("sql_max_connection_lifetime", maxConnLifetime),
				xlog.Err(err))
		return
	}

	if v := query.Get("max_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			l.Warn(fmt.Sprintf(`SQL DSN query parameter "max_conns" value %v could not be parsed to int, falling back to default value %d`, v, maxConns),
				xlog.Err(err))
		} else {
			maxConns = int(s)
		}
		query.Del("max_conns")
	}

	if v := query.Get("max_idle_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			l.Warn(fmt.Sprintf(`SQL DSN query parameter "max_idle_conns" value %v could not be parsed to int, falling back to default value %d`, v, maxIdleConns),
				xlog.Err(err))
		} else {
			maxIdleConns = int(s)
		}
		query.Del("max_idle_conns")
	}

	if v := query.Get("max_conn_lifetime"); v != "" {
		s, err := time.ParseDuration(v)
		if err != nil {
			l.Warn(fmt.Sprintf(`SQL DSN query parameter "max_conn_lifetime" value %v could not be parsed to duration, falling back to default value %d`, v, maxConnLifetime),
				xlog.Err(err))
		} else {
			maxConnLifetime = s
		}
		query.Del("max_conn_lifetime")
	}

	cleanedDSN = fmt.Sprintf("%s?%s", parts[0], query.Encode())

	return
}

func FinalizeDSN(l *xlog.Logger, dsn string) string {
	if strings.HasPrefix(dsn, "mysql://") {
		var q url.Values
		parts := strings.SplitN(dsn, "?", 2)

		if len(parts) == 1 {
			q = make(url.Values)
		} else {
			var err error
			q, err = url.ParseQuery(parts[1])
			if err != nil {
				l.Warn("Unable to parse SQL DSN query, could not finalize the DSN URI.", xlog.Err(err))
				return dsn
			}
		}

		q.Set("multiStatements", "true")
		q.Set("parseTime", "true")

		return fmt.Sprintf("%s?%s", parts[0], q.Encode())
	}

	return dsn
}
