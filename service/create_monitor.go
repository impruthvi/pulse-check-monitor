package service

import (
	"context"
	"errors"

	monitor "github.com/impruthvi/pulse-check-apis/monitor/v1"
	"github.com/impruthvi/pulse-check-monitor/db"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *service) CreateMonitor(ctx context.Context, req *monitor.CreateMonitorRequest) (*monitor.CreateMonitorResponse, error) {
	tracer := otel.Tracer("monitord")
	ctx, span := tracer.Start(ctx, "CreateMonitor")
	defer span.End()

	url := req.GetUrl()
	interval := req.GetIntervalSeconds()

	span.SetAttributes(
		attribute.String("monitor.url", url),
		attribute.Int64("monitor.interval_seconds", int64(interval)),
	)

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
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(attribute.String("monitor.id", resMonitor.ID))

	return &monitor.CreateMonitorResponse{
		Monitor: resMonitor.AsApiMonitor(),
	}, nil
}
