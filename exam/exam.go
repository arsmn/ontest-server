package exam

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/arsmn/ontest-server/module/generate"
	"github.com/arsmn/ontest-server/module/structs"
	"github.com/arsmn/ontest-server/question"
	"github.com/arsmn/ontest-server/settings"
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
	ID           uint64        `xorm:"pk 'id'" json:"id,omitempty" field:"id"`
	Examiner     uint64        `xorm:"not null" json:"examiner,omitempty" field:"examiner"`
	Title        string        `xorm:"varchar(50) not null" json:"title,omitempty" field:"title"`
	State        State         `xorm:"not null" json:"state,omitempty" field:"state"`
	Duration     time.Duration `xorm:"not null" json:"duration,omitempty" field:"duration"`
	StartAt      time.Time     `xorm:"not null" json:"start_at,omitempty" field:"start_at"`
	Deadline     *time.Time    `xorm:"null" json:"deadline,omitempty" field:"deadline"`
	FreeMovement bool          `xorm:"not null" json:"free_movement,omitempty" field:"free_movement"`

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

func (e *Exam) Cover() string {
	ok, err := afero.Exists(e.Fs(), "cover")
	if !ok || err != nil {
		return ""
	}
	return fmt.Sprintf("%s/files/exams/%d/cover", settings.APIURL(), e.ID)
}

func (e *Exam) Fs() afero.Fs {
	if e.fs == nil {
		path := filepath.Join("files", "exams", strconv.FormatUint(e.ID, 10))
		e.fs = afero.NewBasePathFs(afero.NewOsFs(), path)
	}
	return e.fs
}

func (e *Exam) Map(excludes ...string) map[string]interface{} {
	m := structs.Map(e)
	m["cover"] = e.Cover()
	m["state"] = e.State.String()

	for _, e := range excludes {
		delete(m, e)
	}

	return m
}

func (e *Exam) SetTitle(t string) *Exam {
	e.Title = t
	return e
}

func (e *Exam) SetStartAt(sa time.Time) *Exam {
	e.StartAt = sa
	return e
}

func (e *Exam) SetDeadline(d *time.Time) *Exam {
	e.Deadline = d
	return e
}

func (e *Exam) SetFreeMovement(fm bool) *Exam {
	e.FreeMovement = fm
	return e
}
