package question

import "context"

type ServiceProvider interface {
	QuestionService() Service
}

type Service interface {
	GetQuestion(context.Context, uint64) (*Question, error)
	GetQuestionList(context.Context, *GetQuestionListRequest) (*GetQuestionListResponse, error)
	CreateQuestion(context.Context, *CreateQuestionRequest) (*Question, error)
	UpdateQuestion(context.Context, *UpdateQuestionRequest) error
}
