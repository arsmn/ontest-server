package exam

import (
	"path/filepath"
	"strconv"
	"time"

	"github.com/arsmn/ontest-server/module/generate"
	"github.com/arsmn/ontest-server/question"
	"github.com/spf13/afero"
)

type State int

const (
	Unknown State = iota
	Draft
	Published
)

func (s State) String() string {
	return [...]string{"unknown", "draft", "published"}[s]
}

type Exam struct {
	ID           uint64        `xorm:"pk 'id'" json:"id,omitempty"`
	Examiner     uint64        `xorm:"not null" json:"examiner,omitempty"`
	Title        string        `xorm:"varchar(50) not null" json:"title,omitempty"`
	State        State         `xorm:"not null" json:"state,omitempty"`
	Duration     time.Duration `xorm:"not null" json:"duration,omitempty"`
	StartAt      time.Time     `xorm:"not null" json:"start_at,omitempty"`
	Deadline     *time.Time    `xorm:"null" json:"deadline,omitempty"`
	FreeMovement bool          `xorm:"not null" json:"free_movement,omitempty"`

	CreatedAt time.Time `xorm:"created" json:"created_at,omitempty" field:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at,omitempty" field:"updated_at"`
	DeletedAt time.Time `xorm:"deleted" json:"-" field:"-"`

	Questions []*question.Question `xorm:"-" json:"questions,omitempty"`

	fs afero.Fs `xorm:"-" json:"-" field:"-"`
}

func NewDraftExam(examiner uint64) *Exam {
	return &Exam{
		ID:       generate.UID(),
		Examiner: examiner,
		State:    Draft,
	}
}

func (e *Exam) Fs() afero.Fs {
	if e.fs == nil {
		path := filepath.Join("files", "exams", strconv.FormatUint(e.ID, 10))
		e.fs = afero.NewBasePathFs(afero.NewOsFs(), path)
	}
	return e.fs
}

func (e *Exam) SetTitle(t string) *Exam {
	e.Title = t
	return e
}

func (e *Exam) SetStartAt(sa time.Time) *Exam {
	e.StartAt = sa
	return e
}

func (e *Exam) SetDeadline(d time.Time) *Exam {
	e.Deadline = &d
	return e
}
