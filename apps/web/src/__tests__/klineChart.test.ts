import { describe, expect, it } from 'vitest'
import {
  buildKLineIndicatorNames,
  toKLineChartData,
  type KLineChartSourceRow,
} from '@/services/klineChart'
import type { IndicatorConfig } from '@/types'

describe('klinecharts helpers', () => {
  it('converts API kline rows into klinecharts data rows', () => {
    const rows: KLineChartSourceRow[] = [
      { date: '2026-06-12', open: 10, close: 11, low: 9.8, high: 11.2, volume: 1200 },
      { time: '2026-06-15 09:35', open: 11, close: 10.5, low: 10.3, high: 11.3, volume: 900 },
      { date: 'bad-row', open: null, close: 12, low: 11, high: 13, volume: 1 },
    ]

    expect(toKLineChartData(rows)).toEqual([
      {
        timestamp: new Date('2026-06-12T00:00:00+08:00').getTime(),
        open: 10,
        close: 11,
        low: 9.8,
        high: 11.2,
        volume: 1200,
      },
      {
        timestamp: new Date('2026-06-15T09:35:00+08:00').getTime(),
        open: 11,
        close: 10.5,
        low: 10.3,
        high: 11.3,
        volume: 900,
      },
    ])
  })

  it('maps selected overlay and oscillator controls to klinecharts indicator names', () => {
    const indicatorConfig: IndicatorConfig = {
      ma: [5, 10, 20],
      macd: { short: 12, long: 26, signal: 9 },
      boll: { period: 20, stdDev: 2 },
      kdj: { period: 9, kPeriod: 3, dPeriod: 3 },
      rsi: [6, 12, 24],
      dmi: { period: 14, adxPeriod: 14 },
      sar: { afStart: 0.02, afIncrement: 0.02, afMax: 0.2 },
      kc: { emaPeriod: 20, atrPeriod: 10, multiplier: 2 },
    }

    expect(buildKLineIndicatorNames({
      overlays: ['ma', 'boll', 'sar', 'kc'],
      oscillator: 'macd',
      indicatorConfig,
    })).toEqual({
      overlays: ['MA', 'BOLL', 'SAR', 'KC'],
      panes: ['VOL', 'MACD'],
    })
  })
})
