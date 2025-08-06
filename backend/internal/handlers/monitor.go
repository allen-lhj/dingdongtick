package handlers

import (
	"net/http"
	"ticktick-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MonitorHandler 监控处理器
type MonitorHandler struct {
	tokenMonitor *services.TokenMonitor
	tokenStore   *services.TokenStore
}

// NewMonitorHandler 创建监控处理器实例
func NewMonitorHandler(tokenMonitor *services.TokenMonitor, tokenStore *services.TokenStore) *MonitorHandler {
	return &MonitorHandler{
		tokenMonitor: tokenMonitor,
		tokenStore:   tokenStore,
	}
}

// GetHealth 获取健康状态
func (h *MonitorHandler) GetHealth(c *gin.Context) {
	status := h.tokenMonitor.GetHealthStatus()

	httpStatus := http.StatusOK
	if !status["overall_healthy"].(bool) {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, gin.H{
		"status": status,
	})
}

// GetStats 获取统计信息
func (h *MonitorHandler) GetStats(c *gin.Context) {
	stats, err := h.tokenStore.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// GetMonitoringData 获取监控数据
func (h *MonitorHandler) GetMonitoringData(c *gin.Context) {
	data, err := h.tokenMonitor.GetMonitoringData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取监控数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// ForceCleanup 强制执行清理任务
func (h *MonitorHandler) ForceCleanup(c *gin.Context) {
	// 这个接口应该有管理员权限验证
	// 简化处理，实际项目中需要验证管理员权限

	if err := h.tokenMonitor.ForceCleanup(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "执行清理任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "清理任务执行成功"})
}

// GetTokenInfo 获取Token详细信息（调试用）
func (h *MonitorHandler) GetTokenInfo(c *gin.Context) {
	tokenID := c.Param("tokenId")
	if tokenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少tokenId参数"})
		return
	}

	// 检查是否在黑名单中
	isBlacklisted, err := h.tokenStore.IsInBlacklist(tokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查黑名单失败"})
		return
	}

	info := map[string]interface{}{
		"tokenId":       tokenID,
		"isBlacklisted": isBlacklisted,
	}

	// 如果在黑名单中，获取黑名单信息
	if isBlacklisted {
		blacklistInfo, err := h.tokenStore.GetBlacklistInfo(tokenID)
		if err == nil && blacklistInfo != nil {
			info["blacklistInfo"] = blacklistInfo
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"tokenInfo": info,
	})
}

// GetUserSessions 获取指定用户的会话信息（管理员功能）
func (h *MonitorHandler) GetUserSessions(c *gin.Context) {
	// 这个接口应该有管理员权限验证
	userIDStr := c.Param("userId")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少userId参数"})
		return
	}

	// 解析用户ID
	userID, err := parseUUID(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID格式"})
		return
	}

	// 获取用户会话信息
	sessions, err := h.tokenStore.GetUserSessionsInfo(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户会话失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId":   userID,
		"sessions": sessions,
		"total":    len(sessions),
	})
}

// RevokeUserAllSessions 撤销指定用户的所有会话（管理员功能）
func (h *MonitorHandler) RevokeUserAllSessions(c *gin.Context) {
	// 这个接口应该有管理员权限验证
	userIDStr := c.Param("userId")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少userId参数"})
		return
	}

	// 解析用户ID
	userID, err := parseUUID(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID格式"})
		return
	}

	// 撤销用户所有Token
	if err := h.tokenStore.RevokeAllUserTokens(userID, "admin_revoke_all"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "撤销用户会话失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户所有会话已撤销",
		"userId":  userID,
	})
}

// GetSystemMetrics 获取系统指标
func (h *MonitorHandler) GetSystemMetrics(c *gin.Context) {
	// 获取基础统计信息
	stats, err := h.tokenStore.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计信息失败"})
		return
	}

	// 获取健康状态
	health := h.tokenMonitor.GetHealthStatus()

	// 组合系统指标
	metrics := map[string]interface{}{
		"stats":  stats,
		"health": health,
		"uptime": "运行中", // 实际项目中可以计算真实的运行时间
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
	})
}

// parseUUID 解析UUID字符串
func parseUUID(uuidStr string) (uuid.UUID, error) {
	return uuid.Parse(uuidStr)
}
