package transport

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type SignupRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (r SignupRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Email, v.Required, is.Email),
		v.Field(&r.LastName, v.Required, v.Length(3, 50)),
		v.Field(&r.Password, v.Required, v.Length(5, 50)),
		v.Field(&r.FirstName, v.Required, v.Length(3, 50)),
	)
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r SigninRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Email, v.Required),
		v.Field(&r.Password, v.Required),
	)
}

type WhoamiRequest struct {
	Token string
}
