import { fileURLToPath, URL } from 'node:url'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    base: process.env.VITE_BASE_URL || env.VITE_BASE_URL || '/',
    plugins: [vue()],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    server: {
      proxy: {
        '/api': {
          target: process.env.VITE_DEV_API_TARGET || env.VITE_DEV_API_TARGET || 'http://127.0.0.1:8080',
          changeOrigin: true,
        },
      },
    },
    test: {
      environment: 'jsdom',
      globals: true,
    },
    build: {
      rollupOptions: {
        output: {
          manualChunks(id) {
            if (!id.includes('node_modules')) return
            if (id.includes('/echarts/') || id.includes('/zrender/') || id.includes('vue-echarts')) {
              return 'echarts-vendor'
            }
            if (id.includes('ant-design-vue') || id.includes('@ant-design')) {
              return 'antdv-vendor'
            }
            if (id.includes('/vue/') || id.includes('vue-router')) {
              return 'vue-vendor'
            }
          },
        },
      },
    },
  }
})
