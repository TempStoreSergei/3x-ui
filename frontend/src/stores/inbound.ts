import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api'
import type { ApiResponse } from '@/api'

export interface Inbound {
  id: number
  userId: number
  up: number
  down: number
  total: number
  remark: string
  enable: boolean
  expiryTime: number
  listen: string
  port: number
  protocol: string
  settings: string
  streamSettings: string
  tag: string
  sniffing: string
  allocate: string
  clientStats: any[] | null
}

export const useInboundStore = defineStore('inbound', () => {
  const inbounds = ref<Inbound[]>([])
  const loading = ref(false)
  const onlineUsers = ref<string[]>([])

  const totalUp = computed(() => inbounds.value.reduce((s, i) => s + i.up, 0))
  const totalDown = computed(() => inbounds.value.reduce((s, i) => s + i.down, 0))

  async function fetchInbounds() {
    loading.value = true
    try {
      const res = await api.get<any, ApiResponse<Inbound[]>>('/panel/api/inbounds/list')
      if (res.success) {
        inbounds.value = res.obj || []
      }
    } finally {
      loading.value = false
    }
  }

  async function deleteInbound(id: number) {
    const res = await api.post<any, ApiResponse>(`/panel/api/inbounds/del/${id}`)
    if (res.success) await fetchInbounds()
    return res
  }

  async function fetchOnlines() {
    try {
      const res = await api.post<any, ApiResponse<string[]>>('/panel/api/inbounds/onlines')
      if (res.success) onlineUsers.value = res.obj || []
    } catch {}
  }

  return { inbounds, loading, onlineUsers, totalUp, totalDown, fetchInbounds, deleteInbound, fetchOnlines }
})
