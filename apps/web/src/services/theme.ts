import type { AppSettings } from '@/types'
import { getSettings } from '@/services/storage'

export type ThemeMode = 'light' | 'dark'

const THEME_KEY = 'app.theme'
export const THEME_MODE_CHANGE_EVENT = 'app-theme-mode-change'

export function applyColorMode(colorMode: AppSettings['colorMode']): void {
  document.documentElement.dataset.colorMode = colorMode
}

export function syncColorModeFromSettings(): void {
  applyColorMode(getSettings().colorMode)
}

export function getThemeMode(): ThemeMode {
  const stored = localStorage.getItem(THEME_KEY)
  return stored === 'light' ? 'light' : 'dark'
}

export function applyThemeMode(themeMode: ThemeMode): void {
  document.documentElement.dataset.theme = themeMode
  localStorage.setItem(THEME_KEY, themeMode)
  window.dispatchEvent(new CustomEvent<ThemeMode>(THEME_MODE_CHANGE_EVENT, { detail: themeMode }))
}

export function toggleThemeMode(): ThemeMode {
  const next = getThemeMode() === 'dark' ? 'light' : 'dark'
  applyThemeMode(next)
  return next
}

export function syncThemeModeFromStorage(): void {
  document.documentElement.dataset.theme = getThemeMode()
}
