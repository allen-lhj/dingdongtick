package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
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

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int

	// 连接池配置
	PoolSize     int
	MinIdleConns int
	MaxRetries   int

	// 超时配置
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// JWT相关配置
	AccessTokenTTL      time.Duration // AccessToken在Redis中的TTL
	RefreshTokenTTL     time.Duration // RefreshToken在Redis中的TTL
	BlacklistMinTTL     time.Duration // 黑名单最小TTL
	BlacklistMaxTTL     time.Duration // 黑名单最大TTL
	BlacklistDefaultTTL time.Duration // 黑名单默认TTL
	MaxSessionsPerUser  int           // 每用户最大会话数
}

// findProjectRoot 查找项目根目录（包含go.mod的目录）
func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	// 向上查找包含go.mod的目录
	for {
		if _, err := os.Stat(dir + "/go.mod"); err == nil {
			return dir
		}

		parent := dir + "/.."
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 查找项目根目录
	projectRoot := findProjectRoot()

	// 尝试从多个位置加载.env文件
	envPaths := []string{
		".env",                        // 当前目录
		"../.env",                     // 上级目录
		"backend/.env",                // backend子目录
		projectRoot + "/.env",         // 项目根目录
		projectRoot + "/backend/.env", // 项目根目录下的backend子目录
	}

	loaded := false
	var loadedPath string
	for _, path := range envPaths {
		if path == "" {
			continue // 跳过空路径
		}
		// 检查文件是否存在并尝试加载
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err == nil {
				loaded = true
				loadedPath = path
				break
			}
		}
	}

	if loaded {
		log.Printf("成功加载配置文件: %s", loadedPath)
	} else {
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
			AccessTokenDuration:  getEnvAsInt("JWT_ACCESS_DURATION", 24),  // 24小时
			RefreshTokenDuration: getEnvAsInt("JWT_REFRESH_DURATION", 30), // 30天
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),

			// 连接池配置
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
			MaxRetries:   getEnvAsInt("REDIS_MAX_RETRIES", 3),

			// 超时配置
			DialTimeout:  getEnvAsDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
			ReadTimeout:  getEnvAsDuration("REDIS_READ_TIMEOUT", 3*time.Second),
			WriteTimeout: getEnvAsDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),

			// JWT相关配置
			AccessTokenTTL:      getEnvAsDuration("REDIS_ACCESS_TOKEN_TTL", 2*time.Hour),
			RefreshTokenTTL:     getEnvAsDuration("REDIS_REFRESH_TOKEN_TTL", 30*24*time.Hour),
			BlacklistMinTTL:     getEnvAsDuration("REDIS_BLACKLIST_MIN_TTL", 5*time.Minute),
			BlacklistMaxTTL:     getEnvAsDuration("REDIS_BLACKLIST_MAX_TTL", 30*24*time.Hour),
			BlacklistDefaultTTL: getEnvAsDuration("REDIS_BLACKLIST_DEFAULT_TTL", 24*time.Hour),
			MaxSessionsPerUser:  getEnvAsInt("REDIS_MAX_SESSIONS_PER_USER", 5),
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
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsDuration 获取环境变量并转换为时间间隔
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
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
