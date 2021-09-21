package session

import (
	"context"
)

type ServiceProvider interface {
	SessionService() Service
}

type Service interface {
	GetSession(context.Context, string) (*Session, error)
	DeleteSession(context.Context, string) error
	IssueSession(context.Context, *SigninRequest) (*Session, error)
	OAuthIssueSession(context.Context, *OAuthSignRequest) (*Session, error)
}
