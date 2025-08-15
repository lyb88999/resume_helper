// 用户相关类型定义

export interface User {
  id: number
  email: string
  nickname: string
  avatar?: string
  createdAt: string
  updatedAt: string
  permissions?: string[]
}

export interface LoginForm {
  email: string
  password: string
}

export interface RegisterForm {
  email: string
  password: string
  nickname: string
  confirmPassword?: string
}

export interface LoginResponse {
  token: string
  expiresAt: number
  user: User
}

export interface RegisterResponse {
  id: number
  email: string
  nickname: string
  createdAt: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface UserProfile {
  nickname: string
  avatar?: string
}
