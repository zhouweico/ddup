package db

import (
	"fmt"
	"time"

	"ddup-apis/internal/config"
	"ddup-apis/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	var db *gorm.DB
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	// 实现重试机制
	for i := 0; i <= cfg.Database.RetryTimes; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		if i < cfg.Database.RetryTimes {
			time.Sleep(cfg.Database.RetryDelay)
			continue
		}
		return fmt.Errorf("数据库连接失败，已重试 %d 次: %w", i, err)
	}

	// 获取底层的 *sql.DB 对象
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.MaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.Database.MaxIdleTime)

	// 自动迁移表结构
	if err := db.AutoMigrate(
		&model.User{},
		&model.Session{},
		&model.Social{},
		&model.Organization{},
		&model.OrganizationMember{},
	); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	DB = db
	return nil
}

// 添加健康检查方法
func Ping() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}
	return sqlDB.Ping()
}
