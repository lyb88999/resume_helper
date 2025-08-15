import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getToken } from '@/utils/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: {
      title: '登录',
      requiresAuth: false
    }
  },
  {
    path: '/register',
    name: 'Register', 
    component: () => import('@/views/auth/Register.vue'),
    meta: {
      title: '注册',
      requiresAuth: false
    }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: {
      title: '仪表板',
      requiresAuth: true
    }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: {
      title: '个人资料',
      requiresAuth: true
    }
  },
  {
    path: '/resume',
    name: 'Resume',
    component: () => import('@/views/Resume.vue'),
    meta: {
      title: '简历管理',
      requiresAuth: true
    }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: {
      title: '页面不存在'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  const token = getToken()
  
  // 设置页面标题
  if (to.meta?.title) {
    document.title = `${to.meta.title} - 智能简历优化系统`
  }
  
  // 检查是否需要认证
  if (to.meta?.requiresAuth) {
    if (!token) {
      // 未登录，跳转到登录页
      next({
        name: 'Login',
        query: { redirect: to.fullPath }
      })
      return
    }
    
    // 已登录但用户信息为空，尝试获取用户信息
    if (!userStore.isAuthenticated) {
      try {
        await userStore.checkAuth()
      } catch (error) {
        console.error('认证检查失败:', error)
        next({
          name: 'Login',
          query: { redirect: to.fullPath }
        })
        return
      }
    }
  }
  
  // 已登录用户访问登录/注册页面，跳转到仪表板
  if (token && userStore.isAuthenticated && (to.name === 'Login' || to.name === 'Register')) {
    next({ name: 'Dashboard' })
    return
  }
  
  next()
})

export default router
