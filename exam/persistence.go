package exam

import "context"

type PersistenceProvider interface {
	ExamPersister() Persister
}

type Persister interface {
	FindExam(context.Context, uint64) (*Exam, error)
	CreateExam(context.Context, *Exam) error
	UpdateExam(context.Context, *Exam, ...string) error
}
