import { mount } from '@vue/test-utils'
import { defineComponent, nextTick } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import Scanner from '@/pages/Scanner/Scanner.vue'

const push = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({ push }),
}))

vi.mock('@/services/api', () => ({
  getAllAShareQuotes: vi.fn(async () => []),
  getConceptConstituents: vi.fn(async () => []),
  getConceptList: vi.fn(async () => []),
  getIndustryConstituents: vi.fn(async () => []),
  getIndustryList: vi.fn(async () => []),
  getStockChanges: vi.fn(async () => []),
  getZTPool: vi.fn(async () => []),
}))

const tableStub = defineComponent({
  name: 'ATable',
  props: {
    customRow: { type: Function, default: undefined },
    dataSource: { type: Array, default: () => [] },
  },
  template: '<div data-testid="scanner-table" />',
})

describe('Scanner page', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('passes scan result row navigation through Ant Design table customRow', async () => {
    const wrapper = mount(Scanner, {
      global: {
        stubs: {
          AAlert: true,
          AButton: true,
          ACard: { template: '<section><slot /><slot name="title" /></section>' },
          ACheckableTag: true,
          AEmpty: true,
          AForm: { template: '<form><slot /></form>' },
          AFormItem: { template: '<div><slot /></div>' },
          AProgress: true,
          ASegmented: true,
          ASelect: true,
          ASpace: { template: '<div><slot /></div>' },
          ATable: tableStub,
          ATag: true,
        },
      },
    })

    const vm = wrapper.vm as unknown as {
      results: Array<{ code: string; routeCode: string; name: string; matchedSignals: string[]; time: string; added: boolean }>
    }
    vm.results = [{ code: '600519', routeCode: 'sh600519', name: '贵州茅台', matchedSignals: ['MA金叉'], time: '2026-06-15', added: false }]
    await nextTick()

    const customRow = wrapper.getComponent(tableStub).props('customRow')
    expect(customRow).toEqual(expect.any(Function))

    customRow?.({ routeCode: 'sh600519' }).onClick()
    expect(push).toHaveBeenCalledWith('/s/sh600519')
  })
})
