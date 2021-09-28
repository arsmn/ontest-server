package exam

import "context"

type ServiceProvider interface {
	ExamService() Service
}

type Service interface {
	CreateExam(context.Context, *CreateExamRequest) (*Exam, error)
}
