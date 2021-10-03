package exam

import (
	"time"

	"github.com/arsmn/ontest-server/shared"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

///// CreateExamRequest

type CreateExamRequest struct {
	Title    string     `json:"title,omitempty"`
	StartAt  time.Time  `json:"start_at,omitempty"`
	Once     bool       `json:"once,omitempty"`
	Deadline *time.Time `json:"deadline,omitempty"`
}

func (r CreateExamRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Title, v.Required, v.Length(5, 50)),
		v.Field(&r.StartAt, v.Required, v.Min(time.Now().UTC())),
		v.Field(&r.Deadline, v.When(!r.Once, v.Required).Else(v.Nil), v.Min(r.StartAt)),
	)
}

///// UpdateExamRequest

type UpdateExamRequest struct {
	Title        string     `json:"title,omitempty"`
	StartAt      time.Time  `json:"start_at,omitempty"`
	Once         bool       `json:"once,omitempty"`
	Deadline     *time.Time `json:"deadline,omitempty"`
	FreeMovement bool       `json:"free_movement,omitempty"`
}

func (r UpdateExamRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Title, v.Required, v.Length(5, 50)),
		v.Field(&r.StartAt, v.Required, v.Min(time.Now().UTC())),
		v.Field(&r.Deadline, v.When(!r.Once, v.Required).Else(v.Nil), v.Min(r.StartAt)),
	)
}

///// SearchExamRequest

type SearchExamRequest struct {
	shared.PaginatedRequest
}

///// SearchExamResponse

type SearchExamResponse struct {
	shared.PaginatedResponse
	Exams []*Exam
}

///// ExamStatsResponse

type ExamStatsResponse struct {
	TotalDuration       int64 `json:"total_duration"`
	TotalQuestions      int64 `json:"total_questions"`
	TotalScores         int64 `json:"total_scores"`
	TotalNegativeScores int64 `json:"total_negative_scores"`
}

///// SubmitAnswerRequest

type SubmitAnswerRequest struct {
	ID              uint64            `json:"-"`
	Text            string            `json:"text,omitempty"`
	SelectedOptions []*SelectedOption `json:"selected_options,omitempty"`
}
