package exam

import "context"

type PersistenceProvider interface {
	ExamPersister() Persister
}

type Persister interface {
	FindExam(context.Context, uint64) (*Exam, error)
	CreateExam(context.Context, *Exam) error
	UpdateExam(context.Context, *Exam, ...string) error
	SearchExam(_ context.Context, q string, page, pageSize int) (int64, []*Exam, error)
	ExamStats(context.Context, uint64) (*ExamStatsResponse, error)
	CreateResult(context.Context, *Result) error
	FindResult(_ context.Context, id uint64) (*Result, error)
	FindResultByExam(_ context.Context, uid, eid uint64) (*Result, error)
	FindResultAnswers(_ context.Context, rid uint64) ([]*Answer, error)
	UpdateAnswer(context.Context, *Answer, ...string) error
}
