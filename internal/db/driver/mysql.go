package driver

import (
	"ddup-apis/internal/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDriver struct{}

func (d *MySQLDriver) Name() string {
	return "mysql"
}

func (d *MySQLDriver) Open(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
