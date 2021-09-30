package session

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

///// SigninRequest

type SigninRequest struct {
	Identifier string `json:"identifier,omitempty"`
	Password   string `json:"password,omitempty"`
	Remember   bool   `json:"remember,omitempty"`
	IP         string `json:"-"`
	UserAgent  string `json:"-"`
}

func (r SigninRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Identifier, v.Required),
		v.Field(&r.Password, v.Required),
	)
}

///// OAuthSignRequest

type OAuthSignRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	IP        string `json:"-"`
	UserAgent string `json:"-"`
}

func (r OAuthSignRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Email, v.Required, is.Email),
	)
}

///// DeleteSessionRequest

type DeleteSessionRequest struct {
	ID uint64 `json:"-"`
}

//// DeleteSessionByTokenRequest

type DeleteSessionByTokenRequest struct {
	Token string `json:"-"`
}

///// GetUserActiveSessionsRequest

type GetUserActiveSessionsRequest struct {
	UserID uint64 `json:"-"`
}

///// GetActiveSessionsResponse

type GetUserActiveSessionsResponse struct {
	Current *Session   `json:"current,omitempty"`
	Others  []*Session `json:"others,omitempty"`
}
