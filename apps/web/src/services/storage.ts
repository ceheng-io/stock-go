import type { AlertRule, AppSettings, ColumnConfig, HeatmapConfig, IndicatorConfig, SearchHistoryItem, WatchlistGroup } from '@/types'
import { normalizeStockCode } from '@/utils/format'

const STORAGE_KEYS = {
  WATCHLIST_GROUPS: 'watchlist.groups',
  ALERTS: 'watchlist.alerts',
  SETTINGS: 'app.settings',
  TABLE_COLUMNS: 'ui.tableColumns',
  HEATMAP_CONFIG: 'ui.heatmapConfig',
  INDICATOR_CONFIG: 'ui.indicatorConfig',
  SEARCH_HISTORY: 'search.recent',
} as const

const DEFAULT_SETTINGS: AppSettings = {
  refreshInterval: { list: 0, detail: 5000, heatmap: 10000 },
  colorMode: 'red-rise',
  heatmapConfig: {
    dimension: 'industry',
    colorField: 'changePercent',
    sizeField: 'totalMarketCap',
    colorMode: 'red-rise',
    topK: 200,
  },
  indicatorConfig: {
    ma: [5, 10, 20, 60],
    macd: { short: 12, long: 26, signal: 9 },
    boll: { period: 20, stdDev: 2 },
    kdj: { period: 9, kPeriod: 3, dPeriod: 3 },
    rsi: [6, 12, 24],
    dmi: { period: 14, adxPeriod: 14 },
    sar: { afStart: 0.02, afIncrement: 0.02, afMax: 0.2 },
    kc: { emaPeriod: 20, atrPeriod: 10, multiplier: 2 },
  },
}

const DEFAULT_REFRESH_INTERVALS: AppSettings['refreshInterval'] = {
  list: 15000,
  detail: 5000,
  heatmap: 10000,
}

const DEFAULT_WATCHLIST_GROUPS: WatchlistGroup[] = [
  { id: 'default', name: '默认分组', codes: [], createdAt: Date.now(), updatedAt: Date.now() },
]

function safeJsonParse<T>(value: string | null, fallback: T): T {
  if (!value) return fallback
  try {
    return JSON.parse(value) as T
  } catch {
    return fallback
  }
}

export function getWatchlistGroups(): WatchlistGroup[] {
  const groups = safeJsonParse(localStorage.getItem(STORAGE_KEYS.WATCHLIST_GROUPS), DEFAULT_WATCHLIST_GROUPS)
  let changed = false
  const normalized = groups.map((group) => {
    const seen = new Set<string>()
    const codes = group.codes
      .map(normalizeStockCode)
      .filter((code) => {
        if (!code || seen.has(code)) {
          changed = true
          return false
        }
        seen.add(code)
        return true
      })
    if (codes.length !== group.codes.length || codes.some((code, index) => code !== group.codes[index])) {
      changed = true
      return { ...group, codes }
    }
    return group
  })
  if (changed) saveWatchlistGroups(normalized)
  return normalized
}

export function saveWatchlistGroups(groups: WatchlistGroup[]) {
  localStorage.setItem(STORAGE_KEYS.WATCHLIST_GROUPS, JSON.stringify(groups))
}

export function addToWatchlist(code: string, groupId = 'default') {
  const normalized = normalizeStockCode(code)
  if (!normalized) return
  const groups = getWatchlistGroups()
  const group = groups.find((item) => item.id === groupId) || groups[0]
  if (!group.codes.includes(normalized)) {
    group.codes.push(normalized)
    group.updatedAt = Date.now()
    saveWatchlistGroups(groups)
  }
}

export function removeFromWatchlist(code: string, groupId?: string) {
  const normalized = normalizeStockCode(code)
  const groups = getWatchlistGroups()
  groups.forEach((group) => {
    if (!groupId || group.id === groupId) {
      group.codes = group.codes.filter((item) => item !== normalized)
      group.updatedAt = Date.now()
    }
  })
  saveWatchlistGroups(groups)
}

export function getAllWatchlistCodes(): string[] {
  return Array.from(new Set(getWatchlistGroups().flatMap((group) => group.codes)))
}

export function isInWatchlist(code: string): boolean {
  return getAllWatchlistCodes().includes(normalizeStockCode(code))
}

export function createWatchlistGroup(name: string): WatchlistGroup {
  const groups = getWatchlistGroups()
  const group = { id: `group_${Date.now()}`, name, codes: [], createdAt: Date.now(), updatedAt: Date.now() }
  groups.push(group)
  saveWatchlistGroups(groups)
  return group
}

export function renameWatchlistGroup(groupId: string, name: string) {
  const trimmed = name.trim()
  if (!trimmed) return
  const groups = getWatchlistGroups()
  const group = groups.find((item) => item.id === groupId)
  if (!group) return
  group.name = trimmed
  group.updatedAt = Date.now()
  saveWatchlistGroups(groups)
}

export function deleteWatchlistGroup(groupId: string) {
  if (groupId === 'default') return
  const groups = getWatchlistGroups().filter((group) => group.id !== groupId)
  saveWatchlistGroups(groups.length > 0 ? groups : DEFAULT_WATCHLIST_GROUPS)
}

export function batchAddToWatchlist(codes: string[], groupId = 'default') {
  const groups = getWatchlistGroups()
  const group = groups.find((item) => item.id === groupId) || groups[0]
  const seen = new Set(group.codes.map(normalizeStockCode))
  let addedCount = 0
  codes.map(normalizeStockCode).forEach((code) => {
    if (code && !seen.has(code)) {
      seen.add(code)
      group.codes.push(code)
      addedCount += 1
    }
  })
  if (addedCount > 0) {
    group.updatedAt = Date.now()
    saveWatchlistGroups(groups)
  }
  return addedCount
}

export function batchRemoveFromWatchlist(codes: string[], groupId?: string) {
  const targets = new Set(codes.map(normalizeStockCode).filter(Boolean))
  const groups = getWatchlistGroups()
  groups.forEach((group) => {
    if (!groupId || group.id === groupId) {
      group.codes = group.codes.filter((code) => !targets.has(normalizeStockCode(code)))
      group.updatedAt = Date.now()
    }
  })
  saveWatchlistGroups(groups)
}

export function reorderWatchlist(groupId: string, orderedCodes: string[]) {
  const groups = getWatchlistGroups()
  const group = groups.find((item) => item.id === groupId)
  if (!group) return
  const existing = new Set(group.codes.map(normalizeStockCode))
  const reordered = orderedCodes.map(normalizeStockCode).filter((code) => code && existing.has(code))
  group.codes = Array.from(new Set(reordered))
  group.updatedAt = Date.now()
  saveWatchlistGroups(groups)
}

export function getAlertRules(): AlertRule[] {
  return safeJsonParse(localStorage.getItem(STORAGE_KEYS.ALERTS), [])
}

export function saveAlertRules(rules: AlertRule[]) {
  localStorage.setItem(STORAGE_KEYS.ALERTS, JSON.stringify(rules))
}

export function getAlertsByCode(code: string): AlertRule[] {
  const normalized = normalizeStockCode(code)
  return getAlertRules().filter((rule) => normalizeStockCode(rule.code) === normalized)
}

export function addAlertRule(rule: Omit<AlertRule, 'id' | 'createdAt' | 'lastTriggeredAt'>): AlertRule {
  const alerts = getAlertRules()
  const next: AlertRule = {
    ...rule,
    code: normalizeStockCode(rule.code),
    id: `alert_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`,
    createdAt: Date.now(),
    lastTriggeredAt: 0,
  }
  alerts.unshift(next)
  saveAlertRules(alerts)
  return next
}

export function updateAlertRule(id: string, updates: Partial<AlertRule>): AlertRule | null {
  const alerts = getAlertRules()
  const index = alerts.findIndex((rule) => rule.id === id)
  if (index < 0) return null
  alerts[index] = { ...alerts[index], ...updates, code: normalizeStockCode(updates.code || alerts[index].code) }
  saveAlertRules(alerts)
  return alerts[index]
}

export function deleteAlertRule(id: string) {
  saveAlertRules(getAlertRules().filter((rule) => rule.id !== id))
}

export function getSettings(): AppSettings {
  const parsed = safeJsonParse<Partial<AppSettings>>(localStorage.getItem(STORAGE_KEYS.SETTINGS), {})
  const legacyHeatmapConfig = safeJsonParse<Partial<HeatmapConfig>>(
    localStorage.getItem(STORAGE_KEYS.HEATMAP_CONFIG),
    {},
  )
  const legacyIndicatorConfig = safeJsonParse<Partial<IndicatorConfig>>(
    localStorage.getItem(STORAGE_KEYS.INDICATOR_CONFIG),
    {},
  )
  return {
    ...DEFAULT_SETTINGS,
    ...parsed,
    refreshInterval: { ...DEFAULT_SETTINGS.refreshInterval, ...parsed.refreshInterval },
    heatmapConfig: { ...DEFAULT_SETTINGS.heatmapConfig, ...legacyHeatmapConfig, ...parsed.heatmapConfig },
    indicatorConfig: {
      ...DEFAULT_SETTINGS.indicatorConfig,
      ...legacyIndicatorConfig,
      ...parsed.indicatorConfig,
      macd: {
        ...DEFAULT_SETTINGS.indicatorConfig.macd,
        ...legacyIndicatorConfig.macd,
        ...parsed.indicatorConfig?.macd,
      },
      boll: {
        ...DEFAULT_SETTINGS.indicatorConfig.boll,
        ...legacyIndicatorConfig.boll,
        ...parsed.indicatorConfig?.boll,
      },
      kdj: {
        ...DEFAULT_SETTINGS.indicatorConfig.kdj,
        ...legacyIndicatorConfig.kdj,
        ...parsed.indicatorConfig?.kdj,
      },
      dmi: {
        ...DEFAULT_SETTINGS.indicatorConfig.dmi,
        ...legacyIndicatorConfig.dmi,
        ...parsed.indicatorConfig?.dmi,
      },
      sar: {
        ...DEFAULT_SETTINGS.indicatorConfig.sar,
        ...legacyIndicatorConfig.sar,
        ...parsed.indicatorConfig?.sar,
      },
      kc: {
        ...DEFAULT_SETTINGS.indicatorConfig.kc,
        ...legacyIndicatorConfig.kc,
        ...parsed.indicatorConfig?.kc,
      },
    },
  }
}

export function saveSettings(settings: AppSettings) {
  localStorage.setItem(STORAGE_KEYS.SETTINGS, JSON.stringify(settings))
  localStorage.setItem(STORAGE_KEYS.HEATMAP_CONFIG, JSON.stringify(settings.heatmapConfig))
  localStorage.setItem(STORAGE_KEYS.INDICATOR_CONFIG, JSON.stringify(settings.indicatorConfig))
}

export function getRefreshInterval(key: keyof AppSettings['refreshInterval']): number {
  const value = getSettings().refreshInterval[key]
  return value > 0 ? value : DEFAULT_REFRESH_INTERVALS[key]
}

type AppSettingsPatch = Partial<{
  refreshInterval: Partial<AppSettings['refreshInterval']>
  colorMode: AppSettings['colorMode']
  heatmapConfig: Partial<HeatmapConfig>
  indicatorConfig: Partial<IndicatorConfig>
}>

export function updateSettings(patch: AppSettingsPatch): AppSettings {
  const current = getSettings()
  const next: AppSettings = {
    ...current,
    ...patch,
    refreshInterval: { ...current.refreshInterval, ...patch.refreshInterval },
    heatmapConfig: { ...current.heatmapConfig, ...patch.heatmapConfig },
    indicatorConfig: { ...current.indicatorConfig, ...patch.indicatorConfig },
  }
  saveSettings(next)
  return next
}

export function getHeatmapConfig(): HeatmapConfig {
  return getSettings().heatmapConfig
}

export function getIndicatorConfig(): IndicatorConfig {
  return getSettings().indicatorConfig
}

export function getTableColumns(pageKey: string): ColumnConfig[] | null {
  const configs = safeJsonParse<Record<string, ColumnConfig[]>>(localStorage.getItem(STORAGE_KEYS.TABLE_COLUMNS), {})
  return configs[pageKey] || null
}

export function saveTableColumns(pageKey: string, columns: ColumnConfig[]) {
  const configs = safeJsonParse<Record<string, ColumnConfig[]>>(localStorage.getItem(STORAGE_KEYS.TABLE_COLUMNS), {})
  configs[pageKey] = columns
  localStorage.setItem(STORAGE_KEYS.TABLE_COLUMNS, JSON.stringify(configs))
}

const MAX_SEARCH_HISTORY = 20

export function getSearchHistory(): SearchHistoryItem[] {
  return safeJsonParse(localStorage.getItem(STORAGE_KEYS.SEARCH_HISTORY), [])
}

export function addSearchHistory(item: Omit<SearchHistoryItem, 'timestamp'>) {
  const history = getSearchHistory().filter((row) => row.code !== item.code || row.market !== item.market)
  history.unshift({ ...item, timestamp: Date.now() })
  localStorage.setItem(STORAGE_KEYS.SEARCH_HISTORY, JSON.stringify(history.slice(0, MAX_SEARCH_HISTORY)))
}

export function clearSearchHistory() {
  localStorage.removeItem(STORAGE_KEYS.SEARCH_HISTORY)
}
