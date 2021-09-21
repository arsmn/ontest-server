package user

import (
	"context"
)

type ServiceProvider interface {
	UserService() Service
}

type Service interface {
	GetUser(context.Context, uint64) (*User, error)
	RegisterUser(context.Context, *SignupRequest) (*User, error)
	ForgotPassword(context.Context, *ForgotPasswordRequest) error
	ResetPassword(context.Context, *ResetPasswordRequest) error
	ChangePassword(context.Context, *ChangePasswordRequest) error
}
