package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type Provider interface {
	CreateMonitor(monitor *Monitor) (*Monitor, error)
	GetMonitor(id string) (*Monitor, error)
	UpdateMonitor(m *Monitor) error
}

type provider struct {
	db *gorm.DB
}

func New(dbURL string) Provider {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		log.Fatalf("Failed to use tracing plugin: %v", err)
	}

	db.AutoMigrate(&Monitor{})

	return &provider{db}
}
