import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'
import type { ApiResponse } from '@/api'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const isLogin = ref(!!localStorage.getItem('isLogin'))
  const loading = ref(false)

  async function login(username: string, password: string, loginSecret: string = '') {
    loading.value = true
    try {
      const res = await api.post<any, ApiResponse>('/login', { username, password, loginSecret })
      if (res.success) {
        isLogin.value = true
        localStorage.setItem('isLogin', 'true')
        router.push('/')
      }
      return res
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    try {
      await api.post('/logout')
    } catch {}
    isLogin.value = false
    localStorage.removeItem('isLogin')
    router.push('/login')
  }

  return { isLogin, loading, login, logout }
})
