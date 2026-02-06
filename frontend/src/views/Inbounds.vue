<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  PlusOutlined,
  ThunderboltOutlined,
  DeleteOutlined,
  ToolOutlined,
  UserOutlined,
  ReloadOutlined,
  SwapOutlined,
  RocketOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useInboundStore } from '@/stores/inbound'
import { useFormat } from '@/composables/useFormat'
import api from '@/api'
import type { ApiResponse } from '@/api'

const store = useInboundStore()
const { sizeFormat } = useFormat()

const isSimpleMode = ref(localStorage.getItem('uiMode') === 'simple')
const showQuickSetup = ref(false)
const quickLoading = ref(false)
const quickOpts = ref({
  generateInbounds: true,
  configureDns: true,
  configureRouting: true,
})
const searchKey = ref('')

const filteredInbounds = computed(() => {
  if (!searchKey.value) return store.inbounds
  const key = searchKey.value.toLowerCase()
  return store.inbounds.filter(
    (i) => i.remark.toLowerCase().includes(key) || String(i.port).includes(key) || i.protocol.includes(key)
  )
})

const totalTraffic = computed(() => ({
  up: store.inbounds.reduce((s, i) => s + i.up, 0),
  down: store.inbounds.reduce((s, i) => s + i.down, 0),
}))

function toggleMode() {
  isSimpleMode.value = !isSimpleMode.value
  localStorage.setItem('uiMode', isSimpleMode.value ? 'simple' : 'advanced')
}

function protocolColor(protocol: string) {
  const colors: Record<string, string> = {
    vless: 'green',
    vmess: 'purple',
    trojan: 'orange',
    shadowsocks: 'cyan',
    wireguard: 'blue',
    socks: 'red',
    http: 'gold',
  }
  return colors[protocol] || 'default'
}

function parseClients(inbound: any): any[] {
  try {
    const settings = JSON.parse(inbound.settings || '{}')
    return settings.clients || []
  } catch {
    return []
  }
}

async function handleDelete(id: number) {
  const res = await store.deleteInbound(id)
  if (res.success) message.success('Inbound deleted')
  else message.error(res.msg || 'Failed to delete')
}

async function quickSetup() {
  quickLoading.value = true
  try {
    const res = await api.post<any, ApiResponse>('/panel/api/generator/quicksetup', quickOpts.value)
    if (res.success) {
      message.success('Quick setup completed!')
      showQuickSetup.value = false
      store.fetchInbounds()
    } else {
      message.error(res.msg || 'Quick setup failed')
    }
  } catch {
    message.error('Connection error')
  } finally {
    quickLoading.value = false
  }
}

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: 'Remark', dataIndex: 'remark', ellipsis: true },
  { title: 'Protocol', dataIndex: 'protocol', width: 120 },
  { title: 'Port', dataIndex: 'port', width: 80 },
  { title: 'Traffic ↑/↓', key: 'traffic', width: 180 },
  { title: 'Status', dataIndex: 'enable', width: 80 },
  { title: 'Clients', key: 'clients', width: 80 },
  { title: 'Actions', key: 'actions', width: 80 },
]

onMounted(() => {
  store.fetchInbounds()
})
</script>

<template>
  <div>
    <!-- Stats bar -->
    <a-row :gutter="[16, 16]" style="margin-bottom: 16px">
      <a-col :xs="12" :md="6">
        <a-card size="small">
          <a-statistic title="Upload / Download" :value-style="{ fontSize: '14px' }">
            <template #formatter>
              <SwapOutlined /> {{ sizeFormat(totalTraffic.up) }} / {{ sizeFormat(totalTraffic.down) }}
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="12" :md="6">
        <a-card size="small">
          <a-statistic title="Total Usage" :value="sizeFormat(totalTraffic.up + totalTraffic.down)" />
        </a-card>
      </a-col>
      <a-col :xs="12" :md="6">
        <a-card size="small">
          <a-statistic title="Inbounds" :value="store.inbounds.length" />
        </a-card>
      </a-col>
      <a-col :xs="12" :md="6">
        <a-card size="small">
          <a-statistic title="Online" :value="store.onlineUsers.length">
            <template #suffix>
              <a-tag color="green" size="small">users</a-tag>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- Quick Setup panel -->
    <a-card v-if="showQuickSetup" size="small" style="margin-bottom: 16px">
      <template #title>
        <RocketOutlined style="color: #722ed1" /> Quick Setup & Generator
      </template>
      <template #extra>
        <a-button size="small" @click="showQuickSetup = false">Close</a-button>
      </template>
      <a-row :gutter="16">
        <a-col :xs="24" :md="12">
          <a-card size="small" :bordered="false" :body-style="{ background: '#f6ffed', borderRadius: '8px' }">
            <h4 style="margin: 0 0 8px; color: #389e0d"><RocketOutlined /> One-Click Setup</h4>
            <p style="font-size: 12px; color: #666; margin-bottom: 12px">
              Generate optimized inbounds, configure anti-censorship DNS, and add Russia bypass routing.
            </p>
            <a-space direction="vertical">
              <a-checkbox v-model:checked="quickOpts.generateInbounds">Generate inbounds (VLESS, VMess, SS, Trojan)</a-checkbox>
              <a-checkbox v-model:checked="quickOpts.configureDns">Configure DNS (Cloudflare DoH + NextDNS)</a-checkbox>
              <a-checkbox v-model:checked="quickOpts.configureRouting">Add Russian site bypass routing</a-checkbox>
            </a-space>
            <a-popconfirm title="Apply all selected configurations?" @confirm="quickSetup">
              <a-button type="primary" block :loading="quickLoading" style="margin-top: 12px">
                <RocketOutlined /> Run Quick Setup
              </a-button>
            </a-popconfirm>
          </a-card>
        </a-col>
        <a-col :xs="24" :md="12">
          <a-card size="small" :bordered="false" :body-style="{ background: '#f0f5ff', borderRadius: '8px' }">
            <h4 style="margin: 0 0 8px; color: #1890ff">What happens</h4>
            <ul style="font-size: 12px; padding-left: 16px; margin: 0; color: #555">
              <li>10 diverse inbound configs are created</li>
              <li>Cloudflare DoH + NextDNS configured</li>
              <li>Russian sites go direct (bypass proxy)</li>
              <li>Ad/tracker domains blocked</li>
              <li>Each client gets a subscription URL</li>
            </ul>
          </a-card>
        </a-col>
      </a-row>
    </a-card>

    <!-- Toolbar -->
    <a-card size="small" style="margin-bottom: 16px">
      <a-space wrap>
        <template v-if="isSimpleMode">
          <a-button type="primary" @click="showQuickSetup = !showQuickSetup" style="background: #722ed1; border-color: #722ed1">
            <ThunderboltOutlined /> Quick Setup
          </a-button>
        </template>
        <template v-else>
          <a-button type="primary">
            <PlusOutlined /> Add Inbound
          </a-button>
          <a-button @click="showQuickSetup = !showQuickSetup" style="background: #722ed1; border-color: #722ed1; color: #fff">
            <ThunderboltOutlined /> Auto Generate
          </a-button>
        </template>
        <a-button @click="toggleMode">
          <component :is="isSimpleMode ? ToolOutlined : UserOutlined" />
          {{ isSimpleMode ? 'Advanced' : 'Simple' }}
        </a-button>
        <a-button @click="store.fetchInbounds()" :loading="store.loading">
          <ReloadOutlined /> Refresh
        </a-button>
        <a-input-search
          v-model:value="searchKey"
          placeholder="Search inbounds..."
          style="width: 200px"
          allow-clear
        />
      </a-space>
    </a-card>

    <!-- Inbound table -->
    <a-card size="small">
      <a-table
        :columns="columns"
        :data-source="filteredInbounds"
        :row-key="(r: any) => r.id"
        :loading="store.loading"
        :pagination="{ pageSize: 20, showSizeChanger: true }"
        size="small"
        :scroll="{ x: 800 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'protocol'">
            <a-tag :color="protocolColor(record.protocol)">{{ record.protocol.toUpperCase() }}</a-tag>
          </template>
          <template v-if="column.key === 'traffic'">
            <span style="color: #52c41a">↑ {{ sizeFormat(record.up) }}</span>
            <span style="margin: 0 4px">/</span>
            <span style="color: #1890ff">↓ {{ sizeFormat(record.down) }}</span>
          </template>
          <template v-if="column.dataIndex === 'enable'">
            <a-tag :color="record.enable ? 'success' : 'default'">
              {{ record.enable ? 'ON' : 'OFF' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'clients'">
            <a-tag color="blue">{{ parseClients(record).length }}</a-tag>
          </template>
          <template v-if="column.key === 'actions'">
            <a-popconfirm title="Delete this inbound?" @confirm="handleDelete(record.id)">
              <a-button type="link" danger size="small">
                <DeleteOutlined />
              </a-button>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>
