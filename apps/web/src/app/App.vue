<template>
  <a-config-provider :theme="themeConfig" :locale="zhCN">
    <router-view />
  </a-config-provider>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'
import { theme } from 'ant-design-vue'
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import { getThemeMode, syncColorModeFromSettings, syncThemeModeFromStorage, THEME_MODE_CHANGE_EVENT, type ThemeMode } from '@/services/theme'

syncColorModeFromSettings()
syncThemeModeFromStorage()

const themeMode = ref<ThemeMode>(getThemeMode())
function handleThemeModeChange(event: Event) {
  themeMode.value = (event as CustomEvent<ThemeMode>).detail
}

window.addEventListener(THEME_MODE_CHANGE_EVENT, handleThemeModeChange)
onBeforeUnmount(() => {
  window.removeEventListener(THEME_MODE_CHANGE_EVENT, handleThemeModeChange)
})

const fontFamily =
  '-apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif'

const lightPalette = {
  bgLayout: '#f4f6fa',
  bgContainer: '#ffffff',
  bgElevated: '#ffffff',
  bgHeader: '#ffffff',
  bgSidebar: '#ffffff',
  border: '#d9e0ea',
  split: '#e7edf5',
  text: '#172033',
  textHeading: '#0f172a',
  textSecondary: '#44546a',
  textTertiary: '#69758a',
  hover: '#eef3fb',
  active: '#e8f0ff',
  activeText: '#1d4ed8',
}

const darkPalette = {
  bgLayout: '#11151c',
  bgContainer: '#171d26',
  bgElevated: '#1d2530',
  bgHeader: '#151b24',
  bgSidebar: '#121820',
  border: '#2a3442',
  split: '#202a36',
  text: '#d7dde7',
  textHeading: '#eef2f7',
  textSecondary: '#aeb8c7',
  textTertiary: '#8d98a8',
  hover: '#202a36',
  active: '#1f2f46',
  activeText: '#dce9ff',
}

const themeConfig = computed(() => {
  const palette = themeMode.value === 'dark' ? darkPalette : lightPalette

  return {
    algorithm: themeMode.value === 'dark' ? theme.darkAlgorithm : theme.defaultAlgorithm,
    token: {
      colorPrimary: '#2563eb',
      colorBgLayout: palette.bgLayout,
      colorBgContainer: palette.bgContainer,
      colorBgElevated: palette.bgElevated,
      colorBorder: palette.border,
      colorBorderSecondary: palette.split,
      colorSplit: palette.split,
      colorText: palette.text,
      colorTextBase: palette.text,
      colorTextHeading: palette.textHeading,
      colorTextLabel: palette.textSecondary,
      colorTextDescription: palette.textTertiary,
      colorTextSecondary: palette.textSecondary,
      colorTextTertiary: palette.textTertiary,
      colorFillAlter: palette.bgElevated,
      colorFillContent: palette.hover,
      controlItemBgHover: palette.hover,
      controlItemBgActive: palette.active,
      borderRadius: 6,
      fontFamily,
    },
    components: {
      Layout: {
        colorBgBody: palette.bgLayout,
        colorBgHeader: palette.bgHeader,
      },
      Menu: {
        colorItemText: palette.textSecondary,
        colorItemTextHover: palette.textHeading,
        colorItemTextSelected: palette.activeText,
        colorItemBg: 'transparent',
        colorItemBgHover: palette.hover,
        colorItemBgSelected: palette.active,
        colorActiveBarBorderSize: 0,
      },
      Table: {
        colorBgContainer: palette.bgContainer,
        colorFillAlter: palette.bgElevated,
        colorBorderSecondary: palette.split,
        colorText: palette.text,
        colorTextHeading: palette.textHeading,
      },
      Card: {
        colorBgContainer: palette.bgContainer,
        colorBorderSecondary: palette.split,
        colorText: palette.text,
        colorTextHeading: palette.textHeading,
      },
    },
  }
})
</script>
