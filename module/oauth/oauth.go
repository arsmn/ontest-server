package oauth

import (
	"context"

	t "github.com/arsmn/ontest-server/transport"
	"golang.org/x/oauth2"
)

type OAuthProviderType int

const (
	UnknownType OAuthProviderType = iota
	GoogleType
	GitHubType
	LinkedInType
)

func (t OAuthProviderType) String() string {
	return [...]string{"unknown", "google", "github", "linkedin"}[t]
}

type OAuther interface {
	Config() *oauth2.Config
	FetchData(ctx context.Context, token string) (*t.OAuthSignRequest, error)
}

type Provider interface {
	OAuther(typ OAuthProviderType) OAuther
}
