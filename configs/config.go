package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Address string
	Mode    string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	Secret    string
	ExpiresIn int64
}

func Load() (*Config, error) {
	// 1. 首先加载配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	config := &Config{}

	// 2. 从配置文件加载默认值
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 3. 使用环境变量覆盖配置
	// Server
	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		config.Server.Address = addr
	}
	if mode := os.Getenv("SERVER_MODE"); mode != "" {
		config.Server.Mode = mode
	}

	// Database
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if portInt, err := strconv.Atoi(port); err == nil {
			config.Database.Port = portInt
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Database.DBName = dbName
	}

	// JWT
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		config.JWT.Secret = secret
	}
	if expiresIn := os.Getenv("JWT_EXPIRES_IN"); expiresIn != "" {
		if exp, err := strconv.ParseInt(expiresIn, 10, 64); err == nil {
			config.JWT.ExpiresIn = exp
		}
	}

	return config, nil
}

// GetDSN 返回数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)
}

// 用于打印配置信息（去除敏感信息）
func (c *Config) String() string {
	return fmt.Sprintf(`
Server:
  Address: %s
  Mode: %s
Database:
  Host: %s
  Port: %d
  User: %s
  DBName: %s
JWT:
  ExpiresIn: %d
`,
		c.Server.Address,
		c.Server.Mode,
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.DBName,
		c.JWT.ExpiresIn,
	)
}
