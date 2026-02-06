import { ref } from 'vue'
import { allTranslations } from '../i18n/translations'

export interface Language {
  name: string
  value: string
  icon: string
}

export const supportedLanguages: Language[] = [
  { name: 'العربية', value: 'ar-EG', icon: '🇪🇬' },
  { name: 'English', value: 'en-US', icon: '🇺🇸' },
  { name: 'فارسی', value: 'fa-IR', icon: '🇮🇷' },
  { name: '简体中文', value: 'zh-CN', icon: '🇨🇳' },
  { name: '繁體中文', value: 'zh-TW', icon: '🇹🇼' },
  { name: '日本語', value: 'ja-JP', icon: '🇯🇵' },
  { name: 'Русский', value: 'ru-RU', icon: '🇷🇺' },
  { name: 'Tiếng Việt', value: 'vi-VN', icon: '🇻🇳' },
  { name: 'Español', value: 'es-ES', icon: '🇪🇸' },
  { name: 'Indonesian', value: 'id-ID', icon: '🇮🇩' },
  { name: 'Українська', value: 'uk-UA', icon: '🇺🇦' },
  { name: 'Türkçe', value: 'tr-TR', icon: '🇹🇷' },
  { name: 'Português', value: 'pt-BR', icon: '🇧🇷' },
]

function getCookie(name: string): string {
  const match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'))
  return match ? decodeURIComponent(match[2]) : ''
}

function setCookie(name: string, value: string, days: number) {
  const d = new Date()
  d.setTime(d.getTime() + days * 24 * 60 * 60 * 1000)
  document.cookie = `${name}=${encodeURIComponent(value)};expires=${d.toUTCString()};path=/`
}

function detectLanguage(): string {
  const cookie = getCookie('lang')
  if (cookie && supportedLanguages.some(l => l.value === cookie)) return cookie

  const nav = navigator.language || ''
  const prefix = nav.split('-')[0].toLowerCase()
  const map: Record<string, string> = {
    ar: 'ar-EG', fa: 'fa-IR', ja: 'ja-JP', ru: 'ru-RU',
    vi: 'vi-VN', es: 'es-ES', id: 'id-ID', uk: 'uk-UA',
    tr: 'tr-TR', pt: 'pt-BR', zh: 'zh-CN',
  }
  if (map[prefix]) return map[prefix]
  if (nav.startsWith('zh-TW') || nav.startsWith('zh-Hant')) return 'zh-TW'
  if (nav.startsWith('en')) return 'en-US'
  return 'en-US'
}

const currentLang = ref(detectLanguage())

const translations = ref<Record<string, Record<string, string>>>({ ...allTranslations })

export function setLanguage(lang: string) {
  if (!supportedLanguages.some(l => l.value === lang)) lang = 'en-US'
  setCookie('lang', lang, 150)
  currentLang.value = lang
  window.location.reload()
}

export function t(key: string, fallback?: string): string {
  const dict = translations.value[currentLang.value]
  if (dict && dict[key]) return dict[key]
  const en = translations.value['en-US']
  if (en && en[key]) return en[key]
  return fallback || key
}

export function useI18n() {
  return {
    currentLang,
    supportedLanguages,
    t,
    setLanguage,
    translations,
  }
}
