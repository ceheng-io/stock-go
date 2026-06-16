import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import EndOfDayPicker from '@/pages/EndOfDayPicker/EndOfDayPicker.vue'
import { analyzeEndOfDayStocks } from '@/services/analysis'
import {
  DEFAULT_END_OF_DAY_FILTERS,
  addEndOfDayRecentUsage,
  deleteEndOfDayScheme,
  getBatchWatchlistCandidates,
  getEndOfDayFilters,
  getEndOfDayRecentUsage,
  getEndOfDaySchemes,
  saveEndOfDayFilters,
  saveEndOfDayScheme,
  sortEndOfDayStocks,
  toggleSelectedCode,
} from '@/services/endOfDayPicker'
import type { EndOfDayStock } from '@/services/analysis'
import { addToWatchlist } from '@/services/storage'

const push = vi.fn()
const messageMocks = vi.hoisted(() => ({
  info: vi.fn(),
  success: vi.fn(),
  warning: vi.fn(),
  error: vi.fn(),
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({ push }),
}))

vi.mock('ant-design-vue', async (importOriginal) => {
  const actual = await importOriginal<typeof import('ant-design-vue')>()
  return {
    ...actual,
    message: {
      ...actual.message,
      info: messageMocks.info,
      success: messageMocks.success,
      warning: messageMocks.warning,
      error: messageMocks.error,
    },
  }
})

vi.mock('@/services/analysis', () => ({
  analyzeEndOfDayStocks: vi.fn(),
  isAnalysisAborted: vi.fn(() => false),
}))

vi.mock('@/services/storage', () => ({
  addToWatchlist: vi.fn(),
  isInWatchlist: vi.fn(() => false),
}))

const mockAnalyzeEndOfDayStocks = vi.mocked(analyzeEndOfDayStocks)
const mockAddToWatchlist = vi.mocked(addToWatchlist)

function stock(partial: Partial<EndOfDayStock>): EndOfDayStock {
  return {
    code: partial.code || '600000',
    routeCode: partial.routeCode || 'sh600000',
    name: partial.name || '测试股票',
    price: partial.price ?? 10,
    changePercent: partial.changePercent ?? 0,
    change: partial.change ?? 0,
    volume: partial.volume ?? 0,
    amount: partial.amount ?? 0,
    turnoverRate: partial.turnoverRate ?? null,
    volumeRatio: partial.volumeRatio ?? null,
    circulatingMarketCap: partial.circulatingMarketCap ?? null,
    totalMarketCap: partial.totalMarketCap ?? null,
    pe: partial.pe ?? null,
    pb: partial.pb ?? null,
    high: partial.high ?? 10.5,
    low: partial.low ?? 9.5,
    open: partial.open ?? 9.8,
    prevClose: partial.prevClose ?? 9.7,
    timeline: partial.timeline,
    timelineAboveAvgRatio: partial.timelineAboveAvgRatio,
  }
}

describe('end of day picker helpers', () => {
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
    push.mockClear()
    mockAnalyzeEndOfDayStocks.mockReset()
    mockAddToWatchlist.mockClear()
    messageMocks.info.mockClear()
    messageMocks.success.mockClear()
    messageMocks.warning.mockClear()
    messageMocks.error.mockClear()
  })

  it('persists filters while falling back to defaults for missing fields', () => {
    saveEndOfDayFilters({ ...DEFAULT_END_OF_DAY_FILTERS, marketCapMin: 80 })

    expect(getEndOfDayFilters()).toMatchObject({
      marketCapMin: 80,
      marketCapMax: DEFAULT_END_OF_DAY_FILTERS.marketCapMax,
      excludeST: true,
    })
  })

  it('saves and deletes named schemes with immutable filter snapshots', () => {
    const scheme = saveEndOfDayScheme('强势回踩', { ...DEFAULT_END_OF_DAY_FILTERS, changePercentMin: 2 })

    expect(getEndOfDaySchemes()).toEqual([scheme])
    expect(getEndOfDaySchemes()[0].filters.changePercentMin).toBe(2)

    deleteEndOfDayScheme(scheme.id)
    expect(getEndOfDaySchemes()).toEqual([])
  })

  it('deduplicates recent usage by filter content without moving existing entries and caps records at five', () => {
    for (let index = 0; index < 7; index += 1) {
      addEndOfDayRecentUsage({ ...DEFAULT_END_OF_DAY_FILTERS, marketCapMin: 50 + index }, 1000 + index)
    }
    addEndOfDayRecentUsage({ ...DEFAULT_END_OF_DAY_FILTERS, marketCapMin: 55 }, 2000)

    const recent = getEndOfDayRecentUsage()
    expect(recent).toHaveLength(5)
    expect(recent[0]).toMatchObject({ usedAt: 1006, filters: { marketCapMin: 56 } })
    expect(recent.filter((item) => item.filters.marketCapMin === 55)).toHaveLength(1)
    expect(recent.find((item) => item.filters.marketCapMin === 55)?.usedAt).toBe(1005)
  })

  it('sorts result stocks by numeric fields with nulls treated as zero', () => {
    const rows = [
      stock({ code: '1', turnoverRate: 8, timelineAboveAvgRatio: 70 }),
      stock({ code: '2', turnoverRate: null, timelineAboveAvgRatio: 95 }),
      stock({ code: '3', turnoverRate: 3, timelineAboveAvgRatio: 80 }),
    ]

    expect(sortEndOfDayStocks(rows, 'timelineAboveAvgRatio', 'desc').map((item) => item.code)).toEqual(['2', '3', '1'])
    expect(sortEndOfDayStocks(rows, 'turnoverRate', 'asc').map((item) => item.code)).toEqual(['2', '3', '1'])
  })

  it('toggles selected codes and filters batch watchlist candidates', () => {
    expect(toggleSelectedCode(['sh600000'], 'sh600001')).toEqual(['sh600000', 'sh600001'])
    expect(toggleSelectedCode(['sh600000'], 'sh600000')).toEqual([])

    expect(
      getBatchWatchlistCandidates(
        [stock({ routeCode: 'sh600000' }), stock({ routeCode: 'sz000001' })],
        ['sh600000', 'sz000001'],
        (code) => code === 'sh600000',
      ),
    ).toEqual(['sz000001'])
  })

  it('renders TypeScript-style start screen before analysis', () => {
    const wrapper = mount(EndOfDayPicker)

    expect(wrapper.find('[data-testid="eod-start"]').exists()).toBe(true)
    expect(wrapper.find('.start-screen').exists()).toBe(true)
    expect(wrapper.text()).toContain('点击编辑')
    expect(wrapper.text()).not.toContain('共筛选出')
  })

  it('shows card results, navigates from cards, and can return to the start screen', async () => {
    mockAnalyzeEndOfDayStocks.mockResolvedValue([
      stock({
        code: '600519',
        routeCode: 'sh600519',
        name: '贵州茅台',
        changePercent: 4.2,
        change: 12.3,
        price: 305.5,
        prevClose: 293.2,
        circulatingMarketCap: 120,
        volumeRatio: 1.8,
        turnoverRate: 6.5,
        amount: 360000000,
        timelineAboveAvgRatio: 100,
        timeline: [
          { time: '14:57', price: 305, avgPrice: 301 },
          { time: '14:58', price: 306, avgPrice: 302 },
        ],
      }),
    ])

    const wrapper = mount(EndOfDayPicker)
    await wrapper.get('[data-testid="eod-start"]').trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('共筛选出 1 只符合条件的股票')
    expect(wrapper.find('.stock-card').exists()).toBe(true)
    expect(wrapper.text()).toContain('强度 100%')

    await wrapper.get('.stock-card').trigger('click')
    expect(push).toHaveBeenCalledWith('/s/sh600519')

    await wrapper.get('[data-testid="eod-reset"]').trigger('click')
    expect(wrapper.find('.start-screen').exists()).toBe(true)
    expect(wrapper.text()).not.toContain('共筛选出 1 只符合条件的股票')
  })

  it('matches TypeScript batch select mode and adds selected cards to the watchlist', async () => {
    mockAnalyzeEndOfDayStocks.mockResolvedValue([
      stock({ code: '600519', routeCode: 'sh600519', name: '贵州茅台', timelineAboveAvgRatio: 100 }),
      stock({ code: '000001', routeCode: 'sz000001', name: '平安银行', timelineAboveAvgRatio: 90 }),
    ])

    const wrapper = mount(EndOfDayPicker)
    await wrapper.get('[data-testid="eod-start"]').trigger('click')
    await flushPromises()

    await wrapper.get('[data-testid="eod-select-mode"]').trigger('click')
    const selectButtons = wrapper.findAll('.select-btn')
    expect(selectButtons).toHaveLength(2)

    await selectButtons[0].trigger('click')
    await selectButtons[1].trigger('click')
    await wrapper.get('[data-testid="eod-batch-add"]').trigger('click')

    expect(mockAddToWatchlist).toHaveBeenCalledWith('sh600519')
    expect(mockAddToWatchlist).toHaveBeenCalledWith('sz000001')
  })
})
