package exam

import (
	"time"

	"github.com/arsmn/ontest-server/question"
)

type Result struct {
	ID       uint64 `xorm:"pk 'id'" json:"id"`
	Examinee uint64 `xorm:"not null" json:"examinee"`
	ExamID   uint64 `xorm:"not null 'exam_id'" json:"exam_id"`

	CreatedAt time.Time `xorm:"created" json:"created_at,omitempty" field:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at,omitempty" field:"updated_at"`
	DeletedAt time.Time `xorm:"deleted" json:"-" field:"-"`

	Answers []*Answer `xorm:"-" json:"answers"`
}

type Answer struct {
	ID              uint64            `xorm:"pk 'id'" json:"id"`
	QuestionID      uint64            `xorm:"not null 'question_id'" json:"question_id"`
	ResultID        uint64            `xorm:"not null 'result_id'" json:"result_id"`
	Text            string            `xorm:"not null varchar(250)" json:"text"`
	SelectedOptions []*SelectedOption `xorm:"json null" json:"selected_options"`

	CreatedAt time.Time `xorm:"created" json:"created_at,omitempty" field:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at,omitempty" field:"updated_at"`
	DeletedAt time.Time `xorm:"deleted" json:"-" field:"-"`

	Active   bool               `xorm:"-" json:"active"`
	Question *question.Question `xorm:"-" json:"question"`
}

type SelectedOption struct {
	OptionID uint64 `json:"option_id"`
	Answer   bool   `json:"answer"`
}
