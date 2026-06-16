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

export interface StockIssueInfo {
  foundDate?: string | null
  listingDate?: string | null
  issueWay?: string
  parValue?: number | null
  totalIssueShares?: number | null
  issuePrice?: number | null
  totalFunds?: number | null
  netRaiseFunds?: number | null
  openPrice?: number | null
  closePrice?: number | null
  turnoverRate?: number | null
}

export interface StockProfile {
  secuCode: string
  code: string
  name: string
  orgName?: string
  orgNameEn?: string
  formerName?: string
  securityType?: string
  industry?: string
  tradeMarket?: string
  csrcIndustry?: string
  president?: string
  legalRepresentative?: string
  secretary?: string
  chairman?: string
  securitiesRepresentative?: string
  tel?: string
  email?: string
  fax?: string
  website?: string
  address?: string
  registeredAddress?: string
  province?: string
  addressPostcode?: string
  registeredCapital?: number | null
  registrationNumber?: string
  employeeCount?: number | null
  lawFirm?: string
  accountingFirm?: string
  profile?: string
  businessScope?: string
  issue?: StockIssueInfo
}

export interface FinancialIndicator {
  secuCode: string
  code: string
  name: string
  reportDate?: string | null
  reportType?: string
  reportDateName?: string
  noticeDate?: string | null
  updateDate?: string | null
  currency?: string
  basicEps?: number | null
  deductBasicEps?: number | null
  dilutedEps?: number | null
  bps?: number | null
  capitalReservePerShare?: number | null
  unassignedProfitPerShare?: number | null
  operatingCashFlowPerShare?: number | null
  totalRevenue?: number | null
  grossProfit?: number | null
  parentNetProfit?: number | null
  deductParentNetProfit?: number | null
  totalRevenueYoY?: number | null
  parentNetProfitYoY?: number | null
  deductParentNetProfitYoY?: number | null
  roeWeighted?: number | null
  roeDeductWeighted?: number | null
  roa?: number | null
  netMargin?: number | null
  grossMargin?: number | null
  assetLiabilityRatio?: number | null
  roic?: number | null
  staffCount?: number | null
}

export interface StockAnnouncementColumn {
  code: string
  name: string
}

export interface StockAnnouncement {
  artCode: string
  title: string
  titleCh?: string
  titleEn?: string
  noticeDate?: string | null
  displayTime?: string | null
  sortDate?: string | null
  columns?: StockAnnouncementColumn[]
}

export interface StockAnnouncementResult {
  list: StockAnnouncement[]
  pageIndex: number
  pageSize: number
  total: number
}

export interface StockAnnouncementAttachment {
  url: string
  type: string
  size?: number | null
  seq?: number | null
}

export interface StockAnnouncementDetail {
  artCode: string
  title?: string
  noticeDate?: string | null
  attachUrl?: string
  attachUrlWeb?: string
  attachSize?: string
  attachType?: string
  noticeContent?: string
  attachments?: StockAnnouncementAttachment[]
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

export interface ZTPoolItem {
  code: string
  name: string
  price?: number | null
  changePercent?: number | null
  limitPrice?: number | null
  amount?: number | null
  floatMarketValue?: number | null
  totalMarketValue?: number | null
  turnoverRate?: number | null
  continuousBoardCount?: number | null
  firstBoardTime?: string | null
  lastBoardTime?: string | null
  boardAmount?: number | null
  sealAmount?: number | null
  failedCount?: number | null
  industry?: string | null
  ztStatistics?: string | null
  limitUpType?: string | null
  reasonType?: string | null
  amplitude?: number | null
  speed?: number | null
}

export interface SearchResult {
  code: string
  name: string
  market: string
  type: string
  category?: string
}
