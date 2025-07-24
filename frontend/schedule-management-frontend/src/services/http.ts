import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'

// 扩展axios配置类型
declare module 'axios' {
  interface InternalAxiosRequestConfig {
    skipAuth?: boolean
    skipErrorHandler?: boolean
    metadata?: {
      startTime: number
    }
  }
}

// API响应基础接口
export interface ApiResponse<T = any> {
  success: boolean
  data: T
  message: string
  code: number
}

// HTTP错误类型
export interface HttpError {
  code: number
  message: string
  details?: any
}

// 请求配置扩展
export interface RequestConfig extends AxiosRequestConfig {
  skipAuth?: boolean // 跳过认证
  skipErrorHandler?: boolean // 跳过错误处理
}

class HttpService {
  private instance: AxiosInstance

  constructor() {
    // 创建axios实例
    this.instance = axios.create({
      baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
      timeout: 10000,
      withCredentials: true, // 支持Cookie
      headers: {
        'Content-Type': 'application/json',
      },
    })

    this.setupInterceptors()
  }

  private setupInterceptors() {
    // 请求拦截器
    this.instance.interceptors.request.use(
      (config) => {
        // 添加JWT token到请求头
        const token = this.getTokenFromStorage()
        if (token && !config.skipAuth) {
          config.headers.Authorization = `Bearer ${token}`
        }

        // 添加请求时间戳
        config.metadata = { startTime: Date.now() }

        console.log(`🚀 [${config.method?.toUpperCase()}] ${config.url}`, {
          headers: config.headers,
          data: config.data,
        })

        return config
      },
      (error) => {
        console.error('❌ Request Error:', error)
        return Promise.reject(error)
      }
    )

    // 响应拦截器
    this.instance.interceptors.response.use(
      (response: AxiosResponse) => {
        const duration = Date.now() - (response.config.metadata?.startTime || 0)
        console.log(
          `✅ [${response.config.method?.toUpperCase()}] ${response.config.url} (${duration}ms)`,
          response.data
        )

        // 统一处理响应格式
        return this.handleResponse(response)
      },
      (error) => {
        return this.handleError(error)
      }
    )
  }

  private handleResponse(response: AxiosResponse): AxiosResponse {
    const { data } = response

    // 如果后端返回的是标准格式
    if (data && typeof data === 'object' && 'success' in data) {
      if (!data.success) {
        throw new Error(data.message || '请求失败')
      }
    }

    return response
  }

  private handleError(error: any): Promise<never> {
    const duration = Date.now() - (error.config?.metadata?.startTime || 0)
    console.error(
      `❌ [${error.config?.method?.toUpperCase()}] ${error.config?.url} (${duration}ms)`,
      error
    )

    // 跳过错误处理的请求
    if (error.config?.skipErrorHandler) {
      return Promise.reject(error)
    }

    let errorMessage = '网络请求失败'
    let errorCode = 0

    if (error.response) {
      // 服务器响应错误
      const { status, data } = error.response
      errorCode = status

      switch (status) {
        case 400:
          errorMessage = data?.message || '请求参数错误'
          break
        case 401:
          errorMessage = '未授权，请重新登录'
          this.handleUnauthorized()
          break
        case 403:
          errorMessage = '拒绝访问'
          break
        case 404:
          errorMessage = '请求的资源不存在'
          break
        case 422:
          errorMessage = data?.message || '数据验证失败'
          break
        case 500:
          errorMessage = '服务器内部错误'
          break
        case 502:
          errorMessage = '网关错误'
          break
        case 503:
          errorMessage = '服务不可用'
          break
        default:
          errorMessage = data?.message || `请求失败 (${status})`
      }
    } else if (error.request) {
      // 网络错误
      errorMessage = '网络连接失败，请检查网络设置'
      errorCode = -1
    } else {
      // 其他错误
      errorMessage = error.message || '未知错误'
      errorCode = -2
    }

    // 错误消息将在组件中处理
    console.error('HTTP Error:', errorMessage)

    const httpError: HttpError = {
      code: errorCode,
      message: errorMessage,
      details: error.response?.data,
    }

    return Promise.reject(httpError)
  }

  private handleUnauthorized() {
    // 清除本地存储的认证信息
    this.clearAuthData()
    
    // 重定向到登录页面
    if (typeof window !== 'undefined') {
      window.location.href = '/login'
    }
  }

  private getTokenFromStorage(): string | null {
    if (typeof window === 'undefined') return null
    return localStorage.getItem('access_token')
  }

  private clearAuthData() {
    if (typeof window === 'undefined') return
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user_info')
  }

  // 公共请求方法
  public async get<T = any>(url: string, config?: RequestConfig): Promise<T> {
    const response = await this.instance.get(url, config)
    return response.data as T
  }

  public async post<T = any>(url: string, data?: any, config?: RequestConfig): Promise<T> {
    const response = await this.instance.post(url, data, config)
    return response.data as T
  }

  public async put<T = any>(url: string, data?: any, config?: RequestConfig): Promise<T> {
    const response = await this.instance.put(url, data, config)
    return response.data as T
  }

  public async delete<T = any>(url: string, config?: RequestConfig): Promise<T> {
    const response = await this.instance.delete(url, config)
    return response.data as T
  }

  // 获取原始axios实例（用于特殊需求）
  public getInstance(): AxiosInstance {
    return this.instance
  }
}

// 创建并导出HTTP服务实例
export const http = new HttpService()
export default http
