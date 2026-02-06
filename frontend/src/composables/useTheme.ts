import { reactive } from 'vue'

interface ThemeState {
  isDarkTheme: boolean
  isUltra: boolean
}

const state = reactive<ThemeState>({
  isDarkTheme: localStorage.getItem('dark-mode') === 'true',
  isUltra: localStorage.getItem('isUltraDarkThemeEnabled') === 'true',
})

function initTheme() {
  if (state.isUltra) {
    document.documentElement.setAttribute('data-theme', 'ultra-dark')
  }
  document.body.setAttribute('class', state.isDarkTheme ? 'dark' : 'light')
}

function toggleTheme() {
  state.isDarkTheme = !state.isDarkTheme
  localStorage.setItem('dark-mode', String(state.isDarkTheme))
  document.body.setAttribute('class', state.isDarkTheme ? 'dark' : 'light')
  if (!state.isDarkTheme) {
    state.isUltra = false
    localStorage.setItem('isUltraDarkThemeEnabled', 'false')
    document.documentElement.removeAttribute('data-theme')
  }
}

function toggleUltra() {
  if (!state.isDarkTheme) return
  state.isUltra = !state.isUltra
  if (state.isUltra) {
    document.documentElement.setAttribute('data-theme', 'ultra-dark')
  } else {
    document.documentElement.removeAttribute('data-theme')
  }
  localStorage.setItem('isUltraDarkThemeEnabled', String(state.isUltra))
}

export function useTheme() {
  return {
    state,
    currentTheme: () => state.isDarkTheme ? 'dark' : 'light',
    initTheme,
    toggleTheme,
    toggleUltra,
  }
}
