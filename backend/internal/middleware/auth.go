package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"ticktick-backend/config"
	"ticktick-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID    uuid.UUID `json:"userId"`
	Email     string    `json:"email"`
	TokenID   string    `json:"jti"`       // JWT ID，用于撤销
	TokenType string    `json:"tokenType"` // "access" 或 "refresh"
	jwt.RegisteredClaims
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware(cfg *config.Config, tokenStore *services.TokenStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Cookie中获取token
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			// 如果Cookie中没有，尝试从Authorization头获取
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
				c.Abort()
				return
			}

			// 检查Bearer前缀
			if !strings.HasPrefix(authHeader, "Bearer ") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌格式"})
				c.Abort()
				return
			}

			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 解析和验证token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWT.SecretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
			c.Abort()
			return
		}

		// 验证token有效性
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证令牌已过期"})
			c.Abort()
			return
		}

		// 提取用户信息
		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌声明"})
			c.Abort()
			return
		}

		// 检查Token是否在黑名单中
		if claims.TokenID != "" {
			isBlacklisted, err := tokenStore.IsInBlacklist(claims.TokenID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "验证令牌状态失败"})
				c.Abort()
				return
			}
			if isBlacklisted {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "认证令牌已被撤销"})
				c.Abort()
				return
			}
		}

		// 将用户信息存储到上下文中
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("tokenID", claims.TokenID)
		c.Set("tokenType", claims.TokenType)

		c.Next()
	}
}

// GenerateTokens 生成访问令牌和刷新令牌
func GenerateTokens(cfg *config.Config, tokenStore *services.TokenStore, userID uuid.UUID, email string, deviceInfo string) (accessToken, refreshToken string, tokenID string, err error) {
	// 生成唯一的TokenID
	tokenID = uuid.New().String()
	now := time.Now()

	// 计算过期时间
	accessExpiresAt := now.Add(time.Duration(cfg.JWT.AccessTokenDuration) * time.Hour)
	refreshExpiresAt := now.Add(time.Duration(cfg.JWT.RefreshTokenDuration) * 24 * time.Hour)

	// 生成访问令牌
	accessClaims := &JWTClaims{
		UserID:    userID,
		Email:     email,
		TokenID:   tokenID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "ticktick-backend",
			Subject:   userID.String(),
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString([]byte(cfg.JWT.SecretKey))
	if err != nil {
		return "", "", "", err
	}

	// 生成刷新令牌
	refreshClaims := &JWTClaims{
		UserID:    userID,
		Email:     email,
		TokenID:   tokenID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "ticktick-backend",
			Subject:   userID.String(),
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(cfg.JWT.SecretKey))
	if err != nil {
		return "", "", "", err
	}

	// 存储RefreshToken信息到Redis
	refreshTokenInfo := &services.RefreshTokenInfo{
		UserID:     userID,
		Email:      email,
		TokenID:    tokenID,
		DeviceInfo: deviceInfo,
		CreatedAt:  now,
		LastUsedAt: now,
		ExpiresAt:  refreshExpiresAt,
	}

	if err := tokenStore.StoreRefreshToken(refreshTokenInfo); err != nil {
		return "", "", "", err
	}

	return accessToken, refreshToken, tokenID, nil
}

// SetTokenCookies 设置JWT令牌到Cookie
func SetTokenCookies(c *gin.Context, accessToken, refreshToken string, cfg *config.Config) {
	// 设置访问令牌Cookie
	c.SetCookie(
		"access_token",
		accessToken,
		cfg.JWT.AccessTokenDuration*3600, // 转换为秒
		"/",
		"",
		false, // 开发环境设为false，生产环境应设为true
		true,  // HttpOnly
	)

	// 设置刷新令牌Cookie
	c.SetCookie(
		"refresh_token",
		refreshToken,
		cfg.JWT.RefreshTokenDuration*24*3600, // 转换为秒
		"/",
		"",
		false, // 开发环境设为false，生产环境应设为true
		true,  // HttpOnly
	)
}

// ClearTokenCookies 清除JWT令牌Cookie
func ClearTokenCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
}

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, false
	}

	if id, ok := userID.(uuid.UUID); ok {
		return id, true
	}

	return uuid.Nil, false
}

// GetTokenIDFromContext 从上下文中获取TokenID
func GetTokenIDFromContext(c *gin.Context) (string, bool) {
	tokenID, exists := c.Get("tokenID")
	if !exists {
		return "", false
	}

	if id, ok := tokenID.(string); ok {
		return id, true
	}

	return "", false
}

// GetDeviceInfo 获取设备信息
func GetDeviceInfo(c *gin.Context) string {
	userAgent := c.GetHeader("User-Agent")
	clientIP := c.ClientIP()

	// 简化的设备信息，实际项目中可以更详细地解析User-Agent
	deviceInfo := "Unknown"
	if strings.Contains(userAgent, "Mobile") {
		deviceInfo = "Mobile"
	} else if strings.Contains(userAgent, "Chrome") {
		deviceInfo = "Chrome Browser"
	} else if strings.Contains(userAgent, "Firefox") {
		deviceInfo = "Firefox Browser"
	} else if strings.Contains(userAgent, "Safari") {
		deviceInfo = "Safari Browser"
	}

	return fmt.Sprintf("%s (%s)", deviceInfo, clientIP)
}

// RevokeToken 撤销Token
func RevokeToken(tokenStore *services.TokenStore, userID uuid.UUID, tokenID string, tokenType services.TokenType, expiresAt time.Time, reason string, deviceInfo string) error {
	blacklistInfo := &services.BlacklistInfo{
		UserID:     userID,
		TokenType:  tokenType,
		RevokedAt:  time.Now(),
		Reason:     reason,
		DeviceInfo: deviceInfo,
	}

	return tokenStore.AddToBlacklist(tokenID, blacklistInfo, expiresAt)
}
