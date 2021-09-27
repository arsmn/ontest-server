package session

import (
	"time"

	"github.com/arsmn/ontest-server/module/generate"
	"github.com/arsmn/ontest-server/module/httplib"
	"github.com/arsmn/ontest-server/module/httplib/ip"
)

type Session struct {
	ID         uint64      `xorm:"pk 'id'" json:"id"`
	UserID     uint64      `xorm:"'user_id'" json:"-"`
	Token      string      `xorm:"varchar(50) not null unique" json:"-"`
	Active     bool        `xorm:"not null" json:"-"`
	IssuedAt   time.Time   `xorm:"not null" json:"issued_at"`
	ExpiresAt  time.Time   `xorm:"not null" json:"-"`
	Properties *Properties `xorm:"json" json:"properties"`

	CreatedAt time.Time `xorm:"created" json:"-"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`
	DeletedAt time.Time `xorm:"deleted" json:"-"`
}

type Properties struct {
	IPLocation *ip.IPLocation  `json:"ip_location,omitempty"`
	UAInfo     *httplib.UAInfo `json:"ua_info,omitempty"`
}

func (s *Session) SetIPLocation(ipl *ip.IPLocation) *Session {
	if s.Properties == nil {
		s.Properties = new(Properties)
	}
	s.Properties.IPLocation = ipl
	return s
}

func (s *Session) SetUAInfo(uai *httplib.UAInfo) *Session {
	if s.Properties == nil {
		s.Properties = new(Properties)
	}
	s.Properties.UAInfo = uai
	return s
}

func NewActiveSession(userID uint64, lifespan time.Duration) *Session {
	return &Session{
		ID:        generate.UID(),
		UserID:    userID,
		Token:     generate.SessionToken(),
		Active:    true,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(lifespan),
	}
}

func (s *Session) IsActive() bool {
	return s.Active && s.ExpiresAt.After(time.Now().UTC())
}
