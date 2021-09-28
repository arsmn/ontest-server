package question

import "time"

type Option struct {
	ID         uint64 `xorm:"pk 'id'" json:"id,omitempty"`
	QuestionID uint64 `xorm:"not null 'question_id'" json:"question_id,omitempty"`
	Text       string `xorm:"varchar(250) not null" json:"text,omitempty"`
	Answer     bool   `xorm:"not null" json:"answer,omitempty"`

	CreatedAt time.Time `xorm:"created" json:"created_at,omitempty" field:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at,omitempty" field:"updated_at"`
	DeletedAt time.Time `xorm:"deleted" json:"-,omitempty" field:"-"`
}
