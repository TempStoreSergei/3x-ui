<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  DashboardOutlined,
  CloudServerOutlined,
  SettingOutlined,
  ThunderboltOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
} from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const collapsed = ref(false)
const isMobile = ref(window.innerWidth < 768)

function handleResize() {
  isMobile.value = window.innerWidth < 768
  if (isMobile.value) collapsed.value = true
}

onMounted(() => window.addEventListener('resize', handleResize))
onUnmounted(() => window.removeEventListener('resize', handleResize))

const selectedKeys = computed(() => {
  const path = route.path
  if (path === '/') return ['dashboard']
  if (path.startsWith('/inbounds')) return ['inbounds']
  if (path.startsWith('/settings')) return ['settings']
  if (path.startsWith('/xray')) return ['xray']
  return ['dashboard']
})

function navigate(key: string) {
  const routes: Record<string, string> = {
    dashboard: '/',
    inbounds: '/inbounds',
    settings: '/settings',
    xray: '/xray',
  }
  if (routes[key]) router.push(routes[key])
}

function handleMenuClick({ key }: { key: string }) {
  if (key === 'logout') {
    auth.logout()
    return
  }
  navigate(key)
}
</script>

<template>
  <a-layout style="min-height: 100vh">
    <a-layout-sider
      v-model:collapsed="collapsed"
      collapsible
      :trigger="null"
      breakpoint="md"
      :collapsed-width="isMobile ? 0 : 80"
      @breakpoint="(broken: boolean) => (collapsed = broken)"
    >
      <div class="sider-logo">
        <span v-if="!collapsed">3X-UI</span>
        <span v-else>3X</span>
      </div>
      <a-menu
        theme="dark"
        mode="inline"
        :selected-keys="selectedKeys"
        @click="handleMenuClick"
      >
        <a-menu-item key="dashboard">
          <DashboardOutlined />
          <span>Dashboard</span>
        </a-menu-item>
        <a-menu-item key="inbounds">
          <CloudServerOutlined />
          <span>Inbounds</span>
        </a-menu-item>
        <a-menu-item key="settings">
          <SettingOutlined />
          <span>Settings</span>
        </a-menu-item>
        <a-menu-item key="xray">
          <ThunderboltOutlined />
          <span>Xray</span>
        </a-menu-item>
        <a-menu-item key="logout" class="logout-item">
          <LogoutOutlined />
          <span>Logout</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-header class="layout-header">
        <component
          :is="collapsed ? MenuUnfoldOutlined : MenuFoldOutlined"
          class="trigger"
          @click="collapsed = !collapsed"
        />
        <span class="header-title">3X-UI Panel</span>
      </a-layout-header>
      <a-layout-content class="layout-content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<style scoped>
.sider-logo {
  height: 48px;
  margin: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 20px;
  font-weight: 700;
  letter-spacing: 1px;
}
.layout-header {
  background: #fff;
  padding: 0 16px;
  display: flex;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}
.trigger {
  font-size: 18px;
  cursor: pointer;
  padding: 0 8px;
  transition: color 0.3s;
}
.trigger:hover {
  color: #1890ff;
}
.header-title {
  font-size: 18px;
  font-weight: 600;
  margin-left: 12px;
}
.layout-content {
  margin: 16px;
  padding: 24px;
  background: #fff;
  border-radius: 8px;
  min-height: 280px;
}
.logout-item {
  margin-top: auto;
}
</style>
