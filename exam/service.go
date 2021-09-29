package exam

import "context"

type ServiceProvider interface {
	ExamService() Service
}

type Service interface {
	GetExam(context.Context, uint64) (*Exam, error)
	CreateExam(context.Context, *CreateExamRequest) (*Exam, error)
	UpdateExam(context.Context, *UpdateExamRequest) error
}
