package sql

import (
	"context"

	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/question"
)

var _ persistence.Persister = new(Persister)

func (p *Persister) FindQuestion(_ context.Context, id uint64) (*question.Question, error) {
	q := new(question.Question)
	has, err := p.engine.ID(id).Get(q)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, persistence.ErrNoRows
	}
	return q, nil
}

func (p *Persister) FindQuestionOptions(_ context.Context, qid uint64) ([]*question.Option, error) {
	res := make([]*question.Option, 0)
	err := p.engine.Where("question_id = ?", qid).Find(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Persister) FindListQuestions(_ context.Context, eid uint64, page, pageSize int) (int64, []*question.Question, error) {
	sess := p.engine.Where("exam_id = ?", eid)
	p.setSessionPagination(sess, page, pageSize)
	res := make([]*question.Question, 0, pageSize)
	count, err := sess.FindAndCount(&res)
	if err != nil {
		return 0, nil, err
	}
	return count, res, nil
}

func (p *Persister) CreateQuestion(_ context.Context, q *question.Question) error {
	sess := p.engine.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	if _, err := sess.Insert(q); err != nil {
		return err
	}

	if len(q.Options) > 0 {
		if _, err := sess.Insert(q.Options); err != nil {
			return err
		}
	}

	return sess.Commit()
}

func (p *Persister) UpdateQuestion(_ context.Context, q *question.Question, fields ...string) error {
	sess := p.engine.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	if _, err := sess.Where("question_id = ?", q.ID).Delete(new(question.Option)); err != nil {
		return err
	}

	if len(q.Options) > 0 {
		if _, err := sess.Insert(q.Options); err != nil {
			return err
		}
	}

	if _, err := p.engine.ID(q.ID).Cols(fields...).Update(q); err != nil {
		return err
	}

	return sess.Commit()
}
