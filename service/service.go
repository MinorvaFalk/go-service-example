package service

import (
	"github.com/MinorvaFalk/go-service-example/datasource"
	servicea "github.com/MinorvaFalk/go-service-example/service/service-a"
	serviceb "github.com/MinorvaFalk/go-service-example/service/service-b"
	"github.com/MinorvaFalk/go-service-example/utils/logger"
)

// Insert services global dependencies here
type Service struct {
	l  *logger.Logger
	db *datasource.DB
}

func NewService(l *logger.Logger, db *datasource.DB) *Service {
	return &Service{
		l:  l,
		db: db,
	}
}

func (s *Service) ServiceA() *servicea.ServiceA {
	return servicea.NewServiceA(s.l, s.db)
}

func (s *Service) ServiceB() *serviceb.ServiceB {
	return serviceb.NewServiceB()
}
