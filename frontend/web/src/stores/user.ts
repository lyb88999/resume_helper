import { defineStore } from 'pinia'
import { ref, readonly } from 'vue'
import { ElMessage } from 'element-plus'
import * as userApi from '@/api/user'
import type { User, LoginForm, RegisterForm } from '@/types/user'
import { getToken, setToken, removeToken } from '@/utils/auth'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const token = ref<string>('')
  const isAuthenticated = ref(false)
  const loading = ref(false)

  // 检查认证状态
  const checkAuth = async () => {
    const savedToken = getToken()
    if (savedToken) {
      token.value = savedToken
      try {
        const userData = await userApi.getUserInfo()
        user.value = userData
        isAuthenticated.value = true
      } catch (error) {
        console.error('检查认证状态失败:', error)
        await logout()
      }
    }
  }

  // 用户登录
  const login = async (loginForm: LoginForm) => {
    loading.value = true
    try {
      const response = await userApi.login(loginForm)
      
      token.value = response.token
      user.value = response.user
      isAuthenticated.value = true
      
      setToken(response.token)
      
      ElMessage.success('登录成功')
      return response
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || error.message || '登录失败'
      ElMessage.error(errorMessage)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 用户注册
  const register = async (registerForm: RegisterForm) => {
    loading.value = true
    try {
      const response = await userApi.register(registerForm)
      
      token.value = response.token
      user.value = response.user
      isAuthenticated.value = true
      
      setToken(response.token)
      
      ElMessage.success('注册成功')
      return response
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || error.message || '注册失败'
      ElMessage.error(errorMessage)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 用户登出
  const logout = async () => {
    try {
      await userApi.logout()
    } catch (error) {
      console.error('登出API调用失败:', error)
    } finally {
      user.value = null
      token.value = ''
      isAuthenticated.value = false
      removeToken()
      ElMessage.success('已退出登录')
    }
  }

  // 更新用户信息
  const updateProfile = async (profileData: Partial<User>) => {
    loading.value = true
    try {
      const updatedUser = await userApi.updateProfile(profileData)
      user.value = { ...user.value, ...updatedUser }
      ElMessage.success('信息更新成功')
      return updatedUser
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || error.message || '更新失败'
      ElMessage.error(errorMessage)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取用户权限
  const getUserPermissions = () => {
    return user.value?.permissions || []
  }

  // 检查用户权限
  const hasPermission = (permission: string) => {
    const permissions = getUserPermissions()
    return permissions.includes(permission)
  }

  return {
    // state
    user: readonly(user),
    token: readonly(token),
    isAuthenticated: readonly(isAuthenticated),
    loading: readonly(loading),

    // actions
    checkAuth,
    login,
    register,
    logout,
    updateProfile,
    getUserPermissions,
    hasPermission
  }
})
