package ip

import "context"

type IPLocation struct {
	IP          string `json:"ip,omitempty"`
	CountryName string `json:"country_name,omitempty"`
	StateProv   string `json:"state_prov,omitempty"`
	City        string `json:"city,omitempty"`
	Latitude    string `json:"latitude,omitempty"`
	Longitude   string `json:"longitude,omitempty"`
	TimeZone    string `json:"time_zone,omitempty"`
	ISP         string `json:"isp,omitempty"`
	CountryFlag string `json:"country_flag,omitempty"`
}

type IP2Location interface {
	FetchData(ctx context.Context, ip string) (*IPLocation, error)
}

type Provider interface {
	IP2Location() IP2Location
}
