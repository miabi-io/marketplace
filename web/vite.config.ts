import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) },
  },
  build: {
    outDir: 'dist',
  },
  server: {
    port: 3100,
    proxy: {
      '/v1': { target: 'http://localhost:8088', changeOrigin: true },
    },
  },
})
