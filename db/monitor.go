package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	monitor "github.com/impruthvi/pulse-check-apis/monitor/v1"
)

type Monitor struct {
	ID              string `gorm:"primaryKey"`
	URL             string
	IntervalSeconds int32
	Status          string
	LastCheckedAt   *time.Time
	ResponseTimeMs  int64
}

func (m *Monitor) AsApiMonitor() *monitor.Monitor {
	var lastCheckedAt int64
	if m.LastCheckedAt != nil {
		lastCheckedAt = m.LastCheckedAt.Unix()
	}
	return &monitor.Monitor{
		Id:              m.ID,
		Url:             m.URL,
		IntervalSeconds: m.IntervalSeconds,
		Status:          m.Status,
		LastCheckedAt:   lastCheckedAt,
		ResponseTimeMs:  m.ResponseTimeMs,
	}
}

func (m *Monitor) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.NewString()
	return
}

func (p *provider) CreateMonitor(m *Monitor) (*Monitor, error) {
	err := p.db.Create(m).Error
	return m, err
}

func (p *provider) GetMonitor(id string) (*Monitor, error) {
	var m Monitor
	err := p.db.Where("id = ?", id).First(&m).Error
	return &m, err
}

func (p *provider) UpdateMonitor(m *Monitor) error {
	return p.db.Save(m).Error
}
