package user

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

///// SignupRequest

type SignupRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}

func (r SignupRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Email, v.Required, is.Email),
		v.Field(&r.LastName, v.Required, v.Length(3, 50)),
		v.Field(&r.Password, v.Required, v.Length(5, 50)),
		v.Field(&r.FirstName, v.Required, v.Length(3, 50)),
	)
}

///// ForgotPasswordRequest

type ForgotPasswordRequest struct {
	Email string `json:"email,omitempty"`
}

func (r ForgotPasswordRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Email, v.Required, is.Email),
	)
}

///// ResetPasswordRequest

type ResetPasswordRequest struct {
	Code     string `json:"code,omitempty"`
	Password string `json:"password,omitempty"`
}

func (r ResetPasswordRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Code, v.Required),
		v.Field(&r.Password, v.Required, v.Length(5, 50)),
	)
}

///// ChangePasswordRequest

type ChangePasswordRequest struct {
	*User
	CurrentPassword string `json:"current_password,omitempty"`
	NewPassword     string `json:"new_password,omitempty"`
}

func (r ChangePasswordRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.CurrentPassword, v.Required, v.Length(5, 50)),
		v.Field(&r.NewPassword, v.Required, v.Length(5, 50)),
	)
}
