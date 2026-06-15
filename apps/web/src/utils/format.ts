export function formatNumber(value: number | null | undefined, decimals = 2): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '--'
  return value.toLocaleString('zh-CN', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  })
}

export function formatPrice(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '--'
  return value.toFixed(2)
}

export function formatPercent(value: number | null | undefined, showSign = true): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '--'
  const sign = showSign && value > 0 ? '+' : ''
  return `${sign}${value.toFixed(2)}%`
}

export function formatChange(value: number | null | undefined, showSign = true): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '--'
  const sign = showSign && value > 0 ? '+' : ''
  return `${sign}${value.toFixed(2)}`
}

export function formatVolume(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '--'
  if (value >= 100000000) return `${(value / 100000000).toFixed(2)}亿`
  if (value >= 10000) return `${(value / 10000).toFixed(2)}万`
  return value.toFixed(0)
}

export function formatAmount(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '--'
  if (value >= 10000) return `${(value / 10000).toFixed(2)}亿`
  if (value >= 1) return `${value.toFixed(2)}万`
  return `${(value * 10000).toFixed(0)}元`
}

export function formatYuanAmount(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '--'
  const absValue = Math.abs(value)
  const sign = value < 0 ? '-' : ''
  if (absValue >= 1000000000000) return `${sign}${(absValue / 1000000000000).toFixed(2)}万亿`
  if (absValue >= 100000000) return `${sign}${(absValue / 100000000).toFixed(2)}亿`
  if (absValue >= 10000) return `${sign}${(absValue / 10000).toFixed(2)}万`
  return `${sign}${absValue.toFixed(0)}元`
}

export function formatCompactNumber(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '--'
  const absValue = Math.abs(value)
  const sign = value < 0 ? '-' : ''
  if (absValue >= 100000000) return `${sign}${(absValue / 100000000).toFixed(2)}亿`
  if (absValue >= 10000) return `${sign}${(absValue / 10000).toFixed(2)}万`
  return `${sign}${absValue.toFixed(0)}`
}

export function formatMarketCap(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '--'
  if (value >= 10000) return `${(value / 10000).toFixed(2)}万亿`
  return `${value.toFixed(2)}亿`
}

export function formatTurnover(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '--'
  return `${value.toFixed(2)}%`
}

export function formatVolumeRatio(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '--'
  return value.toFixed(2)
}

export function formatRatio(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '--'
  if (value < 0) return '亏损'
  return value.toFixed(2)
}

export function formatTime(time: string | undefined): string {
  if (!time) return '--'
  if (time.length === 14) return `${time.slice(8, 10)}:${time.slice(10, 12)}:${time.slice(12, 14)}`
  return time
}

export function formatDate(date: string | undefined): string {
  if (!date) return '--'
  if (date.includes('-')) return date
  if (date.length === 8) return `${date.slice(0, 4)}-${date.slice(4, 6)}-${date.slice(6, 8)}`
  return date
}

export function getChangeColorClass(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value) || value === 0) return 'text-flat'
  return value > 0 ? 'text-rise' : 'text-fall'
}

export function getChangeColor(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value) || value === 0) return 'var(--color-flat)'
  return value > 0 ? 'var(--color-rise)' : 'var(--color-fall)'
}

export function parseStockCode(code: string): { market: string; symbol: string } {
  const trimmed = code.trim()
  if (!trimmed) return { market: '', symbol: '' }

  const prefixMatch = trimmed.match(/^(sh|sz|bj)\.?(\d{6})$/i)
  if (prefixMatch) return { market: prefixMatch[1].toLowerCase(), symbol: prefixMatch[2] }

  const suffixMatch = trimmed.match(/^(\d{6})\.(sh|sz|bj)$/i)
  if (suffixMatch) return { market: suffixMatch[2].toLowerCase(), symbol: suffixMatch[1] }

  if (/^\d{6}$/.test(trimmed)) {
    if (trimmed.startsWith('6')) return { market: 'sh', symbol: trimmed }
    if (trimmed.startsWith('0') || trimmed.startsWith('3')) return { market: 'sz', symbol: trimmed }
    if (trimmed.startsWith('4') || trimmed.startsWith('8')) return { market: 'bj', symbol: trimmed }
  }

  return { market: '', symbol: trimmed }
}

export function normalizeStockCode(code: string): string {
  const trimmed = code.trim()
  if (!trimmed) return ''
  const { market, symbol } = parseStockCode(trimmed)
  return market ? `${market}${symbol}` : trimmed
}
