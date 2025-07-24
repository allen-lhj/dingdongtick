package handlers

import (
	"net/http"
	"strings"
	"ticktick-backend/config"
	"ticktick-backend/internal/middleware"
	"ticktick-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userService *services.UserService
	config      *config.Config
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(userService *services.UserService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userService: userService,
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

	// 生成JWT令牌
	accessToken, refreshToken, err := middleware.GenerateTokens(h.config, user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败"})
		return
	}

	// 设置Cookie
	middleware.SetTokenCookies(c, accessToken, refreshToken, h.config)

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user":    user,
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

	// 生成JWT令牌
	accessToken, refreshToken, err := middleware.GenerateTokens(h.config, user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败"})
		return
	}

	// 设置Cookie
	middleware.SetTokenCookies(c, accessToken, refreshToken, h.config)

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user":    user,
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

	// 验证用户是否仍然存在
	user, err := h.userService.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	// 生成新的访问令牌和刷新令牌
	newAccessToken, newRefreshToken, err := middleware.GenerateTokens(h.config, user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败"})
		return
	}

	// 设置新的Cookie
	middleware.SetTokenCookies(c, newAccessToken, newRefreshToken, h.config)

	c.JSON(http.StatusOK, gin.H{
		"message": "令牌刷新成功",
		"user":    user,
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
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
