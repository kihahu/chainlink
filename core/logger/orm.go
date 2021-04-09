package logger

import (
	"gorm.io/gorm"
)

type ORM interface {
	GetServiceLogLevel(serviceName string) (string, error)
}

type orm struct {
	db *gorm.DB
}

// NewORM initializes a new ORM
func NewORM(db *gorm.DB) *orm {
	return &orm{db}
}

// GetServiceLogLevel returns the log level for a configured service
func (orm *orm) GetServiceLogLevel(serviceName string) (string, error) {
	config := LogConfig{}
	if err := orm.db.First(&config, "service_name = ?", serviceName).Error; err != nil {
		return "", err
	}
	return config.LogLevel, nil
}
