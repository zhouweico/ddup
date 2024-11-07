package db

import (
	"fmt"
	"time"

	"ddup-apis/internal/config"
	"ddup-apis/internal/db/driver"
	"ddup-apis/internal/model"

	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	factory := driver.NewFactory()
	factory.Register(&driver.PostgresDriver{})
	factory.Register(&driver.MySQLDriver{})
	factory.Register(&driver.SQLiteDriver{})

	driver, err := factory.Get(cfg.Database.Driver)
	if err != nil {
		return err
	}

	var db *gorm.DB
	// 实现重试机制
	for i := 0; i <= cfg.Database.Retry.Attempts; i++ {
		db, err = driver.Open(cfg)
		if err == nil {
			break
		}

		if i < cfg.Database.Retry.Attempts {
			time.Sleep(cfg.Database.Retry.Interval)
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
	sqlDB.SetMaxOpenConns(cfg.Database.Pool.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.Database.Pool.MaxIdle)
	sqlDB.SetConnMaxLifetime(cfg.Database.Pool.Lifetime)
	sqlDB.SetConnMaxIdleTime(cfg.Database.Pool.IdleTime)

	// 自动迁移表结构
	if err := db.AutoMigrate(
		&model.User{},
		&model.Session{},
		&model.Profile{},
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
