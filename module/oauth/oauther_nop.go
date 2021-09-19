package oauth

import (
	"context"

	t "github.com/arsmn/ontest-server/transport"
	"golang.org/x/oauth2"
)

type Nop struct{}

func NewOAutherNop() *Nop {
	return &Nop{}
}

func (n *Nop) Config() *oauth2.Config {
	return nil
}

func (n *Nop) FetchData(_ context.Context, token string) (*t.OAuthSignRequest, error) {
	return nil, nil
}
