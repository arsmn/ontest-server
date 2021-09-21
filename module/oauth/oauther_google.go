package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arsmn/ontest-server/module/xlog"
	"github.com/arsmn/ontest-server/settings"
	t "github.com/arsmn/ontest-server/transport"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type (
	googleDependencies interface {
		settings.Provider
		xlog.Provider
	}
	Google struct {
		dx     googleDependencies
		client *http.Client
	}
)

func NewOAutherGoogle(dx googleDependencies) *Google {
	return &Google{
		dx:     dx,
		client: new(http.Client),
	}
}

func (g *Google) Config() *oauth2.Config {
	c := g.dx.Settings().OAuth.Google
	c.Endpoint = google.Endpoint
	return &c
}

func (g *Google) FetchData(_ context.Context, token string) (*t.OAuthSignRequest, error) {
	resp, err := g.client.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s responded with a %d trying to fetch user profile", "Google", resp.StatusCode)
	}

	var data struct {
		ID         string `json:"id,omitempty"`
		Name       string `json:"name,omitempty"`
		GivenName  string `json:"given_name,omitempty"`
		FamilyName string `json:"family_name,omitempty"`
		Email      string `json:"email,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &t.OAuthSignRequest{
		Email:     data.Email,
		FirstName: data.GivenName,
		LastName:  data.FamilyName,
	}, nil
}
