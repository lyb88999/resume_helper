import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const title = ref('智能简历优化系统')
  const theme = ref<'light' | 'dark'>('light')
  const language = ref('zh-CN')
  const loading = ref(false)
  const sidebarCollapse = ref(false)
  
  // 初始化应用
  const initApp = async () => {
    loading.value = true
    try {
      // 从localStorage恢复设置
      const savedTheme = localStorage.getItem('app_theme')
      if (savedTheme && (savedTheme === 'light' || savedTheme === 'dark')) {
        theme.value = savedTheme
      }
      
      const savedLanguage = localStorage.getItem('app_language')
      if (savedLanguage) {
        language.value = savedLanguage
      }
      
      const savedSidebarCollapse = localStorage.getItem('sidebar_collapse')
      if (savedSidebarCollapse) {
        sidebarCollapse.value = JSON.parse(savedSidebarCollapse)
      }
      
      // 应用主题
      applyTheme()
    } catch (error) {
      console.error('初始化应用失败:', error)
    } finally {
      loading.value = false
    }
  }
  
  // 切换主题
  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
    localStorage.setItem('app_theme', theme.value)
    applyTheme()
  }
  
  // 应用主题
  const applyTheme = () => {
    if (theme.value === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }
  
  // 设置语言
  const setLanguage = (lang: string) => {
    language.value = lang
    localStorage.setItem('app_language', lang)
  }
  
  // 切换侧边栏
  const toggleSidebar = () => {
    sidebarCollapse.value = !sidebarCollapse.value
    localStorage.setItem('sidebar_collapse', JSON.stringify(sidebarCollapse.value))
  }
  
  // 设置全局加载状态
  const setLoading = (status: boolean) => {
    loading.value = status
  }
  
  return {
    // state
    title,
    theme,
    language,
    loading,
    sidebarCollapse,
    
    // actions
    initApp,
    toggleTheme,
    setLanguage,
    toggleSidebar,
    setLoading,
    applyTheme
  }
})
