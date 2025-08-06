#!/bin/bash

# 基于Redis的JWT验证方案测试脚本

BASE_URL="http://localhost:8080/api/v1"

echo "=== 基于Redis的JWT验证方案测试 ==="
echo

# 1. 健康检查
echo "1. 健康检查"
curl -s "$BASE_URL/health" | jq '.'
echo

# 2. 用户注册
echo "2. 用户注册"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "firstName": "Test",
    "lastName": "User"
  }')

echo "$REGISTER_RESPONSE" | jq '.'

# 提取访问令牌和刷新令牌
ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.accessToken // empty')
REFRESH_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.refreshToken // empty')
TOKEN_ID=$(echo "$REGISTER_RESPONSE" | jq -r '.tokenID // empty')

if [ -z "$ACCESS_TOKEN" ]; then
    echo "注册失败，无法获取访问令牌"
    exit 1
fi

echo "获取到访问令牌: ${ACCESS_TOKEN:0:50}..."
echo "获取到刷新令牌: ${REFRESH_TOKEN:0:50}..."
echo "Token ID: $TOKEN_ID"
echo

# 3. 获取用户信息
echo "3. 获取用户信息"
curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/profile" | jq '.'
echo

# 4. 获取用户会话列表
echo "4. 获取用户会话列表"
curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/sessions" | jq '.'
echo

# 5. 获取监控统计信息
echo "5. 获取监控统计信息"
curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/monitor/stats" | jq '.'
echo

# 6. 获取Token信息
echo "6. 获取Token信息"
curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/monitor/token/$TOKEN_ID" | jq '.'
echo

# 7. 刷新Token
echo "7. 刷新Token"
REFRESH_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/refresh" \
  -H "Cookie: refresh_token=$REFRESH_TOKEN")

echo "$REFRESH_RESPONSE" | jq '.'

# 提取新的访问令牌
NEW_ACCESS_TOKEN=$(echo "$REFRESH_RESPONSE" | jq -r '.accessToken // empty')
NEW_TOKEN_ID=$(echo "$REFRESH_RESPONSE" | jq -r '.tokenID // empty')

if [ -n "$NEW_ACCESS_TOKEN" ]; then
    echo "获取到新的访问令牌: ${NEW_ACCESS_TOKEN:0:50}..."
    echo "新的Token ID: $NEW_TOKEN_ID"
    ACCESS_TOKEN="$NEW_ACCESS_TOKEN"
    TOKEN_ID="$NEW_TOKEN_ID"
fi
echo

# 8. 再次获取会话列表（应该看到新的会话）
echo "8. 再次获取会话列表"
curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/sessions" | jq '.'
echo

# 9. 获取系统指标
echo "9. 获取系统指标"
curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/monitor/metrics" | jq '.'
echo

# 10. 登出所有设备
echo "10. 登出所有设备"
curl -s -X POST -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/logout-all" | jq '.'
echo

# 11. 尝试使用已撤销的Token访问（应该失败）
echo "11. 尝试使用已撤销的Token访问（应该失败）"
curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/profile" | jq '.'
echo

# 12. 最终健康检查
echo "12. 最终健康检查"
curl -s "$BASE_URL/health" | jq '.'
echo

echo "=== 测试完成 ==="
