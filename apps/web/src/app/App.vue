<template>
  <a-config-provider :theme="themeConfig" :locale="zhCN">
    <router-view />
  </a-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { theme } from 'ant-design-vue'
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import { syncColorModeFromSettings } from '@/services/theme'

syncColorModeFromSettings()

const themeConfig = computed(() => ({
  algorithm: theme.defaultAlgorithm,
  token: {
    colorPrimary: '#2563eb',
    borderRadius: 6,
    fontFamily:
      '-apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif',
  },
}))
</script>
