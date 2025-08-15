<template>
  <div class="resume-page">
    <div class="page-header">
      <h1>简历管理</h1>
      <el-button type="primary" @click="handleUpload">
        <el-icon><Upload /></el-icon>
        上传简历
      </el-button>
    </div>
    
    <div class="resume-content">
      <el-empty v-if="!resumes.length" description="暂无简历，快去上传第一份简历吧！">
        <el-button type="primary" @click="handleUpload">上传简历</el-button>
      </el-empty>
      
      <div v-else class="resume-grid">
        <div v-for="resume in resumes" :key="resume.id" class="resume-card">
          <div class="resume-header">
            <h3>{{ resume.name }}</h3>
            <el-tag :type="getStatusType(resume.status)">
              {{ getStatusText(resume.status) }}
            </el-tag>
          </div>
          
          <div class="resume-info">
            <p><strong>上传时间:</strong> {{ formatDate(resume.uploadTime) }}</p>
            <p><strong>文件类型:</strong> {{ resume.fileType.toUpperCase() }}</p>
            <p v-if="resume.score"><strong>优化分数:</strong> {{ resume.score }}分</p>
          </div>
          
          <div class="resume-actions">
            <el-button size="small" @click="handleView(resume)">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button
              v-if="resume.status === 'pending'"
              size="small"
              type="primary"
              @click="handleAnalyze(resume)"
            >
              <el-icon><Operation /></el-icon>
              分析
            </el-button>
            <el-button
              v-else-if="resume.status === 'completed'"
              size="small"
              type="success"
              @click="handleViewReport(resume)"
            >
              <el-icon><Document /></el-icon>
              查看报告
            </el-button>
            <el-button
              size="small"
              type="danger"
              @click="handleDelete(resume)"
            >
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Upload,
  View,
  Operation,
  Document,
  Delete
} from '@element-plus/icons-vue'
import { useResumeStore } from '@/stores/resume'
import type { Resume } from '@/types/resume'

const resumeStore = useResumeStore()
const resumes = ref<Resume[]>([])
const loading = ref(false)

onMounted(() => {
  loadResumes()
})

const loadResumes = async () => {
  loading.value = true
  try {
    await resumeStore.fetchResumes()
    resumes.value = resumeStore.resumes
  } catch (error) {
    console.error('加载简历列表失败:', error)
  } finally {
    loading.value = false
  }
}

const getStatusType = (status: string) => {
  const typeMap: Record<string, string> = {
    pending: 'info',
    processing: 'warning',
    completed: 'success',
    failed: 'danger'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    pending: '待分析',
    processing: '分析中',
    completed: '已完成',
    failed: '分析失败'
  }
  return textMap[status] || '未知'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

const handleUpload = () => {
  ElMessage.info('简历上传功能开发中...')
}

const handleView = (resume: Resume) => {
  ElMessage.info(`查看简历: ${resume.name}`)
}

const handleAnalyze = async (resume: Resume) => {
  try {
    await resumeStore.parseResume(resume.id)
    ElMessage.success(`开始解析简历: ${resume.name}`)
    // 重新加载简历列表以更新状态
    await loadResumes()
  } catch (error) {
    ElMessage.error('解析简历失败')
  }
}

const handleViewReport = (resume: Resume) => {
  ElMessage.info(`查看分析报告: ${resume.name}`)
}

const handleDelete = async (resume: Resume) => {
  try {
    await ElMessageBox.confirm(`确定要删除简历 "${resume.name}" 吗？`, '确认删除', {
      type: 'warning'
    })
    
    await resumeStore.deleteResume(resume.id)
    // 重新加载简历列表
    await loadResumes()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}
</script>

<style lang="scss" scoped>
.resume-page {
  min-height: 100vh;
  background-color: var(--el-bg-color-page);
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  
  h1 {
    font-size: 24px;
    font-weight: 600;
    margin: 0;
    color: var(--el-text-color-primary);
  }
}

.resume-content {
  max-width: 1200px;
}

.resume-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.resume-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  border: 1px solid var(--el-border-color-light);
  
  .resume-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 12px;
    
    h3 {
      font-size: 16px;
      font-weight: 600;
      margin: 0;
      color: var(--el-text-color-primary);
      flex: 1;
      margin-right: 12px;
    }
  }
  
  .resume-info {
    margin-bottom: 16px;
    
    p {
      margin: 4px 0;
      font-size: 14px;
      color: var(--el-text-color-regular);
      
      strong {
        color: var(--el-text-color-primary);
      }
    }
  }
  
  .resume-actions {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }
}

@media (max-width: 768px) {
  .resume-grid {
    grid-template-columns: 1fr;
  }
  
  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }
}
</style>
