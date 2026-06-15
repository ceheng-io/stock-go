import { mount } from '@vue/test-utils'
import { defineComponent, nextTick } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import EndOfDayPicker from '@/pages/EndOfDayPicker/EndOfDayPicker.vue'
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

const push = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({ push }),
}))

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

  it('deduplicates recent usage by filter content and caps records at five', () => {
    for (let index = 0; index < 7; index += 1) {
      addEndOfDayRecentUsage({ ...DEFAULT_END_OF_DAY_FILTERS, marketCapMin: 50 + index }, 1000 + index)
    }
    addEndOfDayRecentUsage({ ...DEFAULT_END_OF_DAY_FILTERS, marketCapMin: 55 }, 2000)

    const recent = getEndOfDayRecentUsage()
    expect(recent).toHaveLength(5)
    expect(recent[0]).toMatchObject({ usedAt: 2000, filters: { marketCapMin: 55 } })
    expect(recent.filter((item) => item.filters.marketCapMin === 55)).toHaveLength(1)
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

  it('passes result row navigation through Ant Design table customRow', async () => {
    const tableStub = defineComponent({
      name: 'ATable',
      props: {
        customRow: { type: Function, default: undefined },
        dataSource: { type: Array, default: () => [] },
      },
      template: '<div data-testid="eod-table" />',
    })

    const wrapper = mount(EndOfDayPicker, {
      global: {
        stubs: {
          AButton: true,
          ACard: { template: '<section><slot /><slot name="extra" /></section>' },
          ACheckbox: true,
          ACol: { template: '<div><slot /></div>' },
          ADescriptions: true,
          ADescriptionsItem: true,
          ADivider: true,
          AEmpty: true,
          AForm: { template: '<form><slot /></form>' },
          AFormItem: { template: '<div><slot /></div>' },
          AInputNumber: true,
          AInputSearch: true,
          AList: true,
          AListItem: true,
          AListItemMeta: true,
          APopover: true,
          AProgress: true,
          ARow: { template: '<div><slot /></div>' },
          ASegmented: true,
          ASelect: true,
          ASpace: { template: '<div><slot /></div>' },
          AStatistic: true,
          ATable: tableStub,
        },
      },
    })

    const vm = wrapper.vm as unknown as { stocks: EndOfDayStock[] }
    vm.stocks = [stock({ routeCode: 'sh600519', name: '贵州茅台' })]
    await nextTick()

    const customRow = wrapper.getComponent(tableStub).props('customRow')
    expect(customRow).toEqual(expect.any(Function))

    customRow?.({ routeCode: 'sh600519' }).onClick()
    expect(push).toHaveBeenCalledWith('/s/sh600519')
  })
})
