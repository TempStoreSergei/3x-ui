import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  base: '/dist/',
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  build: {
    outDir: '../web/dist/',
    emptyOutDir: true
  },
  server: {
    proxy: {
      '/panel/api': {
        target: 'http://localhost:2053',
        changeOrigin: true
      },
      '/ws': {
        target: 'http://localhost:2053',
        ws: true
      }
    }
  }
})
