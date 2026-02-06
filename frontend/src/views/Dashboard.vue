<script setup lang="ts">
import { onMounted, onUnmounted, computed } from 'vue'
import {
  DashboardOutlined,
  CloudServerOutlined,
  HddOutlined,
  WifiOutlined,
  ClockCircleOutlined,
  ApiOutlined,
} from '@ant-design/icons-vue'
import { useServerStore } from '@/stores/server'
import { useFormat } from '@/composables/useFormat'

const server = useServerStore()
const { sizeFormat, cpuFormat, timeFormat } = useFormat()

let timer: ReturnType<typeof setInterval> | null = null

const memPercent = computed(() => {
  if (!server.status?.mem.total) return 0
  return Number(((server.status.mem.current / server.status.mem.total) * 100).toFixed(1))
})

const diskPercent = computed(() => {
  if (!server.status?.disk.total) return 0
  return Number(((server.status.disk.current / server.status.disk.total) * 100).toFixed(1))
})

const xrayRunning = computed(() => server.status?.xray?.state === 'running')

onMounted(() => {
  server.fetchStatus()
  timer = setInterval(() => server.fetchStatus(), 2000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <a-spin :spinning="server.loading && !server.status">
    <a-row :gutter="[16, 16]">
      <!-- CPU -->
      <a-col :xs="24" :sm="12" :md="8" :lg="6">
        <a-card size="small">
          <a-statistic title="CPU Usage" :value="server.status ? cpuFormat(server.status.cpu) : '-'">
            <template #prefix><DashboardOutlined /></template>
          </a-statistic>
          <a-progress
            :percent="server.status?.cpu ?? 0"
            :stroke-color="(server.status?.cpu ?? 0) > 80 ? '#ff4d4f' : '#1890ff'"
            :show-info="false"
            size="small"
            style="margin-top: 8px"
          />
        </a-card>
      </a-col>

      <!-- Memory -->
      <a-col :xs="24" :sm="12" :md="8" :lg="6">
        <a-card size="small">
          <a-statistic
            title="Memory"
            :value="server.status ? sizeFormat(server.status.mem.current) : '-'"
            :suffix="server.status ? '/ ' + sizeFormat(server.status.mem.total) : ''"
          >
            <template #prefix><CloudServerOutlined /></template>
          </a-statistic>
          <a-progress :percent="memPercent" :show-info="false" size="small" style="margin-top: 8px" />
        </a-card>
      </a-col>

      <!-- Disk -->
      <a-col :xs="24" :sm="12" :md="8" :lg="6">
        <a-card size="small">
          <a-statistic
            title="Disk"
            :value="server.status ? sizeFormat(server.status.disk.current) : '-'"
            :suffix="server.status ? '/ ' + sizeFormat(server.status.disk.total) : ''"
          >
            <template #prefix><HddOutlined /></template>
          </a-statistic>
          <a-progress
            :percent="diskPercent"
            :stroke-color="diskPercent > 90 ? '#ff4d4f' : '#52c41a'"
            :show-info="false"
            size="small"
            style="margin-top: 8px"
          />
        </a-card>
      </a-col>

      <!-- Uptime -->
      <a-col :xs="24" :sm="12" :md="8" :lg="6">
        <a-card size="small">
          <a-statistic
            title="System Uptime"
            :value="server.status ? timeFormat(server.status.uptime) : '-'"
          >
            <template #prefix><ClockCircleOutlined /></template>
          </a-statistic>
        </a-card>
      </a-col>

      <!-- Xray Status -->
      <a-col :xs="24" :sm="12" :md="8" :lg="6">
        <a-card size="small">
          <a-statistic title="Xray Status">
            <template #prefix><ApiOutlined /></template>
            <template #formatter>
              <a-tag :color="xrayRunning ? 'success' : 'error'">
                {{ xrayRunning ? 'Running' : 'Stopped' }}
              </a-tag>
              <span v-if="server.status?.xray?.version" style="font-size: 12px; color: #888; margin-left: 4px">
                v{{ server.status.xray.version }}
              </span>
            </template>
          </a-statistic>
        </a-card>
      </a-col>

      <!-- Network I/O -->
      <a-col :xs="24" :sm="12" :md="8" :lg="6">
        <a-card size="small">
          <a-statistic title="Network I/O">
            <template #prefix><WifiOutlined /></template>
            <template #formatter>
              <div style="font-size: 14px">
                <span style="color: #52c41a">↑ {{ server.status ? sizeFormat(server.status.netIO.up) + '/s' : '-' }}</span>
                <span style="margin: 0 8px">|</span>
                <span style="color: #1890ff">↓ {{ server.status ? sizeFormat(server.status.netIO.down) + '/s' : '-' }}</span>
              </div>
            </template>
          </a-statistic>
        </a-card>
      </a-col>

      <!-- Connections -->
      <a-col :xs="24" :sm="12" :md="8" :lg="6">
        <a-card size="small">
          <a-statistic title="Connections">
            <template #formatter>
              <a-space>
                <a-tag color="blue">TCP: {{ server.status?.tcpCount ?? 0 }}</a-tag>
                <a-tag color="purple">UDP: {{ server.status?.udpCount ?? 0 }}</a-tag>
              </a-space>
            </template>
          </a-statistic>
        </a-card>
      </a-col>

      <!-- Public IP -->
      <a-col :xs="24" :sm="12" :md="8" :lg="6">
        <a-card size="small">
          <a-statistic title="Public IP">
            <template #formatter>
              <div style="font-size: 12px; word-break: break-all">
                <div v-if="server.status?.publicIP?.ipv4">IPv4: {{ server.status.publicIP.ipv4 }}</div>
                <div v-if="server.status?.publicIP?.ipv6">IPv6: {{ server.status.publicIP.ipv6 }}</div>
                <span v-if="!server.status?.publicIP?.ipv4 && !server.status?.publicIP?.ipv6">-</span>
              </div>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>
  </a-spin>
</template>
