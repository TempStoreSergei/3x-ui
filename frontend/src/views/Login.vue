<script setup lang="ts">
import { reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const form = reactive({
  username: '',
  password: '',
  loginSecret: ''
})
const showSecret = ref(false)
const loading = ref(false)

async function handleLogin() {
  if (!form.username || !form.password) {
    message.warning('Please enter username and password')
    return
  }
  loading.value = true
  try {
    const res = await auth.login(form.username, form.password, form.loginSecret)
    if (!res.success) {
      message.error(res.msg || 'Login failed')
    }
  } catch {
    message.error('Connection error')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-container">
    <a-card class="login-card">
      <div class="login-header">
        <h1>3X-UI</h1>
        <p style="color: #888">Panel Management</p>
      </div>
      <a-form layout="vertical" @finish="handleLogin">
        <a-form-item label="Username">
          <a-input
            v-model:value="form.username"
            size="large"
            placeholder="Enter username"
            @pressEnter="handleLogin"
          >
            <template #prefix><UserOutlined /></template>
          </a-input>
        </a-form-item>
        <a-form-item label="Password">
          <a-input-password
            v-model:value="form.password"
            size="large"
            placeholder="Enter password"
            @pressEnter="handleLogin"
          >
            <template #prefix><LockOutlined /></template>
          </a-input-password>
        </a-form-item>
        <a-form-item v-if="showSecret" label="Secret Token">
          <a-input
            v-model:value="form.loginSecret"
            size="large"
            placeholder="Login secret"
          />
        </a-form-item>
        <div style="margin-bottom: 16px">
          <a-checkbox v-model:checked="showSecret">Use login secret</a-checkbox>
        </div>
        <a-button
          type="primary"
          html-type="submit"
          size="large"
          block
          :loading="loading"
        >
          Login
        </a-button>
      </a-form>
      <div style="text-align: center; margin-top: 16px">
        <a-select
          default-value="en"
          style="width: 120px"
          size="small"
        >
          <a-select-option value="en">English</a-select-option>
          <a-select-option value="zh">中文</a-select-option>
          <a-select-option value="fa">فارسی</a-select-option>
          <a-select-option value="ru">Русский</a-select-option>
        </a-select>
      </div>
    </a-card>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.login-card {
  width: 420px;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}
.login-header {
  text-align: center;
  margin-bottom: 24px;
}
.login-header h1 {
  margin: 0;
  font-size: 32px;
  font-weight: 700;
  color: #1890ff;
}
</style>
