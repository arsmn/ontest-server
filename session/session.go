package session

import "time"

type Session struct {
	ID        uint64    `xorm:"pk 'id'"`
	UserID    uint64    `xorm:"'user_id'"`
	Token     string    `xorm:"varchar(50) not null unique"`
	Active    bool      `xorm:"not null"`
	IssuedAt  time.Time `xorm:"not null"`
	ExpiresAt time.Time `xorm:"not null"`

	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func (s *Session) IsActive() bool {
	return s.Active && s.ExpiresAt.After(time.Now().UTC())
}
