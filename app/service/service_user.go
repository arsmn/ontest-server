package service

import "github.com/arsmn/ontest/app"

var _ app.App = new(Service)

func (s *Service) Register() error {
	panic("not implemented")
}
