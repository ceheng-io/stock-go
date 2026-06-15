import { beforeEach, describe, expect, it, vi } from 'vitest'
import { applyColorMode, syncColorModeFromSettings } from '@/services/theme'

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
})
