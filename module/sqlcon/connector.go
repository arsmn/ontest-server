package sqlcon

import (
	"fmt"
	"net/url"
	"runtime"
	"strings"
)

func cleanURLQuery(in url.Values) (out url.Values) {
	out, _ = url.ParseQuery(in.Encode())
	out.Del("max_conns")
	out.Del("max_idle_conns")
	out.Del("max_conn_lifetime")
	out.Del("parseTime")
	return out
}

func GetDriverName(dsn string) string {
	return strings.Split(dsn, "://")[0]
}

func classifyDSN(dsn string) string {
	scheme := strings.Split(dsn, "://")[0]
	parts := strings.Split(dsn, "@")
	host := parts[len(parts)-1]
	return fmt.Sprintf("%s://*:*@%s", scheme, host)
}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func parseQuery(dsn string) (clean string, query url.Values, err error) {
	query = url.Values{}
	parts := strings.Split(dsn, "?")
	clean = parts[0]
	if len(parts) == 2 {
		if query, err = url.ParseQuery(parts[1]); err != nil {
			return "", query, err
		}
	}
	return
}
