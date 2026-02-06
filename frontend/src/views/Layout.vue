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
  BulbOutlined,
  BulbFilled,
  GlobalOutlined,
} from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useTheme } from '@/composables/useTheme'
import { useI18n, supportedLanguages, setLanguage } from '@/composables/useI18n'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const theme = useTheme()
const { t, currentLang } = useI18n()
const lang = ref(currentLang.value)

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

function onLangChange(val: string) {
  setLanguage(val)
}
</script>

<template>
  <a-layout :class="theme.currentTheme()" style="min-height: 100vh">
    <a-layout-sider
      v-model:collapsed="collapsed"
      collapsible
      :trigger="null"
      breakpoint="md"
      :collapsed-width="isMobile ? 0 : 80"
      :theme="theme.state.isDarkTheme ? 'dark' : 'dark'"
      @breakpoint="(broken: boolean) => (collapsed = broken)"
    >
      <div class="sider-logo">
        <span v-if="!collapsed" style="color: var(--color-primary-100, #008771); font-size: 22px; font-weight: 700;">3X-UI</span>
        <span v-else style="color: var(--color-primary-100, #008771); font-size: 18px; font-weight: 700;">3X</span>
      </div>
      <a-menu
        theme="dark"
        mode="inline"
        :selected-keys="selectedKeys"
        @click="handleMenuClick"
      >
        <a-menu-item key="dashboard">
          <DashboardOutlined />
          <span>{{ t('pages.index.title') }}</span>
        </a-menu-item>
        <a-menu-item key="inbounds">
          <CloudServerOutlined />
          <span>{{ t('pages.inbounds.title') }}</span>
        </a-menu-item>
        <a-menu-item key="settings">
          <SettingOutlined />
          <span>{{ t('pages.settings.title') }}</span>
        </a-menu-item>
        <a-menu-item key="xray">
          <ThunderboltOutlined />
          <span>{{ t('pages.xray.title') }}</span>
        </a-menu-item>
        <!-- Theme sub-menu -->
        <a-sub-menu key="theme">
          <template #title>
            <component :is="theme.state.isDarkTheme ? BulbFilled : BulbOutlined" />
            <span>{{ t('menu.theme') }}</span>
          </template>
          <a-menu-item key="toggle-dark" @click.stop="theme.toggleTheme()">
            <span>{{ t('menu.dark') }}</span>
            <a-switch size="small" :checked="theme.state.isDarkTheme" style="margin-left: 8px" />
          </a-menu-item>
          <a-menu-item v-if="theme.state.isDarkTheme" key="toggle-ultra" @click.stop="theme.toggleUltra()">
            <span>{{ t('menu.ultraDark') }}</span>
            <a-checkbox :checked="theme.state.isUltra" style="margin-left: 8px" />
          </a-menu-item>
        </a-sub-menu>
        <!-- Language -->
        <a-sub-menu key="language">
          <template #title>
            <GlobalOutlined />
            <span>{{ t('pages.settings.language') }}</span>
          </template>
          <a-menu-item v-for="l in supportedLanguages" :key="'lang-' + l.value" @click="onLangChange(l.value)">
            <span>{{ l.icon }}&nbsp;{{ l.name }}</span>
          </a-menu-item>
        </a-sub-menu>
        <a-menu-item key="logout" style="margin-top: 16px;">
          <LogoutOutlined />
          <span>{{ t('login', 'Logout') }}</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-header class="layout-header" :style="{ background: theme.state.isDarkTheme ? 'var(--dark-color-surface-100, #151f31)' : '#fff' }">
        <component
          :is="collapsed ? MenuUnfoldOutlined : MenuFoldOutlined"
          class="trigger"
          @click="collapsed = !collapsed"
        />
        <span class="header-title">3X-UI Panel</span>
      </a-layout-header>
      <a-layout-content class="layout-content" :style="{ background: theme.state.isDarkTheme ? 'var(--dark-color-background, #0a1222)' : '#f0f2f5' }">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<style>
.sider-logo {
  height: 48px;
  margin: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.layout-header {
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
  color: var(--color-primary-100, #008771);
}
.header-title {
  font-size: 18px;
  font-weight: 600;
  margin-left: 12px;
}
.layout-content {
  margin: 16px;
  padding: 24px;
  border-radius: 8px;
  min-height: 280px;
}
</style>
