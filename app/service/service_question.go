package service

import (
	"context"
	"time"

	"github.com/arsmn/ontest-server/app"
	c "github.com/arsmn/ontest-server/module/context"
	v "github.com/arsmn/ontest-server/module/validation"
	"github.com/arsmn/ontest-server/question"
	"github.com/arsmn/ontest-server/shared"
)

var _ app.App = new(Service)

func (s *Service) GetQuestion(ctx context.Context, id uint64) (*question.Question, error) {
	return s.dx.Persister().FindQuestion(ctx, id)
}

func (s *Service) GetQuestionOptions(ctx context.Context, id uint64) ([]*question.Option, error) {
	return s.dx.Persister().FindQuestionOptions(ctx, id)
}

func (s *Service) GetQuestionList(ctx context.Context, req *question.GetQuestionListRequest) (*question.GetQuestionListResponse, error) {
	total, res, err := s.dx.Persister().FindListQuestions(ctx, req.ExamID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &question.GetQuestionListResponse{
		Questions:         res,
		PaginatedResponse: shared.NewPaginatedResponse(int(total), req.PageSize, req.Page),
	}, nil
}

func (s *Service) CreateQuestion(ctx context.Context, req *question.CreateQuestionRequest) (*question.Question, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	u := c.User(ctx)
	e := c.Exam(ctx)

	q := question.NewQuestion(u.ID, e.ID, req.Text).
		SetType(req.Type).
		SetScore(req.Score).
		SetNegativeScore(req.NegativeScore).
		SetDuration(time.Duration(req.Duration))

	if q.Type == question.SingleChoice || q.Type == question.MultipleChoice {
		for _, opt := range req.Options {
			q.AddOption(question.NewOption(q.ID, opt.Text, opt.Answer))
		}
	}

	if err := s.dx.Persister().CreateQuestion(ctx, q); err != nil {
		return nil, err
	}

	return q, nil
}

func (s *Service) UpdateQuestion(ctx context.Context, req *question.CreateQuestionRequest) error {
	if err := v.Validate(req); err != nil {
		return err
	}

	q := c.Question(ctx)

	q.SetText(req.Text).
		SetType(req.Type).
		SetScore(req.Score).
		SetNegativeScore(req.NegativeScore).
		SetDuration(time.Duration(req.Duration))

	if q.Type == question.Descriptive {
		q.SetOptions(nil)
	}

	if q.Type == question.SingleChoice || q.Type == question.MultipleChoice {
		for _, opt := range req.Options {
			q.AddOption(question.NewOption(q.ID, opt.Text, opt.Answer))
		}
	}

	return s.dx.Persister().UpdateQuestion(ctx, q, "text", "type", "score", "negative_score", "duration")
}

func (s *Service) DeleteQuestion(ctx context.Context, id uint64) error {
	return s.dx.Persister().RemoveQuestion(ctx, id)
}
