package service

import (
	checker "github.com/impruthvi/pulse-check-apis/checker/v1"
	monitor "github.com/impruthvi/pulse-check-apis/monitor/v1"

	"github.com/impruthvi/pulse-check-monitor/db"
)

type Dependencies struct {
	DBProvider    db.Provider
	CheckerClient checker.CheckerServiceClient
}

type Service interface {
	monitor.MonitorServiceServer
}

type service struct {
	Dependencies
}

func New(deps Dependencies) Service {
	return &service{
		Dependencies: deps,
	}
}
