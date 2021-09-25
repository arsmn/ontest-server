package service

import (
	"context"
	stderr "errors"
	"fmt"
	"time"

	"github.com/arsmn/ontest-server/app"
	"github.com/arsmn/ontest-server/module/cache"
	"github.com/arsmn/ontest-server/module/errors"
	"github.com/arsmn/ontest-server/module/generate"
	v "github.com/arsmn/ontest-server/module/validation"
	"github.com/arsmn/ontest-server/module/xlog"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/user"
)

var _ app.App = new(Service)

func (s *Service) GetUser(ctx context.Context, id uint64) (*user.User, error) {
	return s.dx.Persister().FindUser(ctx, id)
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	return s.dx.Persister().FindUserByUsername(ctx, username)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	return s.dx.Persister().FindUserByEmail(ctx, email)
}

func (s *Service) createUser(ctx context.Context, user *user.User) (*user.User, error) {
	return user, s.dx.Persister().CreateUser(ctx, user)
}

func (s *Service) RegisterUser(ctx context.Context, req *user.SignupRequest) (*user.User, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	pswd, err := s.dx.Hasher().Generate(ctx, []byte(req.Password))
	if err != nil {
		return nil, err
	}

	u := &user.User{
		ID:        generate.UID(),
		Username:  generate.HFUID(),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
		Password:  string(pswd),
		Rands:     generate.UserRandCode(),
	}

	return s.createUser(ctx, u)
}

func (s *Service) SendResetPassword(ctx context.Context, req *user.SendResetPasswordRequest) error {
	if err := v.Validate(req); err != nil {
		return err
	}

	u, err := s.dx.Persister().FindUserByEmail(ctx, req.Email)
	if err != nil {
		if stderr.Is(err, persistence.ErrNoRows) {
			return nil
		}
		return err
	}

	code := generate.ResetPasswordCode(u.Email)
	if err := s.dx.Cacher().Set(ctx, &cache.Item{
		Key:   fmt.Sprintf("rpc_%s", code),
		Value: u.ID,
		TTL:   30 * time.Minute,
	}); err != nil {
		return err
	}

	s.dx.Mailer().SendResetPassword(ctx, u, code)

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, req *user.ResetPasswordRequest) error {
	if err := v.Validate(req); err != nil {
		return err
	}

	var uid uint64
	key := fmt.Sprintf("rpc_%s", req.Code)
	if err := s.dx.Cacher().Get(ctx, key, &uid); err != nil {
		return err
	}

	u, err := s.dx.Persister().FindUser(ctx, uid)
	if err != nil {
		return err
	}

	if !generate.VerifyResetPasswordCode(req.Code, u.Email) {
		return errors.ErrBadRequest
	}

	pswd, err := s.dx.Hasher().Generate(ctx, []byte(req.Password))
	if err != nil {
		return err
	}

	if err := s.dx.Cacher().Delete(ctx, key); err != nil {
		s.dx.Logger().Error(fmt.Sprintf("error while deleting cache key: %s", key), xlog.Err(err))
	}

	u.Password = string(pswd)
	u.Rands = generate.UserRandCode()

	return s.dx.Persister().UpdateUser(ctx, u, "password", "rands")
}

func (s *Service) ChangePassword(ctx context.Context, req *user.ChangePasswordRequest) error {
	if err := v.Validate(req); err != nil {
		return err
	}

	u := req.SignedUser()

	if err := s.dx.Hasher().Compare(ctx, []byte(req.CurrentPassword), []byte(u.Password)); err != nil {
		return app.ErrInvalidCredentials
	}

	pswd, err := s.dx.Hasher().Generate(ctx, []byte(req.NewPassword))
	if err != nil {
		return err
	}

	u.Password = string(pswd)
	err = s.dx.Persister().UpdateUser(ctx, u, "password")
	if err != nil {
		return err
	}

	if req.Terminate {
		err = s.dx.Persister().RemoveUserSessions(ctx, u.ID, req.Token())
		if err != nil {
			s.dx.Logger().Error("error while removing user sessions", xlog.Err(err))
		}
	}

	return nil
}

func (s *Service) SetPassword(ctx context.Context, req *user.SetPasswordRequest) error {
	if err := v.Validate(req); err != nil {
		return err
	}

	u := req.SignedUser()

	if u.PasswordSet() {
		return nil
	}

	pswd, err := s.dx.Hasher().Generate(ctx, []byte(req.Password))
	if err != nil {
		return err
	}

	u.Password = string(pswd)
	err = s.dx.Persister().UpdateUser(ctx, u, "password")
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateProfile(ctx context.Context, req *user.UpdateProfileRequest) error {
	if err := v.Validate(req); err != nil {
		return err
	}

	u := req.SignedUser()
	u.FirstName = req.FirstName
	u.LastName = req.LastName
	u.Username = req.Username

	return s.dx.Persister().UpdateUser(ctx, req.SignedUser(), "first_name", "last_name", "username")
}

func (s *Service) SendVerification(ctx context.Context, req *user.SendVerificationRequest) error {
	if err := v.Validate(req); err != nil {
		return err
	}

	u := req.SignedUser()

	if u.EmailVerified {
		return nil
	}

	code := generate.ResetPasswordCode(u.Email)
	if err := s.dx.Cacher().Set(ctx, &cache.Item{
		Key:   fmt.Sprintf("vc_%s", code),
		Value: u.ID,
		TTL:   30 * time.Minute,
	}); err != nil {
		return err
	}

	s.dx.Mailer().SendVerification(ctx, u, code)

	return nil
}

func (s *Service) Verify(ctx context.Context, req *user.VerificationRequest) error {
	if err := v.Validate(req); err != nil {
		return err
	}

	u := req.SignedUser()
	if u.EmailVerified {
		return nil
	}

	var uid uint64
	key := fmt.Sprintf("vc_%s", req.Code)
	if err := s.dx.Cacher().Get(ctx, key, &uid); err != nil {
		return err
	}

	if uid != u.ID {
		return errors.ErrBadRequest
	}

	if err := s.dx.Cacher().Delete(ctx, key); err != nil {
		s.dx.Logger().Error(fmt.Sprintf("error while deleting cache key: %s", key), xlog.Err(err))
	}

	u.EmailVerified = true

	return s.dx.Persister().UpdateUser(ctx, u, "email_verified")
}

func (s *Service) SetPreference(ctx context.Context, req *user.SetPreferenceRequest) error {
	if err := v.Validate(req); err != nil {
		return err
	}

	u := req.SignedUser().SetPreference(req.Key, req.Value)

	return s.dx.Persister().UpdateUser(ctx, u, "preferences")
}
