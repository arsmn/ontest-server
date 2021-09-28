package ip

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/arsmn/ontest-server/module/cache"
	"github.com/arsmn/ontest-server/module/xlog"
	"github.com/arsmn/ontest-server/settings"
)

type (
	geoDependencies interface {
		settings.Provider
		cache.Provider
		xlog.Provider
	}
	Geo struct {
		dx     geoDependencies
		client *http.Client
	}
)

func NewIP2LocationGeo(dx geoDependencies) *Geo {
	return &Geo{
		dx:     dx,
		client: new(http.Client),
	}
}

func (g *Geo) FetchData(ctx context.Context, ip string) (*IPLocation, error) {
	var data IPLocation

	key := fmt.Sprintf("ip_%s", ip)
	if err := g.dx.Cacher().Get(ctx, key, &data); err == nil {
		return &data, nil
	}

	q := make(url.Values)
	q.Set("ip", ip)
	q.Set("apiKey", g.dx.Settings().External.IPGeoLocation.APIKey)

	resp, err := g.client.Get("https://api.ipgeolocation.io/ipgeo?" + q.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s responded with a %d trying to fetch ip geo location", "ipgeolocation", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if err := g.dx.Cacher().Set(ctx, &cache.Item{
		Key:   key,
		Value: data,
		TTL:   24 * time.Hour,
	}); err != nil {
		g.dx.Logger().Warn("error while caching ip geo location", xlog.Err(err))
	}

	return &data, nil
}
