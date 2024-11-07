package driver

import (
	"ddup-apis/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteDriver struct{}

func (d *SQLiteDriver) Name() string {
	return "sqlite"
}

func (d *SQLiteDriver) Open(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
}
