import { describe, expect, it } from 'vitest'
import { calculateTimelineStrength, scanSignalPool } from '@/services/analysis'

describe('analysis service', () => {
  it('calculates percent of timeline points above average price', () => {
    const result = calculateTimelineStrength({
      code: 'sh600519',
      date: '2026-06-15',
      data: [
        { time: '09:30', price: 10, avgPrice: 10, volume: 100, amount: 1000 },
        { time: '09:31', price: 11, avgPrice: 10.5, volume: 100, amount: 1100 },
        { time: '09:32', price: 9, avgPrice: 10, volume: 100, amount: 900 },
      ],
    })

    expect(result.ratio).toBeCloseTo(66.666, 2)
    expect(result.points).toHaveLength(3)
  })

  it('matches only requested technical scanner signals from indicator kline data', async () => {
    const rows = [
      {
        date: '2026-06-11',
        close: 10,
        ma: { 5: 9.8, 10: 10.2 },
        macd: { dif: -0.2, dea: -0.1 },
        rsi: { 6: 45 },
        boll: { upper: 11, lower: 9 },
      },
      {
        date: '2026-06-12',
        close: 8.8,
        ma: { 5: 10.5, 10: 10.1 },
        macd: { dif: 0.1, dea: -0.05 },
        rsi: { 6: 24 },
        boll: { upper: 11, lower: 9 },
      },
    ]

    const results = await scanSignalPool(
      [{ code: '600519', routeCode: 'sh600519', name: '贵州茅台' }],
      ['ma_golden', 'rsi_oversold', 'boll_upper'],
      {
        loadKline: async () => rows,
      },
    )

    expect(results).toEqual([
      {
        code: '600519',
        routeCode: 'sh600519',
        name: '贵州茅台',
        matchedSignals: ['MA金叉', 'RSI超卖'],
      },
    ])
  })

  it('omits stocks when selected scanner signals are not matched', async () => {
    const results = await scanSignalPool(
      [{ code: '000001', routeCode: 'sz000001', name: '平安银行' }],
      ['ma_death', 'macd_death'],
      {
        loadKline: async () => [
          { date: '2026-06-11', close: 10, ma: { 5: 10, 10: 9 }, macd: { dif: 0.2, dea: 0.1 } },
          { date: '2026-06-12', close: 10.2, ma: { 5: 10.3, 10: 9.5 }, macd: { dif: 0.3, dea: 0.15 } },
        ],
      },
    )

    expect(results).toEqual([])
  })
})
