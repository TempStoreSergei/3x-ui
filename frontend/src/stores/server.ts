import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'
import type { ApiResponse } from '@/api'

export interface ServerStatus {
  cpu: number
  mem: { current: number; total: number }
  swap: { current: number; total: number }
  disk: { current: number; total: number }
  xray: { state: string; errorMsg: string; version: string }
  uptime: number
  loads: number[]
  tcpCount: number
  udpCount: number
  netIO: { up: number; down: number }
  netTraffic: { sent: number; recv: number }
  publicIP: { ipv4: string; ipv6: string }
  appStats: { threads: number; mem: number; uptime: number }
}

export const useServerStore = defineStore('server', () => {
  const status = ref<ServerStatus | null>(null)
  const loading = ref(false)

  async function fetchStatus() {
    loading.value = true
    try {
      const res = await api.post<any, ApiResponse<ServerStatus>>('/panel/api/server/status')
      if (res.success) {
        status.value = res.obj
      }
    } finally {
      loading.value = false
    }
  }

  return { status, loading, fetchStatus }
})
