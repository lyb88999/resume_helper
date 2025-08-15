<template>
  <div id="app" class="app-container">
    <router-view />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'

const appStore = useAppStore()
const userStore = useUserStore()

onMounted(async () => {
  // 初始化应用配置
  await appStore.initApp()
  
  // 检查用户登录状态
  await userStore.checkAuth()
})
</script>

<style lang="scss">
.app-container {
  min-height: 100vh;
  background-color: var(--el-bg-color-page);
}

// 全局滚动条样式
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: var(--el-fill-color-light);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: var(--el-fill-color);
  border-radius: 3px;
  
  &:hover {
    background: var(--el-fill-color-dark);
  }
}

// 响应式断点
@media (max-width: 768px) {
  .app-container {
    padding: 0;
  }
}
</style>
