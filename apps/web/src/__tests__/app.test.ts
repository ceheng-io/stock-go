import { mount } from '@vue/test-utils'
import { defineComponent, nextTick } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { theme } from 'ant-design-vue'
import App from '@/app/App.vue'
import { applyThemeMode } from '@/services/theme'

const configProviderStub = defineComponent({
  name: 'AConfigProvider',
  props: {
    theme: { type: Object, default: undefined },
    locale: { type: Object, default: undefined },
  },
  template: '<div><slot /></div>',
})

function mountApp() {
  return mount(App, {
    global: {
      stubs: {
        AConfigProvider: configProviderStub,
        RouterView: true,
      },
    },
  })
}

describe('App theme provider', () => {
  const store = new Map<string, string>()

  beforeEach(() => {
    store.clear()
    vi.stubGlobal('localStorage', {
      getItem: (key: string) => store.get(key) ?? null,
      setItem: (key: string, value: string) => store.set(key, value),
      removeItem: (key: string) => store.delete(key),
      clear: () => store.clear(),
    })
    document.documentElement.removeAttribute('data-theme')
    document.documentElement.removeAttribute('data-color-mode')
    localStorage.clear()
  })

  it('uses Ant Design dark algorithm from persisted theme mode', () => {
    localStorage.setItem('app.theme', 'dark')

    const wrapper = mountApp()

    const providerTheme = wrapper.getComponent(configProviderStub).props('theme') as { algorithm: unknown }
    expect(providerTheme.algorithm).toBe(theme.darkAlgorithm)
  })

  it('uses coordinated dark tokens for layout, content, and navigation surfaces', () => {
    localStorage.setItem('app.theme', 'dark')

    const wrapper = mountApp()

    const providerTheme = wrapper.getComponent(configProviderStub).props('theme') as {
      token: Record<string, string>
      components: Record<string, Record<string, string>>
    }
    expect(providerTheme.token).toMatchObject({
      colorBgLayout: '#11151c',
      colorBgContainer: '#171d26',
      colorBgElevated: '#1d2530',
      colorText: '#d7dde7',
      colorTextHeading: '#eef2f7',
      colorTextSecondary: '#aeb8c7',
      colorBorder: '#2a3442',
    })
    expect(providerTheme.components.Layout).toMatchObject({
      colorBgBody: '#11151c',
      colorBgHeader: '#151b24',
    })
    expect(providerTheme.components.Menu).toMatchObject({
      colorItemText: '#aeb8c7',
      colorItemTextSelected: '#dce9ff',
      colorItemBgSelected: '#1f2f46',
      colorItemBgHover: '#202a36',
    })
  })

  it('updates Ant Design algorithm when theme mode changes in the same session', async () => {
    const wrapper = mountApp()

    applyThemeMode('dark')
    await nextTick()

    const providerTheme = wrapper.getComponent(configProviderStub).props('theme') as { algorithm: unknown }
    expect(providerTheme.algorithm).toBe(theme.darkAlgorithm)
  })
})
