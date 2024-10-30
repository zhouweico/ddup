package db

import (
	"fmt"

	"ddup-apis/internal/config"
	"ddup-apis/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	// 自动迁移表结构
	if err := DB.AutoMigrate(
		&model.User{},
		&model.UserSession{},
	); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	return nil
}
