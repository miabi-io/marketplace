import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// The dev server proxies the catalog API to a locally running `marketplace
// server` (make run), so the SPA talks to the same endpoints it will in prod.
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
