package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arsmn/ontest-server/settings"
	t "github.com/arsmn/ontest-server/transport"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type (
	githubDependencies interface {
		settings.Provider
	}
	GitHub struct {
		dx     githubDependencies
		client *http.Client
	}
)

func NewOAutherGitHub(dx githubDependencies) *GitHub {
	return &GitHub{
		dx:     dx,
		client: new(http.Client),
	}
}

func (g *GitHub) Config() *oauth2.Config {
	c := g.dx.Settings().OAuth().GitHub
	c.Endpoint = github.Endpoint
	return &c
}

func (g *GitHub) FetchData(_ context.Context, token string) (*t.OAuthSignRequest, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+token)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s responded with a %d trying to fetch user profile", "GitHub", resp.StatusCode)
	}

	var data struct {
		ID    int    `json:"id,omitempty"`
		Login string `json:"login,omitempty"`
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &t.OAuthSignRequest{
		Email:     data.Email,
		FirstName: data.Name,
	}, nil
}
