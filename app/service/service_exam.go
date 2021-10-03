package service

import (
	"context"
	stderr "errors"
	"fmt"
	"time"

	"github.com/arsmn/ontest-server/app"
	"github.com/arsmn/ontest-server/exam"
	"github.com/arsmn/ontest-server/module/cache"
	c "github.com/arsmn/ontest-server/module/context"
	"github.com/arsmn/ontest-server/module/errors"
	"github.com/arsmn/ontest-server/module/generate"
	v "github.com/arsmn/ontest-server/module/validation"
	"github.com/arsmn/ontest-server/module/xlog"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/shared"
	"github.com/arsmn/ontest-server/user"
)

var _ app.App = new(Service)

func (s *Service) GetExam(ctx context.Context, id uint64) (*exam.Exam, error) {
	return s.dx.Persister().FindExam(ctx, id)
}

func (s *Service) GetExamStats(ctx context.Context, id uint64) (*exam.ExamStatsResponse, error) {
	return s.dx.Persister().ExamStats(ctx, id)
}

func (s *Service) CreateExam(ctx context.Context, req *exam.CreateExamRequest) (*exam.Exam, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	u := c.User(ctx)

	exam := exam.NewDraftExam(u.ID).
		SetTitle(req.Title).
		SetStartAt(req.StartAt).
		SetDeadline(req.Deadline)

	if err := s.dx.Persister().CreateExam(ctx, exam); err != nil {
		return nil, err
	}

	return exam, nil
}

func (s *Service) UpdateExam(ctx context.Context, req *exam.UpdateExamRequest) error {
	if err := v.Validate(req); err != nil {
		return err
	}

	e := c.Exam(ctx)

	if e.State == exam.Published {
		return nil
	}

	e.SetTitle(req.Title).
		SetStartAt(req.StartAt).
		SetDeadline(req.Deadline).
		SetFreeMovement(req.FreeMovement)

	if err := s.dx.Persister().UpdateExam(ctx, e, "title", "start_at", "deadline", "free_movement"); err != nil {
		return err
	}

	return nil
}

func (s *Service) PublishExam(ctx context.Context, e *exam.Exam) error {
	if e.State == exam.Published {
		return nil
	}

	if e.StartAt.Before(time.Now().UTC()) {
		return errors.ErrBadRequest
	}

	if e.Deadline != nil {
		if e.Deadline.Before(e.StartAt) {
			return errors.ErrBadRequest
		}
	}

	e.SetState(exam.Published)

	if err := s.dx.Persister().UpdateExam(ctx, e, "state"); err != nil {
		return err
	}

	return nil
}

func (s *Service) SearchExam(ctx context.Context, req *exam.SearchExamRequest) (*exam.SearchExamResponse, error) {
	count, res, err := s.dx.Persister().SearchExam(ctx, req.Query, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &exam.SearchExamResponse{
		Exams:             res,
		PaginatedResponse: shared.NewPaginatedResponse(int(count), req.PageSize, req.Page),
	}, nil
}

func (s *Service) Participate(ctx context.Context, u *user.User, e *exam.Exam) (*exam.Result, error) {
	if e.StartAt.After(time.Now().UTC()) {
		return nil, errors.ErrBadRequest.WithError("Exam has not started!")
	}

	ress, err := s.dx.Persister().FindResultByExam(ctx, u.ID, e.ID)
	if err != nil {
		if !stderr.Is(err, persistence.ErrNoRows) {
			return nil, err
		}
	} else {
		return ress, nil
	}

	_, questions, err := s.dx.Persister().FindListQuestions(ctx, e.ID, 1, 100)
	if err != nil {
		return nil, err
	}

	res := &exam.Result{
		ID:       generate.UID(),
		Examinee: u.ID,
		ExamID:   e.ID,
		Answers:  make([]*exam.Answer, 0),
	}

	for _, q := range questions {
		ans := &exam.Answer{
			ID:              generate.UID(),
			QuestionID:      q.ID,
			ResultID:        res.ID,
			SelectedOptions: []*exam.SelectedOption{},
		}
		opts, err := s.dx.Persister().FindQuestionOptions(ctx, q.ID)
		if err != nil {
			return nil, err
		}
		for _, opt := range opts {
			ans.SelectedOptions = append(ans.SelectedOptions, &exam.SelectedOption{
				OptionID: opt.ID,
				Answer:   false,
			})
		}
		res.Answers = append(res.Answers, ans)
	}

	if err := s.dx.Persister().CreateResult(ctx, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) GetResult(ctx context.Context, id uint64) (*exam.Result, error) {
	res, err := s.dx.Persister().FindResult(ctx, id)
	if err != nil {
		return nil, err
	}

	ans, err := s.dx.Persister().FindResultAnswers(ctx, id)
	if err != nil {
		return nil, err
	}

	res.Answers = ans

	var active uint64
	key := fmt.Sprintf("uaq_%d", c.User(ctx).ID)
	s.dx.Cacher().Get(ctx, key, &active)

	for _, a := range res.Answers {
		q, err := s.dx.Persister().FindQuestion(ctx, a.QuestionID)
		if err != nil {
			return nil, err
		}
		opts, err := s.dx.Persister().FindQuestionOptions(ctx, a.QuestionID)
		if err != nil {
			return nil, err
		}
		a.Active = a.ID == active
		a.Question = q.SetOptions(opts)
	}

	return res, nil
}

func (s *Service) SubmitAnswer(ctx context.Context, req *exam.SubmitAnswerRequest) error {
	if err := s.dx.Persister().UpdateAnswer(ctx, &exam.Answer{
		ID:              req.ID,
		Text:            req.Text,
		SelectedOptions: req.SelectedOptions,
	}); err != nil {
		return err
	}

	key := fmt.Sprintf("uaq_%d", c.User(ctx).ID)
	if err := s.dx.Cacher().Set(ctx, &cache.Item{
		Key:   key,
		Value: req.ID,
		TTL:   24 * time.Hour,
	}); err != nil {
		s.dx.Logger().Error("error while saving UAQ", xlog.Err(err))
	}

	return nil
}
