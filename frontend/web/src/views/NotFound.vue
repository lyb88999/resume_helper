<template>
  <div class="not-found">
    <div class="not-found-content">
      <div class="error-code">404</div>
      <h1>页面不存在</h1>
      <p>抱歉，您访问的页面不存在或已被移除。</p>
      <el-button type="primary" @click="goHome">
        <el-icon><HomeFilled /></el-icon>
        返回首页
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { HomeFilled } from '@element-plus/icons-vue'

const router = useRouter()

const goHome = () => {
  router.push('/dashboard')
}
</script>

<style lang="scss" scoped>
.not-found {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: var(--el-bg-color-page);
}

.not-found-content {
  text-align: center;
  
  .error-code {
    font-size: 120px;
    font-weight: 700;
    color: var(--el-color-primary);
    margin-bottom: 20px;
    line-height: 1;
  }
  
  h1 {
    font-size: 32px;
    font-weight: 600;
    color: var(--el-text-color-primary);
    margin: 0 0 16px 0;
  }
  
  p {
    font-size: 16px;
    color: var(--el-text-color-regular);
    margin: 0 0 32px 0;
  }
}

@media (max-width: 768px) {
  .not-found-content {
    .error-code {
      font-size: 80px;
    }
    
    h1 {
      font-size: 24px;
    }
  }
}
</style>
