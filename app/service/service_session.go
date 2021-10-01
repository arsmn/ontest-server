package service

import (
	"context"
	stderr "errors"

	"github.com/arsmn/ontest-server/app"
	c "github.com/arsmn/ontest-server/module/context"
	"github.com/arsmn/ontest-server/module/httplib"
	v "github.com/arsmn/ontest-server/module/validation"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/session"
	"github.com/arsmn/ontest-server/user"
)

var _ app.App = new(Service)

func (s *Service) GetSession(ctx context.Context, token string) (*session.Session, error) {
	return s.dx.Persister().FindSessionByToken(ctx, token)
}

func (s *Service) IssueSession(ctx context.Context, req *session.SigninRequest) (*session.Session, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	find := func(ctx context.Context, typ, value string) (*user.User, error) {
		switch typ {
		case "email":
			return s.dx.Persister().FindUserByEmail(ctx, value)
		case "username":
			return s.dx.Persister().FindUserByUsername(ctx, value)
		default:
			return nil, nil
		}
	}

	var user *user.User
	for _, t := range []string{"email", "username"} {
		u, err := find(ctx, t, req.Identifier)
		if err != nil {
			if stderr.Is(err, persistence.ErrNoRows) {
				continue
			}
			return nil, err
		}
		user = u
		break
	}

	if user == nil {
		return nil, app.ErrInvalidCredentials
	}

	if !user.PasswordSet() {
		return nil, app.ErrInvalidCredentials
	}

	if err := s.dx.Hasher().Compare(ctx, []byte(req.Password), []byte(user.Password)); err != nil {
		return nil, app.ErrInvalidCredentials
	}

	hv := c.HTTPValues(ctx)
	sess := session.NewActiveSession(user.ID, s.dx.Settings().Session.Lifespan)

	if len(hv.IP) != 0 {
		ipl, err := s.dx.IP2Location().FetchData(ctx, hv.IP)
		if err == nil {
			sess.SetIPLocation(ipl)
		}
	}

	if len(hv.UserAgent) != 0 {
		uai := httplib.ParseUserAgent(hv.UserAgent)
		sess.SetUAInfo(uai)
	}

	if err := s.dx.Persister().CreateSession(ctx, sess); err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *Service) OAuthIssueSession(ctx context.Context, req *session.OAuthSignRequest) (*session.Session, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	var newuser bool
	u, err := s.dx.Persister().FindUserByEmail(ctx, req.Email)
	if err != nil {
		if stderr.Is(err, persistence.ErrNoRows) {
			newuser = true
		} else {
			return nil, err
		}
	}

	if newuser {
		u = user.NewActiveUser(req.FirstName, req.LastName, req.Email)
		u.EmailVerified = true
		if _, err := s.createUser(ctx, u); err != nil {
			return nil, err
		}
	}

	hv := c.HTTPValues(ctx)
	sess := session.NewActiveSession(u.ID, s.dx.Settings().Session.Lifespan)

	if len(hv.IP) != 0 {
		ipl, err := s.dx.IP2Location().FetchData(ctx, hv.IP)
		if err == nil {
			sess.SetIPLocation(ipl)
		}
	}

	if len(hv.UserAgent) != 0 {
		uai := httplib.ParseUserAgent(hv.UserAgent)
		sess.SetUAInfo(uai)
	}

	if err := s.dx.Persister().CreateSession(ctx, sess); err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *Service) GetUserActiveSessions(ctx context.Context, req *session.GetUserActiveSessionsRequest) (*session.GetUserActiveSessionsResponse, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	sessions, err := s.dx.Persister().FindUserSessions(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	t := c.Session(ctx).Token

	response := new(session.GetUserActiveSessionsResponse)
	for _, sess := range sessions {
		if sess.Token == t {
			response.Current = sess
		} else {
			response.Others = append(response.Others, sess)
		}
	}

	return response, nil
}

func (s *Service) deleteSession(ctx context.Context, sess *session.Session, err error) error {
	if err != nil {
		if stderr.Is(err, persistence.ErrNoRows) {
			return nil
		}
		return err
	}

	if !sess.IsActive() {
		return nil
	}

	return s.dx.Persister().RemoveSession(ctx, sess.ID)
}

func (s *Service) DeleteSessionByToken(ctx context.Context, req *session.DeleteSessionByTokenRequest) error {
	sess, err := s.dx.Persister().FindSessionByToken(ctx, req.Token)
	return s.deleteSession(ctx, sess, err)
}

func (s *Service) DeleteSession(ctx context.Context, req *session.DeleteSessionRequest) error {
	sess, err := s.dx.Persister().FindSession(ctx, req.ID)
	return s.deleteSession(ctx, sess, err)
}
