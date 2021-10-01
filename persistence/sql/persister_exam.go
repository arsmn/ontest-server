package sql

import (
	"context"

	"github.com/arsmn/ontest-server/exam"
	"github.com/arsmn/ontest-server/persistence"
)

var _ persistence.Persister = new(Persister)

func (p *Persister) FindExam(_ context.Context, id uint64) (*exam.Exam, error) {
	s := new(exam.Exam)
	has, err := p.engine.ID(id).Get(s)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, persistence.ErrNoRows
	}
	return s, nil
}

func (p *Persister) CreateExam(_ context.Context, e *exam.Exam) error {
	_, err := p.engine.InsertOne(e)
	return handleError(err)
}

func (p *Persister) UpdateExam(_ context.Context, e *exam.Exam, fields ...string) error {
	_, err := p.engine.ID(e.ID).Cols(fields...).Update(e)
	return err
}