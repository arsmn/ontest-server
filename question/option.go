package question

import (
	"time"

	"github.com/arsmn/ontest-server/module/generate"
)

type Option struct {
	ID         uint64 `xorm:"pk 'id'" json:"id,omitempty"`
	QuestionID uint64 `xorm:"not null 'question_id'" json:"question_id,omitempty"`
	Text       string `xorm:"varchar(250) not null" json:"text,omitempty"`
	Answer     bool   `xorm:"not null" json:"answer,omitempty"`

	CreatedAt time.Time `xorm:"created" json:"created_at,omitempty" field:"created_at"`
}

func NewOption(qid uint64, text string, answer bool) *Option {
	return &Option{
		ID:         generate.UID(),
		QuestionID: qid,
		Text:       text,
		Answer:     answer,
	}
}
