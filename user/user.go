package user

import (
	"time"
)

type User struct {
	ID              uint64 `xorm:"pk 'id'"`
	Username        string `xorm:"varchar(50) not null unique"`
	FirstName       string `xorm:"varchar(50) not null"`
	LastName        string `xorm:"varchar(50) not null"`
	Email           string `xorm:"varchar(100) not null unique"`
	Password        string `xorm:"varchar(250) not null"`
	IsEmailVerified bool   `xorm:"not null"`
	IsVerified      bool   `xorm:"not null"`
	IsActive        bool   `xorm:"not null"`

	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}
