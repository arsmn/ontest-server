package question

import (
	"path/filepath"
	"strconv"
	"time"

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
	ID            uint64        `xorm:"pk 'id'" json:"id,omitempty"`
	ExamID        uint64        `xorm:"not null 'exam_id'" json:"exam_id,omitempty"`
	Text          string        `xorm:"varchar(250) not null" json:"text,omitempty"`
	Type          Type          `xorm:"not null" json:"type,omitempty"`
	Duration      time.Duration `xorm:"null" json:"duration,omitempty"`
	Score         int           `xorm:"not null" json:"score,omitempty"`
	NegativeScore int           `xorm:"not null" json:"negative_score,omitempty"`

	CreatedAt time.Time `xorm:"created" json:"created_at,omitempty" field:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at,omitempty" field:"updated_at"`
	DeletedAt time.Time `xorm:"deleted" json:"-,omitempty" field:"-"`

	Options []*Option `xorm:"-" json:"options,omitempty"`

	fs afero.Fs `xorm:"-" json:"-,omitempty" field:"-"`
}

func (q *Question) Fs() afero.Fs {
	if q.fs == nil {
		path := filepath.Join("files", "questions", strconv.FormatUint(q.ID, 10))
		q.fs = afero.NewBasePathFs(afero.NewOsFs(), path)
	}
	return q.fs
}
