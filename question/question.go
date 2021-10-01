package question

import (
	"path/filepath"
	"strconv"
	"time"

	"github.com/arsmn/ontest-server/module/generate"
	"github.com/spf13/afero"
)

type Type int

const (
	Unknown Type = iota
	Descriptive
	SingleChoice
	MultipleChoice
)

func (t Type) String() string {
	return [...]string{"unknown", "descriptive", "single_choice", "multiple_choice"}[t]
}

type Question struct {
	ID            uint64        `xorm:"pk 'id'" json:"id" field:"id"`
	Examiner      uint64        `xorm:"not null" json:"examiner" field:"examiner"`
	ExamID        uint64        `xorm:"not null 'exam_id'" json:"exam_id" field:"exam_id"`
	Text          string        `xorm:"varchar(250) not null" json:"text" field:"text"`
	Type          Type          `xorm:"not null" json:"type" field:"type"`
	Duration      time.Duration `xorm:"null" json:"duration" field:"duration"`
	Score         int           `xorm:"not null" json:"score" field:"score"`
	NegativeScore int           `xorm:"not null" json:"negative_score" field:"negative_score"`

	CreatedAt time.Time `xorm:"created" json:"created_at" field:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at" field:"updated_at"`
	DeletedAt time.Time `xorm:"deleted" json:"-" field:"-"`

	Options []*Option `xorm:"-" json:"options" field:"options"`

	fs afero.Fs `xorm:"-" json:"-,omitempty" field:"-"`
}

func NewQuestion(uid, eid uint64, text string) *Question {
	return &Question{
		ID:       generate.UID(),
		Examiner: uid,
		ExamID:   eid,
		Text:     text,
	}
}

func (q *Question) Fs() afero.Fs {
	if q.fs == nil {
		path := filepath.Join("files", "questions", strconv.FormatUint(q.ID, 10))
		q.fs = afero.NewBasePathFs(afero.NewOsFs(), path)
	}
	return q.fs
}

func (q *Question) SetType(typ Type) *Question {
	q.Type = typ
	return q
}

func (q *Question) SetText(text string) *Question {
	q.Text = text
	return q
}

func (q *Question) SetDuration(d time.Duration) *Question {
	q.Duration = d
	return q
}

func (q *Question) SetScore(s int) *Question {
	q.Score = s
	return q
}

func (q *Question) SetNegativeScore(ns int) *Question {
	q.NegativeScore = ns
	return q
}

func (q *Question) SetOptions(opts []*Option) *Question {
	q.Options = opts
	return q
}

func (q *Question) AddOption(opt *Option) *Question {
	if q.Options == nil {
		q.Options = make([]*Option, 0)
	}
	q.Options = append(q.Options, opt)
	return q
}
