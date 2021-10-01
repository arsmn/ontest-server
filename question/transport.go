package question

import (
	"github.com/arsmn/ontest-server/shared"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

///// GetQuestionListRequest

type GetQuestionListRequest struct {
	shared.PaginatedRequest
	ExamID uint64 `json:"-"`
}

///// GetQuestionListResponse

type GetQuestionListResponse struct {
	shared.PaginatedResponse
	Questions []*Question `json:"questions"`
}

///// CreateQuestionRequest

type CreateQuestionRequest struct {
	ExamID        uint64 `json:"-"`
	Text          string `json:"text,omitempty"`
	Type          Type   `json:"type,omitempty"`
	Duration      int64  `json:"duration,omitempty"`
	Score         int    `json:"score,omitempty"`
	NegativeScore int    `json:"negative_score,omitempty"`
	Options       []struct {
		Text   string `json:"text,omitempty"`
		Answer bool   `json:"answer,omitempty"`
	} `json:"options,omitempty"`
}

func (r CreateQuestionRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Text, v.Required, v.Length(1, 250)),
		v.Field(&r.Score, v.Required, v.Min(0)),
		v.Field(&r.NegativeScore, v.Required, v.Min(0)),
		v.Field(&r.Options, v.When(r.Type == SingleChoice ||
			r.Type == MultipleChoice, v.Length(2, 10))),
	)
}

///// UpdateQuestionRequest

type UpdateQuestionRequest struct {
	QuestionID    uint64 `json:"-"`
	Text          string `json:"text,omitempty"`
	Type          Type   `json:"type,omitempty"`
	Duration      int64  `json:"duration,omitempty"`
	Score         int    `json:"score,omitempty"`
	NegativeScore int    `json:"negative_score,omitempty"`
	Options       []struct {
		Text   string `json:"text,omitempty"`
		Answer bool   `json:"answer,omitempty"`
	} `json:"options,omitempty"`
}

func (r UpdateQuestionRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Text, v.Required, v.Length(1, 250)),
		v.Field(&r.Score, v.Required, v.Min(0)),
		v.Field(&r.NegativeScore, v.Required, v.Min(0)),
		v.Field(&r.Options, v.When(r.Type == SingleChoice ||
			r.Type == MultipleChoice, v.Length(2, 10))),
	)
}
