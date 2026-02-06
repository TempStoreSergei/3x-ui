import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'
import type { ApiResponse } from '@/api'

export interface AllSetting {
  webListen: string
  webDomain: string
  webPort: number
  webCertFile: string
  webKeyFile: string
  webBasePath: string
  sessionMaxAge: number
  pageSize: number
  expireDiff: number
  trafficDiff: number
  remarkModel: string
  datepicker: string
  tgBotEnable: boolean
  tgBotToken: string
  tgBotChatId: string
  tgBotProxy: string
  tgBotAPIServer: string
  tgRunTime: string
  tgBotBackup: boolean
  tgBotLoginNotify: boolean
  tgCpu: number
  tgLang: string
  subEnable: boolean
  subListen: string
  subPort: number
  subPath: string
  subDomain: string
  subCertFile: string
  subKeyFile: string
  subUpdates: number
  subEncrypt: boolean
  subShowInfo: boolean
  subURI: string
  subJsonEnable: boolean
  subJsonURI: string
  subTitle: string
  subSupportUrl: string
  subProfileUrl: string
  cloudflareAPIToken: string
  cloudflareZoneID: string
  nextDnsProfileId: string
  timeLocation: string
  [key: string]: any
}

export const useSettingStore = defineStore('setting', () => {
  const allSetting = ref<AllSetting | null>(null)
  const loading = ref(false)

  async function fetchSettings() {
    loading.value = true
    try {
      const res = await api.post<any, ApiResponse<AllSetting>>('/panel/setting/all')
      if (res.success) {
        allSetting.value = res.obj
      }
    } finally {
      loading.value = false
    }
  }

  async function saveSettings() {
    if (!allSetting.value) return
    loading.value = true
    try {
      return await api.post<any, ApiResponse>('/panel/setting/update', allSetting.value)
    } finally {
      loading.value = false
    }
  }

  return { allSetting, loading, fetchSettings, saveSettings }
})
