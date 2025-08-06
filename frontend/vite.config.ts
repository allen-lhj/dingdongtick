import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), vueDevTools()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    port: 5273,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8080/api/v1',
        changeOrigin: true,
        ws: true,
        rewrite: (path) => path.replace(new RegExp(`^/api`), ''),
      },
    },
  },
})
