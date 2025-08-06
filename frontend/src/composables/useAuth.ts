import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/modules/auth/store'
import type { LoginRequest, RegisterRequest } from '@/modules/auth/types'

/**
 * 认证相关的组合式函数
 * 提供认证状态和操作的便捷访问
 */
export function useAuth() {
  const authStore = useAuthStore()
  const router = useRouter()

  // ==================== 状态 ====================
  
  const user = computed(() => authStore.user)
  const isAuthenticated = computed(() => authStore.isAuthenticated)
  const isGuest = computed(() => authStore.isGuest)
  const isLoading = computed(() => authStore.isLoadingState)
  const error = computed(() => authStore.error)
  const userDisplayName = computed(() => authStore.userDisplayName)

  // ==================== 方法 ====================
  
  /**
   * 登录
   */
  const login = async (loginData: LoginRequest, redirectTo?: string) => {
    try {
      console.log('开始登录...')
      await authStore.login(loginData)

      console.log('登录成功，认证状态:', authStore.isAuthenticated)
      console.log('用户信息:', authStore.user)

      // 登录成功后重定向
      const redirect = redirectTo || router.currentRoute.value.query.redirect as string || '/'
      console.log('准备跳转到:', redirect)

      await router.push(redirect)
      console.log('跳转完成')
    } catch (error) {
      console.error('登录失败:', error)
      throw error
    }
  }

  /**
   * 注册
   */
  const register = async (registerData: RegisterRequest, redirectTo?: string) => {
    try {
      await authStore.register(registerData)
      
      // 注册成功后重定向
      const redirect = redirectTo || '/'
      await router.push(redirect)
    } catch (error) {
      throw error
    }
  }

  /**
   * 登出
   */
  const logout = async (redirectTo?: string) => {
    try {
      await authStore.logout()
      
      // 登出后重定向到登录页
      const redirect = redirectTo || '/login'
      await router.push(redirect)
    } catch (error) {
      console.error('Logout error:', error)
      // 即使出错也要重定向
      await router.push('/login')
    }
  }

  /**
   * 刷新token
   */
  const refreshToken = async () => {
    return await authStore.refreshAccessToken()
  }

  /**
   * 获取当前用户信息
   */
  const fetchCurrentUser = async () => {
    return await authStore.fetchCurrentUser()
  }

  /**
   * 检查是否需要登录
   */
  const requireAuth = () => {
    if (!isAuthenticated.value) {
      router.push({
        path: '/login',
        query: { redirect: router.currentRoute.value.fullPath }
      })
      return false
    }
    return true
  }

  /**
   * 检查是否需要访客状态（已登录用户不能访问）
   */
  const requireGuest = () => {
    if (isAuthenticated.value) {
      router.push('/')
      return false
    }
    return true
  }

  /**
   * 清除错误信息
   */
  const clearError = () => {
    authStore.setError(null)
  }

  return {
    // 状态
    user,
    isAuthenticated,
    isGuest,
    isLoading,
    error,
    userDisplayName,
    
    // 方法
    login,
    register,
    logout,
    refreshToken,
    fetchCurrentUser,
    requireAuth,
    requireGuest,
    clearError,
  }
}

/**
 * 路由守卫相关的组合式函数
 */
export function useAuthGuard() {
  const authStore = useAuthStore()

  /**
   * 认证路由守卫
   */
  const authGuard = (to: any, from: any, next: any) => {
    console.log('sss')
    const requiresAuth = to.meta?.requiresAuth ?? false
    const guestOnly = to.meta?.guestOnly ?? false

    if (requiresAuth && !authStore.isAuthenticated) {
      // 需要认证但未登录，重定向到登录页
      console.log('t')
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
    } else if (guestOnly && authStore.isAuthenticated) {
      console.log('this is to re1')

      // 只允许访客访问但已登录，重定向到首页
      next('/')
    } else {
      console.log('this is to re')

      // 允许访问
      next()
    }
  }

  return {
    authGuard
  }
}
