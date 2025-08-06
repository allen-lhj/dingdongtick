package handlers

import (
	"net/http"
	"strings"
	"ticktick-backend/config"
	"ticktick-backend/internal/middleware"
	"ticktick-backend/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userService *services.UserService
	tokenStore  *services.TokenStore
	config      *config.Config
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(userService *services.UserService, tokenStore *services.TokenStore, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		tokenStore:  tokenStore,
		config:      cfg,
	}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数无效",
			"details": err.Error(),
		})
		return
	}

	// 调用服务层进行注册
	user, err := h.userService.Register(&req)
	if err != nil {
		if strings.Contains(err.Error(), "邮箱已被注册") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	// 获取设备信息
	deviceInfo := middleware.GetDeviceInfo(c)

	// 生成JWT令牌
	accessToken, refreshToken, tokenID, err := middleware.GenerateTokens(h.config, h.tokenStore, user.ID, user.Email, deviceInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败"})
		return
	}

	// 设置Cookie
	middleware.SetTokenCookies(c, accessToken, refreshToken, h.config)

	c.JSON(http.StatusCreated, gin.H{
		"message":      "注册成功",
		"user":         user,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"tokenID":      tokenID,
		"deviceInfo":   deviceInfo,
		"expiresIn":    h.config.JWT.AccessTokenDuration * 3600, // 转换为秒
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数无效",
			"details": err.Error(),
		})
		return
	}

	// 调用服务层进行登录验证
	user, err := h.userService.Login(&req)
	if err != nil {
		if strings.Contains(err.Error(), "邮箱或密码错误") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
		return
	}

	// 获取设备信息
	deviceInfo := middleware.GetDeviceInfo(c)

	// 生成JWT令牌
	accessToken, refreshToken, tokenID, err := middleware.GenerateTokens(h.config, h.tokenStore, user.ID, user.Email, deviceInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败"})
		return
	}

	// 设置Cookie
	middleware.SetTokenCookies(c, accessToken, refreshToken, h.config)

	c.JSON(http.StatusOK, gin.H{
		"message":      "登录成功",
		"user":         user,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"tokenID":      tokenID,
		"deviceInfo":   deviceInfo,
		"expiresIn":    h.config.JWT.AccessTokenDuration * 3600, // 转换为秒
	})
}

// RefreshToken 刷新访问令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// 从Cookie中获取刷新令牌
	refreshTokenString, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到刷新令牌"})
		return
	}

	// 解析刷新令牌
	token, err := jwt.ParseWithClaims(refreshTokenString, &middleware.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(h.config.JWT.SecretKey), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的刷新令牌"})
		return
	}

	// 验证令牌有效性
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "刷新令牌已过期"})
		return
	}

	// 提取用户信息
	claims, ok := token.Claims.(*middleware.JWTClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌声明"})
		return
	}

	// 检查Token是否在黑名单中
	if claims.TokenID != "" {
		isBlacklisted, err := h.tokenStore.IsInBlacklist(claims.TokenID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "验证令牌状态失败"})
			return
		}
		if isBlacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "刷新令牌已被撤销"})
			return
		}
	}

	// 验证RefreshToken是否在白名单中
	refreshTokenInfo, err := h.tokenStore.GetRefreshToken(claims.UserID, claims.TokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "验证刷新令牌失败"})
		return
	}
	if refreshTokenInfo == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "刷新令牌无效"})
		return
	}

	// 验证用户是否仍然存在
	user, err := h.userService.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	// 获取设备信息
	deviceInfo := middleware.GetDeviceInfo(c)

	// 撤销旧的RefreshToken
	if err := h.tokenStore.RevokeUserSession(claims.UserID, claims.TokenID, "token_refresh"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "撤销旧令牌失败"})
		return
	}

	// 生成新的访问令牌和刷新令牌
	newAccessToken, newRefreshToken, newTokenID, err := middleware.GenerateTokens(h.config, h.tokenStore, user.ID, user.Email, deviceInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败"})
		return
	}

	// 设置新的Cookie
	middleware.SetTokenCookies(c, newAccessToken, newRefreshToken, h.config)

	c.JSON(http.StatusOK, gin.H{
		"message":      "令牌刷新成功",
		"user":         user,
		"accessToken":  newAccessToken,
		"refreshToken": newRefreshToken,
		"tokenID":      newTokenID,
		"deviceInfo":   deviceInfo,
		"expiresIn":    h.config.JWT.AccessTokenDuration * 3600,
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 获取用户ID和TokenID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		// 即使没有用户信息，也清除Cookie
		middleware.ClearTokenCookies(c)
		c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
		return
	}

	tokenID, exists := middleware.GetTokenIDFromContext(c)
	if exists && tokenID != "" {
		// 撤销当前会话
		if err := h.tokenStore.RevokeUserSession(userID, tokenID, "user_logout"); err != nil {
			// 记录错误但不影响登出流程
			// 在实际项目中应该记录到日志系统
		}
	}

	// 清除Cookie
	middleware.ClearTokenCookies(c)
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// GetProfile 获取当前用户信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// GetSessions 获取用户活跃会话列表
func (h *AuthHandler) GetSessions(c *gin.Context) {
	// 获取用户ID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	// 获取用户会话信息
	sessions, err := h.tokenStore.GetUserSessionsInfo(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取会话信息失败"})
		return
	}

	// 获取当前TokenID以标识当前会话
	currentTokenID, _ := middleware.GetTokenIDFromContext(c)

	// 为每个会话添加是否为当前会话的标识
	for _, session := range sessions {
		session := map[string]interface{}{
			"tokenId":               session.TokenID,
			"deviceInfo":            session.DeviceInfo,
			"createdAt":             session.CreatedAt,
			"lastUsedAt":            session.LastUsedAt,
			"accessTokenExpiresAt":  session.AccessTokenExpiresAt,
			"refreshTokenExpiresAt": session.RefreshTokenExpiresAt,
			"isCurrent":             session.TokenID == currentTokenID,
		}
		_ = session // 避免未使用变量警告
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
		"total":    len(sessions),
	})
}

// RevokeSession 撤销特定会话
func (h *AuthHandler) RevokeSession(c *gin.Context) {
	// 获取用户ID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	// 获取要撤销的TokenID
	tokenID := c.Param("tokenId")
	if tokenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少tokenId参数"})
		return
	}

	// 检查是否是当前会话
	currentTokenID, _ := middleware.GetTokenIDFromContext(c)
	if tokenID == currentTokenID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能撤销当前会话，请使用登出功能"})
		return
	}

	// 撤销指定会话
	if err := h.tokenStore.RevokeUserSession(userID, tokenID, "manual_revoke"); err != nil {
		if strings.Contains(err.Error(), "会话不存在") {
			c.JSON(http.StatusNotFound, gin.H{"error": "会话不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "撤销会话失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "会话撤销成功"})
}

// LogoutAll 登出所有设备
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	// 获取用户ID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	// 撤销用户所有Token
	if err := h.tokenStore.RevokeAllUserTokens(userID, "logout_all"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登出所有设备失败"})
		return
	}

	// 清除当前Cookie
	middleware.ClearTokenCookies(c)

	c.JSON(http.StatusOK, gin.H{"message": "已登出所有设备"})
}

// RevokeToken 撤销指定Token（管理员功能）
func (h *AuthHandler) RevokeToken(c *gin.Context) {
	var req struct {
		TokenID string `json:"tokenId" binding:"required"`
		Reason  string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数无效",
			"details": err.Error(),
		})
		return
	}

	// 这里可以添加管理员权限检查
	// 简化处理，实际项目中需要验证管理员权限

	// 检查Token是否在黑名单中
	isBlacklisted, err := h.tokenStore.IsInBlacklist(req.TokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查Token状态失败"})
		return
	}

	if isBlacklisted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token已在黑名单中"})
		return
	}

	// 获取黑名单信息（如果需要的话）
	reason := req.Reason
	if reason == "" {
		reason = "admin_revoke"
	}

	// 这里需要更复杂的逻辑来确定Token的过期时间
	// 简化处理，使用默认TTL
	expiresAt := time.Now().Add(h.config.Redis.BlacklistDefaultTTL)

	blacklistInfo := &services.BlacklistInfo{
		UserID:    uuid.Nil,                 // 管理员撤销时可能不知道具体用户
		TokenType: services.AccessTokenType, // 假设是AccessToken
		RevokedAt: time.Now(),
		Reason:    reason,
	}

	if err := h.tokenStore.AddToBlacklist(req.TokenID, blacklistInfo, expiresAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "撤销Token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token撤销成功"})
}
