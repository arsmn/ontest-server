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
	"github.com/arsmn/ontest-server/persistence"
	t "github.com/arsmn/ontest-server/transport"
	"github.com/arsmn/ontest-server/user"
)

var _ app.App = new(Service)

func (s *Service) GetUser(ctx context.Context, id uint64) (*user.User, error) {
	return s.dx.Persister().FindUser(ctx, id)
}

func (s *Service) createUser(ctx context.Context, user *user.User) (*user.User, error) {
	return user, s.dx.Persister().CreateUser(ctx, user)
}

func (s *Service) RegisterUser(ctx context.Context, req *t.SignupRequest) (*user.User, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	pswd, err := s.dx.Hasher().Generate(ctx, []byte(req.Password))
	if err != nil {
		return nil, err
	}

	u := &user.User{
		ID:       generate.UID(),
		Username: generate.HFUID(),
		Email:    req.Email,
		IsActive: true,
		Password: string(pswd),
		Rands:    generate.UserRandCode(),
	}

	return s.createUser(ctx, u)
}

func (s *Service) ForgotPassword(ctx context.Context, req *t.ForgotPasswordRequest) error {
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

	// send code

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, req *t.ResetPasswordRequest) error {
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

	u.Password = string(pswd)
	u.Rands = generate.UserRandCode()
	if err := s.dx.Persister().UpdateUser(ctx, u, "password", "rands"); err != nil {
		return err
	}

	return s.dx.Cacher().Delete(ctx, key)
}
