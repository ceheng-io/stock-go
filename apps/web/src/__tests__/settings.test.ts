import { mount } from '@vue/test-utils'
import { defineComponent } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import Settings from '@/pages/Settings/Settings.vue'

vi.mock('ant-design-vue', async (importOriginal) => {
  const actual = await importOriginal<typeof import('ant-design-vue')>()
  return {
    ...actual,
    message: {
      ...actual.message,
      success: vi.fn(),
    },
  }
})

const selectStub = defineComponent({
  name: 'ASelect',
  props: {
    value: { type: [String, Number], default: undefined },
  },
  template: '<select><slot /></select>',
})

const selectOptionStub = defineComponent({
  name: 'ASelectOption',
  props: {
    value: { type: [String, Number], default: undefined },
  },
  template: '<option :value="value"><slot /></option>',
})

function mountSettings() {
  return mount(Settings, {
    global: {
      stubs: {
        AButton: true,
        ACard: { template: '<section><slot name="title" /><slot /></section>' },
        ACol: { template: '<div><slot /></div>' },
        ADescriptions: { template: '<dl><slot /></dl>' },
        ADescriptionsItem: { template: '<div><slot /></div>' },
        AForm: { template: '<form><slot /></form>' },
        AFormItem: { template: '<label><slot /></label>' },
        AInput: true,
        AInputNumber: true,
        ARadioButton: { template: '<button><slot /></button>' },
        ARadioGroup: { template: '<div><slot /></div>' },
        ARow: { template: '<div><slot /></div>' },
        ASelect: selectStub,
        ASelectOption: selectOptionStub,
        ASpace: { template: '<div><slot /></div>' },
      },
    },
  })
}

describe('Settings page', () => {
  const store = new Map<string, string>()

  beforeEach(() => {
    store.clear()
    vi.stubGlobal('localStorage', {
      getItem: (key: string) => store.get(key) ?? null,
      setItem: (key: string, value: string) => store.set(key, value),
      removeItem: (key: string) => store.delete(key),
      clear: () => store.clear(),
    })
    localStorage.clear()
  })

  it('keeps the default label for zero list refresh interval', () => {
    const wrapper = mountSettings()

    expect(wrapper.text()).toContain('默认')
    expect(wrapper.text()).not.toContain('手动刷新')
  })

  it('shows data source and unit notes', () => {
    const wrapper = mountSettings()
    const text = wrapper.text()

    expect(text).toContain('策衡 A 股看板')
    expect(text).toContain('apps/server')
    expect(text).toContain('Go SDK')
    expect(text).toContain('成交量单位：手')
    expect(text).toContain('成交额单位：万元')
    expect(text).toContain('资金流、北向、龙虎榜等数据默认使用元级展示')
    expect(text).toContain('市值单位：亿元')
  })
})
