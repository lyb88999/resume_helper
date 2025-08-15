<template>
  <div class="profile-page">
    <div class="page-header">
      <h1>个人资料</h1>
      <p>管理您的账户信息</p>
    </div>
    
    <div class="profile-content">
      <el-card class="profile-card">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
          </div>
        </template>
        
        <el-form
          ref="profileFormRef"
          :model="profileForm"
          :rules="profileRules"
          label-width="80px"
          class="profile-form"
        >
          <el-form-item label="邮箱">
            <el-input :value="userStore.user?.email" disabled />
          </el-form-item>
          
          <el-form-item label="昵称" prop="nickname">
            <el-input
              v-model="profileForm.nickname"
              placeholder="请输入昵称"
              clearable
            />
          </el-form-item>
          
          <el-form-item label="头像" prop="avatar">
            <el-input
              v-model="profileForm.avatar"
              placeholder="请输入头像URL"
              clearable
            />
          </el-form-item>
          
          <el-form-item>
            <el-button
              type="primary"
              :loading="userStore.loading"
              @click="handleUpdateProfile"
            >
              保存修改
            </el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </el-card>
      
      <el-card class="account-info">
        <template #header>
          <div class="card-header">
            <span>账户信息</span>
          </div>
        </template>
        
        <div class="info-item">
          <label>注册时间</label>
          <span>{{ formatDate(userStore.user?.createdAt) }}</span>
        </div>
        
        <div class="info-item">
          <label>最后更新</label>
          <span>{{ formatDate(userStore.user?.updatedAt) }}</span>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { ElForm, ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import type { UserProfile } from '@/types/user'

const userStore = useUserStore()
const profileFormRef = ref<InstanceType<typeof ElForm>>()

const profileForm = reactive<UserProfile>({
  nickname: '',
  avatar: ''
})

const profileRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '昵称长度在2到20个字符', trigger: 'blur' }
  ]
}

onMounted(() => {
  initForm()
})

const initForm = () => {
  if (userStore.user) {
    profileForm.nickname = userStore.user.nickname
    profileForm.avatar = userStore.user.avatar || ''
  }
}

const handleUpdateProfile = async () => {
  if (!profileFormRef.value) return
  
  try {
    await profileFormRef.value.validate()
    await userStore.updateProfile(profileForm)
  } catch (error) {
    console.error('更新失败:', error)
  }
}

const handleReset = () => {
  initForm()
}

const formatDate = (dateString?: string) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>

<style lang="scss" scoped>
.profile-page {
  min-height: 100vh;
  background-color: var(--el-bg-color-page);
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
  
  h1 {
    font-size: 24px;
    font-weight: 600;
    margin: 0 0 8px 0;
    color: var(--el-text-color-primary);
  }
  
  p {
    color: var(--el-text-color-regular);
    margin: 0;
  }
}

.profile-content {
  max-width: 800px;
  display: grid;
  gap: 24px;
  grid-template-columns: 2fr 1fr;
}

.profile-card {
  .card-header {
    font-weight: 600;
  }
  
  .profile-form {
    .el-form-item {
      margin-bottom: 20px;
    }
  }
}

.account-info {
  .card-header {
    font-weight: 600;
  }
  
  .info-item {
    display: flex;
    justify-content: space-between;
    padding: 12px 0;
    border-bottom: 1px solid var(--el-border-color-lighter);
    
    &:last-child {
      border-bottom: none;
    }
    
    label {
      color: var(--el-text-color-regular);
      font-weight: 500;
    }
    
    span {
      color: var(--el-text-color-primary);
    }
  }
}

@media (max-width: 768px) {
  .profile-content {
    grid-template-columns: 1fr;
  }
}
</style>
