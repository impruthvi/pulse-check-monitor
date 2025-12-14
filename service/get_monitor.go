package service

import (
	"context"
	"errors"
	"time"

	checker "github.com/impruthvi/pulse-check-apis/checker/v1"
	monitor "github.com/impruthvi/pulse-check-apis/monitor/v1"
)

func (s *service) GetMonitor(
	ctx context.Context,
	req *monitor.GetMonitorRequest,
) (*monitor.GetMonitorResponse, error) {

	monitorID := req.GetId()
	if monitorID == "" {
		return nil, errors.New("monitor id is required")
	}

	resMonitor, err := s.DBProvider.GetMonitor(monitorID)
	if err != nil {
		return nil, err
	}

	checkResp, err := s.CheckerClient.CheckURL(ctx, &checker.CheckURLRequest{
		MonitorId: resMonitor.ID,
		Url:       resMonitor.URL,
	})

	if err != nil {
		return nil, err
	}

	// 4. Update monitor with checker result
	now := time.Unix(checkResp.CheckedAt, 0)

	resMonitor.Status = checkResp.Status
	resMonitor.ResponseTimeMs = checkResp.ResponseTimeMs
	resMonitor.LastCheckedAt = &now

	if err := s.DBProvider.UpdateMonitor(resMonitor); err != nil {
		return nil, err
	}

	return &monitor.GetMonitorResponse{
		Monitor: resMonitor.AsApiMonitor(),
	}, nil
}
