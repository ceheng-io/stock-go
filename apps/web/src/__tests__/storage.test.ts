import { beforeEach, describe, expect, it, vi } from 'vitest'
import {
  addSearchHistory,
  addToWatchlist,
  batchAddToWatchlist,
  batchRemoveFromWatchlist,
  createWatchlistGroup,
  deleteWatchlistGroup,
  getAllWatchlistCodes,
  getRefreshInterval,
  getSearchHistory,
  getSettings,
  getHeatmapConfig,
  getWatchlistGroups,
  getIndicatorConfig,
  renameWatchlistGroup,
  reorderWatchlist,
  updateSettings,
} from '@/services/storage'

describe('storage service', () => {
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

  it('creates a default watchlist group', () => {
    expect(getWatchlistGroups()).toHaveLength(1)
    expect(getWatchlistGroups()[0].name).toBe('默认分组')
  })

  it('normalizes and deduplicates watchlist codes', () => {
    addToWatchlist('600519')
    addToWatchlist('sh600519')

    expect(getAllWatchlistCodes()).toEqual(['sh600519'])
  })

  it('keeps recent search history capped at 20 rows', () => {
    for (let index = 0; index < 25; index += 1) {
      addSearchHistory({ code: `6005${index.toString().padStart(2, '0')}`, name: `股票${index}`, market: 'sh', type: '股票' })
    }

    expect(getSearchHistory()).toHaveLength(20)
    expect(getSearchHistory()[0].name).toBe('股票24')
  })

  it('renames and deletes custom watchlist groups while keeping default group', () => {
    const group = createWatchlistGroup('短线')

    renameWatchlistGroup(group.id, '观察池')
    expect(getWatchlistGroups().find((item) => item.id === group.id)?.name).toBe('观察池')

    deleteWatchlistGroup(group.id)
    expect(getWatchlistGroups().some((item) => item.id === group.id)).toBe(false)

    deleteWatchlistGroup('default')
    expect(getWatchlistGroups().some((item) => item.id === 'default')).toBe(true)
  })

  it('batch adds removes and reorders watchlist codes in one group', () => {
    const group = createWatchlistGroup('策略')

    const addedCount = batchAddToWatchlist(['600519', '000001', 'sh600519'], group.id)

    expect(addedCount).toBe(2)
    expect(batchAddToWatchlist(['600519', 'sh600519'], group.id)).toBe(0)
    expect(getWatchlistGroups().find((item) => item.id === group.id)?.codes).toEqual(['sh600519', 'sz000001'])

    reorderWatchlist(group.id, ['sz000001', 'sh600519', 'sz300750'])
    expect(getWatchlistGroups().find((item) => item.id === group.id)?.codes).toEqual(['sz000001', 'sh600519'])

    batchRemoveFromWatchlist(['000001', '600000'], group.id)
    expect(getWatchlistGroups().find((item) => item.id === group.id)?.codes).toEqual(['sh600519'])
  })

  it('updates nested heatmap settings without resetting unrelated settings', () => {
    updateSettings({
      refreshInterval: { list: 5000 },
      heatmapConfig: { dimension: 'watchlist', topK: 50 },
    })

    const settings = getSettings()
    expect(settings.refreshInterval).toEqual({ list: 5000, detail: 5000, heatmap: 10000 })
    expect(settings.heatmapConfig).toMatchObject({
      dimension: 'watchlist',
      colorField: 'changePercent',
      sizeField: 'totalMarketCap',
      colorMode: 'red-rise',
      topK: 50,
    })
    expect(settings.indicatorConfig.ma).toEqual([5, 10, 20, 60])
    expect(settings.indicatorConfig.sar).toEqual({ afStart: 0.02, afIncrement: 0.02, afMax: 0.2 })
    expect(settings.indicatorConfig.kc).toEqual({ emaPeriod: 20, atrPeriod: 10, multiplier: 2 })
  })

  it('treats zero refresh interval as the default list interval', () => {
    updateSettings({ refreshInterval: { list: 0 } })

    expect(getSettings().refreshInterval.list).toBe(0)
    expect(getRefreshInterval('list')).toBe(15000)
    expect(getRefreshInterval('detail')).toBe(5000)
    expect(getRefreshInterval('heatmap')).toBe(10000)
  })

  it('loads existing heatmap and indicator config keys', () => {
    localStorage.setItem('ui.heatmapConfig', JSON.stringify({
      dimension: 'watchlist',
      colorMode: 'green-rise',
      topK: 80,
    }))
    localStorage.setItem('ui.indicatorConfig', JSON.stringify({
      macd: { short: 8 },
      sar: { afMax: 0.3 },
    }))

    const settings = getSettings()

    expect(settings.heatmapConfig).toMatchObject({
      dimension: 'watchlist',
      colorField: 'changePercent',
      sizeField: 'totalMarketCap',
      colorMode: 'green-rise',
      topK: 80,
    })
    expect(settings.indicatorConfig.macd).toEqual({ short: 8, long: 26, signal: 9 })
    expect(settings.indicatorConfig.sar).toEqual({ afStart: 0.02, afIncrement: 0.02, afMax: 0.3 })
    expect(getHeatmapConfig().dimension).toBe('watchlist')
    expect(getIndicatorConfig().macd.short).toBe(8)
  })

  it('updates SAR and KC indicator settings without resetting other indicator settings', () => {
    updateSettings({
      indicatorConfig: {
        sar: { afStart: 0.03, afIncrement: 0.04, afMax: 0.3 },
        kc: { emaPeriod: 22, atrPeriod: 11, multiplier: 1.8 },
      },
    })

    const settings = getSettings()
    expect(settings.indicatorConfig.sar).toEqual({ afStart: 0.03, afIncrement: 0.04, afMax: 0.3 })
    expect(settings.indicatorConfig.kc).toEqual({ emaPeriod: 22, atrPeriod: 11, multiplier: 1.8 })
    expect(settings.indicatorConfig.ma).toEqual([5, 10, 20, 60])
  })
})
