<template>
  <a-layout style="min-height: 100vh">
    <a-layout-sider collapsible>
      <div style="height: 32px; margin: 16px; color: #fff; text-align: center; font-weight: bold">
        3X-UI
      </div>
      <a-menu theme="dark" mode="inline">
        <a-menu-item key="dashboard" @click="$router.push('/')">Dashboard</a-menu-item>
        <a-menu-item key="inbounds" @click="$router.push('/inbounds')">Inbounds</a-menu-item>
        <a-menu-item key="settings" @click="$router.push('/settings')">Settings</a-menu-item>
        <a-menu-item key="xray" @click="$router.push('/xray')">Xray</a-menu-item>
      </a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-header style="background: #fff; padding: 0 16px">
        <span style="font-size: 18px; font-weight: bold">3X-UI Panel</span>
      </a-layout-header>
      <a-layout-content style="margin: 16px">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
