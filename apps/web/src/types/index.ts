export interface WatchlistGroup {
  id: string
  name: string
  codes: string[]
  createdAt: number
  updatedAt: number
}

export type AlertType =
  | 'price_gte'
  | 'price_lte'
  | 'change_percent_gte'
  | 'change_percent_lte'
  | 'amount_gte'
  | 'near_limit_up'
  | 'near_limit_down'

export interface AlertRule {
  id: string
  code: string
  name: string
  type: AlertType
  value: number
  cooldownSec: number
  enabled: boolean
  lastTriggeredAt: number
  createdAt: number
}

export interface HeatmapConfig {
  dimension: 'industry' | 'concept' | 'stock' | 'watchlist'
  colorField: 'changePercent' | 'change' | 'volumeRatio' | 'turnoverRate'
  sizeField: 'totalMarketCap' | 'amount' | 'volume'
  colorMode: 'red-rise' | 'green-rise'
  topK: number
}

export interface IndicatorConfig {
  ma: number[]
  macd: { short: number; long: number; signal: number }
  boll: { period: number; stdDev: number }
  kdj: { period: number; kPeriod: number; dPeriod: number }
  rsi: number[]
  dmi: { period: number; adxPeriod: number }
  sar: { afStart: number; afIncrement: number; afMax: number }
  kc: { emaPeriod: number; atrPeriod: number; multiplier: number }
}

export interface ColumnConfig {
  key: string
  label: string
  visible: boolean
  width?: number
}

export interface AppSettings {
  refreshInterval: {
    list: number
    detail: number
    heatmap: number
  }
  colorMode: 'red-rise' | 'green-rise'
  heatmapConfig: HeatmapConfig
  indicatorConfig: IndicatorConfig
}

export interface SearchHistoryItem {
  code: string
  name: string
  market: string
  type: string
  timestamp: number
}

export interface CacheItem<T> {
  data: T
  timestamp: number
  ttl: number
}

export type RefreshStatus = 'idle' | 'loading' | 'success' | 'error'
export type MarketStatus = 'pre' | 'trading' | 'break' | 'closed'
export type KlinePeriod = 'daily' | 'weekly' | 'monthly'
export type MinutePeriod = '1' | '5' | '15' | '30' | '60'
export type AdjustType = '' | 'qfq' | 'hfq'
export type SortDirection = 'asc' | 'desc'

export interface SortConfig {
  field: string
  direction: SortDirection
}

export interface FullQuote {
  code: string
  name: string
  price: number
  prevClose: number
  open: number
  high: number
  low: number
  change: number
  changePercent: number
  volume: number
  amount: number
  turnoverRate?: number | null
  pe?: number | null
  pb?: number | null
  volumeRatio?: number | null
  totalMarketCap?: number | null
  circulatingMarketCap?: number | null
  limitUp?: number | null
  limitDown?: number | null
  bid?: Array<{ price: number; volume: number }>
  ask?: Array<{ price: number; volume: number }>
  time?: string
  source?: string
  assetType?: string
  market?: string
}

export interface Board {
  rank: number
  name: string
  code: string
  price?: number | null
  change?: number | null
  changePercent?: number | null
  totalMarketCap?: number | null
  turnoverRate?: number | null
  riseCount?: number | null
  fallCount?: number | null
  leadingStock?: string | null
  leadingStockChangePercent?: number | null
}

export interface SearchResult {
  code: string
  name: string
  market: string
  type: string
  category?: string
}
