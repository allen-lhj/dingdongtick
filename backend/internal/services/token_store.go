package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// TokenType Token类型
type TokenType string

const (
	AccessTokenType  TokenType = "access"
	RefreshTokenType TokenType = "refresh"
)

// RefreshTokenInfo RefreshToken信息
type RefreshTokenInfo struct {
	UserID     uuid.UUID `json:"userId"`
	Email      string    `json:"email"`
	TokenID    string    `json:"tokenId"`
	DeviceInfo string    `json:"deviceInfo"`
	CreatedAt  time.Time `json:"createdAt"`
	LastUsedAt time.Time `json:"lastUsedAt"`
	ExpiresAt  time.Time `json:"expiresAt"`
}

// BlacklistInfo 黑名单信息
type BlacklistInfo struct {
	UserID     uuid.UUID `json:"userId"`
	TokenType  TokenType `json:"tokenType"`
	RevokedAt  time.Time `json:"revokedAt"`
	Reason     string    `json:"reason"`
	DeviceInfo string    `json:"deviceInfo,omitempty"`
}

// SessionInfo 会话信息
type SessionInfo struct {
	TokenID               string    `json:"tokenId"`
	DeviceInfo            string    `json:"deviceInfo"`
	CreatedAt             time.Time `json:"createdAt"`
	LastUsedAt            time.Time `json:"lastUsedAt"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
}

// UserOnlineInfo 用户在线信息
type UserOnlineInfo struct {
	LastActivity   time.Time `json:"lastActivity"`
	ActiveTokens   int       `json:"activeTokens"`
	DeviceCount    int       `json:"deviceCount"`
	LastDeviceInfo string    `json:"lastDeviceInfo"`
}

// TokenStore Token存储服务
type TokenStore struct {
	redis *RedisService
}

// NewTokenStore 创建Token存储服务
func NewTokenStore(redisService *RedisService) *TokenStore {
	return &TokenStore{
		redis: redisService,
	}
}

// Redis键名生成函数
func (ts *TokenStore) refreshTokenKey(userID uuid.UUID, tokenID string) string {
	return fmt.Sprintf("refresh_token:%s:%s", userID.String(), tokenID)
}

func (ts *TokenStore) blacklistKey(tokenID string) string {
	return fmt.Sprintf("blacklist:%s", tokenID)
}

func (ts *TokenStore) userSessionsKey(userID uuid.UUID) string {
	return fmt.Sprintf("user_sessions:%s", userID.String())
}

func (ts *TokenStore) userOnlineKey(userID uuid.UUID) string {
	return fmt.Sprintf("user_online:%s", userID.String())
}

// StoreRefreshToken 存储RefreshToken
func (ts *TokenStore) StoreRefreshToken(info *RefreshTokenInfo) error {
	key := ts.refreshTokenKey(info.UserID, info.TokenID)

	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("序列化RefreshToken信息失败: %w", err)
	}

	// 设置TTL为RefreshToken的剩余有效期
	ttl := time.Until(info.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("RefreshToken已过期")
	}

	if err := ts.redis.Set(key, data, ttl); err != nil {
		return fmt.Errorf("存储RefreshToken失败: %w", err)
	}

	// 添加到用户会话列表
	if err := ts.addToUserSessions(info.UserID, info.TokenID); err != nil {
		return fmt.Errorf("添加到用户会话失败: %w", err)
	}

	// 更新用户在线状态
	if err := ts.updateUserOnlineStatus(info.UserID, info.DeviceInfo); err != nil {
		return fmt.Errorf("更新用户在线状态失败: %w", err)
	}

	return nil
}

// GetRefreshToken 获取RefreshToken信息
func (ts *TokenStore) GetRefreshToken(userID uuid.UUID, tokenID string) (*RefreshTokenInfo, error) {
	key := ts.refreshTokenKey(userID, tokenID)

	data, err := ts.redis.Get(key)
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Token不存在
		}
		return nil, fmt.Errorf("获取RefreshToken失败: %w", err)
	}

	var info RefreshTokenInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return nil, fmt.Errorf("反序列化RefreshToken信息失败: %w", err)
	}

	return &info, nil
}

// DeleteRefreshToken 删除RefreshToken
func (ts *TokenStore) DeleteRefreshToken(userID uuid.UUID, tokenID string) error {
	key := ts.refreshTokenKey(userID, tokenID)

	if err := ts.redis.Del(key); err != nil {
		return fmt.Errorf("删除RefreshToken失败: %w", err)
	}

	// 从用户会话列表中移除
	if err := ts.removeFromUserSessions(userID, tokenID); err != nil {
		return fmt.Errorf("从用户会话移除失败: %w", err)
	}

	return nil
}

// AddToBlacklist 添加Token到黑名单
func (ts *TokenStore) AddToBlacklist(tokenID string, info *BlacklistInfo, expiresAt time.Time) error {
	key := ts.blacklistKey(tokenID)

	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("序列化黑名单信息失败: %w", err)
	}

	// 计算TTL
	ttl := ts.calculateBlacklistTTL(expiresAt)
	if ttl <= 0 {
		return nil // Token已过期，无需加入黑名单
	}

	if err := ts.redis.Set(key, data, ttl); err != nil {
		return fmt.Errorf("添加到黑名单失败: %w", err)
	}

	return nil
}

// IsInBlacklist 检查Token是否在黑名单中
func (ts *TokenStore) IsInBlacklist(tokenID string) (bool, error) {
	key := ts.blacklistKey(tokenID)

	exists, err := ts.redis.Exists(key)
	if err != nil {
		return false, fmt.Errorf("检查黑名单失败: %w", err)
	}

	return exists > 0, nil
}

// GetBlacklistInfo 获取黑名单信息
func (ts *TokenStore) GetBlacklistInfo(tokenID string) (*BlacklistInfo, error) {
	key := ts.blacklistKey(tokenID)

	data, err := ts.redis.Get(key)
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 不在黑名单中
		}
		return nil, fmt.Errorf("获取黑名单信息失败: %w", err)
	}

	var info BlacklistInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return nil, fmt.Errorf("反序列化黑名单信息失败: %w", err)
	}

	return &info, nil
}

// addToUserSessions 添加到用户会话列表
func (ts *TokenStore) addToUserSessions(userID uuid.UUID, tokenID string) error {
	key := ts.userSessionsKey(userID)

	if err := ts.redis.SAdd(key, tokenID); err != nil {
		return err
	}

	// 设置会话列表的TTL
	config := ts.redis.GetConfig()
	if err := ts.redis.Expire(key, config.RefreshTokenTTL); err != nil {
		return err
	}

	return nil
}

// removeFromUserSessions 从用户会话列表中移除
func (ts *TokenStore) removeFromUserSessions(userID uuid.UUID, tokenID string) error {
	key := ts.userSessionsKey(userID)
	return ts.redis.SRem(key, tokenID)
}

// GetUserSessions 获取用户所有会话
func (ts *TokenStore) GetUserSessions(userID uuid.UUID) ([]string, error) {
	key := ts.userSessionsKey(userID)
	return ts.redis.SMembers(key)
}

// updateUserOnlineStatus 更新用户在线状态
func (ts *TokenStore) updateUserOnlineStatus(userID uuid.UUID, deviceInfo string) error {
	key := ts.userOnlineKey(userID)

	// 获取当前活跃Token数量
	sessionsKey := ts.userSessionsKey(userID)
	activeTokens, err := ts.redis.SCard(sessionsKey)
	if err != nil {
		activeTokens = 0
	}

	info := UserOnlineInfo{
		LastActivity:   time.Now(),
		ActiveTokens:   int(activeTokens),
		DeviceCount:    1, // 简化处理，实际可以统计不同设备
		LastDeviceInfo: deviceInfo,
	}

	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("序列化用户在线信息失败: %w", err)
	}

	// 设置1小时TTL
	if err := ts.redis.Set(key, data, time.Hour); err != nil {
		return fmt.Errorf("更新用户在线状态失败: %w", err)
	}

	return nil
}

// GetUserOnlineStatus 获取用户在线状态
func (ts *TokenStore) GetUserOnlineStatus(userID uuid.UUID) (*UserOnlineInfo, error) {
	key := ts.userOnlineKey(userID)

	data, err := ts.redis.Get(key)
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 用户不在线
		}
		return nil, fmt.Errorf("获取用户在线状态失败: %w", err)
	}

	var info UserOnlineInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return nil, fmt.Errorf("反序列化用户在线信息失败: %w", err)
	}

	return &info, nil
}

// calculateBlacklistTTL 计算黑名单TTL
func (ts *TokenStore) calculateBlacklistTTL(expiresAt time.Time) time.Duration {
	config := ts.redis.GetConfig()
	remaining := time.Until(expiresAt)

	// 如果已经过期，返回0
	if remaining <= 0 {
		return 0
	}

	// 应用最小TTL限制
	if remaining < config.BlacklistMinTTL {
		return config.BlacklistMinTTL
	}

	// 应用最大TTL限制
	if remaining > config.BlacklistMaxTTL {
		return config.BlacklistMaxTTL
	}

	return remaining
}

// RevokeAllUserTokens 撤销用户所有Token
func (ts *TokenStore) RevokeAllUserTokens(userID uuid.UUID, reason string) error {
	// 获取用户所有会话
	sessions, err := ts.GetUserSessions(userID)
	if err != nil {
		return fmt.Errorf("获取用户会话失败: %w", err)
	}

	// 批量撤销所有Token
	pipe := ts.redis.Pipeline()

	for _, tokenID := range sessions {
		// 获取RefreshToken信息以确定过期时间
		refreshInfo, err := ts.GetRefreshToken(userID, tokenID)
		if err != nil || refreshInfo == nil {
			continue
		}

		// 添加RefreshToken到黑名单
		blacklistInfo := &BlacklistInfo{
			UserID:     userID,
			TokenType:  RefreshTokenType,
			RevokedAt:  time.Now(),
			Reason:     reason,
			DeviceInfo: refreshInfo.DeviceInfo,
		}

		blacklistKey := ts.blacklistKey(tokenID)
		data, _ := json.Marshal(blacklistInfo)
		ttl := ts.calculateBlacklistTTL(refreshInfo.ExpiresAt)
		if ttl > 0 {
			pipe.Set(ts.redis.GetContext(), blacklistKey, data, ttl)
		}

		// 删除RefreshToken
		refreshKey := ts.refreshTokenKey(userID, tokenID)
		pipe.Del(ts.redis.GetContext(), refreshKey)
	}

	// 清空用户会话列表
	sessionsKey := ts.userSessionsKey(userID)
	pipe.Del(ts.redis.GetContext(), sessionsKey)

	// 执行批量操作
	if _, err := ts.redis.ExecutePipeline(pipe); err != nil {
		return fmt.Errorf("批量撤销Token失败: %w", err)
	}

	return nil
}

// RevokeUserSession 撤销用户特定会话
func (ts *TokenStore) RevokeUserSession(userID uuid.UUID, tokenID string, reason string) error {
	// 获取RefreshToken信息
	refreshInfo, err := ts.GetRefreshToken(userID, tokenID)
	if err != nil {
		return fmt.Errorf("获取RefreshToken信息失败: %w", err)
	}
	if refreshInfo == nil {
		return fmt.Errorf("会话不存在")
	}

	// 添加到黑名单
	blacklistInfo := &BlacklistInfo{
		UserID:     userID,
		TokenType:  RefreshTokenType,
		RevokedAt:  time.Now(),
		Reason:     reason,
		DeviceInfo: refreshInfo.DeviceInfo,
	}

	if err := ts.AddToBlacklist(tokenID, blacklistInfo, refreshInfo.ExpiresAt); err != nil {
		return fmt.Errorf("添加到黑名单失败: %w", err)
	}

	// 删除RefreshToken
	if err := ts.DeleteRefreshToken(userID, tokenID); err != nil {
		return fmt.Errorf("删除RefreshToken失败: %w", err)
	}

	return nil
}

// UpdateRefreshTokenLastUsed 更新RefreshToken最后使用时间
func (ts *TokenStore) UpdateRefreshTokenLastUsed(userID uuid.UUID, tokenID string) error {
	// 获取当前信息
	info, err := ts.GetRefreshToken(userID, tokenID)
	if err != nil {
		return fmt.Errorf("获取RefreshToken信息失败: %w", err)
	}
	if info == nil {
		return fmt.Errorf("RefreshToken不存在")
	}

	// 更新最后使用时间
	info.LastUsedAt = time.Now()

	// 重新存储
	key := ts.refreshTokenKey(userID, tokenID)
	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("序列化RefreshToken信息失败: %w", err)
	}

	// 保持原有的TTL
	ttl, err := ts.redis.TTL(key)
	if err != nil {
		return fmt.Errorf("获取TTL失败: %w", err)
	}

	if err := ts.redis.Set(key, data, ttl); err != nil {
		return fmt.Errorf("更新RefreshToken失败: %w", err)
	}

	// 更新用户在线状态
	if err := ts.updateUserOnlineStatus(userID, info.DeviceInfo); err != nil {
		return fmt.Errorf("更新用户在线状态失败: %w", err)
	}

	return nil
}

// GetUserSessionsInfo 获取用户会话详细信息
func (ts *TokenStore) GetUserSessionsInfo(userID uuid.UUID) ([]*SessionInfo, error) {
	sessions, err := ts.GetUserSessions(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户会话失败: %w", err)
	}

	var sessionsInfo []*SessionInfo
	for _, tokenID := range sessions {
		refreshInfo, err := ts.GetRefreshToken(userID, tokenID)
		if err != nil || refreshInfo == nil {
			continue
		}

		sessionInfo := &SessionInfo{
			TokenID:               tokenID,
			DeviceInfo:            refreshInfo.DeviceInfo,
			CreatedAt:             refreshInfo.CreatedAt,
			LastUsedAt:            refreshInfo.LastUsedAt,
			RefreshTokenExpiresAt: refreshInfo.ExpiresAt,
			// AccessToken过期时间需要根据配置计算
			AccessTokenExpiresAt: refreshInfo.LastUsedAt.Add(ts.redis.GetConfig().AccessTokenTTL),
		}

		sessionsInfo = append(sessionsInfo, sessionInfo)
	}

	return sessionsInfo, nil
}

// CleanupExpiredTokens 清理过期Token
func (ts *TokenStore) CleanupExpiredTokens() error {
	// 这个方法主要用于清理一些边缘情况下的过期数据
	// Redis的TTL机制会自动清理大部分过期数据

	// 清理没有设置TTL的异常黑名单记录
	script := `
		local keys = redis.call('KEYS', 'blacklist:*')
		local cleaned = 0
		for i=1,#keys do
			local ttl = redis.call('TTL', keys[i])
			if ttl == -1 then
				redis.call('DEL', keys[i])
				cleaned = cleaned + 1
			end
		end
		return cleaned
	`

	result, err := ts.redis.Eval(script, []string{})
	if err != nil {
		return fmt.Errorf("清理过期黑名单失败: %w", err)
	}

	if cleaned, ok := result.(int64); ok && cleaned > 0 {
		fmt.Printf("清理了 %d 个异常黑名单记录\n", cleaned)
	}

	return nil
}

// GetStats 获取统计信息
func (ts *TokenStore) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 统计黑名单数量
	blacklistKeys, err := ts.redis.Keys("blacklist:*")
	if err == nil {
		stats["blacklist_count"] = len(blacklistKeys)
	}

	// 统计RefreshToken数量
	refreshTokenKeys, err := ts.redis.Keys("refresh_token:*")
	if err == nil {
		stats["refresh_token_count"] = len(refreshTokenKeys)
	}

	// 统计用户会话数量
	sessionKeys, err := ts.redis.Keys("user_sessions:*")
	if err == nil {
		stats["user_sessions_count"] = len(sessionKeys)
	}

	// 统计在线用户数量
	onlineKeys, err := ts.redis.Keys("user_online:*")
	if err == nil {
		stats["online_users_count"] = len(onlineKeys)
	}

	return stats, nil
}
