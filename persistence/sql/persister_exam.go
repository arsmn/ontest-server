package sql

import (
	"context"

	"github.com/arsmn/ontest-server/exam"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/question"
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

func (p *Persister) SearchExam(_ context.Context, q string, page, pageSize int) (int64, []*exam.Exam, error) {
	sess := p.engine.Where("state = ?", exam.Published)
	p.setSessionPagination(sess, page, pageSize)

	if len(q) > 0 {
		sess = sess.And("title LIKE '%" + q + "%'")
	}

	sess = sess.Desc("updated_at")

	res := make([]*exam.Exam, 0, pageSize)
	count, err := sess.FindAndCount(&res)
	if err != nil {
		return 0, nil, err
	}

	return count, res, nil
}

func (p *Persister) ExamStats(_ context.Context, id uint64) (*exam.ExamStatsResponse, error) {
	res, err := p.engine.Where("exam_id = ?", id).SumsInt(new(question.Question), "duration", "score", "negative_score")
	if err != nil {
		return nil, err
	}

	count, err := p.engine.Where("exam_id = ?", id).Count(new(question.Question))
	if err != nil {
		return nil, err
	}

	return &exam.ExamStatsResponse{
		TotalQuestions:      count,
		TotalDuration:       res[0],
		TotalScores:         res[1],
		TotalNegativeScores: res[2],
	}, nil
}

func (p *Persister) CreateResult(_ context.Context, r *exam.Result) error {
	sess := p.engine.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	if _, err := sess.Insert(r); err != nil {
		return err
	}

	if len(r.Answers) > 0 {
		if _, err := sess.Insert(r.Answers); err != nil {
			return err
		}
	}

	return sess.Commit()
}

func (p *Persister) FindResult(_ context.Context, id uint64) (*exam.Result, error) {
	s := new(exam.Result)
	has, err := p.engine.ID(id).Get(s)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, persistence.ErrNoRows
	}
	return s, nil
}

func (p *Persister) FindResultByExam(_ context.Context, uid, eid uint64) (*exam.Result, error) {
	s := new(exam.Result)
	has, err := p.engine.Where("exam_id = ?", eid).And("examinee = ?", uid).Get(s)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, persistence.ErrNoRows
	}
	return s, nil
}

func (p *Persister) FindResultAnswers(_ context.Context, rid uint64) ([]*exam.Answer, error) {
	res := make([]*exam.Answer, 0)
	err := p.engine.Where("result_id = ?", rid).Asc("id").Find(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Persister) UpdateAnswer(_ context.Context, a *exam.Answer, fields ...string) error {
	_, err := p.engine.ID(a.ID).Cols(fields...).Update(a)
	return err
}
