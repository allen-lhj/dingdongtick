import { defineStore } from 'pinia'
import { ref, computed, watch, readonly } from 'vue'
import authService from './service'
import type {
  User,
  UserStatusType,
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  AuthEventTypeType,
} from './types'
import { UserStatus, AuthEventType } from './types'
import { DEFAULT_AUTH_CONFIG } from './types'

/**
 * 认证状态管理Store
 * 使用Composition API风格的Pinia store
 */
export const useAuthStore = defineStore('auth', () => {
  // ==================== 状态定义 ====================

  // 用户信息
  const user = ref<User | null>(null)

  // 认证状态
  const status = ref<UserStatusType>(UserStatus.GUEST)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const tokenExpiresAt = ref<number | null>(null)

  // 加载状态
  const isLoading = ref(false)
  const isRefreshing = ref(false)

  // 错误信息
  const error = ref<string | null>(null)

  // 自动刷新定时器
  let refreshTimer: NodeJS.Timeout | null = null

  // ==================== 计算属性 ====================

  // 是否已认证
  const isAuthenticated = computed(() => {
    return status.value === UserStatus.AUTHENTICATED && !!user.value && !!accessToken.value
  })

  // 是否为访客
  const isGuest = computed(() => {
    return status.value === UserStatus.GUEST
  })

  // 是否正在加载
  const isLoadingState = computed(() => {
    return status.value === UserStatus.LOADING || isLoading.value
  })

  // token是否即将过期
  const isTokenExpiringSoon = computed(() => {
    if (!tokenExpiresAt.value) return false
    const now = Date.now()
    const threshold = DEFAULT_AUTH_CONFIG.refreshThreshold * 1000
    return tokenExpiresAt.value - now <= threshold
  })

  // 用户显示名称
  const userDisplayName = computed(() => {
    if (!user.value) return ''
    return `${user.value.firstName} ${user.value.lastName}`.trim() || user.value.email
  })

  // ==================== 状态管理方法 ====================

  /**
   * 设置认证状态
   */
  const setAuthState = (authData: AuthResponse) => {
    user.value = authData.user
    accessToken.value = authData.accessToken
    refreshToken.value = authData.refreshToken
    tokenExpiresAt.value = Date.now() + authData.expiresIn * 1000
    status.value = UserStatus.AUTHENTICATED
    error.value = null

    // 启动自动刷新
    startAutoRefresh()
  }

  /**
   * 清除认证状态
   */
  const clearAuthState = () => {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    tokenExpiresAt.value = null
    status.value = UserStatus.GUEST
    error.value = null

    // 停止自动刷新
    stopAutoRefresh()

    // 触发登出事件
    emitAuthEvent(AuthEventType.LOGOUT)
  }

  /**
   * 设置加载状态
   */
  const setLoading = (loading: boolean) => {
    isLoading.value = loading
    if (loading) {
      status.value = UserStatus.LOADING
    }
  }

  /**
   * 设置错误信息
   */
  const setError = (errorMessage: string | null) => {
    error.value = errorMessage
    // 错误信息将在组件中通过computed属性获取并显示
  }

  // ==================== 认证操作方法 ====================

  /**
   * 用户登录
   */
  const login = async (loginData: LoginRequest): Promise<void> => {
    try {
      setLoading(true)
      setError(null)

      const response = await authService.login(loginData)
      setAuthState(response)

      // 触发登录成功事件
      emitAuthEvent(AuthEventType.LOGIN_SUCCESS, { user: response.user })

      // 登录成功消息将在组件中显示
    } catch (err: any) {
      const errorMessage = err.message || '登录失败'
      setError(errorMessage)
      emitAuthEvent(AuthEventType.LOGIN_FAILED, { error: errorMessage })
      throw err
    } finally {
      setLoading(false)
    }
  }

  /**
   * 用户注册
   */
  const register = async (registerData: RegisterRequest): Promise<void> => {
    try {
      setLoading(true)
      setError(null)

      const response = await authService.register(registerData)
      setAuthState(response)

      // 注册成功消息将在组件中显示
      emitAuthEvent(AuthEventType.REGISTRATION_SUCCESS, { user: response.user })
    } catch (err: any) {
      const errorMessage = err.message || '注册失败'
      setError(errorMessage)
      emitAuthEvent(AuthEventType.REGISTRATION_FAILED, { error: errorMessage })
      throw err
    } finally {
      setLoading(false)
    }
  }

  /**
   * 用户登出
   */
  const logout = async (): Promise<void> => {
    try {
      setLoading(true)

      await authService.logout()
      clearAuthState()

      // 登出成功消息将在组件中显示
    } catch (err: any) {
      console.warn('Logout error:', err)
      // 即使登出API失败，也要清除本地状态
      clearAuthState()
    } finally {
      setLoading(false)
    }
  }

  /**
   * 刷新访问令牌
   */
  const refreshAccessToken = async (): Promise<boolean> => {
    if (isRefreshing.value) return false

    try {
      isRefreshing.value = true

      const response = await authService.refreshToken()

      // 更新token信息
      accessToken.value = response.accessToken
      refreshToken.value = response.refreshToken
      tokenExpiresAt.value = Date.now() + response.expiresIn * 1000

      emitAuthEvent(AuthEventType.TOKEN_REFRESHED)
      return true
    } catch (err: any) {
      console.error('Token refresh failed:', err)

      // 刷新失败，清除认证状态
      clearAuthState()
      emitAuthEvent(AuthEventType.TOKEN_EXPIRED)

      return false
    } finally {
      isRefreshing.value = false
    }
  }

  /**
   * 获取当前用户信息
   */
  const fetchCurrentUser = async (): Promise<void> => {
    try {
      setLoading(true)

      const userData = await authService.getCurrentUser()
      user.value = userData
    } catch (err: any) {
      console.error('Fetch current user failed:', err)
      setError('获取用户信息失败')
      throw err
    } finally {
      setLoading(false)
    }
  }

  /**
   * 初始化认证状态
   * 从本地存储恢复认证信息
   */
  const initializeAuth = async (): Promise<void> => {
    try {
      setLoading(true)

      // 检查是否有有效的认证信息
      if (!authService.isAuthenticated()) {
        status.value = UserStatus.GUEST
        return
      }

      // 恢复认证状态
      const storedUser = authService.getStoredUser()
      const storedToken = authService.getAccessToken()

      if (storedUser && storedToken) {
        user.value = storedUser
        accessToken.value = storedToken
        refreshToken.value = localStorage.getItem('refresh_token')
        tokenExpiresAt.value = parseInt(localStorage.getItem('token_expires_at') || '0')
        status.value = UserStatus.AUTHENTICATED

        // 检查token是否需要刷新
        if (isTokenExpiringSoon.value) {
          await refreshAccessToken()
        }

        // 启动自动刷新
        startAutoRefresh()

        // 验证用户信息是否仍然有效
        try {
          await fetchCurrentUser()
        } catch {
          // 如果获取用户信息失败，清除认证状态
          clearAuthState()
        }
      }
    } catch (err: any) {
      console.error('Initialize auth failed:', err)
      clearAuthState()
    } finally {
      setLoading(false)
    }
  }

  // ==================== 自动刷新逻辑 ====================

  /**
   * 启动自动刷新定时器
   */
  const startAutoRefresh = () => {
    stopAutoRefresh() // 先清除现有定时器

    refreshTimer = setInterval(async () => {
      if (isTokenExpiringSoon.value && !isRefreshing.value) {
        await refreshAccessToken()
      }
    }, DEFAULT_AUTH_CONFIG.refreshInterval)
  }

  /**
   * 停止自动刷新定时器
   */
  const stopAutoRefresh = () => {
    if (refreshTimer) {
      clearInterval(refreshTimer)
      refreshTimer = null
    }
  }

  // ==================== 事件处理 ====================

  /**
   * 触发认证事件
   */
  const emitAuthEvent = (type: AuthEventTypeType, payload?: any) => {
    const event = {
      type,
      payload,
      timestamp: Date.now(),
    }

    console.log('Auth Event:', event)

    // 可以在这里添加事件监听器或发送到分析服务
    // eventBus.emit('auth:event', event
  }

  // ==================== 监听器设置 ====================

  // 监听认证状态变化
  watch(isAuthenticated, (newValue) => {
    if (newValue) {
      startAutoRefresh()
    } else {
      stopAutoRefresh()
    }
  })

  // 页面可见性变化时检查token状态
  if (typeof window !== 'undefined') {
    document.addEventListener('visibilitychange', () => {
      if (!document.hidden && isAuthenticated.value && isTokenExpiringSoon.value) {
        refreshAccessToken()
      }
    })
  }

  // ==================== 返回store接口 ====================

  return {
    // 状态
    user: readonly(user),
    status: readonly(status),
    accessToken: readonly(accessToken),
    refreshToken: readonly(refreshToken),
    tokenExpiresAt: readonly(tokenExpiresAt),
    isLoading: readonly(isLoading),
    isRefreshing: readonly(isRefreshing),
    error: readonly(error),

    // 计算属性
    isAuthenticated,
    isGuest,
    isLoadingState,
    isTokenExpiringSoon,
    userDisplayName,

    // 方法
    login,
    register,
    logout,
    refreshAccessToken,
    fetchCurrentUser,
    initializeAuth,
    setError,
    clearAuthState,
  }
})

// 导出store类型
export type AuthStore = ReturnType<typeof useAuthStore>
