import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const Layout = () => import('@/views/Layout.vue')

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: Layout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue')
      },
      {
        path: 'inbounds',
        name: 'Inbounds',
        component: () => import('@/views/Inbounds.vue')
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/views/Settings.vue')
      },
      {
        path: 'xray',
        name: 'Xray',
        component: () => import('@/views/Xray.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  const isLogin = localStorage.getItem('isLogin')
  if (to.matched.some(record => record.meta.requiresAuth !== false) && !isLogin) {
    next({ name: 'Login' })
  } else {
    next()
  }
})

export default router
