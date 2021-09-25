package user

import (
	"context"
)

type ServiceProvider interface {
	UserService() Service
}

type Service interface {
	GetUser(context.Context, uint64) (*User, error)
	GetUserByUsername(context.Context, string) (*User, error)
	GetUserByEmail(context.Context, string) (*User, error)
	RegisterUser(context.Context, *SignupRequest) (*User, error)
	SendResetPassword(context.Context, *SendResetPasswordRequest) error
	ResetPassword(context.Context, *ResetPasswordRequest) error
	ChangePassword(context.Context, *ChangePasswordRequest) error
	SetPassword(context.Context, *SetPasswordRequest) error
	UpdateProfile(context.Context, *UpdateProfileRequest) error
	SendVerification(context.Context, *SendVerificationRequest) error
	Verify(context.Context, *VerificationRequest) error
	SetPreference(context.Context, *SetPreferenceRequest) error
}
