package user

import "time"

type User struct {
	ID        int64
	Username  string
	FirstName string
	LastName  string
	Email     string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
