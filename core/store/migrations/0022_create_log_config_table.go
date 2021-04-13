package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const up22 = `
	CREATE TYPE log_level AS ENUM (
		'debug',
		'info',
		'warn',
		'error',
		'panic'
	);
	
	CREATE TABLE log_configs (
		"id" BIGSERIAL PRIMARY KEY,
		"service_name" text NOT NULL,
		"log_level" log_level NOT NULL,
		"created_at" timestamp with time zone,
		"updated_at" timestamp with time zone
	);
`

const down22 = `
	DROP TABLE IF EXISTS log_conf;

	DROP TABLE IF EXISTS log_services;
`

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "0022_create_log_config_table",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up22).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down22).Error
		},
	})
}
