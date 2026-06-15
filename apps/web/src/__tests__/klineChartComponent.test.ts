import { flushPromises, mount } from '@vue/test-utils'
import { afterEach, describe, expect, it, vi } from 'vitest'
import KLineChart from '@/components/charts/KLineChart.vue'

const chartMock = vi.hoisted(() => ({
  applyNewData: vi.fn(),
  createIndicator: vi.fn(),
  removeIndicator: vi.fn(),
  setPriceVolumePrecision: vi.fn(),
  setStyles: vi.fn(),
  resize: vi.fn(),
}))
const initMock = vi.hoisted(() => vi.fn(() => chartMock))
const disposeMock = vi.hoisted(() => vi.fn())

vi.mock('klinecharts', () => ({
  init: initMock,
  dispose: disposeMock,
}))

const rows = [
  { date: '2026-06-12', open: 10, close: 11, low: 9.8, high: 11.2, volume: 1200 },
  { date: '2026-06-15', open: 11, close: 10.5, low: 10.3, high: 11.3, volume: 900 },
]

describe('KLineChart component', () => {
  afterEach(() => {
    vi.clearAllMocks()
    vi.unstubAllGlobals()
  })

  it('initializes klinecharts and applies converted kline data plus indicators', async () => {
    vi.stubGlobal('ResizeObserver', class {
      observe = vi.fn()
      disconnect = vi.fn()
    })

    mount(KLineChart, {
      props: {
        rows,
        overlays: ['ma', 'boll'],
        oscillator: 'macd',
        emptyText: '暂无 K 线',
      },
      attachTo: document.body,
    })
    await flushPromises()

    expect(initMock).toHaveBeenCalledTimes(1)
    expect(chartMock.applyNewData).toHaveBeenCalledWith([
      expect.objectContaining({ open: 10, high: 11.2, low: 9.8, close: 11, volume: 1200 }),
      expect.objectContaining({ open: 11, high: 11.3, low: 10.3, close: 10.5, volume: 900 }),
    ])
    expect(chartMock.createIndicator).toHaveBeenCalledWith('MA', true, { id: 'candle_pane' })
    expect(chartMock.createIndicator).toHaveBeenCalledWith('BOLL', true, { id: 'candle_pane' })
    expect(chartMock.createIndicator).toHaveBeenCalledWith('VOL', false, { id: 'volume_pane', height: 72 })
    expect(chartMock.createIndicator).toHaveBeenCalledWith('MACD', false, { id: 'indicator_pane', height: 92 })
  })

  it('renders empty text and disposes chart instance', async () => {
    vi.stubGlobal('ResizeObserver', class {
      observe = vi.fn()
      disconnect = vi.fn()
    })

    const wrapper = mount(KLineChart, {
      props: {
        rows: [],
        emptyText: '暂无板块 K 线',
      },
      attachTo: document.body,
    })
    await flushPromises()

    expect(wrapper.text()).toContain('暂无板块 K 线')

    wrapper.unmount()

    expect(disposeMock).toHaveBeenCalledTimes(1)
  })
})
