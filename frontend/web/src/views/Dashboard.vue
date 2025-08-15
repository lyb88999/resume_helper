<template>
  <div class="dashboard">
    <header class="dashboard-header">
      <div class="header-content">
        <h1 class="page-title">仪表板</h1>
        <div class="user-info">
          <span>欢迎回来，{{ userStore.user?.nickname || '用户' }}！</span>
          <el-dropdown @command="handleDropdownCommand">
            <el-avatar :size="32" :src="userStore.user?.avatar">
              {{ userStore.user?.nickname?.charAt(0) || 'U' }}
            </el-avatar>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人资料</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </header>

    <main class="dashboard-main">
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon><Document /></el-icon>
          </div>
          <div class="stat-content">
            <h3>简历总数</h3>
            <p class="stat-number">{{ stats.resumeCount }}</p>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon><Check /></el-icon>
          </div>
          <div class="stat-content">
            <h3>已优化</h3>
            <p class="stat-number">{{ stats.optimizedCount }}</p>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon><Clock /></el-icon>
          </div>
          <div class="stat-content">
            <h3>待处理</h3>
            <p class="stat-number">{{ stats.pendingCount }}</p>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon><TrendCharts /></el-icon>
          </div>
          <div class="stat-content">
            <h3>平均分数</h3>
            <p class="stat-number">{{ stats.averageScore }}分</p>
          </div>
        </div>
      </div>

      <div class="action-cards">
        <div class="action-card">
          <h3>上传新简历</h3>
          <p>上传您的简历，让AI为您提供专业的优化建议</p>
          <el-button type="primary" @click="handleUploadResume">
            <el-icon><Upload /></el-icon>
            上传简历
          </el-button>
        </div>
        
        <div class="action-card">
          <h3>简历管理</h3>
          <p>查看和管理您已上传的所有简历</p>
          <el-button @click="handleManageResumes">
            <el-icon><Folder /></el-icon>
            管理简历
          </el-button>
        </div>
        
        <div class="action-card">
          <h3>知识库</h3>
          <p>浏览简历优化的专业知识和建议</p>
          <el-button @click="handleKnowledgeBase">
            <el-icon><Reading /></el-icon>
            浏览知识库
          </el-button>
        </div>
      </div>

      <div class="recent-activity">
        <h2>最近活动</h2>
        <el-empty v-if="!recentActivities.length" description="暂无活动记录" />
        <div v-else class="activity-list">
          <div v-for="activity in recentActivities" :key="activity.id" class="activity-item">
            <div class="activity-icon">
              <el-icon><Document /></el-icon>
            </div>
            <div class="activity-content">
              <h4>{{ activity.title }}</h4>
              <p>{{ activity.description }}</p>
              <span class="activity-time">{{ activity.time }}</span>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import {
  Document,
  Check,
  Clock, 
  TrendCharts,
  Upload,
  Folder,
  Reading
} from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const stats = reactive({
  resumeCount: 0,
  optimizedCount: 0,
  pendingCount: 0,
  averageScore: 0
})

const recentActivities = ref<Array<{
  id: number
  title: string
  description: string
  time: string
}>>([])

onMounted(() => {
  loadDashboardData()
})

const loadDashboardData = async () => {
  // 模拟加载数据
  stats.resumeCount = 5
  stats.optimizedCount = 3
  stats.pendingCount = 2
  stats.averageScore = 85

  recentActivities.value = [
    {
      id: 1,
      title: '简历优化完成',
      description: '您的软件工程师简历已完成优化分析',
      time: '2小时前'
    },
    {
      id: 2,
      title: '新简历上传',
      description: '产品经理简历上传成功',
      time: '1天前'
    }
  ]
}

const handleDropdownCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '确认', {
          type: 'warning'
        })
        await userStore.logout()
        router.push('/login')
      } catch (error) {
        // 用户取消
      }
      break
  }
}

const handleUploadResume = () => {
  ElMessage.info('简历上传功能开发中...')
}

const handleManageResumes = () => {
  router.push('/resume')
}

const handleKnowledgeBase = () => {
  ElMessage.info('知识库功能开发中...')
}
</script>

<style lang="scss" scoped>
.dashboard {
  min-height: 100vh;
  background-color: var(--el-bg-color-page);
}

.dashboard-header {
  background: white;
  border-bottom: 1px solid var(--el-border-color-light);
  padding: 16px 24px;
  
  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    .page-title {
      font-size: 24px;
      font-weight: 600;
      margin: 0;
      color: var(--el-text-color-primary);
    }
    
    .user-info {
      display: flex;
      align-items: center;
      gap: 12px;
      
      span {
        color: var(--el-text-color-regular);
      }
    }
  }
}

.dashboard-main {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.stat-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  display: flex;
  align-items: center;
  gap: 16px;
  
  .stat-icon {
    width: 48px;
    height: 48px;
    background: var(--el-color-primary-light-9);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--el-color-primary);
    font-size: 20px;
  }
  
  .stat-content {
    h3 {
      font-size: 14px;
      color: var(--el-text-color-regular);
      margin: 0 0 4px 0;
      font-weight: normal;
    }
    
    .stat-number {
      font-size: 24px;
      font-weight: 600;
      color: var(--el-text-color-primary);
      margin: 0;
    }
  }
}

.action-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.action-card {
  background: white;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  text-align: center;
  
  h3 {
    font-size: 18px;
    margin: 0 0 8px 0;
    color: var(--el-text-color-primary);
  }
  
  p {
    color: var(--el-text-color-regular);
    margin: 0 0 20px 0;
    line-height: 1.5;
  }
}

.recent-activity {
  background: white;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  
  h2 {
    font-size: 18px;
    margin: 0 0 16px 0;
    color: var(--el-text-color-primary);
  }
}

.activity-list {
  .activity-item {
    display: flex;
    gap: 12px;
    padding: 12px 0;
    border-bottom: 1px solid var(--el-border-color-lighter);
    
    &:last-child {
      border-bottom: none;
    }
    
    .activity-icon {
      width: 32px;
      height: 32px;
      background: var(--el-color-info-light-9);
      border-radius: 6px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: var(--el-color-info);
      flex-shrink: 0;
    }
    
    .activity-content {
      flex: 1;
      
      h4 {
        font-size: 14px;
        margin: 0 0 4px 0;
        color: var(--el-text-color-primary);
      }
      
      p {
        font-size: 13px;
        color: var(--el-text-color-regular);
        margin: 0 0 4px 0;
      }
      
      .activity-time {
        font-size: 12px;
        color: var(--el-text-color-placeholder);
      }
    }
  }
}

@media (max-width: 768px) {
  .dashboard-main {
    padding: 16px;
  }
  
  .stats-grid,
  .action-cards {
    grid-template-columns: 1fr;
  }
}
</style>
