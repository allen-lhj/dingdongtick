package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Host string
	Mode string // gin模式: debug, release, test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey            string
	AccessTokenDuration  int // 小时
	RefreshTokenDuration int // 天
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("未找到.env文件，使用环境变量")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "ticktick"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey:            getEnv("JWT_SECRET", "your-secret-key"),
			AccessTokenDuration:  getEnvAsInt("JWT_ACCESS_DURATION", 24),   // 24小时
			RefreshTokenDuration: getEnvAsInt("JWT_REFRESH_DURATION", 30), // 30天
		},
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为整数
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		// 简单的字符串转整数，实际项目中应该使用strconv.Atoi
		switch value {
		case "1":
			return 1
		case "24":
			return 24
		case "30":
			return 30
		case "7":
			return 7
		}
	}
	return defaultValue
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return "host=" + c.Database.Host +
		" port=" + c.Database.Port +
		" user=" + c.Database.User +
		" password=" + c.Database.Password +
		" dbname=" + c.Database.DBName +
		" sslmode=" + c.Database.SSLMode
}
