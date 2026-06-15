import type { IndicatorConfig } from '@/types'

export type ChartAxis = Record<string, unknown>
export type ChartSeries = Record<string, unknown>
export type OverlayIndicatorKey = 'ma' | 'boll' | 'sar' | 'kc'
export type OscillatorIndicatorKey = 'macd' | 'kdj' | 'rsi' | 'obv' | 'roc' | 'dmi'

export interface ChartOption extends Record<string, unknown> {
  animation?: boolean
  title?: Record<string, unknown>
  tooltip?: Record<string, unknown>
  legend?: Record<string, unknown>
  grid?: ChartAxis | ChartAxis[]
  xAxis: ChartAxis[]
  yAxis: ChartAxis[]
  dataZoom?: Array<Record<string, unknown>>
  series: ChartSeries[]
}

export interface KlineVolumeRow {
  date: string
  open?: number | null
  close?: number | null
  low?: number | null
  high?: number | null
  volume?: number | null
}

export interface IndicatorKlineRow extends KlineVolumeRow {
  ma?: Record<string, number | null | undefined>
  macd?: { dif?: number | null; dea?: number | null; macd?: number | null }
  boll?: { upper?: number | null; mid?: number | null; lower?: number | null }
  kdj?: { k?: number | null; d?: number | null; j?: number | null }
  rsi?: Record<string, number | null | undefined>
  obv?: { obv?: number | null; obvMa?: number | null }
  roc?: { roc?: number | null; signal?: number | null }
  dmi?: { pdi?: number | null; mdi?: number | null; adx?: number | null }
  sar?: { sar?: number | null; trend?: number | null }
  kc?: { upper?: number | null; mid?: number | null; lower?: number | null }
}

export interface BoardSpotRow {
  item: string
  value?: number | string | null
}

export interface FundFlowRow {
  date: string
  mainNetInflow?: number | null
  mainNetInflowPercent?: number | null
  superLargeNetInflow?: number | null
  smallNetInflow?: number | null
}

export interface TimelineChartPoint {
  time: string
  price?: number | null
  avgPrice?: number | null
}

export interface TimelineChartPayload {
  preClose?: number | null
  data?: TimelineChartPoint[]
}

export interface MinuteKlineRow {
  time: string
  open?: number | null
  close?: number | null
  low?: number | null
  high?: number | null
  volume?: number | null
}

export function buildEmptyChartOption(text: string): ChartOption {
  return {
    title: { text, left: 'center', top: 'middle', textStyle: { color: '#8c8c8c', fontSize: 14, fontWeight: 400 } },
    xAxis: [{ show: false }],
    yAxis: [{ show: false }],
    series: [],
  }
}

export function buildMinuteChartOption(options: {
  period: string
  timeline: TimelineChartPayload | null
  minuteKline: MinuteKlineRow[]
  emptyText: string
}): ChartOption {
  if (options.period === '1') {
    const rows = options.timeline?.data || []
    if (rows.length === 0) return buildEmptyChartOption(options.emptyText)
    const prices = rows.map((item) => item.price ?? null)
    const finitePrices = prices.filter((item): item is number => typeof item === 'number' && Number.isFinite(item))
    const base = options.timeline?.preClose ?? finitePrices[0] ?? 0
    const spread = Math.max(...finitePrices.map((price) => Math.abs(price - base)), 0.01) * 1.2
    return {
      animation: false,
      tooltip: { trigger: 'axis' },
      grid: { left: 52, right: 18, top: 24, bottom: 34 },
      xAxis: [{ type: 'category', data: rows.map((item) => item.time), boundaryGap: false }],
      yAxis: [{ type: 'value', min: base - spread, max: base + spread, axisLabel: { formatter: (value: number) => value.toFixed(2) } }],
      series: [
        { name: '价格', type: 'line', symbol: 'none', data: prices, lineStyle: { color: '#1677ff', width: 1.5 } },
        { name: '均价', type: 'line', symbol: 'none', data: rows.map((item) => item.avgPrice ?? null), lineStyle: { color: '#faad14', width: 1, type: 'dashed' } },
      ],
    }
  }

  if (options.minuteKline.length === 0) return buildEmptyChartOption(options.emptyText)
  const rows = options.minuteKline
  const times = rows.map((item) => item.time)
  const startPercent = rows.length > 80 ? Math.max(0, ((rows.length - 80) / rows.length) * 100) : 0
  return {
    animation: false,
    tooltip: { trigger: 'axis' },
    legend: { top: 0, data: [`${options.period}分K`, '成交量'] },
    grid: [{ left: 52, right: 18, top: 32, height: 260 }, { left: 52, right: 18, top: 318, height: 78 }],
    xAxis: [
      { type: 'category', data: times, boundaryGap: true },
      { type: 'category', data: times, gridIndex: 1, boundaryGap: true, axisLabel: { show: false } },
    ],
    yAxis: [{ scale: true }, { gridIndex: 1, axisLabel: { show: false }, splitLine: { show: false } }],
    dataZoom: [
      { type: 'inside', xAxisIndex: [0, 1], start: startPercent, end: 100 },
      { type: 'slider', xAxisIndex: [0, 1], start: startPercent, end: 100, bottom: 0, height: 18 },
    ],
    series: [
      {
        name: `${options.period}分K`,
        type: 'candlestick',
        data: rows.map((item) => [item.open, item.close, item.low, item.high]),
        itemStyle: { color: '#cf1322', color0: '#389e0d', borderColor: '#cf1322', borderColor0: '#389e0d' },
      },
      {
        name: '成交量',
        type: 'bar',
        xAxisIndex: 1,
        yAxisIndex: 1,
        data: rows.map((item) => ({
          value: item.volume || 0,
          itemStyle: { color: (item.close ?? 0) >= (item.open ?? 0) ? '#cf1322' : '#389e0d' },
        })),
      },
    ],
  }
}

export function buildKlineVolumeOption(rows: KlineVolumeRow[], options: { emptyText: string }): ChartOption {
  if (rows.length === 0) return buildEmptyChartOption(options.emptyText)

  const dates = rows.map((item) => item.date)
  return {
    animation: false,
    tooltip: { trigger: 'axis' },
    legend: { top: 0, data: ['K线', '成交量'] },
    grid: [{ left: 52, right: 18, top: 32, height: 250 }, { left: 52, right: 18, top: 310, height: 80 }],
    xAxis: [
      { type: 'category', data: dates, boundaryGap: true },
      { type: 'category', data: dates, gridIndex: 1, boundaryGap: true, axisLabel: { show: false } },
    ],
    yAxis: [{ scale: true }, { gridIndex: 1 }],
    dataZoom: [{ type: 'inside', xAxisIndex: [0, 1] }, { type: 'slider', xAxisIndex: [0, 1], bottom: 0 }],
    series: [
      {
        name: 'K线',
        type: 'candlestick',
        data: rows.map((item) => [item.open, item.close, item.low, item.high]),
        itemStyle: { color: '#cf1322', color0: '#389e0d', borderColor: '#cf1322', borderColor0: '#389e0d' },
      },
      {
        name: '成交量',
        type: 'bar',
        xAxisIndex: 1,
        yAxisIndex: 1,
        data: rows.map((item) => item.volume || 0),
        itemStyle: { color: '#8c8c8c' },
      },
    ],
  }
}

export function buildIndicatorKlineOption(
  rows: IndicatorKlineRow[],
  options: {
    emptyText: string
    overlays: OverlayIndicatorKey[]
    oscillator: OscillatorIndicatorKey
    indicatorConfig: IndicatorConfig
  },
): ChartOption {
  if (rows.length === 0) return buildEmptyChartOption(options.emptyText)

  const dates = rows.map((item) => item.date)
  const startPercent = rows.length > 80 ? Math.max(0, ((rows.length - 80) / rows.length) * 100) : 0
  const series: ChartSeries[] = [
    {
      name: 'K线',
      type: 'candlestick',
      data: rows.map((item) => [item.open, item.close, item.low, item.high]),
      itemStyle: { color: '#cf1322', color0: '#389e0d', borderColor: '#cf1322', borderColor0: '#389e0d' },
    },
    {
      name: '成交量',
      type: 'bar',
      xAxisIndex: 1,
      yAxisIndex: 1,
      data: rows.map((item) => ({
        value: item.volume || 0,
        itemStyle: { color: (item.close ?? 0) >= (item.open ?? 0) ? '#cf1322' : '#389e0d' },
      })),
    },
  ]

  addOverlaySeries(series, rows, options.overlays, options.indicatorConfig)
  addOscillatorSeries(series, rows, options.oscillator, options.indicatorConfig)

  return {
    animation: false,
    tooltip: { trigger: 'axis' },
    legend: { type: 'scroll', top: 0, data: series.map((item) => item.name).filter(Boolean) },
    grid: [
      { left: 52, right: 18, top: 38, height: 230 },
      { left: 52, right: 18, top: 294, height: 72 },
      { left: 52, right: 18, top: 394, height: 92 },
    ],
    xAxis: [
      { type: 'category', data: dates, boundaryGap: true, axisLabel: { show: false } },
      { type: 'category', data: dates, gridIndex: 1, boundaryGap: true, axisLabel: { show: false } },
      { type: 'category', data: dates, gridIndex: 2, boundaryGap: true },
    ],
    yAxis: [
      { scale: true },
      { gridIndex: 1, axisLabel: { show: false }, splitLine: { show: false } },
      { gridIndex: 2, scale: true },
    ],
    dataZoom: [
      { type: 'inside', xAxisIndex: [0, 1, 2], start: startPercent, end: 100 },
      { type: 'slider', xAxisIndex: [0, 1, 2], start: startPercent, end: 100, bottom: 0, height: 18 },
    ],
    series,
  }
}

export function buildFundFlowOption(rows: FundFlowRow[], options: { emptyText: string }): ChartOption {
  if (rows.length === 0) return buildEmptyChartOption(options.emptyText)

  return {
    animation: false,
    tooltip: { trigger: 'axis' },
    legend: { top: 0, data: ['主力净流入', '净占比'] },
    grid: { left: 58, right: 48, top: 36, bottom: 34 },
    xAxis: [{ type: 'category', data: rows.map((item) => item.date), boundaryGap: true }],
    yAxis: [
      {
        type: 'value',
        axisLabel: { formatter: (value: number) => `${value.toFixed(1)}亿` },
      },
      {
        type: 'value',
        axisLabel: { formatter: (value: number) => `${value.toFixed(1)}%` },
        splitLine: { show: false },
      },
    ],
    series: [
      {
        name: '主力净流入',
        type: 'bar',
        data: rows.map((item) => ({
          value: (item.mainNetInflow ?? 0) / 100000000,
          itemStyle: { color: (item.mainNetInflow ?? 0) >= 0 ? '#cf1322' : '#389e0d' },
        })),
      },
      {
        name: '净占比',
        type: 'line',
        yAxisIndex: 1,
        data: rows.map((item) => item.mainNetInflowPercent ?? null),
        symbol: 'none',
        lineStyle: { width: 1.5, color: '#1677ff' },
      },
    ],
  }
}

function addOverlaySeries(series: ChartSeries[], rows: IndicatorKlineRow[], overlays: OverlayIndicatorKey[], indicatorConfig: IndicatorConfig) {
  if (overlays.includes('ma')) {
    const colors = ['#faad14', '#1677ff', '#eb2f96', '#722ed1', '#13c2c2', '#fa8c16']
    indicatorConfig.ma.forEach((period, index) => {
      series.push({
        name: `MA${period}`,
        type: 'line',
        data: rows.map((item) => indicatorMapValue(item.ma, 'ma', period)),
        symbol: 'none',
        lineStyle: { width: 1, color: colors[index % colors.length] },
      })
    })
  }
  if (overlays.includes('boll')) {
    series.push(
      overlayLine('BOLL上轨', rows.map((item) => item.boll?.upper ?? null), '#faad14', 'dashed'),
      overlayLine('BOLL中轨', rows.map((item) => item.boll?.mid ?? null), '#722ed1'),
      overlayLine('BOLL下轨', rows.map((item) => item.boll?.lower ?? null), '#faad14', 'dashed'),
    )
  }
  if (overlays.includes('sar')) {
    series.push({
      name: 'SAR',
      type: 'scatter',
      data: rows.map((item) => item.sar?.sar ?? null),
      symbolSize: 5,
      itemStyle: { color: '#13c2c2' },
    })
  }
  if (overlays.includes('kc')) {
    series.push(
      overlayLine('KC上轨', rows.map((item) => item.kc?.upper ?? null), '#52c41a'),
      overlayLine('KC中轨', rows.map((item) => item.kc?.mid ?? null), '#1677ff'),
      overlayLine('KC下轨', rows.map((item) => item.kc?.lower ?? null), '#52c41a'),
    )
  }
}

function addOscillatorSeries(series: ChartSeries[], rows: IndicatorKlineRow[], oscillator: OscillatorIndicatorKey, indicatorConfig: IndicatorConfig) {
  if (oscillator === 'kdj') {
    series.push(
      oscillatorLine('K', rows.map((item) => item.kdj?.k ?? null), '#faad14'),
      oscillatorLine('D', rows.map((item) => item.kdj?.d ?? null), '#1677ff'),
      oscillatorLine('J', rows.map((item) => item.kdj?.j ?? null), '#eb2f96'),
    )
    return
  }
  if (oscillator === 'rsi') {
    const colors = ['#faad14', '#1677ff', '#eb2f96', '#13c2c2']
    indicatorConfig.rsi.forEach((period, index) => {
      series.push(oscillatorLine(`RSI${period}`, rows.map((item) => indicatorMapValue(item.rsi, 'rsi', period)), colors[index % colors.length]))
    })
    return
  }
  if (oscillator === 'obv') {
    series.push(
      oscillatorLine('OBV', rows.map((item) => item.obv?.obv ?? null), '#1677ff'),
      oscillatorLine('OBV MA', rows.map((item) => item.obv?.obvMa ?? null), '#faad14'),
    )
    return
  }
  if (oscillator === 'roc') {
    series.push(
      oscillatorLine('ROC', rows.map((item) => item.roc?.roc ?? null), '#52c41a'),
      oscillatorLine('ROC Signal', rows.map((item) => item.roc?.signal ?? null), '#faad14'),
    )
    return
  }
  if (oscillator === 'dmi') {
    series.push(
      oscillatorLine('+DI', rows.map((item) => item.dmi?.pdi ?? null), '#52c41a'),
      oscillatorLine('-DI', rows.map((item) => item.dmi?.mdi ?? null), '#cf1322'),
      oscillatorLine('ADX', rows.map((item) => item.dmi?.adx ?? null), '#1677ff'),
    )
    return
  }
  series.push(
    {
      name: 'MACD',
      type: 'bar',
      xAxisIndex: 2,
      yAxisIndex: 2,
      data: rows.map((item) => ({
        value: item.macd?.macd ?? 0,
        itemStyle: { color: (item.macd?.macd ?? 0) >= 0 ? '#cf1322' : '#389e0d' },
      })),
    },
    oscillatorLine('DIF', rows.map((item) => item.macd?.dif ?? null), '#faad14'),
    oscillatorLine('DEA', rows.map((item) => item.macd?.dea ?? null), '#1677ff'),
  )
}

function overlayLine(name: string, data: Array<number | null>, color: string, type?: string): ChartSeries {
  return {
    name,
    type: 'line',
    data,
    symbol: 'none',
    lineStyle: { width: 1, color, ...(type ? { type } : {}) },
  }
}

function oscillatorLine(name: string, data: Array<number | null | undefined>, color: string): ChartSeries {
  return {
    name,
    type: 'line',
    xAxisIndex: 2,
    yAxisIndex: 2,
    data,
    symbol: 'none',
    lineStyle: { width: 1, color },
  }
}

function indicatorMapValue(map: Record<string, number | null | undefined> | undefined, prefix: string, period: number) {
  return map?.[`${prefix}${period}`] ?? map?.[String(period)] ?? null
}

export function normalizeBoardSpotRows(value: unknown): BoardSpotRow[] {
  if (Array.isArray(value)) {
    return value
      .map((item) => normalizeSpotRow(item))
      .filter((item): item is BoardSpotRow => item !== null)
  }

  if (value && typeof value === 'object') {
    const record = value as Record<string, unknown>
    const wrapped = record.data ?? record.rows ?? record.list ?? record.items
    if (Array.isArray(wrapped)) return normalizeBoardSpotRows(wrapped)

    return Object.entries(record)
      .filter(([, item]) => typeof item === 'number')
      .map(([key, item]) => ({ item: key, value: item as number }))
  }

  return []
}

function normalizeSpotRow(value: unknown): BoardSpotRow | null {
  if (!value || typeof value !== 'object') return null
  const record = value as Record<string, unknown>
  const item = record.item ?? record.name ?? record.label ?? record.metric
  if (typeof item !== 'string' || item.length === 0) return null
  const rawValue = record.value ?? record.amount ?? record.price
  if (rawValue !== null && rawValue !== undefined && typeof rawValue !== 'number' && typeof rawValue !== 'string') return null
  return { item, value: rawValue as number | string | null | undefined }
}
