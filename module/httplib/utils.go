package httplib

import (
	"net/http"
	"strings"

	"github.com/mssola/user_agent"
)

func ParseAddr(raw string) (host, port string) {
	if i := strings.LastIndex(raw, ":"); i != -1 {
		return raw[:i], raw[i+1:]
	}
	return raw, ""
}

type UAInfo struct {
	UA            string `json:"ua,omitempty"`
	OSName        string `json:"os_name,omitempty"`
	OSVersion     string `json:"os_version,omitempty"`
	ClientName    string `json:"client_name,omitempty"`
	ClientVersion string `json:"client_version,omitempty"`
}

func ParseUserAgent(uas string) *UAInfo {
	info := UAInfo{UA: uas}
	ua := user_agent.New(uas)

	osInfo := ua.OSInfo()
	info.OSName = osInfo.Name
	info.OSVersion = osInfo.Version

	bName, bVersion := ua.Browser()
	info.ClientName = bName
	info.ClientVersion = bVersion

	return &info
}

func FetchIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
