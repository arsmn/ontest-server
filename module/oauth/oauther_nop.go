package oauth

import (
	"context"

	"github.com/arsmn/ontest-server/session"
	"golang.org/x/oauth2"
)

type Nop struct{}

func NewOAutherNop() *Nop {
	return &Nop{}
}

func (n *Nop) Config() *oauth2.Config {
	return nil
}

func (n *Nop) FetchData(_ context.Context, token string) (*session.OAuthSignRequest, error) {
	return nil, nil
}
