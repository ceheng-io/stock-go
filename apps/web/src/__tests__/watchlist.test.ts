import { mount } from '@vue/test-utils'
import { defineComponent, nextTick } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { moveWatchlistCodeBefore } from '@/services/watchlist'
import { saveWatchlistGroups } from '@/services/storage'
import Watchlist from '@/pages/Watchlist/Watchlist.vue'
import type { FullQuote } from '@/types'

const push = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({ push }),
}))

vi.mock('@/services/api', () => ({
  getAllQuotesByCodes: vi.fn(async (): Promise<FullQuote[]> => [
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
}))

const tableStub = defineComponent({
  name: 'ATable',
  props: {
    columns: { type: Array, default: () => [] },
    customRow: { type: Function, default: undefined },
    dataSource: { type: Array, default: () => [] },
  },
  template: '<div data-testid="watchlist-table" />',
})

const emptyStub = defineComponent({
  name: 'AEmpty',
  template: '<div><slot /></div>',
})

const cardStub = defineComponent({
  name: 'ACard',
  template: '<section><slot name="title" /><slot name="extra" /><slot /></section>',
})

describe('watchlist helpers', () => {
  const store = new Map<string, string>()

  beforeEach(() => {
    store.clear()
    vi.stubGlobal('localStorage', {
      getItem: (key: string) => store.get(key) ?? null,
      setItem: (key: string, value: string) => store.set(key, value),
      removeItem: (key: string) => store.delete(key),
      clear: () => store.clear(),
    })
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('moves a dragged code before a target code without losing other codes', () => {
    expect(moveWatchlistCodeBefore(['sh600519', 'sz000001', 'sh600000'], 'sh600000', 'sz000001')).toEqual([
      'sh600519',
      'sh600000',
      'sz000001',
    ])
  })

  it('keeps order unchanged for missing or identical drag targets', () => {
    const rows = ['sh600519', 'sz000001']

    expect(moveWatchlistCodeBefore(rows, 'sh600519', 'sh600519')).toEqual(rows)
    expect(moveWatchlistCodeBefore(rows, 'sh600000', 'sh600519')).toEqual(rows)
    expect(moveWatchlistCodeBefore(rows, 'sh600519', 'sh600000')).toEqual(rows)
  })

  it('passes drag row props to Ant Design Vue table customRow in default sort mode', async () => {
    saveWatchlistGroups([
      {
        id: 'default',
        name: '默认分组',
        codes: ['sh600519'],
        createdAt: 1,
        updatedAt: 1,
      },
    ])

    const wrapper = mount(Watchlist, {
      global: {
        stubs: {
          AAlert: true,
          ABadge: true,
          AButton: true,
          ACard: cardStub,
          ACheckbox: true,
          ACheckboxGroup: true,
          ADrawer: true,
          ADropdown: true,
          AEmpty: emptyStub,
          AInput: true,
          AInputSearch: true,
          AList: true,
          AListItem: true,
          AListItemMeta: true,
          AMenu: true,
          AMenuItem: true,
          ASegmented: true,
          ASelect: true,
          ASelectOption: true,
          ASpace: true,
          ASwitch: true,
          ATable: tableStub,
          ATag: true,
          ATextarea: true,
        },
      },
    })

    await nextTick()
    await nextTick()

    const table = wrapper.getComponent(tableStub)
    const customRow = table.props('customRow')
    expect(customRow).toEqual(expect.any(Function))
    const rowProps = customRow?.({
      code: 'sh600519',
      name: '贵州茅台',
    } as FullQuote)

    expect(rowProps?.draggable).toBe(true)
    expect(rowProps?.onDragstart).toEqual(expect.any(Function))
    expect(rowProps?.onDragover).toEqual(expect.any(Function))
    expect(rowProps?.onDrop).toEqual(expect.any(Function))
    expect(rowProps?.onDragend).toEqual(expect.any(Function))
  })
})
