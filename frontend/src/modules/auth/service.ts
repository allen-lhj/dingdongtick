import http from '@/services/http'
import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  RefreshTokenRequest,
  RefreshTokenResponse,
  ResetPasswordRequest,
  ResetPasswordConfirmRequest,
  ChangePasswordRequest,
  UpdateUserRequest,
  User,
} from './types'

/**
 * 认证服务类
 * 负责处理所有与认证相关的API调用
 */
class AuthService {
  /**
   * 用户注册
   * @param data 注册数据
   * @returns 认证响应
   */
  async register(data: RegisterRequest): Promise<AuthResponse> {
    try {
      const response = await http.post<AuthResponse>('/auth/register', data)

      // 存储认证信息
      this.storeAuthData(response)

      return response
    } catch (error) {
      console.error('Registration failed:', error)
      throw error
    }
  }

  /**
   * 用户登录
   * @param data 登录数据
   * @returns 认证响应
   */
  async login(data: LoginRequest): Promise<AuthResponse> {
    try {
      const response = await http.post<AuthResponse>('/auth/login', data)

      // 存储认证信息
      this.storeAuthData(response)

      return response
    } catch (error) {
      console.error('Login failed:', error)
      throw error
    }
  }

  /**
   * 刷新访问令牌
   * @param data 刷新令牌数据
   * @returns 新的令牌信息
   */
  async refreshToken(data?: RefreshTokenRequest): Promise<RefreshTokenResponse> {
    try {
      // 如果没有提供refreshToken，从本地存储获取
      const refreshToken = data?.refreshToken || this.getRefreshToken()

      if (!refreshToken) {
        throw new Error('No refresh token available')
      }

      const response = await http.post<RefreshTokenResponse>(
        '/auth/refresh',
        { refreshToken },
        { skipAuth: true }, // 刷新token时跳过认证
      )

      // 更新本地存储的token信息
      this.updateTokens(response)

      return response
    } catch (error) {
      console.error('Token refresh failed:', error)
      // 刷新失败时清除所有认证数据
      this.clearAuthData()
      throw error
    }
  }

  /**
   * 用户登出
   */
  async logout(): Promise<void> {
    try {
      // 调用后端登出接口
      await http.post('/auth/logout', {}, { skipErrorHandler: true })
    } catch (error) {
      console.warn('Logout API call failed:', error)
      // 即使API调用失败，也要清除本地数据
    } finally {
      // 清除本地存储的认证数据
      this.clearAuthData()
    }
  }

  /**
   * 获取当前用户信息
   * @returns 用户信息
   */
  async getCurrentUser(): Promise<User> {
    try {
      return await http.get<User>('/profile')
    } catch (error) {
      console.error('Get current user failed:', error)
      throw error
    }
  }

  /**
   * 更新用户信息
   * @param data 更新数据
   * @returns 更新后的用户信息
   */
  async updateUser(data: UpdateUserRequest): Promise<User> {
    try {
      return await http.put<User>('/auth/profile', data)
    } catch (error) {
      console.error('Update user failed:', error)
      throw error
    }
  }

  /**
   * 修改密码
   * @param data 密码修改数据
   */
  async changePassword(data: ChangePasswordRequest): Promise<void> {
    try {
      await http.post('/auth/change-password', data)
    } catch (error) {
      console.error('Change password failed:', error)
      throw error
    }
  }

  /**
   * 请求密码重置
   * @param data 重置请求数据
   */
  async requestPasswordReset(data: ResetPasswordRequest): Promise<void> {
    try {
      await http.post('/auth/reset-password', data, { skipAuth: true })
    } catch (error) {
      console.error('Password reset request failed:', error)
      throw error
    }
  }

  /**
   * 确认密码重置
   * @param data 重置确认数据
   */
  async confirmPasswordReset(data: ResetPasswordConfirmRequest): Promise<void> {
    try {
      await http.post('/auth/reset-password/confirm', data, { skipAuth: true })
    } catch (error) {
      console.error('Password reset confirmation failed:', error)
      throw error
    }
  }

  /**
   * 重置密码（使用token）
   * @param data 重置密码数据
   */
  async resetPassword(data: { token: string; password: string }): Promise<void> {
    try {
      await http.post('/auth/reset-password/confirm', data, { skipAuth: true })
    } catch (error) {
      console.error('Password reset failed:', error)
      throw error
    }
  }

  /**
   * 更新用户资料
   * @param data 更新数据
   * @returns 更新后的用户信息
   */
  async updateProfile(data: Partial<User>): Promise<User> {
    try {
      return await http.put<User>('/auth/profile', data)
    } catch (error) {
      console.error('Update profile failed:', error)
      throw error
    }
  }

  /**
   * 删除账户
   */
  async deleteAccount(): Promise<void> {
    try {
      await http.delete('/auth/account')
      // 删除成功后清除本地数据
      this.clearAuthData()
    } catch (error) {
      console.error('Delete account failed:', error)
      throw error
    }
  }

  // ==================== 私有辅助方法 ====================

  /**
   * 存储认证数据到本地存储
   * @param authData 认证响应数据
   */
  private storeAuthData(authData: AuthResponse): void {
    if (typeof window === 'undefined') return

    const { user, accessToken, refreshToken, expiresIn } = authData
    const expiresAt = Date.now() + expiresIn * 1000

    localStorage.setItem('access_token', accessToken)
    localStorage.setItem('refresh_token', refreshToken)
    localStorage.setItem('user_info', JSON.stringify(user))
    localStorage.setItem('token_expires_at', expiresAt.toString())
  }

  /**
   * 更新令牌信息
   * @param tokenData 令牌响应数据
   */
  private updateTokens(tokenData: RefreshTokenResponse): void {
    if (typeof window === 'undefined') return

    const { accessToken, refreshToken, expiresIn } = tokenData
    const expiresAt = Date.now() + expiresIn * 1000

    localStorage.setItem('access_token', accessToken)
    localStorage.setItem('refresh_token', refreshToken)
    localStorage.setItem('token_expires_at', expiresAt.toString())
  }

  /**
   * 清除所有认证数据
   */
  private clearAuthData(): void {
    if (typeof window === 'undefined') return

    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user_info')
    localStorage.removeItem('token_expires_at')
    localStorage.removeItem('remember_me')
  }

  /**
   * 从本地存储获取刷新令牌
   * @returns 刷新令牌
   */
  private getRefreshToken(): string | null {
    if (typeof window === 'undefined') return null
    return localStorage.getItem('refresh_token')
  }

  // ==================== 公共辅助方法 ====================

  /**
   * 检查是否已认证
   * @returns 是否已认证
   */
  isAuthenticated(): boolean {
    if (typeof window === 'undefined') return false

    const token = localStorage.getItem('access_token')
    const expiresAt = localStorage.getItem('token_expires_at')

    if (!token || !expiresAt) return false

    // 检查token是否过期
    return Date.now() < parseInt(expiresAt)
  }

  /**
   * 获取存储的用户信息
   * @returns 用户信息
   */
  getStoredUser(): User | null {
    if (typeof window === 'undefined') return null

    const userInfo = localStorage.getItem('user_info')
    if (!userInfo) return null

    try {
      return JSON.parse(userInfo)
    } catch {
      return null
    }
  }

  /**
   * 获取访问令牌
   * @returns 访问令牌
   */
  getAccessToken(): string | null {
    if (typeof window === 'undefined') return null
    return localStorage.getItem('access_token')
  }

  /**
   * 检查token是否即将过期
   * @param threshold 阈值（秒）
   * @returns 是否即将过期
   */
  isTokenExpiringSoon(threshold: number = 300): boolean {
    if (typeof window === 'undefined') return false

    const expiresAt = localStorage.getItem('token_expires_at')
    if (!expiresAt) return true

    const expirationTime = parseInt(expiresAt)
    const currentTime = Date.now()
    const timeUntilExpiry = expirationTime - currentTime

    return timeUntilExpiry <= threshold * 1000
  }

  /**
   * 获取token剩余有效时间（秒）
   * @returns 剩余时间
   */
  getTokenRemainingTime(): number {
    if (typeof window === 'undefined') return 0

    const expiresAt = localStorage.getItem('token_expires_at')
    if (!expiresAt) return 0

    const expirationTime = parseInt(expiresAt)
    const currentTime = Date.now()
    const remainingTime = Math.max(0, expirationTime - currentTime)

    return Math.floor(remainingTime / 1000)
  }
}

// 创建并导出认证服务实例
export const authService = new AuthService()
export default authService
