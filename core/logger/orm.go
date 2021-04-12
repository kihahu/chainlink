package logger

import (
	"context"

	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type ORM interface {
	GetServiceLogLevel(serviceName string) (string, error)
	SetServiceLogLevel(ctx context.Context, serviceName string, level zapcore.Level) error
}

type Orm struct {
	DB *gorm.DB
}

// NewORM initializes a new ORM
func NewORM(db *gorm.DB) *Orm {
	return &Orm{db}
}

// GetServiceLogLevel returns the log level for a configured service
func (orm *Orm) GetServiceLogLevel(serviceName string) (string, error) {
	config := LogConfig{}
	if err := orm.DB.First(&config, "service_name = ?", serviceName).Error; err != nil {
		return "", err
	}
	return config.LogLevel, nil
}

func (orm *Orm) SetServiceLogLevel(ctx context.Context, serviceName string, level zapcore.Level) error {
	return orm.DB.WithContext(ctx).Where(LogConfig{ServiceName: serviceName}).
		Assign(LogConfig{ServiceName: serviceName, LogLevel: level.String()}).
		FirstOrCreate(&LogConfig{ServiceName: serviceName, LogLevel: level.String()}).Error
}
