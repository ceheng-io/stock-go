import { describe, expect, it } from 'vitest'
import {
  buildEmptyChartOption,
  buildFundFlowOption,
  buildIndicatorKlineOption,
  buildKlineVolumeOption,
  buildMinuteChartOption,
  normalizeBoardSpotRows,
} from '@/services/charts'
import type { IndicatorConfig } from '@/types'

describe('chart helpers', () => {
  it('keeps empty chart axes as arrays so stale multi-axis series cannot reference missing xAxis 1', () => {
    const option = buildEmptyChartOption('暂无数据')

    expect(Array.isArray(option.xAxis)).toBe(true)
    expect(Array.isArray(option.yAxis)).toBe(true)
    expect(option.series).toEqual([])
  })

  it('builds K line and volume option with matching secondary axes', () => {
    const option = buildKlineVolumeOption(
      [
        { date: '2026-06-12', open: 10, close: 11, low: 9.8, high: 11.2, volume: 1200 },
        { date: '2026-06-15', open: 11, close: 10.5, low: 10.3, high: 11.3, volume: 900 },
      ],
      { emptyText: '暂无 K 线数据' },
    )

    expect(option.xAxis).toHaveLength(2)
    expect(option.yAxis).toHaveLength(2)
    expect(option.series).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ name: 'K线', type: 'candlestick' }),
        expect.objectContaining({ name: '成交量', type: 'bar', xAxisIndex: 1, yAxisIndex: 1 }),
      ]),
    )
  })

  it('normalizes board spot arrays, wrapped payloads and object payloads', () => {
    expect(normalizeBoardSpotRows([{ item: '最新', value: 123 }])).toEqual([{ item: '最新', value: 123 }])
    expect(normalizeBoardSpotRows({ data: [{ item: '涨幅', value: 1.2 }] })).toEqual([{ item: '涨幅', value: 1.2 }])
    expect(normalizeBoardSpotRows({ price: 12.3, changePercent: 2.1, name: '板块' })).toEqual([
      { item: 'price', value: 12.3 },
      { item: 'changePercent', value: 2.1 },
    ])
  })

  it('builds indicator K line option with overlays and oscillator panels', () => {
    const indicatorConfig: IndicatorConfig = {
      ma: [5, 10],
      macd: { short: 12, long: 26, signal: 9 },
      boll: { period: 20, stdDev: 2 },
      kdj: { period: 9, kPeriod: 3, dPeriod: 3 },
      rsi: [6, 12],
      dmi: { period: 14, adxPeriod: 14 },
      sar: { afStart: 0.02, afIncrement: 0.02, afMax: 0.2 },
      kc: { emaPeriod: 20, atrPeriod: 10, multiplier: 2 },
    }

    const option = buildIndicatorKlineOption(
      [
        {
          date: '2026-06-12',
          open: 10,
          close: 11,
          low: 9.8,
          high: 11.2,
          volume: 1200,
          ma: { ma5: 10.4, ma10: 10.1 },
          boll: { upper: 12, mid: 10, lower: 8 },
          kc: { upper: 11.8, mid: 10.2, lower: 8.6 },
          macd: { dif: 0.1, dea: 0.05, macd: 0.1 },
        },
        {
          date: '2026-06-15',
          open: 11,
          close: 10.5,
          low: 10.3,
          high: 11.3,
          volume: 900,
          ma: { ma5: 10.6, ma10: 10.2 },
          boll: { upper: 12.1, mid: 10.1, lower: 8.1 },
          kc: { upper: 11.9, mid: 10.3, lower: 8.7 },
          macd: { dif: 0.2, dea: 0.08, macd: 0.24 },
        },
      ],
      {
        emptyText: '暂无 K 线数据',
        overlays: ['ma', 'boll', 'kc'],
        oscillator: 'macd',
        indicatorConfig,
      },
    )

    expect(option.xAxis).toHaveLength(3)
    expect(option.yAxis).toHaveLength(3)
    expect(option.dataZoom?.[0]).toMatchObject({ xAxisIndex: [0, 1, 2] })
    expect(option.series).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ name: 'MA5', type: 'line' }),
        expect.objectContaining({ name: 'BOLL上轨', type: 'line' }),
        expect.objectContaining({ name: 'KC下轨', type: 'line' }),
        expect.objectContaining({ name: 'MACD', type: 'bar', xAxisIndex: 2, yAxisIndex: 2 }),
        expect.objectContaining({ name: 'DIF', type: 'line', xAxisIndex: 2, yAxisIndex: 2 }),
      ]),
    )
  })

  it('builds sector fund flow option with amount bar and ratio line axes', () => {
    const option = buildFundFlowOption(
      [
        { date: '2026-06-12', mainNetInflow: 120000000, mainNetInflowPercent: 3.2 },
        { date: '2026-06-15', mainNetInflow: -50000000, mainNetInflowPercent: -1.1 },
      ],
      { emptyText: '暂无资金流数据' },
    )

    expect(option.xAxis).toHaveLength(1)
    expect(option.yAxis).toHaveLength(2)
    expect(option.series).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ name: '主力净流入', type: 'bar' }),
        expect.objectContaining({ name: '净占比', type: 'line', yAxisIndex: 1 }),
      ]),
    )
  })

  it('builds one minute timeline option with price and average lines', () => {
    const option = buildMinuteChartOption({
      period: '1',
      timeline: {
        preClose: 10,
        data: [
          { time: '09:30', price: 10.1, avgPrice: 10.05 },
          { time: '09:31', price: 10.2, avgPrice: 10.08 },
        ],
      },
      minuteKline: [],
      emptyText: '暂无分时数据',
    })

    expect(option.xAxis).toHaveLength(1)
    expect(option.yAxis).toHaveLength(1)
    expect(option.series).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ name: '价格', type: 'line' }),
        expect.objectContaining({ name: '均价', type: 'line' }),
      ]),
    )
  })

  it('builds minute K line option with volume axes for non-one-minute periods', () => {
    const option = buildMinuteChartOption({
      period: '5',
      timeline: null,
      minuteKline: [
        { time: '09:35', open: 10, close: 10.2, low: 9.9, high: 10.3, volume: 1200 },
        { time: '09:40', open: 10.2, close: 10.1, low: 10, high: 10.4, volume: 900 },
      ],
      emptyText: '暂无分钟 K 数据',
    })

    expect(option.xAxis).toHaveLength(2)
    expect(option.yAxis).toHaveLength(2)
    expect(option.dataZoom?.[0]).toMatchObject({ xAxisIndex: [0, 1] })
    expect(option.series).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ name: '5分K', type: 'candlestick' }),
        expect.objectContaining({ name: '成交量', type: 'bar', xAxisIndex: 1, yAxisIndex: 1 }),
      ]),
    )
  })
})
