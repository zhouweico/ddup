package config

import (
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
		Secret             string
		ExpiresIn          time.Duration
		RefreshGracePeriod time.Duration
	}
}

var globalConfig Config

func GetConfig() *Config {
	if globalConfig.JWT.Secret == "" {
		log.Fatal("配置未初始化，请确保已调用 SetConfig()")
	}
	return &globalConfig
}

// 在程序初始化时设置配置
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
	config.Server.Port = viper.GetString("APP_PORT")
	config.Server.Mode = viper.GetString("APP_ENV")

	config.Database.Host = viper.GetString("DB_HOST")
	config.Database.Port = viper.GetString("DB_PORT")
	config.Database.Name = viper.GetString("DB_NAME")
	config.Database.User = viper.GetString("DB_USER")
	config.Database.Password = viper.GetString("DB_PASSWORD")

	config.JWT.Secret = viper.GetString("JWT_SECRET")

	// 将 JWT_EXPIRES_IN 从秒转换为 time.Duration
	viper.SetDefault("JWT_EXPIRES_IN", 86400)
	config.JWT.ExpiresIn = time.Duration(viper.GetInt("JWT_EXPIRES_IN")) * time.Second
	if gracePeriod, err := time.ParseDuration(viper.GetString("JWT_REFRESH_GRACE_PERIOD")); err == nil {
		config.JWT.RefreshGracePeriod = gracePeriod
	} else {
		config.JWT.RefreshGracePeriod = 6 * time.Hour // 默认值
	}

	return &config, nil
}
