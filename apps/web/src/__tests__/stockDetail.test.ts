import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import StockDetail from '@/pages/StockDetail/StockDetail.vue'
import type { FullQuote } from '@/types'

vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { code: 'sh600519' } }),
  useRouter: () => ({ back: vi.fn(), push: vi.fn() }),
}))

vi.mock('@/services/api', () => ({
  getDividendDetail: vi.fn(async () => []),
  getFullQuotes: vi.fn(async (): Promise<FullQuote[]> => [
    {
      code: 'sh600519',
      name: '贵州茅台',
      price: 1271.1,
      prevClose: 1291.9,
      open: 1288,
      high: 1290,
      low: 1268,
      change: -20.8,
      changePercent: -1.61,
      volume: 430000,
      amount: 5304000000,
      turnoverRate: 0.33,
    },
  ]),
  getIndividualFundFlow: vi.fn(async () => []),
  getKlineWithIndicators: vi.fn(async () => []),
  getMinuteKline: vi.fn(async () => []),
  getNorthboundIndividual: vi.fn(async () => []),
  getPanelLargeOrder: vi.fn(async () => []),
  getQuoteFundFlow: vi.fn(async () => []),
  getTodayTimeline: vi.fn(async () => ({ preClose: 1291.9, data: [] })),
}))

function mountStockDetail() {
  return mount(StockDetail, {
    global: {
      stubs: {
        AAlert: true,
        AButton: true,
        ACard: { template: '<section><slot name="title" /><slot /></section>' },
        ACheckboxGroup: true,
        ACol: { template: '<div><slot /></div>' },
        ADescriptions: { template: '<dl><slot /></dl>' },
        ADescriptionsItem: { template: '<div><slot /></div>' },
        AEmpty: true,
        AForm: { template: '<form><slot /></form>' },
        AFormItem: { template: '<div><slot /></div>' },
        AInputNumber: true,
        AList: true,
        AListItem: true,
        AListItemMeta: true,
        ARow: { template: '<div><slot /></div>' },
        ASegmented: true,
        ASelect: { template: '<select><slot /></select>' },
        ASelectOption: { template: '<option><slot /></option>' },
        ASpace: { template: '<div><slot /></div>' },
        ATabPane: { template: '<div><slot /></div>' },
        ATable: true,
        ATabs: { template: '<div><slot /></div>' },
        ATag: true,
        Echarts: true,
        KLineChart: true,
      },
    },
  })
}

describe('StockDetail page', () => {
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

  it('exposes near limit-up and limit-down alert types from the migrated detail page', () => {
    const wrapper = mountStockDetail()

    expect(wrapper.text()).toContain('接近涨停')
    expect(wrapper.text()).toContain('接近跌停')
  })
})
