import axios, { type AxiosInstance, type AxiosResponse, type AxiosError } from 'axios'
import { ElMessage } from 'element-plus'
import { getToken, removeToken } from './auth'

// 创建axios实例
const request: AxiosInstance = axios.create({
  baseURL: '', // 使用相对路径，让vite代理处理
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error: AxiosError) => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const { data } = response
    
    // 直接返回kratos HTTP响应的数据
    return data
  },
  (error: AxiosError) => {
    console.error('响应错误:', error)
    
    if (error.response) {
      const { status, data, config } = error.response
      const isLoginRequest = config.url?.includes('/login')
      const isRegisterRequest = config.url?.includes('/register')
      
      switch (status) {
        case 401:
          // 如果是登录或注册请求的401错误，不显示通用错误信息，让具体页面处理
          if (!isLoginRequest && !isRegisterRequest) {
            ElMessage.error('登录已过期，请重新登录')
            removeToken()
            window.location.href = '/login'
          }
          break
        case 403:
          ElMessage.error('权限不足')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          // 对于登录和注册请求，不显示通用错误，让页面自己处理
          if (!isLoginRequest && !isRegisterRequest) {
            ElMessage.error((data as any)?.message || '请求失败')
          }
      }
    } else if (error.request) {
      ElMessage.error('网络连接失败，请检查网络')
    } else {
      ElMessage.error('请求配置错误')
    }
    
    return Promise.reject(error)
  }
)

export default request
