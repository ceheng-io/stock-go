import type { AppSettings } from '@/types'
import { getSettings } from '@/services/storage'

export function applyColorMode(colorMode: AppSettings['colorMode']): void {
  document.documentElement.dataset.colorMode = colorMode
}

export function syncColorModeFromSettings(): void {
  applyColorMode(getSettings().colorMode)
}
