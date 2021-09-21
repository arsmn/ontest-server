package session

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

///// SigninRequest

type SigninRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (r SigninRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Email, v.Required),
		v.Field(&r.Password, v.Required),
	)
}

///// OAuthSignRequest

type OAuthSignRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
}

func (r OAuthSignRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Email, v.Required, is.Email),
	)
}
