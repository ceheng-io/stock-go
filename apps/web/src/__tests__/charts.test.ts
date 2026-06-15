import { describe, expect, it } from 'vitest'
import {
  buildEmptyChartOption,
  buildFundFlowOption,
  buildMinuteChartOption,
  normalizeBoardSpotRows,
} from '@/services/charts'

describe('chart helpers', () => {
  it('keeps empty chart axes as arrays so stale multi-axis series cannot reference missing xAxis 1', () => {
    const option = buildEmptyChartOption('暂无数据')

    expect(Array.isArray(option.xAxis)).toBe(true)
    expect(Array.isArray(option.yAxis)).toBe(true)
    expect(option.series).toEqual([])
  })

  it('normalizes board spot arrays, wrapped payloads and object payloads', () => {
    expect(normalizeBoardSpotRows([{ item: '最新', value: 123 }])).toEqual([{ item: '最新', value: 123 }])
    expect(normalizeBoardSpotRows({ data: [{ item: '涨幅', value: 1.2 }] })).toEqual([{ item: '涨幅', value: 1.2 }])
    expect(normalizeBoardSpotRows({ price: 12.3, changePercent: 2.1, name: '板块' })).toEqual([
      { item: 'price', value: 12.3 },
      { item: 'changePercent', value: 2.1 },
    ])
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
