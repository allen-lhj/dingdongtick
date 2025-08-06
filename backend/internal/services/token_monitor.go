package services

import (
	"log"
	"time"
)

// TokenMonitor Token监控服务
type TokenMonitor struct {
	tokenStore   *TokenStore
	redisService *RedisService
	stopChan     chan struct{}
	isRunning    bool
}

// NewTokenMonitor 创建Token监控服务
func NewTokenMonitor(tokenStore *TokenStore, redisService *RedisService) *TokenMonitor {
	return &TokenMonitor{
		tokenStore:   tokenStore,
		redisService: redisService,
		stopChan:     make(chan struct{}),
		isRunning:    false,
	}
}

// Start 启动监控服务
func (tm *TokenMonitor) Start() {
	if tm.isRunning {
		return
	}

	tm.isRunning = true
	log.Println("Token监控服务启动")

	// 启动清理任务
	go tm.startCleanupTask()

	// 启动健康检查任务
	go tm.startHealthCheckTask()

	// 启动统计任务
	go tm.startStatsTask()
}

// Stop 停止监控服务
func (tm *TokenMonitor) Stop() {
	if !tm.isRunning {
		return
	}

	tm.isRunning = false
	close(tm.stopChan)
	log.Println("Token监控服务停止")
}

// startCleanupTask 启动清理任务
func (tm *TokenMonitor) startCleanupTask() {
	ticker := time.NewTicker(1 * time.Hour) // 每小时执行一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tm.performCleanup()
		case <-tm.stopChan:
			return
		}
	}
}

// startHealthCheckTask 启动健康检查任务
func (tm *TokenMonitor) startHealthCheckTask() {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟检查一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tm.performHealthCheck()
		case <-tm.stopChan:
			return
		}
	}
}

// startStatsTask 启动统计任务
func (tm *TokenMonitor) startStatsTask() {
	ticker := time.NewTicker(10 * time.Minute) // 每10分钟统计一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tm.collectStats()
		case <-tm.stopChan:
			return
		}
	}
}

// performCleanup 执行清理任务
func (tm *TokenMonitor) performCleanup() {
	log.Println("开始执行Token清理任务")

	// 清理过期的Token
	if err := tm.tokenStore.CleanupExpiredTokens(); err != nil {
		log.Printf("清理过期Token失败: %v", err)
	}

	// 清理孤立的用户会话记录
	tm.cleanupOrphanedSessions()

	log.Println("Token清理任务完成")
}

// cleanupOrphanedSessions 清理孤立的用户会话记录
func (tm *TokenMonitor) cleanupOrphanedSessions() {
	// 获取所有用户会话键
	sessionKeys, err := tm.redisService.Keys("user_sessions:*")
	if err != nil {
		log.Printf("获取用户会话键失败: %v", err)
		return
	}

	cleanedCount := 0
	for _, sessionKey := range sessionKeys {
		// 获取会话中的TokenID列表
		tokenIDs, err := tm.redisService.SMembers(sessionKey)
		if err != nil {
			continue
		}

		// 检查每个TokenID对应的RefreshToken是否还存在
		validTokenIDs := make([]string, 0)
		for _, tokenID := range tokenIDs {
			// 从会话键中提取用户ID
			userID := extractUserIDFromSessionKey(sessionKey)
			if userID == "" {
				continue
			}

			// 检查RefreshToken是否存在
			refreshKey := "refresh_token:" + userID + ":" + tokenID
			exists, err := tm.redisService.Exists(refreshKey)
			if err != nil || exists == 0 {
				// RefreshToken不存在，从会话中移除
				tm.redisService.SRem(sessionKey, tokenID)
				cleanedCount++
			} else {
				validTokenIDs = append(validTokenIDs, tokenID)
			}
		}

		// 如果会话中没有有效的Token，删除整个会话记录
		if len(validTokenIDs) == 0 {
			tm.redisService.Del(sessionKey)
			cleanedCount++
		}
	}

	if cleanedCount > 0 {
		log.Printf("清理了 %d 个孤立的会话记录", cleanedCount)
	}
}

// extractUserIDFromSessionKey 从会话键中提取用户ID
func extractUserIDFromSessionKey(sessionKey string) string {
	// user_sessions:uuid -> uuid
	if len(sessionKey) > 14 { // len("user_sessions:") = 14
		return sessionKey[14:]
	}
	return ""
}

// performHealthCheck 执行健康检查
func (tm *TokenMonitor) performHealthCheck() {
	// 检查Redis连接
	if err := tm.redisService.HealthCheck(); err != nil {
		log.Printf("Redis健康检查失败: %v", err)
		return
	}

	// 检查Token存储服务状态
	stats, err := tm.tokenStore.GetStats()
	if err != nil {
		log.Printf("获取Token统计信息失败: %v", err)
		return
	}

	// 检查异常情况
	tm.detectAnomalies(stats)
}

// detectAnomalies 检测异常情况
func (tm *TokenMonitor) detectAnomalies(stats map[string]interface{}) {
	// 检查黑名单数量是否异常
	if blacklistCount, ok := stats["blacklist_count"].(int); ok {
		if blacklistCount > 10000 { // 阈值可配置
			log.Printf("警告: 黑名单数量异常 - %d", blacklistCount)
		}
	}

	// 检查RefreshToken数量是否异常
	if refreshTokenCount, ok := stats["refresh_token_count"].(int); ok {
		if refreshTokenCount > 50000 { // 阈值可配置
			log.Printf("警告: RefreshToken数量异常 - %d", refreshTokenCount)
		}
	}

	// 检查在线用户数量
	if onlineUsersCount, ok := stats["online_users_count"].(int); ok {
		log.Printf("当前在线用户数: %d", onlineUsersCount)
	}
}

// collectStats 收集统计信息
func (tm *TokenMonitor) collectStats() {
	stats, err := tm.tokenStore.GetStats()
	if err != nil {
		log.Printf("收集统计信息失败: %v", err)
		return
	}

	// 记录统计信息到日志
	log.Printf("Token统计信息: %+v", stats)

	// 这里可以将统计信息发送到监控系统
	// 例如: Prometheus, InfluxDB, CloudWatch等
	tm.sendStatsToMonitoring(stats)
}

// sendStatsToMonitoring 发送统计信息到监控系统
func (tm *TokenMonitor) sendStatsToMonitoring(stats map[string]interface{}) {
	// 实际项目中可以集成具体的监控系统
	// 这里只是示例

	// 示例: 发送到Prometheus
	// prometheus.GaugeVec.WithLabelValues("blacklist").Set(float64(stats["blacklist_count"]))

	// 示例: 发送到InfluxDB
	// influxClient.WritePoint(context.Background(), "token_stats", stats)

	// 目前只记录到日志
	log.Printf("统计信息已记录: %+v", stats)
}

// GetMonitoringData 获取监控数据
func (tm *TokenMonitor) GetMonitoringData() (map[string]interface{}, error) {
	data := make(map[string]interface{})

	// 获取基础统计信息
	stats, err := tm.tokenStore.GetStats()
	if err != nil {
		return nil, err
	}
	data["stats"] = stats

	// 获取Redis状态
	data["redis_status"] = "healthy"
	if err := tm.redisService.Ping(); err != nil {
		data["redis_status"] = "unhealthy"
		data["redis_error"] = err.Error()
	}

	// 获取服务运行状态
	data["monitor_running"] = tm.isRunning
	data["last_check"] = time.Now()

	return data, nil
}

// ForceCleanup 强制执行清理任务
func (tm *TokenMonitor) ForceCleanup() error {
	log.Println("强制执行清理任务")
	tm.performCleanup()
	return nil
}

// GetHealthStatus 获取健康状态
func (tm *TokenMonitor) GetHealthStatus() map[string]interface{} {
	status := make(map[string]interface{})

	// 检查Redis连接
	redisHealthy := true
	if err := tm.redisService.Ping(); err != nil {
		redisHealthy = false
		status["redis_error"] = err.Error()
	}
	status["redis_healthy"] = redisHealthy

	// 检查监控服务状态
	status["monitor_running"] = tm.isRunning

	// 整体健康状态
	status["overall_healthy"] = redisHealthy && tm.isRunning

	return status
}
