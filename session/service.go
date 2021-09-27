package session

import (
	"context"
)

type ServiceProvider interface {
	SessionService() Service
}

type Service interface {
	GetSession(context.Context, string) (*Session, error)
	GetUserActiveSessions(context.Context, *GetUserActiveSessionsRequest) (*GetUserActiveSessionsResponse, error)
	DeleteSession(context.Context, *DeleteSessionRequest) error
	DeleteSessionByToken(context.Context, *DeleteSessionByTokenRequest) error
	IssueSession(context.Context, *SigninRequest) (*Session, error)
	OAuthIssueSession(context.Context, *OAuthSignRequest) (*Session, error)
}
