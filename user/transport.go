package user

import (
	"github.com/arsmn/ontest-server/module/validation"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

///// SignedRequest
type SignedRequest struct {
	t string `json:"-"`
	u *User  `json:"-"`
}

func (r *SignedRequest) WithUser(u *User) *SignedRequest {
	r.u = u
	return r
}

func (r *SignedRequest) WithToken(t string) *SignedRequest {
	r.t = t
	return r
}

func (r *SignedRequest) SignedUser() *User {
	return r.u
}

func (r *SignedRequest) Token() string {
	return r.t
}

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

///// SendResetPasswordRequest

type SendResetPasswordRequest struct {
	Email string `json:"email,omitempty"`
}

func (r SendResetPasswordRequest) Validate() error {
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
	SignedRequest
	CurrentPassword string `json:"current_password,omitempty"`
	NewPassword     string `json:"new_password,omitempty"`
	Terminate       bool   `json:"terminate,omitempty"`
}

func (r ChangePasswordRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.CurrentPassword, v.Required),
		v.Field(&r.NewPassword, v.Required, v.Length(5, 50)),
	)
}

///// SetPasswordRequest

type SetPasswordRequest struct {
	SignedRequest
	Password string `json:"password,omitempty"`
}

func (r SetPasswordRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Password, v.Required, v.Length(5, 50)),
	)
}

///// UpdateProfileRequest

type UpdateProfileRequest struct {
	SignedRequest
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

func (r UpdateProfileRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.LastName, v.Required, v.Length(3, 50)),
		v.Field(&r.FirstName, v.Required, v.Length(3, 50)),
		v.Field(&r.Username, v.Required, v.Length(3, 30), v.Match(validation.UsernameRegex)),
	)
}

///// SendVerificationRequest

type SendVerificationRequest struct {
	SignedRequest
}

///// VerificationRequest

type VerificationRequest struct {
	SignedRequest
	Code string `json:"code,omitempty"`
}

func (r VerificationRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Code, v.Required),
	)
}

///// SetPreferenceRequest

type SetPreferenceRequest struct {
	SignedRequest
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func (r SetPreferenceRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Key, v.Required),
	)
}
