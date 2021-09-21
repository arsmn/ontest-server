package oauth

import (
	"context"

	"github.com/arsmn/ontest-server/session"
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
	FetchData(ctx context.Context, token string) (*session.OAuthSignRequest, error)
}

type Provider interface {
	OAuther(typ OAuthProviderType) OAuther
}
