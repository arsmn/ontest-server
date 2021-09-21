package user

import (
	"time"
)

type User struct {
	ID            uint64 `xorm:"pk 'id'" json:"id,omitempty"`
	Username      string `xorm:"varchar(50) not null unique" json:"username,omitempty"`
	FirstName     string `xorm:"varchar(50) not null" json:"first_name,omitempty"`
	LastName      string `xorm:"varchar(50) not null" json:"last_name,omitempty"`
	Email         string `xorm:"varchar(100) not null unique" json:"email,omitempty"`
	Password      string `xorm:"varchar(250) not null" json:"-"`
	EmailVerified bool   `xorm:"not null" json:"email_verified,omitempty"`
	Verified      bool   `xorm:"not null" json:"verified,omitempty"`
	IsActive      bool   `xorm:"not null" json:"is_active,omitempty"`
	Rands         string `xorm:"varchar(10) not null" json:"-"`

	CreatedAt time.Time `xorm:"created" json:"created_at,omitempty"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at,omitempty"`
	DeletedAt time.Time `xorm:"deleted" json:"-"`
}

var defaultSanitizeFields = []string{
	"Password",
	"EmailVerified",
	"IsActive",
	"UpdatedAt",
	"DeletedAt",
}

func (u *User) CopySanitize(fields ...string) *User {
	if len(fields) == 0 {
		fields = defaultSanitizeFields
	}

	var uu = *u

	for _, f := range fields {
		switch f {
		case "ID":
			uu.ID = 0
		case "FirstName":
			uu.FirstName = ""
		case "LastName":
			uu.LastName = ""
		case "Email":
			uu.Email = ""
		case "Password":
			uu.Password = ""
		case "EmailVerified":
			uu.EmailVerified = false
		case "Verified":
			uu.Verified = false
		case "IsActive":
			uu.IsActive = false
		case "CreatedAt":
			uu.CreatedAt = time.Time{}
		case "UpdatedAt":
			uu.UpdatedAt = time.Time{}
		case "DeletedAt":
			uu.DeletedAt = time.Time{}
		}
	}

	return &uu
}
