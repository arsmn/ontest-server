package question

import "context"

type PersistenceProvider interface {
	QuestionPersister() Persister
}

type Persister interface {
	FindQuestion(context.Context, uint64) (*Question, error)
	FindQuestionOptions(_ context.Context, qid uint64) ([]*Option, error)
	FindListQuestions(ctx context.Context, eid uint64, page, pageSize int) (int64, []*Question, error)
	CreateQuestion(context.Context, *Question) error
	UpdateQuestion(context.Context, *Question, ...string) error
	RemoveQuestion(context.Context, uint64) error
}
