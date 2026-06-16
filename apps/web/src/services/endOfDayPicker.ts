import type { EndOfDayFilters, EndOfDayStock } from '@/services/analysis'

export type EndOfDaySortField =
  | 'changePercent'
  | 'timelineAboveAvgRatio'
  | 'turnoverRate'
  | 'circulatingMarketCap'
  | 'volumeRatio'
export type EndOfDaySortOrder = 'asc' | 'desc'

export interface EndOfDayScheme {
  id: string
  name: string
  filters: EndOfDayFilters
  createdAt: number
}

export interface EndOfDayRecentUsage {
  filters: EndOfDayFilters
  usedAt: number
}

const FILTERS_STORAGE_KEY = 'end-of-day-picker-settings'
const SCHEMES_STORAGE_KEY = 'end-of-day-picker-schemes'
const RECENT_USAGE_STORAGE_KEY = 'end-of-day-picker-recent'
const MAX_RECENT_USAGE = 5

export const DEFAULT_END_OF_DAY_FILTERS: EndOfDayFilters = {
  marketCapMin: 50,
  marketCapMax: 200,
  volumeRatioMin: 1.2,
  changePercentMin: 3,
  changePercentMax: 5,
  turnoverRateMin: 5,
  turnoverRateMax: 10,
  excludeST: true,
  timelineAboveAvgRatio: 80,
}

function safeJsonParse<T>(value: string | null, fallback: T): T {
  if (!value) return fallback
  try {
    return JSON.parse(value) as T
  } catch {
    return fallback
  }
}

function cloneFilters(filters: EndOfDayFilters): EndOfDayFilters {
  return { ...DEFAULT_END_OF_DAY_FILTERS, ...filters }
}

function filtersKey(filters: EndOfDayFilters): string {
  return JSON.stringify(cloneFilters(filters))
}

export function getEndOfDayFilters(): EndOfDayFilters {
  return cloneFilters(safeJsonParse<Partial<EndOfDayFilters>>(localStorage.getItem(FILTERS_STORAGE_KEY), {} as Partial<EndOfDayFilters>) as EndOfDayFilters)
}

export function saveEndOfDayFilters(filters: EndOfDayFilters): void {
  localStorage.setItem(FILTERS_STORAGE_KEY, JSON.stringify(cloneFilters(filters)))
}

export function getEndOfDaySchemes(): EndOfDayScheme[] {
  return safeJsonParse(localStorage.getItem(SCHEMES_STORAGE_KEY), [])
}

export function saveEndOfDaySchemes(schemes: EndOfDayScheme[]): void {
  localStorage.setItem(SCHEMES_STORAGE_KEY, JSON.stringify(schemes))
}

export function saveEndOfDayScheme(name: string, filters: EndOfDayFilters, now = Date.now()): EndOfDayScheme {
  const scheme: EndOfDayScheme = {
    id: `scheme_${now}_${Math.random().toString(36).slice(2, 8)}`,
    name: name.trim(),
    filters: cloneFilters(filters),
    createdAt: now,
  }
  saveEndOfDaySchemes([...getEndOfDaySchemes(), scheme])
  return scheme
}

export function deleteEndOfDayScheme(id: string): void {
  saveEndOfDaySchemes(getEndOfDaySchemes().filter((scheme) => scheme.id !== id))
}

export function getEndOfDayRecentUsage(): EndOfDayRecentUsage[] {
  return safeJsonParse(localStorage.getItem(RECENT_USAGE_STORAGE_KEY), [])
}

export function saveEndOfDayRecentUsage(items: EndOfDayRecentUsage[]): void {
  localStorage.setItem(RECENT_USAGE_STORAGE_KEY, JSON.stringify(items.slice(0, MAX_RECENT_USAGE)))
}

export function addEndOfDayRecentUsage(filters: EndOfDayFilters, now = Date.now()): EndOfDayRecentUsage {
  const normalized = cloneFilters(filters)
  const key = filtersKey(normalized)
  const recent = getEndOfDayRecentUsage()
  const existing = recent.find((item) => filtersKey(item.filters) === key)
  if (existing) return existing
  const entry = { filters: normalized, usedAt: now }
  saveEndOfDayRecentUsage([entry, ...recent])
  return entry
}

function numericValue(stock: EndOfDayStock, field: EndOfDaySortField): number {
  const value = stock[field]
  return typeof value === 'number' && Number.isFinite(value) ? value : 0
}

export function sortEndOfDayStocks(
  stocks: EndOfDayStock[],
  field: EndOfDaySortField,
  order: EndOfDaySortOrder,
): EndOfDayStock[] {
  return [...stocks].sort((a, b) => {
    const left = numericValue(a, field)
    const right = numericValue(b, field)
    return order === 'desc' ? right - left : left - right
  })
}

export function toggleSelectedCode(selected: string[], code: string): string[] {
  return selected.includes(code) ? selected.filter((item) => item !== code) : [...selected, code]
}

export function getBatchWatchlistCandidates(
  stocks: EndOfDayStock[],
  selectedRouteCodes: string[],
  isInWatchlist: (code: string) => boolean,
): string[] {
  const selected = new Set(selectedRouteCodes)
  return stocks
    .filter((stock) => selected.has(stock.routeCode) && !isInWatchlist(stock.routeCode))
    .map((stock) => stock.routeCode)
}
