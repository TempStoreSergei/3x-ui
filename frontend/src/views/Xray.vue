<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { SaveOutlined, ReloadOutlined, CodeOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import api from '@/api'
import type { ApiResponse } from '@/api'

const xrayConfig = ref('')
const loading = ref(false)
const saving = ref(false)

async function loadConfig() {
  loading.value = true
  try {
    const res = await api.post<any, ApiResponse>('/panel/xray/setting/all')
    if (res.success && res.obj) {
      // The config is the xrayTemplateConfig field
      const config = res.obj.xrayTemplateConfig || res.obj
      if (typeof config === 'string') {
        try {
          xrayConfig.value = JSON.stringify(JSON.parse(config), null, 2)
        } catch {
          xrayConfig.value = config
        }
      } else {
        xrayConfig.value = JSON.stringify(config, null, 2)
      }
    }
  } catch {
    message.error('Failed to load Xray config')
  } finally {
    loading.value = false
  }
}

async function saveConfig() {
  saving.value = true
  try {
    // Validate JSON
    JSON.parse(xrayConfig.value)
    const res = await api.post<any, ApiResponse>('/panel/xray/setting/update', {
      xrayTemplateConfig: xrayConfig.value
    })
    if (res.success) {
      message.success('Xray config saved')
    } else {
      message.error(res.msg || 'Failed to save config')
    }
  } catch (e: any) {
    message.error('Invalid JSON: ' + (e.message || ''))
  } finally {
    saving.value = false
  }
}

async function restartXray() {
  try {
    const res = await api.post<any, ApiResponse>('/panel/api/server/restartXray')
    if (res.success) {
      message.success('Xray restarted')
    } else {
      message.error(res.msg || 'Failed to restart Xray')
    }
  } catch {
    message.error('Failed to restart Xray')
  }
}

function formatJson() {
  try {
    const parsed = JSON.parse(xrayConfig.value)
    xrayConfig.value = JSON.stringify(parsed, null, 2)
    message.success('Formatted')
  } catch (e: any) {
    message.error('Invalid JSON: ' + (e.message || ''))
  }
}

onMounted(loadConfig)
</script>

<template>
  <a-spin :spinning="loading">
    <a-card size="small" style="margin-bottom: 16px">
      <a-space>
        <a-button type="primary" @click="saveConfig" :loading="saving">
          <SaveOutlined /> Save Config
        </a-button>
        <a-popconfirm title="Restart Xray?" @confirm="restartXray">
          <a-button>
            <ReloadOutlined /> Restart Xray
          </a-button>
        </a-popconfirm>
        <a-button @click="formatJson">
          <CodeOutlined /> Format JSON
        </a-button>
      </a-space>
    </a-card>

    <a-card title="Xray Template Config">
      <a-textarea
        v-model:value="xrayConfig"
        :auto-size="{ minRows: 20, maxRows: 50 }"
        style="font-family: 'Courier New', monospace; font-size: 13px"
        placeholder="Loading Xray configuration..."
      />
    </a-card>
  </a-spin>
</template>
