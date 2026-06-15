import { describe, expect, it } from 'vitest'
import dashboardSource from '@/pages/Dashboard/Dashboard.vue?raw'
import {
  buildDashboardFailureMessage,
  getLatestMarketFundFlow,
  loadDashboardSections,
  pickNorthboundSnapshot,
  rankDashboardQuotes,
  summarizeMarketBreadth,
} from '@/services/dashboard'
import type { FullQuote } from '@/types'

function quote(partial: Partial<FullQuote>): FullQuote {
  return {
    code: partial.code || 'sh600000',
    name: partial.name || '测试股票',
    price: partial.price ?? 10,
    prevClose: partial.prevClose ?? 10,
    open: partial.open ?? 10,
    high: partial.high ?? 10,
    low: partial.low ?? 10,
    change: partial.change ?? 0,
    changePercent: partial.changePercent ?? 0,
    volume: partial.volume ?? 0,
    amount: partial.amount ?? 0,
    turnoverRate: partial.turnoverRate,
  }
}

describe('dashboard helpers', () => {
  it('summarizes market breadth with rise/fall/flat and limit counts', () => {
    const summary = summarizeMarketBreadth([
      quote({ changePercent: 10.01, amount: 20_000 }),
      quote({ changePercent: 1.2, amount: 30_000 }),
      quote({ changePercent: 0, amount: 10_000 }),
      quote({ changePercent: -9.91, amount: 40_000 }),
    ])

    expect(summary).toEqual({
      riseCount: 2,
      fallCount: 1,
      flatCount: 1,
      limitUpCount: 1,
      limitDownCount: 1,
      totalAmount: 100_000,
    })
  })

  it('ranks dashboard quotes by selected list metric', () => {
    const rows = [
      quote({ code: 'sh600001', changePercent: 3, amount: 20, turnoverRate: 2 }),
      quote({ code: 'sh600002', changePercent: -5, amount: 50, turnoverRate: 1 }),
      quote({ code: 'sh600003', changePercent: 1, amount: 10, turnoverRate: 9 }),
    ]

    expect(rankDashboardQuotes(rows, 'rise').map((item) => item.code)).toEqual(['sh600001', 'sh600003', 'sh600002'])
    expect(rankDashboardQuotes(rows, 'fall').map((item) => item.code)).toEqual(['sh600002', 'sh600003', 'sh600001'])
    expect(rankDashboardQuotes(rows, 'amount').map((item) => item.code)).toEqual(['sh600002', 'sh600001', 'sh600003'])
    expect(rankDashboardQuotes(rows, 'turnover').map((item) => item.code)).toEqual(['sh600003', 'sh600001', 'sh600002'])
  })

  it('prefers northbound aggregate rows when selecting a capital snapshot', () => {
    const snapshot = pickNorthboundSnapshot([
      { direction: '沪股通', boardName: '沪股通', netInflow: 100 },
      { direction: '北向资金', boardName: '北向汇总', netInflow: 300 },
      { direction: '沪深港通', boardName: '港股通', netInflow: 200 },
    ])

    expect(snapshot).toMatchObject({ direction: '北向资金', netInflow: 300 })
  })

  it('returns the latest market fund-flow row from history-like payloads', () => {
    expect(
      getLatestMarketFundFlow([
        { date: '2026-06-12', mainNetInflow: 10 },
        { date: '2026-06-15', mainNetInflow: -5 },
      ]),
    ).toMatchObject({ date: '2026-06-15', mainNetInflow: -5 })
    expect(getLatestMarketFundFlow([])).toBeNull()
  })

  it('commits dashboard sections as each one resolves instead of waiting for slow sections', async () => {
    const events: string[] = []
    let releaseSlow: ((value: string) => void) | undefined
    const slow = new Promise<string>((resolve) => {
      releaseSlow = resolve
    })

    const running = loadDashboardSections([
      { name: 'fast', load: async () => 'fast-value', commit: (value) => events.push(`fast:${value}`) },
      { name: 'slow', load: async () => slow, commit: (value) => events.push(`slow:${value}`) },
    ])

    await Promise.resolve()

    expect(events).toEqual(['fast:fast-value'])
    releaseSlow?.('slow-value')
    await running
    expect(events).toEqual(['fast:fast-value', 'slow:slow-value'])
  })

  it('reports a failed dashboard section while keeping other section commits', async () => {
    const events: string[] = []

    await loadDashboardSections(
      [
        { name: 'ok', load: async () => 1, commit: (value) => events.push(`ok:${value}`) },
        {
          name: 'bad',
          load: async () => {
            throw new Error('network')
          },
          commit: () => events.push('bad:commit'),
        },
      ],
      (name, error) => events.push(`${name}:${error instanceof Error ? error.message : 'unknown'}`),
    )

    expect(events).toEqual(['ok:1', 'bad:network'])
  })

  it('builds friendly dashboard failure messages from failed section names', () => {
    expect(buildDashboardFailureMessage([])).toBe('')
    expect(buildDashboardFailureMessage(['板块行情', '全市场行情'])).toBe(
      '板块行情、全市场行情加载失败，其他区域已保留可用数据',
    )
  })

  it('uses themed hover background for dashboard list rows', () => {
    expect(dashboardSource).toContain('background: var(--color-hover);')
    expect(dashboardSource).not.toContain('background: #f8fafc;')
  })
})
