import axios from 'axios'
import router from '@/router'

const api = axios.create({
  timeout: 30000,
})

// Request interceptor - convert data to form-urlencoded (matching Go backend expectations)
api.interceptors.request.use((config) => {
  if (config.data && !(config.data instanceof FormData)) {
    config.headers['Content-Type'] = 'application/x-www-form-urlencoded'
    config.data = new URLSearchParams(
      Object.entries(config.data).reduce((acc, [key, val]) => {
        acc[key] = val == null ? '' : typeof val === 'object' ? JSON.stringify(val) : String(val)
        return acc
      }, {} as Record<string, string>)
    ).toString()
  }
  return config
})

// Response interceptor
api.interceptors.response.use(
  (response) => {
    if (response.data && typeof response.data === 'object') {
      return response.data
    }
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('isLogin')
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

export default api

// Type for standard API response
export interface ApiResponse<T = any> {
  success: boolean
  msg: string
  obj: T
}
