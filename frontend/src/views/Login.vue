<template>
  <div style="display: flex; justify-content: center; align-items: center; min-height: 100vh">
    <a-card title="Login" style="width: 400px">
      <a-form @finish="handleLogin" layout="vertical">
        <a-form-item label="Username" name="username">
          <a-input v-model:value="form.username" placeholder="Username" />
        </a-form-item>
        <a-form-item label="Password" name="password">
          <a-input-password v-model:value="form.password" placeholder="Password" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" block>Login</a-button>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const form = reactive({ username: '', password: '' })

function handleLogin() {
  localStorage.setItem('isLogin', 'true')
  router.push('/')
}
</script>
