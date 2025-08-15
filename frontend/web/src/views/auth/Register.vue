<template>
  <div class="register-container">
    <div class="register-card">
      <div class="header">
        <h1 class="title">智能简历优化系统</h1>
        <p class="subtitle">创建您的账户</p>
      </div>
      
      <el-form
        ref="registerFormRef"
        :model="registerForm"
        :rules="registerRules"
        label-width="0"
        size="large"
        class="register-form"
      >
        <el-form-item prop="nickname">
          <el-input
            v-model="registerForm.nickname"
            placeholder="请输入昵称"
            prefix-icon="User"
            clearable
          />
        </el-form-item>
        
        <el-form-item prop="email">
          <el-input
            v-model="registerForm.email"
            placeholder="请输入邮箱"
            prefix-icon="Message"
            clearable
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="registerForm.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="Lock"
            show-password
            clearable
          />
        </el-form-item>
        
        <el-form-item prop="confirmPassword">
          <el-input
            v-model="registerForm.confirmPassword"
            type="password"
            placeholder="请确认密码"
            prefix-icon="Lock"
            show-password
            clearable
            @keyup.enter="handleRegister"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            class="register-btn"
            :loading="userStore.loading"
            @click="handleRegister"
          >
            注册
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="footer">
        <p>
          已有账户？
          <router-link to="/login" class="link">立即登录</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElForm, ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import type { RegisterForm } from '@/types/user'

const router = useRouter()
const userStore = useUserStore()

const registerFormRef = ref<InstanceType<typeof ElForm>>()

const registerForm = reactive<RegisterForm>({
  nickname: '',
  email: '',
  password: '',
  confirmPassword: ''
})

// 确认密码验证
const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (value !== registerForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const registerRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '昵称长度在2到20个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  if (!registerFormRef.value) return
  
  try {
    await registerFormRef.value.validate()
    
    await userStore.register(registerForm)
    
    // 注册成功后跳转到仪表板
    router.push('/dashboard')
  } catch (error) {
    console.error('注册失败:', error)
  }
}
</script>

<style lang="scss" scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.register-card {
  width: 100%;
  max-width: 400px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  padding: 40px;
}

.header {
  text-align: center;
  margin-bottom: 32px;
  
  .title {
    font-size: 24px;
    font-weight: 600;
    color: #2c3e50;
    margin: 0 0 8px 0;
  }
  
  .subtitle {
    font-size: 14px;
    color: #7f8c8d;
    margin: 0;
  }
}

.register-form {
  .el-form-item {
    margin-bottom: 20px;
  }
  
  .register-btn {
    width: 100%;
    height: 44px;
    font-size: 16px;
    font-weight: 500;
  }
}

.footer {
  text-align: center;
  margin-top: 24px;
  
  p {
    font-size: 14px;
    color: #7f8c8d;
    margin: 0;
  }
  
  .link {
    color: #409eff;
    text-decoration: none;
    font-weight: 500;
    
    &:hover {
      text-decoration: underline;
    }
  }
}

@media (max-width: 480px) {
  .register-card {
    padding: 24px;
  }
}
</style>
