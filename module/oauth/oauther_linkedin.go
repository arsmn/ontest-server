package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arsmn/ontest-server/settings"
	t "github.com/arsmn/ontest-server/transport"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

type (
	linkedinDependencies interface {
		settings.Provider
	}
	LinkedIn struct {
		dx     githubDependencies
		client *http.Client
	}
)

func NewOAutherLinkedIn(dx linkedinDependencies) *LinkedIn {
	return &LinkedIn{
		dx:     dx,
		client: new(http.Client),
	}
}

func (l *LinkedIn) Config() *oauth2.Config {
	c := l.dx.Settings().OAuth.LinkedIn
	c.Endpoint = linkedin.Endpoint
	return &c
}

func (l *LinkedIn) FetchData(ctx context.Context, token string) (*t.OAuthSignRequest, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.linkedin.com/v2/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s responded with a %d trying to fetch user profile", "LinkedIn", resp.StatusCode)
	}

	var data struct {
		ID        string `json:"id,omitempty"`
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	email, err := l.GetEmail(ctx, token)
	if err != nil {
		return nil, err
	}

	return &t.OAuthSignRequest{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     email,
	}, nil
}

func (l *LinkedIn) GetEmail(_ context.Context, token string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.linkedin.com/v2/clientAwareMemberHandles?q=members&projection=(elements*(primary,type,handle~))", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := l.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s responded with a %d trying to fetch user email", "LinkedIn", resp.StatusCode)
	}

	var data struct {
		Elements []struct {
			Handle struct {
				EmailAddress string `json:"emailAddress,omitempty"`
			} `json:"handle~"`
		} `json:"elements"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	var email string
	if len(data.Elements) > 0 {
		email = data.Elements[0].Handle.EmailAddress
	}

	return email, nil
}
