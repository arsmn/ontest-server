package user

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/arsmn/ontest-server/module/structs"
	"github.com/arsmn/ontest-server/settings"
	"github.com/spf13/afero"
)

type User struct {
	ID            uint64            `xorm:"pk 'id'" json:"id" field:"id"`
	Username      string            `xorm:"varchar(50) not null unique" json:"username" field:"username"`
	FirstName     string            `xorm:"varchar(50) not null" json:"first_name" field:"first_name"`
	LastName      string            `xorm:"varchar(50) not null" json:"last_name" field:"last_name"`
	Email         string            `xorm:"varchar(100) not null unique" json:"email" field:"email"`
	Password      string            `xorm:"varchar(250) not null" json:"-" field:"-"`
	EmailVerified bool              `xorm:"not null" json:"email_verified" field:"email_verified"`
	Verified      bool              `xorm:"not null" json:"verified" field:"verified"`
	IsActive      bool              `xorm:"not null" json:"is_active" field:"is_active"`
	Rands         string            `xorm:"varchar(10) not null" json:"-" field:"-"`
	Preferences   map[string]string `xorm:"json" json:"preferences" field:"preferences"`

	CreatedAt time.Time `xorm:"created" json:"created_at" field:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at" field:"updated_at"`
	DeletedAt time.Time `xorm:"deleted" json:"-" field:"-"`

	fs afero.Fs `xorm:"-" json:"-" field:"-"`
}

func (u *User) FullName() string {
	if u.FirstName == "" || u.LastName == "" {
		return u.FirstName + u.LastName
	}
	return u.FirstName + " " + u.LastName
}

func (u *User) PasswordSet() bool {
	return u.Password != ""
}

func (u *User) Fs() afero.Fs {
	if u.fs == nil {
		path := filepath.Join("files", strconv.FormatUint(u.ID, 10))
		u.fs = afero.NewBasePathFs(afero.NewOsFs(), path)
	}
	return u.fs
}

func (u *User) Avatar() string {
	ok, err := afero.Exists(u.Fs(), "avatar")
	if !ok || err != nil {
		return ""
	}
	return fmt.Sprintf("%s/files/%d/avatar", settings.APIURL(), u.ID)
}

func (u *User) SetPreference(key, value string) *User {
	if u.Preferences == nil {
		u.Preferences = make(map[string]string)
	}
	u.Preferences[key] = value
	return u
}

var PrivateFields = []string{
	"email",
	"password",
	"email_verified",
	"is_active",
	"rands",
	"updated_at",
	"deleted_at",
}

func (u *User) CopySanitize(fields ...string) *User {
	var uu = *u

	for _, f := range fields {
		switch f {
		case "id":
			uu.ID = 0
		case "username":
			uu.Username = ""
		case "first_name":
			uu.FirstName = ""
		case "last_name":
			uu.LastName = ""
		case "email":
			uu.Email = ""
		case "password":
			uu.Password = ""
		case "email_verified":
			uu.EmailVerified = false
		case "verified":
			uu.Verified = false
		case "is_active":
			uu.IsActive = false
		case "rands":
			uu.Rands = ""
		case "created_at":
			uu.CreatedAt = time.Time{}
		case "updated_at":
			uu.UpdatedAt = time.Time{}
		case "deleted_at":
			uu.DeletedAt = time.Time{}
		}
	}

	return &uu
}

func (u *User) Map(excludes ...string) map[string]interface{} {
	m := structs.Map(u)
	m["avatar"] = u.Avatar()
	m["full_name"] = u.FullName()
	m["password_set"] = u.PasswordSet()

	for _, e := range excludes {
		delete(m, e)
	}

	return m
}
