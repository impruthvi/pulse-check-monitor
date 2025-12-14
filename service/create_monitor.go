package service

import (
	"context"
	"errors"

	monitor "github.com/impruthvi/pulse-check-apis/monitor/v1"
	"github.com/impruthvi/pulse-check-monitor/db"
)

func (s *service) CreateMonitor(ctx context.Context, req *monitor.CreateMonitorRequest) (*monitor.CreateMonitorResponse, error) {
	url := req.GetUrl()
	interval := req.GetIntervalSeconds()

	if url == "" {
		return nil, errors.New("url is required")
	}

	if interval <= 0 {
		return nil, errors.New("interval_seconds must be greater than 0")
	}

	resMonitor, err := s.DBProvider.CreateMonitor(&db.Monitor{
		URL:             url,
		IntervalSeconds: interval,
		Status:          "PENDING",
	})
	if err != nil {
		return nil, err
	}

	return &monitor.CreateMonitorResponse{
		Monitor: resMonitor.AsApiMonitor(),
	}, nil
}
