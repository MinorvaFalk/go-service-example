package service

import (
	"github.com/MinorvaFalk/go-service-example/datasource"
	servicea "github.com/MinorvaFalk/go-service-example/service/service-a"
	serviceb "github.com/MinorvaFalk/go-service-example/service/service-b"
	"github.com/MinorvaFalk/go-service-example/utils/logger"
)

type Services struct {
	ServiceA *servicea.ServiceA
	ServiceB *serviceb.ServiceB
}

func New(l *logger.Logger, db *datasource.DB) *Service {
	return &Service{
		l:  l,
		db: db,
	}
}

func (s *Service) NewServices() *Services {
	return &Services{
		ServiceA: servicea.NewServiceA(s.l, s.db),
		ServiceB: serviceb.NewServiceB(),
	}
}
