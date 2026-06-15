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
