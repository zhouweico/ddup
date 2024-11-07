package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port" yaml:"port" default:"8080"`
		Mode string `mapstructure:"mode" yaml:"mode" default:"development"`
	} `mapstructure:"server" yaml:"server"`

	Database struct {
		Driver   string `mapstructure:"driver" yaml:"driver" default:"postgres"`
		Host     string `mapstructure:"host" yaml:"host" default:"localhost"`
		Port     string `mapstructure:"port" yaml:"port" default:"5432"`
		User     string `mapstructure:"user" yaml:"user" default:"ddup"`
		Password string `mapstructure:"password" yaml:"password"`
		Name     string `mapstructure:"name" yaml:"name" default:"ddup"`
		Pool     struct {
			MaxOpen  int           `mapstructure:"max_open" yaml:"max_open" default:"50"`
			MaxIdle  int           `mapstructure:"max_idle" yaml:"max_idle" default:"10"`
			Lifetime time.Duration `mapstructure:"lifetime" yaml:"lifetime" default:"1h"`
			IdleTime time.Duration `mapstructure:"idle_time" yaml:"idle_time" default:"15m"`
		} `mapstructure:"pool" yaml:"pool"`
		Retry struct {
			Attempts int           `mapstructure:"attempts" yaml:"attempts" default:"3"`
			Interval time.Duration `mapstructure:"interval" yaml:"interval" default:"5s"`
		} `mapstructure:"retry" yaml:"retry"`
	} `mapstructure:"database" yaml:"database"`

	JWT struct {
		Secret    string        `mapstructure:"secret" yaml:"secret"`
		ExpiresIn time.Duration `mapstructure:"expires_in" yaml:"expires_in" default:"24h"`
	} `mapstructure:"jwt" yaml:"jwt"`

	HealthCheck struct {
		Interval time.Duration `mapstructure:"interval" yaml:"interval" default:"5m"`
	} `mapstructure:"health_check" yaml:"health_check"`

	Log struct {
		Level      zapcore.Level
		Filename   string
		MaxSize    int
		MaxBackups int
		MaxAge     int
		Compress   bool
	} `mapstructure:"log" yaml:"log"`

	Swagger struct {
		Host    string   `mapstructure:"host" yaml:"host" default:"localhost:8080"`
		Schemes []string `mapstructure:"schemes" yaml:"schemes" default:"[\"https\"]"`
	} `mapstructure:"swagger" yaml:"swagger"`
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
	config.Server.Port = viper.GetString("SERVER_PORT")
	config.Server.Mode = viper.GetString("SERVER_MODE")

	// 数据库配置
	config.Database.Driver = viper.GetString("DB_DRIVER")
	config.Database.Host = viper.GetString("DB_HOST")
	config.Database.Port = viper.GetString("DB_PORT")
	config.Database.User = viper.GetString("DB_USER")
	config.Database.Password = viper.GetString("DB_PASSWORD")
	config.Database.Name = viper.GetString("DB_NAME")

	// JWT配置
	config.JWT.Secret = viper.GetString("JWT_SECRET")
	config.JWT.ExpiresIn = viper.GetDuration("JWT_EXPIRES_IN")

	// 数据库连接池配置
	config.Database.Pool.MaxOpen = viper.GetInt("DB_POOL_MAX_OPEN")
	config.Database.Pool.MaxIdle = viper.GetInt("DB_POOL_MAX_IDLE")
	config.Database.Pool.Lifetime = viper.GetDuration("DB_POOL_LIFETIME")
	config.Database.Pool.IdleTime = viper.GetDuration("DB_POOL_IDLE_TIME")
	config.Database.Retry.Attempts = viper.GetInt("DB_RETRY_ATTEMPTS")
	config.Database.Retry.Interval = viper.GetDuration("DB_RETRY_INTERVAL")

	// 健康检查配置
	config.HealthCheck.Interval = viper.GetDuration("HEALTH_CHECK_INTERVAL")

	// 日志配置
	config.Log.Level = zapcore.Level(viper.GetInt("LOG_LEVEL"))
	config.Log.Filename = viper.GetString("LOG_FILENAME")
	config.Log.MaxSize = viper.GetInt("LOG_MAX_SIZE")
	config.Log.MaxBackups = viper.GetInt("LOG_MAX_BACKUPS")
	config.Log.MaxAge = viper.GetInt("LOG_MAX_AGE")
	config.Log.Compress = viper.GetBool("LOG_COMPRESS")

	// Swagger 配置
	config.Swagger.Host = viper.GetString("SWAGGER_HOST")
	config.Swagger.Schemes = viper.GetStringSlice("SWAGGER_SCHEMES")

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
		return errors.New("database host and port are required")
	}
	return nil
}
