import { getAllAShareQuotes, getKlineWithIndicators, getTodayTimeline } from '@/services/api'
import type { FullQuote } from '@/types'
import { normalizeStockCode, parseStockCode } from '@/utils/format'

export interface AnalysisProgress {
  completed: number
  total: number
  stage: string
}

export interface TimelinePoint {
  time: string
  price: number
  avgPrice: number
}

export interface TodayTimelineResponse {
  code: string
  date: string
  data: Array<{ time: string; price: number; avgPrice: number; volume: number; amount: number }>
}

export interface EndOfDayFilters {
  marketCapMin: number
  marketCapMax: number
  volumeRatioMin: number
  changePercentMin: number
  changePercentMax: number
  turnoverRateMin: number
  turnoverRateMax: number
  excludeST: boolean
  timelineAboveAvgRatio: number
}

export interface EndOfDayStock {
  code: string
  routeCode: string
  name: string
  price: number
  changePercent: number
  change: number
  volume: number
  amount: number
  turnoverRate: number | null
  volumeRatio: number | null
  circulatingMarketCap: number | null
  totalMarketCap: number | null
  pe: number | null
  pb: number | null
  high: number
  low: number
  open: number
  prevClose: number
  timeline?: TimelinePoint[]
  timelineAboveAvgRatio?: number
}

export type ScannerSignalKey =
  | 'ma_golden'
  | 'ma_death'
  | 'macd_golden'
  | 'macd_death'
  | 'rsi_oversold'
  | 'rsi_overbought'
  | 'boll_upper'
  | 'boll_lower'

export interface ScannerStockPoolItem {
  code: string
  routeCode: string
  name: string
}

export interface ScannerSignalResult {
  code: string
  routeCode: string
  name: string
  matchedSignals: string[]
}

interface ScannerKlineRow {
  close?: number | null
  ma?: Record<string, number | null | undefined>
  macd?: { dif?: number | null; dea?: number | null }
  rsi?: Record<string, number | null | undefined>
  boll?: { upper?: number | null; lower?: number | null }
}

const SCANNER_SIGNAL_LABELS: Record<ScannerSignalKey, string> = {
  ma_golden: 'MA金叉',
  ma_death: 'MA死叉',
  macd_golden: 'MACD金叉',
  macd_death: 'MACD死叉',
  rsi_oversold: 'RSI超卖',
  rsi_overbought: 'RSI超买',
  boll_upper: 'BOLL上轨',
  boll_lower: 'BOLL下轨',
}

const DEFAULT_SCAN_CONCURRENCY = 4

class AnalysisAbortError extends Error {
  code = 'ANALYSIS_ABORTED'
}

function throwIfAborted(signal?: AbortSignal) {
  if (signal?.aborted) throw new AnalysisAbortError('Analysis aborted')
}

export function isAnalysisAborted(error: unknown): boolean {
  return error instanceof AnalysisAbortError
}

export function calculateTimelineStrength(timeline: TodayTimelineResponse): { ratio: number; points: TimelinePoint[] } {
  if (!timeline.data?.length) return { ratio: 0, points: [] }
  const points = timeline.data.map((item) => ({ time: item.time, price: item.price, avgPrice: item.avgPrice }))
  const aboveAvgCount = points.filter((point) => point.price >= point.avgPrice).length
  return { ratio: (aboveAvgCount / points.length) * 100, points }
}

function toRouteCode(code: string) {
  return normalizeStockCode(code) || code
}

function toDisplayCode(code: string) {
  const normalized = toRouteCode(code)
  return parseStockCode(normalized).symbol || normalized
}

async function mapWithConcurrency<T, R>(
  items: T[],
  mapper: (item: T, index: number) => Promise<R>,
  options?: { concurrency?: number; signal?: AbortSignal; onProgress?: (completed: number, total: number) => void },
) {
  const concurrency = Math.max(1, options?.concurrency ?? DEFAULT_SCAN_CONCURRENCY)
  const results: R[] = []
  let cursor = 0
  let completed = 0
  async function worker() {
    while (true) {
      throwIfAborted(options?.signal)
      const index = cursor
      cursor += 1
      if (index >= items.length) return
      const result = await mapper(items[index], index)
      throwIfAborted(options?.signal)
      results.push(result)
      completed += 1
      options?.onProgress?.(completed, items.length)
    }
  }
  await Promise.all(Array.from({ length: Math.min(concurrency, items.length) }, () => worker()))
  return results
}

function filterBasicQuotes(quotes: FullQuote[], filters: EndOfDayFilters): EndOfDayStock[] {
  return quotes
    .filter((quote) => {
      const marketCap = quote.circulatingMarketCap
      const volumeRatio = quote.volumeRatio
      const changePercent = quote.changePercent
      const turnoverRate = quote.turnoverRate

      if (filters.excludeST && (quote.name.includes('ST') || quote.name.includes('*ST'))) return false
      if (marketCap === null || marketCap === undefined || marketCap < filters.marketCapMin || marketCap > filters.marketCapMax) return false
      if (volumeRatio === null || volumeRatio === undefined || volumeRatio < filters.volumeRatioMin) return false
      if (changePercent < filters.changePercentMin || changePercent > filters.changePercentMax) return false
      if (turnoverRate === null || turnoverRate === undefined || turnoverRate < filters.turnoverRateMin || turnoverRate > filters.turnoverRateMax) return false

      return true
    })
    .map((quote) => ({
      code: toDisplayCode(quote.code),
      routeCode: toRouteCode(quote.code),
      name: quote.name,
      price: quote.price,
      changePercent: quote.changePercent,
      change: quote.change,
      volume: quote.volume,
      amount: quote.amount,
      turnoverRate: quote.turnoverRate ?? null,
      volumeRatio: quote.volumeRatio ?? null,
      circulatingMarketCap: quote.circulatingMarketCap ?? null,
      totalMarketCap: quote.totalMarketCap ?? null,
      pe: quote.pe ?? null,
      pb: quote.pb ?? null,
      high: quote.high,
      low: quote.low,
      open: quote.open,
      prevClose: quote.prevClose,
    }))
    .sort((a, b) => b.changePercent - a.changePercent)
}

export async function analyzeEndOfDayStocks(
  filters: EndOfDayFilters,
  options?: { signal?: AbortSignal; onProgress?: (progress: AnalysisProgress) => void; timelineConcurrency?: number },
) {
  options?.onProgress?.({ completed: 0, total: 0, stage: '获取行情数据' })
  const quotes = await getAllAShareQuotes({
    batchSize: 500,
    concurrency: 4,
    onProgress: (completed, total) => {
      options?.onProgress?.({ completed, total, stage: '获取行情数据' })
    },
  })
  throwIfAborted(options?.signal)

  const basicStocks = filterBasicQuotes(quotes, filters)
  if (basicStocks.length === 0) return []

  options?.onProgress?.({ completed: 0, total: basicStocks.length, stage: '分时结构筛选' })
  const results = await mapWithConcurrency(
    basicStocks,
    async (stock) => {
      const timeline = await getTodayTimeline(stock.routeCode) as TodayTimelineResponse
      const { ratio, points } = calculateTimelineStrength(timeline)
      if (ratio < filters.timelineAboveAvgRatio) return null
      return { ...stock, timeline: points, timelineAboveAvgRatio: ratio }
    },
    { concurrency: options?.timelineConcurrency ?? DEFAULT_SCAN_CONCURRENCY, signal: options?.signal, onProgress: (completed, total) => options?.onProgress?.({ completed, total, stage: '分时结构筛选' }) },
  )
  return results
    .filter((item) => item !== null)
    .map((item) => item as EndOfDayStock)
    .sort((a, b) => (b.timelineAboveAvgRatio ?? 0) - (a.timelineAboveAvgRatio ?? 0))
}

export async function scanSignalPool(
  pool: ScannerStockPoolItem[],
  signals: ScannerSignalKey[],
  options?: {
    signal?: AbortSignal
    onProgress?: (progress: AnalysisProgress) => void
    onResult?: (result: ScannerSignalResult) => void
    concurrency?: number
    loadKline?: (stock: ScannerStockPoolItem) => Promise<ScannerKlineRow[]>
  },
) {
  const results = await mapWithConcurrency(
    pool,
    async (stock) => {
      throwIfAborted(options?.signal)
      const rows = await (options?.loadKline ? options.loadKline(stock) : getKlineWithIndicators(stock.routeCode, { period: 'daily', adjust: 'qfq' }) as Promise<ScannerKlineRow[]>)
      const matchedSignals = matchScannerSignals(rows, signals)
      if (matchedSignals.length === 0) return null
      const result = { code: stock.code, routeCode: stock.routeCode, name: stock.name, matchedSignals }
      options?.onResult?.(result)
      return result
    },
    { concurrency: options?.concurrency ?? 4, signal: options?.signal, onProgress: (completed, total) => options?.onProgress?.({ completed, total, stage: '技术信号扫描' }) },
  )
  return results.filter((item): item is ScannerSignalResult => item !== null)
}

function matchScannerSignals(rows: ScannerKlineRow[], signals: ScannerSignalKey[]) {
  if (rows.length < 2) return []
  const prev = rows[rows.length - 2]
  const latest = rows[rows.length - 1]
  return signals.filter((signal) => isSignalMatched(signal, prev, latest)).map((signal) => SCANNER_SIGNAL_LABELS[signal])
}

function isSignalMatched(signal: ScannerSignalKey, prev: ScannerKlineRow, latest: ScannerKlineRow) {
  switch (signal) {
    case 'ma_golden':
      return crossedUp(maValue(prev, 5), maValue(prev, 10), maValue(latest, 5), maValue(latest, 10))
    case 'ma_death':
      return crossedDown(maValue(prev, 5), maValue(prev, 10), maValue(latest, 5), maValue(latest, 10))
    case 'macd_golden':
      return crossedUp(prev.macd?.dif, prev.macd?.dea, latest.macd?.dif, latest.macd?.dea)
    case 'macd_death':
      return crossedDown(prev.macd?.dif, prev.macd?.dea, latest.macd?.dif, latest.macd?.dea)
    case 'rsi_oversold':
      return minRsi(latest) !== null && minRsi(latest)! <= 30
    case 'rsi_overbought':
      return maxRsi(latest) !== null && maxRsi(latest)! >= 70
    case 'boll_upper':
      return latest.close !== null && latest.close !== undefined && latest.boll?.upper !== null && latest.boll?.upper !== undefined && latest.close >= latest.boll.upper
    case 'boll_lower':
      return latest.close !== null && latest.close !== undefined && latest.boll?.lower !== null && latest.boll?.lower !== undefined && latest.close <= latest.boll.lower
    default:
      return false
  }
}

function maValue(row: ScannerKlineRow, period: number) {
  return row.ma?.[period] ?? row.ma?.[String(period)]
}

function crossedUp(prevLeft?: number | null, prevRight?: number | null, latestLeft?: number | null, latestRight?: number | null) {
  if (!isFiniteNumber(prevLeft) || !isFiniteNumber(prevRight) || !isFiniteNumber(latestLeft) || !isFiniteNumber(latestRight)) return false
  return prevLeft <= prevRight && latestLeft > latestRight
}

function crossedDown(prevLeft?: number | null, prevRight?: number | null, latestLeft?: number | null, latestRight?: number | null) {
  if (!isFiniteNumber(prevLeft) || !isFiniteNumber(prevRight) || !isFiniteNumber(latestLeft) || !isFiniteNumber(latestRight)) return false
  return prevLeft >= prevRight && latestLeft < latestRight
}

function isFiniteNumber(value: unknown): value is number {
  return typeof value === 'number' && Number.isFinite(value)
}

function rsiValues(row: ScannerKlineRow) {
  return Object.values(row.rsi || {}).filter(isFiniteNumber)
}

function minRsi(row: ScannerKlineRow) {
  const values = rsiValues(row)
  return values.length > 0 ? Math.min(...values) : null
}

function maxRsi(row: ScannerKlineRow) {
  const values = rsiValues(row)
  return values.length > 0 ? Math.max(...values) : null
}
