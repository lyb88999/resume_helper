<template>
  <div class="login-container">
    <div class="login-card">
      <div class="header">
        <h1 class="title">智能简历优化系统</h1>
        <p class="subtitle">登录您的账户</p>
      </div>
      
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        label-width="0"
        size="large"
        class="login-form"
      >
        <el-form-item prop="email">
          <el-input
            v-model="loginForm.email"
            placeholder="请输入邮箱"
            prefix-icon="Message"
            clearable
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="Lock"
            show-password
            clearable
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            class="login-btn"
            :loading="userStore.loading"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="footer">
        <p>
          还没有账户？
          <router-link to="/register" class="link">立即注册</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElForm, ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import type { LoginForm } from '@/types/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loginFormRef = ref<InstanceType<typeof ElForm>>()

const loginForm = reactive<LoginForm>({
  email: '',
  password: ''
})

const loginRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    
    await userStore.login(loginForm)
    
    // 登录成功后跳转
    const redirectPath = (route.query.redirect as string) || '/dashboard'
    router.push(redirectPath)
  } catch (error) {
    console.error('登录失败:', error)
  }
}
</script>

<style lang="scss" scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
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

.login-form {
  .el-form-item {
    margin-bottom: 20px;
  }
  
  .login-btn {
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
  .login-card {
    padding: 24px;
  }
}
</style>
