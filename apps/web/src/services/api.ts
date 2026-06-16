import type { Board, FinancialIndicator, FullQuote, SearchResult, StockAnnouncementDetail, StockAnnouncementResult, StockProfile, ZTPoolItem } from '@/types'
import { normalizeBoardSpotRows, type BoardSpotRow } from '@/services/charts'

const API_BASE = (import.meta.env.VITE_API_BASE_URL || '/api').replace(/\/$/, '')
type QueryValue = string | number | boolean | undefined | null
type QuoteBatchOptions = { batchSize?: number; concurrency?: number; onProgress?: (completed: number, total: number) => void }

const inFlightKlineRequests = new Map<string, Promise<unknown>>()

const INITIALISM_PREFIXES = [
  'MACD',
  'BOLL',
  'KDJ',
  'DMI',
  'ADX',
  'RSI',
  'OBV',
  'ROC',
  'SAR',
  'JSON',
  'HTML',
  'HTTP',
  'HTTPS',
  'URL',
  'API',
  'SDK',
  'EPS',
  'BPS',
  'SH',
  'SZ',
  'HK',
  'US',
  'CN',
  'ZT',
  'PE',
  'PB',
  'TZ',
  'ID',
] as const

export class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

export async function apiRequest<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: {
      Accept: 'application/json',
      ...init?.headers,
    },
  })

  if (!response.ok) {
    let message = response.statusText || `HTTP ${response.status}`
    try {
      const body = (await response.json()) as { error?: { message?: string } }
      message = body.error?.message || message
    } catch {
      // Keep status text when the server does not return JSON.
    }
    throw new ApiError(message, response.status)
  }

  return camelizeKeys(await response.json()) as T
}

function requestKline<T>(path: string): Promise<T> {
  const existing = inFlightKlineRequests.get(path) as Promise<T> | undefined
  if (existing) return existing

  const request = apiRequest<T>(path).finally(() => {
    inFlightKlineRequests.delete(path)
  })
  inFlightKlineRequests.set(path, request)
  return request
}

function toCamelKey(key: string): string {
  if (!key || /^[a-z]/.test(key)) return key
  const prefix = INITIALISM_PREFIXES.find((item) => key.startsWith(item))
  if (prefix) {
    return `${prefix.toLowerCase()}${key.slice(prefix.length)}`
  }
  return `${key.charAt(0).toLowerCase()}${key.slice(1)}`.replace(/ID$/, 'Id').replace(/URL$/, 'Url')
}

function camelizeKeys(value: unknown): unknown {
  if (Array.isArray(value)) return value.map((item) => camelizeKeys(item))
  if (value === null || typeof value !== 'object') return value
  return Object.fromEntries(
    Object.entries(value).map(([key, item]) => [toCamelKey(key), camelizeKeys(item)]),
  )
}

function query(params: Record<string, QueryValue>): string {
  const search = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      search.set(key, String(value))
    }
  })
  const value = search.toString()
  return value ? `?${value}` : ''
}

function codesQuery(codes: string[]) {
  return query({ codes: codes.join(',') })
}

function quoteSymbol(code: string): string {
  const trimmed = code.trim()
  const prefixMatch = trimmed.match(/^(sh|sz|bj)\.?(\d{6})$/i)
  if (prefixMatch) return prefixMatch[2]
  const suffixMatch = trimmed.match(/^(\d{6})\.(sh|sz|bj)$/i)
  if (suffixMatch) return suffixMatch[1]
  return trimmed
}

function normalizeRequestedQuoteCode(code: string): string {
  const trimmed = code.trim()
  const prefixMatch = trimmed.match(/^(sh|sz|bj)\.?(\d{6})$/i)
  if (prefixMatch) return `${prefixMatch[1].toLowerCase()}${prefixMatch[2]}`
  const suffixMatch = trimmed.match(/^(\d{6})\.(sh|sz|bj)$/i)
  if (suffixMatch) return `${suffixMatch[2].toLowerCase()}${suffixMatch[1]}`
  return trimmed
}

function restoreRequestedQuoteCodes(rows: FullQuote[], requestedCodes: string[]): FullQuote[] {
  const requestedBySymbol = new Map<string, string[]>()
  requestedCodes.map(normalizeRequestedQuoteCode).forEach((code) => {
    const symbol = quoteSymbol(code)
    if (!symbol) return
    const codes = requestedBySymbol.get(symbol) || []
    codes.push(code)
    requestedBySymbol.set(symbol, codes)
  })

  return rows.map((row) => {
    const requestedCode = requestedBySymbol.get(quoteSymbol(row.code))?.shift()
    return requestedCode && requestedCode !== row.code ? { ...row, code: requestedCode } : row
  })
}

export function getFullQuotes(codes: string[]) {
  return apiRequest<FullQuote[]>(`/quotes/full${codesQuery(codes)}`).then((rows) => restoreRequestedQuoteCodes(rows, codes))
}

function quoteBatchQuery(options?: QuoteBatchOptions) {
  return query({
    batchSize: options?.batchSize,
    concurrency: options?.concurrency,
  })
}

export function getAllQuotesByCodes(codes: string[], options?: QuoteBatchOptions) {
  return apiRequest<FullQuote[]>(`/quotes/batch${query({ codes: codes.join(','), batchSize: options?.batchSize, concurrency: options?.concurrency })}`)
}

export function getAllAShareQuotes(options?: QuoteBatchOptions) {
  return apiRequest<FullQuote[]>(`/quotes/a-share${quoteBatchQuery(options)}`)
}

export function search(keyword: string) {
  return apiRequest<SearchResult[]>(`/search${query({ keyword })}`)
}

export function getIndustryList() {
  return apiRequest<Board[]>('/boards/industry')
}

export function getConceptList() {
  return apiRequest<Board[]>('/boards/concept')
}

export function getIndustryConstituents(code: string) {
  return apiRequest(`/boards/industry/${encodeURIComponent(code)}/constituents`)
}

export function getConceptConstituents(code: string) {
  return apiRequest(`/boards/concept/${encodeURIComponent(code)}/constituents`)
}

export function getHistoryKline(symbol: string, options?: Record<string, string>) {
  return requestKline(`/kline/history${query({ symbol, ...options })}`)
}

export function getMinuteKline(symbol: string, options?: Record<string, string>) {
  return requestKline(`/kline/minute${query({ symbol, ...options })}`)
}

export function getKlineWithIndicators(symbol: string, options?: Record<string, string | boolean>) {
  return requestKline(`/kline/indicators${query({ symbol, ...options })}`)
}

export function getTodayTimeline(code: string) {
  return apiRequest(`/timeline/today${query({ code })}`)
}

export function getBoardSpot(type: 'industry' | 'concept', code: string) {
  return apiRequest<unknown>(`/boards/${type}/${encodeURIComponent(code)}/spot`).then((value): BoardSpotRow[] => normalizeBoardSpotRows(value))
}

export function getBoardKline(type: 'industry' | 'concept', code: string, options?: Record<string, QueryValue>) {
  return requestKline(`/boards/${type}/${encodeURIComponent(code)}/kline${query(options || {})}`)
}

export function getBoardMinuteKline(type: 'industry' | 'concept', code: string, options?: Record<string, QueryValue>) {
  return requestKline(`/boards/${type}/${encodeURIComponent(code)}/minute${query(options || {})}`)
}

export function getQuoteFundFlow(codes: string[]) {
  return apiRequest(`/fund-flow/quotes${codesQuery(codes)}`)
}

export function getPanelLargeOrder(codes: string[]) {
  return apiRequest(`/panel-large-order${codesQuery(codes)}`)
}

export function getIndividualFundFlow(symbol: string, options?: Record<string, QueryValue>) {
  return apiRequest(`/fund-flow/individual${query({ symbol, ...options })}`)
}

export function getMarketFundFlow() {
  return apiRequest('/fund-flow/market')
}

export function getFundFlowRank(options?: Record<string, QueryValue>) {
  return apiRequest(`/fund-flow/rank${query(options || {})}`)
}

export function getSectorFundFlowRank(options?: Record<string, QueryValue>) {
  return apiRequest(`/fund-flow/sector-rank${query(options || {})}`)
}

export function getSectorFundFlowHistory(symbol: string, options?: Record<string, QueryValue>) {
  return apiRequest(`/fund-flow/sector-history${query({ symbol, ...options })}`)
}

export function getNorthboundMinute(direction: 'north' | 'south' = 'north') {
  return apiRequest(`/northbound/minute${query({ direction })}`)
}

export function getNorthboundFlowSummary() {
  return apiRequest('/northbound/summary')
}

export function getNorthboundHoldingRank(options?: Record<string, QueryValue>) {
  return apiRequest(`/northbound/holding-rank${query(options || {})}`)
}

export function getNorthboundHistory(direction: 'north' | 'south' = 'north', options?: Record<string, QueryValue>) {
  return apiRequest(`/northbound/history${query({ direction, ...options })}`)
}

export function getNorthboundIndividual(symbol: string, options?: Record<string, QueryValue>) {
  return apiRequest(`/northbound/individual${query({ symbol, ...options })}`)
}

export function getZTPool(type = 'zt', date?: string) {
  return apiRequest<unknown>(`/market-event/zt-pool${query({ type, date })}`).then(normalizeZTPoolRows)
}

export function getTHSLimitUpPool(options?: { date?: string; page?: number; limit?: number }) {
  return apiRequest<unknown>(`/market-event/ths-limit-up-pool${query(options || {})}`).then(normalizeZTPoolRows)
}

function normalizeZTPoolRows(payload: unknown): ZTPoolItem[] {
  const rows = Array.isArray(payload)
    ? payload
    : payload && typeof payload === 'object' && Array.isArray((payload as { items?: unknown[] }).items)
      ? (payload as { items: unknown[] }).items
      : []

  return rows.map((row) => normalizeZTPoolItem(row)).filter((row) => row.code || row.name)
}

function normalizeZTPoolItem(row: unknown): ZTPoolItem {
  const item = row && typeof row === 'object' ? row as Record<string, unknown> : {}
  return {
    code: stringValue(item.code),
    name: stringValue(item.name),
    price: numberValue(item.price ?? item.latest),
    changePercent: numberValue(item.changePercent ?? item.changeRate),
    limitPrice: numberValue(item.limitPrice),
    amount: numberValue(item.amount),
    floatMarketValue: numberValue(item.floatMarketValue ?? item.currencyValue),
    totalMarketValue: numberValue(item.totalMarketValue),
    turnoverRate: numberValue(item.turnoverRate),
    continuousBoardCount: boardCountValue(item),
    firstBoardTime: stringValue(item.firstBoardTime ?? item.firstLimitUpTimeText),
    lastBoardTime: stringValue(item.lastBoardTime ?? item.lastLimitUpTimeText),
    boardAmount: numberValue(item.boardAmount),
    sealAmount: numberValue(item.sealAmount ?? item.orderAmount),
    failedCount: numberValue(item.failedCount ?? item.openNum),
    industry: stringValue(item.industry),
    ztStatistics: stringValue(item.ztStatistics ?? item.highDays),
    limitUpType: stringValue(item.limitUpType),
    reasonType: stringValue(item.reasonType),
    amplitude: numberValue(item.amplitude),
    speed: numberValue(item.speed),
  }
}

function stringValue(value: unknown): string {
  return typeof value === 'string' ? value : ''
}

function numberValue(value: unknown): number | null {
  return typeof value === 'number' && Number.isFinite(value) ? value : null
}

function boardCountValue(item: Record<string, unknown>): number | null {
  const continuousBoardCount = numberValue(item.continuousBoardCount)
  if (continuousBoardCount !== null) return continuousBoardCount

  const highDaysCount = boardCountFromText(item.highDays)
  if (highDaysCount !== null) return highDaysCount

  const highDaysValue = numberValue(item.highDaysValue)
  return highDaysValue !== null && highDaysValue > 0 && highDaysValue < 100 ? highDaysValue : null
}

function boardCountFromText(value: unknown): number | null {
  if (typeof value !== 'string') return null
  if (value.includes('首板')) return 1

  const matches = Array.from(value.matchAll(/(\d+(?:\.\d+)?)\s*板/g))
  if (matches.length === 0) return null
  const count = Number(matches[matches.length - 1][1])
  return Number.isFinite(count) ? count : null
}

export function getStockChanges(type = 'large_buy') {
  return apiRequest(`/market-event/stock-changes${query({ type })}`)
}

export function getBoardChanges() {
  return apiRequest('/market-event/board-changes')
}

export function getDragonTigerDetail(options?: Record<string, QueryValue>) {
  return apiRequest(`/dragon-tiger/detail${query(options || {})}`)
}

export function getBlockTradeDetail(options?: Record<string, QueryValue>) {
  return apiRequest(`/block-trade/detail${query(options || {})}`)
}

export function getMarginAccountInfo() {
  return apiRequest('/margin/account')
}

export function getDividendDetail(symbol: string) {
  return apiRequest(`/dividends${query({ symbol })}`)
}

export function getStockProfile(symbol: string) {
  return apiRequest<StockProfile>(`/stocks/profile${query({ symbol })}`)
}

export function getFinancialIndicators(symbol: string, options?: { period?: 'all' | 'annual' }) {
  return apiRequest<FinancialIndicator[]>(`/stocks/financial-indicators${query({ symbol, ...options })}`)
}

export function getStockAnnouncements(symbol: string, options?: { pageSize?: number; pageIndex?: number }) {
  return apiRequest<StockAnnouncementResult>(`/stocks/announcements${query({ symbol, ...options })}`)
}

export function getStockAnnouncementDetail(artCode: string) {
  return apiRequest<StockAnnouncementDetail>(`/stocks/announcements/${encodeURIComponent(artCode)}`)
}

export function getTradingCalendar() {
  return apiRequest<string[]>('/trading-calendar')
}
