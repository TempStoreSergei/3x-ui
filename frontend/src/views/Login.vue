<script setup lang="ts">
import { reactive, ref, onMounted, nextTick } from 'vue'
import { SettingOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useI18n, supportedLanguages, setLanguage, t } from '@/composables/useI18n'
import { useTheme } from '@/composables/useTheme'
import api from '@/api'
import type { ApiResponse } from '@/api'
import { useRouter } from 'vue-router'

const router = useRouter()
const { currentLang } = useI18n()
const theme = useTheme()

const spinning = ref(false)
const fetched = ref(false)
const twoFactorEnable = ref(false)
const animationStarted = ref(false)
const lang = ref(currentLang.value)

const user = reactive({
  username: '',
  password: '',
  twoFactorCode: '',
})

async function getTwoFactorEnable() {
  try {
    const res = await api.post<any, ApiResponse>('/getTwoFactorEnable')
    if (res.success) {
      twoFactorEnable.value = res.obj
      fetched.value = true
      await nextTick()
      if (!animationStarted.value) {
        animationStarted.value = true
        initHeadline()
      }
    }
  } catch {
    fetched.value = true
  }
}

async function handleLogin() {
  spinning.value = true
  try {
    const res = await api.post<any, ApiResponse>('/login', user)
    if (res.success) {
      localStorage.setItem('isLogin', 'true')
      router.push('/')
    } else {
      message.error(res.msg || 'Login failed')
    }
  } catch {
    message.error('Connection error')
  } finally {
    spinning.value = false
  }
}

function onLangChange(val: string) {
  setLanguage(val)
}

// Text animation - matching original exactly
function initHeadline() {
  const headlines = document.querySelectorAll('.headline')
  headlines.forEach((headline) => {
    const first = headline.querySelector('.is-visible')
    if (!first) return
    setTimeout(() => hideWord(first as HTMLElement, 2000), 2000)
  })
}

function hideWord(word: HTMLElement, delay: number) {
  const next = takeNext(word)
  switchWord(word, next)
  setTimeout(() => hideWord(next, delay), delay)
}

function takeNext(word: HTMLElement): HTMLElement {
  return (word.nextElementSibling || word.parentElement!.firstElementChild) as HTMLElement
}

function switchWord(oldWord: HTMLElement, newWord: HTMLElement) {
  oldWord.classList.remove('is-visible')
  oldWord.classList.add('is-hidden')
  newWord.classList.remove('is-hidden')
  newWord.classList.add('is-visible')
}

onMounted(() => {
  getTwoFactorEnable()
})
</script>

<template>
  <a-layout :class="[theme.currentTheme(), 'login-app']">
    <a-layout-content class="under min-h-0">
      <div class="waves-header">
        <div class="waves-inner-header"></div>
        <svg class="waves" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
          viewBox="0 24 150 28" preserveAspectRatio="none" shape-rendering="auto">
          <defs>
            <path id="gentle-wave" d="M-160 44c30 0 58-18 88-18s 58 18 88 18 58-18 88-18 58 18 88 18 v44h-352z" />
          </defs>
          <g class="parallax">
            <use xlink:href="#gentle-wave" x="48" y="0" fill="rgba(0, 135, 113, 0.08)" />
            <use xlink:href="#gentle-wave" x="48" y="3" fill="rgba(0, 135, 113, 0.08)" />
            <use xlink:href="#gentle-wave" x="48" y="5" fill="rgba(0, 135, 113, 0.08)" />
            <use xlink:href="#gentle-wave" x="48" y="7" fill="#c7ebe2" />
          </g>
        </svg>
      </div>
      <a-row type="flex" justify="center" align="middle" class="h-100 overflow-y-auto overflow-x-hidden">
        <a-col :xs="22" :sm="12" :md="10" :lg="8" :xl="6" :xxl="5" id="login" class="my-3rem">
          <template v-if="!fetched">
            <div class="text-center">
              <a-spin size="large" />
            </div>
          </template>
          <template v-else>
            <div class="setting-section">
              <a-popover :overlay-class-name="theme.currentTheme()" :title="t('menu.settings')" placement="bottomRight" trigger="click">
                <template #content>
                  <a-space direction="vertical" :size="10">
                    <!-- Theme switch -->
                    <a-space direction="horizontal" size="small">
                      <a-switch size="small" :checked="theme.state.isDarkTheme" @change="theme.toggleTheme()" />
                      <span>{{ t('menu.dark') }}</span>
                    </a-space>
                    <a-space v-if="theme.state.isDarkTheme" direction="horizontal" size="small">
                      <a-checkbox :checked="theme.state.isUltra" @click="theme.toggleUltra()" />
                      <span>{{ t('menu.ultraDark') }}</span>
                    </a-space>
                    <span>{{ t('pages.settings.language') }}</span>
                    <a-select v-model:value="lang" style="width: 100%" :popup-class-name="theme.currentTheme()" @change="onLangChange">
                      <a-select-option v-for="l in supportedLanguages" :key="l.value" :value="l.value">
                        <span>{{ l.icon }}</span>&nbsp;&nbsp;<span>{{ l.name }}</span>
                      </a-select-option>
                    </a-select>
                  </a-space>
                </template>
                <a-button shape="circle">
                  <template #icon><SettingOutlined /></template>
                </a-button>
              </a-popover>
            </div>
            <a-row type="flex" justify="center">
              <a-col :style="{ width: '100%' }">
                <h2 class="title headline zoom">
                  <span class="words-wrapper">
                    <b class="is-visible">{{ t('pages.login.hello') }}</b>
                    <b>{{ t('pages.login.title') }}</b>
                  </span>
                </h2>
              </a-col>
            </a-row>
            <a-row type="flex" justify="center">
              <a-col :span="24">
                <form @submit.prevent="handleLogin">
                  <a-space direction="vertical" size="middle" style="width: 100%">
                    <a-form-item>
                      <a-input
                        v-model:value="user.username"
                        autocomplete="username"
                        :placeholder="t('username')"
                        autofocus
                        size="large"
                      >
                        <template #prefix>
                          <span class="anticon fs-1rem">👤</span>
                        </template>
                      </a-input>
                    </a-form-item>
                    <a-form-item>
                      <a-input-password
                        v-model:value="user.password"
                        autocomplete="current-password"
                        :placeholder="t('password')"
                        size="large"
                      >
                        <template #prefix>
                          <span class="anticon fs-1rem">🔒</span>
                        </template>
                      </a-input-password>
                    </a-form-item>
                    <a-form-item v-if="twoFactorEnable">
                      <a-input
                        v-model:value="user.twoFactorCode"
                        autocomplete="one-time-code"
                        :placeholder="t('twoFactorCode')"
                        size="large"
                      >
                        <template #prefix>
                          <span class="anticon fs-1rem">🔑</span>
                        </template>
                      </a-input>
                    </a-form-item>
                    <a-form-item>
                      <a-row justify="center" class="centered">
                        <div class="wave-btn-bg wave-btn-bg-cl h-50px mt-1rem"
                          :style="spinning ? 'width: 52px' : 'display: inline-block'">
                          <a-button
                            class="ant-btn-primary-login"
                            type="primary"
                            :loading="spinning"
                            html-type="submit"
                          >
                            {{ spinning ? '' : t('login') }}
                          </a-button>
                        </div>
                      </a-row>
                    </a-form-item>
                  </a-space>
                </form>
              </a-col>
            </a-row>
          </template>
        </a-col>
      </a-row>
    </a-layout-content>
  </a-layout>
</template>

<style>
/* Import original 3x-ui styles for login page */
.login-app {
  overflow: hidden;
  min-height: 100vh;
}
.login-app * {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
.login-app .under {
  background-color: #c7ebe2;
  z-index: 0;
}
.login-app.dark .under {
  background-color: var(--dark-color-login-wave, #0f2d32);
}
.login-app .waves-header {
  position: fixed;
  width: 100%;
  text-align: center;
  background-color: #dbf5ed;
  color: #fff;
  z-index: -1;
}
.login-app.dark .waves-header {
  background-color: var(--dark-color-login-background, #0a2227);
}
.login-app .waves {
  position: relative;
  width: 100%;
  margin-bottom: -7px;
  min-height: 100px;
  max-height: 150px;
}
.login-app .parallax > use {
  animation: wave-move 25s cubic-bezier(.55,.5,.45,.5) infinite;
}
.login-app .parallax > use:nth-child(1) { animation-delay: -2s; animation-duration: 7s; }
.login-app .parallax > use:nth-child(2) { animation-delay: -3s; animation-duration: 10s; }
.login-app .parallax > use:nth-child(3) { animation-delay: -4s; animation-duration: 13s; }
.login-app .parallax > use:nth-child(4) { animation-delay: -5s; animation-duration: 20s; }
@keyframes wave-move {
  0% { transform: translate3d(-90px, 0, 0); }
  100% { transform: translate3d(85px, 0, 0); }
}
#login {
  animation: charge 0.5s both;
  background-color: #fff;
  border-radius: 2rem;
  padding: 4rem 3rem;
  transition: all 0.3s;
  user-select: none;
}
#login:hover {
  box-shadow: 0 2px 8px rgb(0 0 0 / .09);
}
.login-app.dark #login {
  background-color: var(--dark-color-surface-100, #151f31);
}
.login-app.dark h1,
.login-app.dark h2 {
  color: #fff;
}
@keyframes charge {
  from { transform: translateY(30px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}
.login-app h1 {
  text-align: center;
  height: 110px;
}
.login-app .ant-input,
.login-app .ant-input-affix-wrapper {
  height: 50px;
  border-radius: 30px;
}
.login-app .ant-input-affix-wrapper .ant-input {
  height: auto;
}
.login-app .centered {
  display: flex;
  text-align: center;
  align-items: center;
  justify-content: center;
  width: 100%;
}
.login-app .title {
  font-size: 2rem;
  margin-block-end: 2rem;
  text-align: center;
}
.login-app .title b {
  font-weight: bold !important;
}
.login-app .wave-btn-bg {
  position: relative;
  border-radius: 25px;
  width: 100%;
  transition: all 0.3s cubic-bezier(.645,.045,.355,1);
}
.login-app .ant-btn-primary-login {
  width: 100%;
  height: 50px;
  border-radius: 25px;
  font-size: 16px;
  background-color: var(--color-primary-100, #008771);
  border-color: var(--color-primary-100, #008771);
}
.login-app .ant-btn-primary-login:hover,
.login-app .ant-btn-primary-login:focus {
  background-color: #065;
  border-color: #065;
}
.login-app .setting-section {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 10;
}
.login-app .h-100 { height: 100vh; }
.login-app .min-h-0 { min-height: 0; }
.login-app .overflow-y-auto { overflow-y: auto; }
.login-app .overflow-x-hidden { overflow-x: hidden; }
.login-app .text-center { text-align: center; }
.login-app .my-3rem { margin-top: 3rem; margin-bottom: 3rem; }
.login-app .fs-1rem { font-size: 1rem; }
.login-app .h-50px { height: 50px; }
.login-app .mt-1rem { margin-top: 1rem; }

/* Animated words */
.login-app .words-wrapper { display: inline-block; position: relative; }
.login-app .words-wrapper b { display: inline-block; position: absolute; white-space: nowrap; left: 0; top: 0; width: 100%; }
.login-app .words-wrapper b.is-visible { position: relative; opacity: 1; }
.login-app .words-wrapper b.is-hidden { opacity: 0; }

/* CSS variables for themes */
:root {
  --color-primary-100: #008771;
  --dark-color-background: #0a1222;
  --dark-color-surface-100: #151f31;
  --dark-color-login-background: #0a2227;
  --dark-color-login-wave: #0f2d32;
}
</style>
