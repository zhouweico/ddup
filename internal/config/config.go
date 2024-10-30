package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
		Mode string
	}
	Database struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
	}
	JWT struct {
		Secret    string
		ExpiresIn time.Duration
	}
}

var globalConfig Config

func GetConfig() *Config {
	if globalConfig.JWT.Secret == "" {
		log.Fatal("配置未初始化，请确保已调用 SetConfig()")
	}
	return &globalConfig
}

func SetConfig(cfg Config) {
	globalConfig = cfg
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config

	// 服务器配置
	config.Server.Port = viper.GetString("APP_PORT")
	config.Server.Mode = viper.GetString("APP_ENV")

	// 数据库配置
	config.Database.Host = viper.GetString("DB_HOST")
	config.Database.Port = viper.GetString("DB_PORT")
	config.Database.Name = viper.GetString("DB_NAME")
	config.Database.User = viper.GetString("DB_USER")
	config.Database.Password = viper.GetString("DB_PASSWORD")

	// JWT配置
	config.JWT.Secret = viper.GetString("JWT_SECRET")

	// 将 JWT_EXPIRES_IN 从秒转换为 time.Duration
	config.JWT.ExpiresIn = time.Duration(viper.GetInt("JWT_EXPIRES_IN")) * time.Second

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(cfg *Config) error {
	if cfg.JWT.Secret == "" {
		return errors.New("JWT secret is required")
	}
	if cfg.JWT.ExpiresIn <= 0 {
		return errors.New("JWT expiration time must be positive")
	}
	if cfg.Database.Host == "" || cfg.Database.Port == "" {
		return errors.New("Database host and port are required")
	}
	return nil
}
