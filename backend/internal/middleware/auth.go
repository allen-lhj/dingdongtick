package middleware

import (
	"net/http"
	"strings"
	"time"

	"ticktick-backend/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
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

		// 将用户信息存储到上下文中
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// GenerateTokens 生成访问令牌和刷新令牌
func GenerateTokens(cfg *config.Config, userID uuid.UUID, email string) (accessToken, refreshToken string, err error) {
	// 生成访问令牌
	accessClaims := &JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWT.AccessTokenDuration) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ticktick-backend",
			Subject:   userID.String(),
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString([]byte(cfg.JWT.SecretKey))
	if err != nil {
		return "", "", err
	}

	// 生成刷新令牌
	refreshClaims := &JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWT.RefreshTokenDuration) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ticktick-backend",
			Subject:   userID.String(),
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(cfg.JWT.SecretKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
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

	id, ok := userID.(uuid.UUID)
	return id, ok
}
