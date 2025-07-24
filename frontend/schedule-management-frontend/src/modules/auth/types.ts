// 用户信息接口
export interface User {
  id: string
  email: string
  firstName: string
  lastName: string
  phone?: string
  bio?: string
  preferences?: {
    emailNotifications: boolean
    desktopNotifications: boolean
    marketingEmails: boolean
  }
  createdAt: string
  updatedAt: string
}

// 用户注册请求接口
export interface RegisterRequest {
  email: string
  password: string
  firstName: string
  lastName: string
}

// 用户登录请求接口
export interface LoginRequest {
  email: string
  password: string
}

// 认证响应接口
export interface AuthResponse {
  user: User
  accessToken: string
  refreshToken: string
  expiresIn: number // token过期时间（秒）
}

// 令牌刷新请求接口
export interface RefreshTokenRequest {
  refreshToken: string
}

// 令牌刷新响应接口
export interface RefreshTokenResponse {
  accessToken: string
  refreshToken: string
  expiresIn: number
}

// 用户状态枚举
export const UserStatus = {
  GUEST: 'guest',
  AUTHENTICATED: 'authenticated',
  LOADING: 'loading',
} as const

export type UserStatusType = typeof UserStatus[keyof typeof UserStatus]

// 认证状态接口
export interface AuthState {
  // 用户信息
  user: User | null

  // 认证状态
  status: UserStatusType
  isAuthenticated: boolean

  // Token信息
  accessToken: string | null
  refreshToken: string | null
  tokenExpiresAt: number | null

  // 加载状态
  isLoading: boolean

  // 错误信息
  error: string | null
}

// 登录表单数据
export interface LoginFormData {
  email: string
  password: string
  remember?: boolean
}

// 注册表单数据
export interface RegisterFormData {
  email: string
  password: string
  confirmPassword: string
  firstName: string
  lastName: string
  agreeToTerms: boolean
}

// 密码重置请求
export interface ResetPasswordRequest {
  email: string
}

// 密码重置确认请求
export interface ResetPasswordConfirmRequest {
  token: string
  newPassword: string
  confirmPassword: string
}

// 修改密码请求
export interface ChangePasswordRequest {
  currentPassword: string
  newPassword: string
  confirmPassword: string
}

// 更新用户信息请求
export interface UpdateUserRequest {
  firstName?: string
  lastName?: string
  email?: string
}

// API错误响应
export interface ApiError {
  code: number
  message: string
  details?: Record<string, any>
  field?: string // 字段级别的错误
}

// 表单验证错误
export interface ValidationError {
  field: string
  message: string
}

// 认证事件类型
export enum AuthEventType {
  LOGIN_SUCCESS = 'login_success',
  LOGIN_FAILED = 'login_failed',
  LOGOUT = 'logout',
  TOKEN_EXPIRED = 'token_expired',
  TOKEN_REFRESHED = 'token_refreshed',
  REGISTRATION_SUCCESS = 'registration_success',
  REGISTRATION_FAILED = 'registration_failed',
}

// 认证事件接口
export interface AuthEvent {
  type: AuthEventType
  payload?: any
  timestamp: number
}

// 本地存储键名
export const AUTH_STORAGE_KEYS = {
  ACCESS_TOKEN: 'access_token',
  REFRESH_TOKEN: 'refresh_token',
  USER_INFO: 'user_info',
  TOKEN_EXPIRES_AT: 'token_expires_at',
  REMEMBER_ME: 'remember_me',
} as const

// 认证配置
export interface AuthConfig {
  // Token刷新阈值（提前多少秒刷新token）
  refreshThreshold: number

  // 自动刷新间隔（毫秒）
  refreshInterval: number

  // 最大重试次数
  maxRetries: number

  // 记住我功能的有效期（天）
  rememberMeDays: number
}

// 默认认证配置
export const DEFAULT_AUTH_CONFIG: AuthConfig = {
  refreshThreshold: 300, // 5分钟
  refreshInterval: 60000, // 1分钟
  maxRetries: 3,
  rememberMeDays: 30,
}