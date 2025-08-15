import request from '@/utils/request'
import type { User, LoginForm, RegisterForm, LoginResponse, RegisterResponse, UserProfile } from '@/types/user'

/**
 * 用户注册
 */
export function register(data: RegisterForm): Promise<RegisterResponse> {
  return request.post('/v1/user/register', {
    email: data.email,
    password: data.password,
    nickname: data.nickname
  })
}

/**
 * 用户登录
 */
export function login(data: LoginForm): Promise<LoginResponse> {
  return request.post('/v1/user/login', {
    email: data.email,
    password: data.password
  })
}

/**
 * 获取用户信息
 */
export function getUserInfo(): Promise<User> {
  // 注意：这里需要用户ID，实际使用时应该从token中解析或从其他方式获取
  // 暂时使用固定ID，后续需要优化
  return request.get('/v1/user/1')
}

/**
 * 更新用户信息
 */
export function updateProfile(data: UserProfile): Promise<User> {
  // 同样需要用户ID
  return request.put('/v1/user/1', data)
}

/**
 * 用户登出
 */
export function logout(): Promise<void> {
  // 目前后端没有登出接口，前端直接清除token
  return Promise.resolve()
}
