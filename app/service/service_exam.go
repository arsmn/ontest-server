package service

import (
	"context"

	"github.com/arsmn/ontest-server/app"
	"github.com/arsmn/ontest-server/exam"
	v "github.com/arsmn/ontest-server/module/validation"
)

var _ app.App = new(Service)

func (s *Service) CreateExam(ctx context.Context, req *exam.CreateExamRequest) (*exam.Exam, error) {
	if err := v.Validate(req); err != nil {
		return nil, err
	}

	exam := exam.NewDraftExam(req.SignedUser().ID).
		SetTitle(req.Title).
		SetStartAt(req.StartAt)

	if req.Deadline != nil {
		exam.SetDeadline(*req.Deadline)
	}

	if err := s.dx.Persister().CreateExam(ctx, exam); err != nil {
		return nil, err
	}

	return exam, nil
}
