package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"ticktick-backend/config"

	"github.com/redis/go-redis/v9"
)

// RedisService Redis服务
type RedisService struct {
	client *redis.Client
	config *config.RedisConfig
	ctx    context.Context
}

// NewRedisService 创建Redis服务实例
func NewRedisService(cfg *config.RedisConfig) (*RedisService, error) {
	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	ctx := context.Background()

	// 测试连接
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis连接失败: %w", err)
	}

	log.Println("Redis连接成功")

	return &RedisService{
		client: rdb,
		config: cfg,
		ctx:    ctx,
	}, nil
}

// Close 关闭Redis连接
func (r *RedisService) Close() error {
	return r.client.Close()
}

// GetClient 获取Redis客户端
func (r *RedisService) GetClient() *redis.Client {
	return r.client
}

// GetContext 获取上下文
func (r *RedisService) GetContext() context.Context {
	return r.ctx
}

// Ping 测试连接
func (r *RedisService) Ping() error {
	return r.client.Ping(r.ctx).Err()
}

// Set 设置键值对
func (r *RedisService) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// Get 获取值
func (r *RedisService) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

// Del 删除键
func (r *RedisService) Del(keys ...string) error {
	return r.client.Del(r.ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *RedisService) Exists(keys ...string) (int64, error) {
	return r.client.Exists(r.ctx, keys...).Result()
}

// TTL 获取键的剩余生存时间
func (r *RedisService) TTL(key string) (time.Duration, error) {
	return r.client.TTL(r.ctx, key).Result()
}

// Expire 设置键的过期时间
func (r *RedisService) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

// SAdd 向集合添加成员
func (r *RedisService) SAdd(key string, members ...interface{}) error {
	return r.client.SAdd(r.ctx, key, members...).Err()
}

// SRem 从集合移除成员
func (r *RedisService) SRem(key string, members ...interface{}) error {
	return r.client.SRem(r.ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func (r *RedisService) SMembers(key string) ([]string, error) {
	return r.client.SMembers(r.ctx, key).Result()
}

// SIsMember 检查成员是否在集合中
func (r *RedisService) SIsMember(key string, member interface{}) (bool, error) {
	return r.client.SIsMember(r.ctx, key, member).Result()
}

// SCard 获取集合成员数量
func (r *RedisService) SCard(key string) (int64, error) {
	return r.client.SCard(r.ctx, key).Result()
}

// HSet 设置哈希字段
func (r *RedisService) HSet(key string, values ...interface{}) error {
	return r.client.HSet(r.ctx, key, values...).Err()
}

// HGet 获取哈希字段值
func (r *RedisService) HGet(key, field string) (string, error) {
	return r.client.HGet(r.ctx, key, field).Result()
}

// HGetAll 获取哈希所有字段
func (r *RedisService) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(r.ctx, key).Result()
}

// HDel 删除哈希字段
func (r *RedisService) HDel(key string, fields ...string) error {
	return r.client.HDel(r.ctx, key, fields...).Err()
}

// HExists 检查哈希字段是否存在
func (r *RedisService) HExists(key, field string) (bool, error) {
	return r.client.HExists(r.ctx, key, field).Result()
}

// Keys 查找匹配模式的键
func (r *RedisService) Keys(pattern string) ([]string, error) {
	return r.client.Keys(r.ctx, pattern).Result()
}

// Pipeline 创建管道
func (r *RedisService) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

// ExecutePipeline 执行管道
func (r *RedisService) ExecutePipeline(pipe redis.Pipeliner) ([]redis.Cmder, error) {
	return pipe.Exec(r.ctx)
}

// Eval 执行Lua脚本
func (r *RedisService) Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(r.ctx, script, keys, args...).Result()
}

// GetConfig 获取Redis配置
func (r *RedisService) GetConfig() *config.RedisConfig {
	return r.config
}

// HealthCheck 健康检查
func (r *RedisService) HealthCheck() error {
	// 执行简单的ping命令
	if err := r.Ping(); err != nil {
		return fmt.Errorf("Redis健康检查失败: %w", err)
	}

	// 检查连接池状态
	stats := r.client.PoolStats()
	log.Printf("Redis连接池状态 - 总连接数: %d, 空闲连接数: %d, 过期连接数: %d", 
		stats.TotalConns, stats.IdleConns, stats.StaleConns)

	return nil
}
