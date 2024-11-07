package driver

import (
	"ddup-apis/internal/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDriver struct{}

func (d *PostgresDriver) Name() string {
	return "postgres"
}

func (d *PostgresDriver) Open(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
