import { mount } from '@vue/test-utils'
import { defineComponent, nextTick } from 'vue'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import AppHeader from '@/layouts/AppHeader.vue'

const push = vi.fn()
const searchMock = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({ push }),
}))

vi.mock('@/services/api', () => ({
  search: (...args: unknown[]) => searchMock(...args),
}))

const autoCompleteStub = defineComponent({
  name: 'AAutoComplete',
  props: {
    value: { type: String, default: '' },
    options: { type: Array, default: () => [] },
    allowClear: { type: Boolean, default: false },
  },
  emits: ['search', 'select', 'focus', 'update:value'],
  template: `
    <div>
      <input
        data-testid="header-search"
        :value="value"
        @input="$emit('update:value', $event.target.value); $emit('search', $event.target.value)"
        @focus="$emit('focus')"
      />
      <div data-testid="header-options">
        <slot
          v-for="option in options"
          name="option"
          :key="option.value"
          :option="option"
        />
      </div>
    </div>
  `,
})

function mountHeader() {
  return mount(AppHeader, {
    global: {
      stubs: {
        AAutoComplete: autoCompleteStub,
        AButton: true,
        ALayoutHeader: { template: '<header><slot /></header>' },
        ASpace: { template: '<div><slot /></div>' },
        ATag: true,
        BulbOutlined: true,
        CheckOutlined: true,
        DatabaseOutlined: true,
        EyeInvisibleOutlined: true,
        GithubOutlined: true,
        MenuOutlined: true,
        StarOutlined: true,
      },
    },
  })
}

function createStorage() {
  const store = new Map<string, string>()
  vi.stubGlobal('localStorage', {
    getItem: (key: string) => store.get(key) ?? null,
    setItem: (key: string, value: string) => store.set(key, value),
    removeItem: (key: string) => store.delete(key),
    clear: () => store.clear(),
  })
}

describe('AppHeader', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    vi.clearAllMocks()
    createStorage()
    localStorage.clear()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('stores selected search result into recent history', async () => {
    searchMock.mockResolvedValue([
      { code: '600519', name: '贵州茅台', market: 'sh', type: '股票', entityType: 'stock', isSupported: true },
    ])

    const wrapper = mountHeader()

    await wrapper.get('[data-testid="header-search"]').setValue('茅台')
    await vi.advanceTimersByTimeAsync(250)
    await nextTick()

    await wrapper.get('[data-testid="header-options"] .search-option').trigger('click')

    expect(push).toHaveBeenCalledWith('/s/sh600519')
    expect(localStorage.getItem('search.recent')).toContain('贵州茅台')
  })

  it('selects a highlighted search result from the autocomplete select event', async () => {
    searchMock.mockResolvedValue([
      { code: '600519', name: '贵州茅台', market: 'sh', type: '股票', entityType: 'stock', isSupported: true },
    ])

    const wrapper = mountHeader()

    await wrapper.get('[data-testid="header-search"]').setValue('茅台')
    await vi.advanceTimersByTimeAsync(250)
    await nextTick()

    const autoComplete = wrapper.getComponent(autoCompleteStub)
    const [option] = autoComplete.props('options') as Array<{
      value: string
      label: string
      route: string
      market: string
      type: string
    }>

    autoComplete.vm.$emit('select', option.value, option)
    await nextTick()

    expect(push).toHaveBeenCalledWith('/s/sh600519')
    expect(localStorage.getItem('search.recent')).toContain('贵州茅台')
  })

  it('enables clearing the search input from the autocomplete control', () => {
    const wrapper = mountHeader()

    expect(wrapper.getComponent(autoCompleteStub).props('allowClear')).toBe(true)
  })

  it('adds stock search result to watchlist without navigating', async () => {
    searchMock.mockResolvedValue([
      { code: '600519', name: '贵州茅台', market: 'sh', type: '股票', entityType: 'stock', isSupported: true },
    ])

    const wrapper = mountHeader()

    await wrapper.get('[data-testid="header-search"]').setValue('茅台')
    await vi.advanceTimersByTimeAsync(250)
    await nextTick()

    await wrapper.get('.quick-add').trigger('click')

    expect(push).not.toHaveBeenCalled()
    expect(localStorage.getItem('watchlist.groups')).toContain('sh600519')
  })

  it('shows and clears recent search history', async () => {
    localStorage.setItem('search.recent', JSON.stringify([
      { code: 'sh600519', name: '贵州茅台', market: 'sh', type: '股票', timestamp: 1 },
    ]))

    const wrapper = mountHeader()

    await wrapper.get('[data-testid="header-search"]').trigger('focus')

    expect(wrapper.text()).toContain('最近搜索')
    expect(wrapper.text()).toContain('贵州茅台')

    await wrapper.get('.clear-history').trigger('click')

    expect(localStorage.getItem('search.recent')).toBeNull()
  })

  it('renders SDK GitHub links and toggles theme mode', async () => {
    const wrapper = mountHeader()

    expect(wrapper.get('a[href="https://stock-sdk.linkdiary.cn/"]').attributes('target')).toBe('_blank')
    expect(wrapper.get('a[href="https://github.com/chengzuopeng/stock-dashboard"]').attributes('target')).toBe('_blank')

    await wrapper.get('.theme-toggle').trigger('click')

    expect(document.documentElement.dataset.theme).toBe('light')
    expect(localStorage.getItem('app.theme')).toBe('light')
  })
})
