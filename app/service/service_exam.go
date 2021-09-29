package service

import (
	"context"

	"github.com/arsmn/ontest-server/app"
	"github.com/arsmn/ontest-server/exam"
	v "github.com/arsmn/ontest-server/module/validation"
)

var _ app.App = new(Service)

func (s *Service) GetExam(ctx context.Context, id uint64) (*exam.Exam, error) {
	return s.dx.Persister().FindExam(ctx, id)
}

func (s *Service) CreateExam(ctx context.Context, req *exam.CreateExamRequest) (*exam.Exam, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	exam := exam.NewDraftExam(req.SignedUser().ID).
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

	req.Exam().
		SetTitle(req.Title).
		SetStartAt(req.StartAt).
		SetDeadline(req.Deadline).
		SetFreeMovement(req.FreeMovement)

	if err := s.dx.Persister().UpdateExam(ctx, req.Exam(), "title", "start_at", "deadline", "free_movement"); err != nil {
		return err
	}

	return nil
}
