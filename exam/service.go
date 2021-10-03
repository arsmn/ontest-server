package exam

import (
	"context"

	"github.com/arsmn/ontest-server/user"
)

type ServiceProvider interface {
	ExamService() Service
}

type Service interface {
	GetExam(context.Context, uint64) (*Exam, error)
	GetExamStats(context.Context, uint64) (*ExamStatsResponse, error)
	CreateExam(context.Context, *CreateExamRequest) (*Exam, error)
	UpdateExam(context.Context, *UpdateExamRequest) error
	PublishExam(context.Context, *Exam) error
	SearchExam(context.Context, *SearchExamRequest) (*SearchExamResponse, error)
	Participate(context.Context, *user.User, *Exam) (*Result, error)
	GetResult(context.Context, uint64) (*Result, error)
	SubmitAnswer(context.Context, *SubmitAnswerRequest) error
}
