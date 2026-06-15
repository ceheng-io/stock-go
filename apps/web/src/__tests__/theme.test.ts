import { beforeEach, describe, expect, it, vi } from 'vitest'
import { applyColorMode, applyThemeMode, getThemeMode, syncColorModeFromSettings, toggleThemeMode } from '@/services/theme'

describe('theme service', () => {
  const store = new Map<string, string>()

  beforeEach(() => {
    store.clear()
    vi.stubGlobal('localStorage', {
      getItem: (key: string) => store.get(key) ?? null,
      setItem: (key: string, value: string) => store.set(key, value),
      removeItem: (key: string) => store.delete(key),
      clear: () => store.clear(),
    })
    document.documentElement.removeAttribute('data-color-mode')
    document.documentElement.removeAttribute('data-theme')
    localStorage.clear()
  })

  it('applies color mode to document root', () => {
    applyColorMode('green-rise')

    expect(document.documentElement.dataset.colorMode).toBe('green-rise')
  })

  it('syncs color mode from persisted app settings', () => {
    localStorage.setItem('app.settings', JSON.stringify({ colorMode: 'green-rise' }))

    syncColorModeFromSettings()

    expect(document.documentElement.dataset.colorMode).toBe('green-rise')
  })

  it('applies and persists light/dark theme mode', () => {
    applyThemeMode('dark')

    expect(document.documentElement.dataset.theme).toBe('dark')
    expect(localStorage.getItem('app.theme')).toBe('dark')
    expect(getThemeMode()).toBe('dark')
  })

  it('defaults to dark theme when no theme mode has been persisted', () => {
    expect(getThemeMode()).toBe('dark')
  })

  it('toggles theme mode from the persisted value', () => {
    localStorage.setItem('app.theme', 'dark')
    document.documentElement.dataset.theme = 'dark'

    expect(toggleThemeMode()).toBe('light')
    expect(document.documentElement.dataset.theme).toBe('light')
    expect(localStorage.getItem('app.theme')).toBe('light')
  })
})
