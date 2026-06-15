import { describe, expect, it } from 'vitest'
import { formatPercent, normalizeStockCode, parseStockCode } from '@/utils/format'

describe('format utilities', () => {
  it('normalizes A-share symbols with exchange prefix', () => {
    expect(normalizeStockCode('600519')).toBe('sh600519')
    expect(normalizeStockCode('000001')).toBe('sz000001')
    expect(normalizeStockCode('430047')).toBe('bj430047')
    expect(normalizeStockCode('bj430047')).toBe('bj430047')
  })

  it('parses prefixed symbols', () => {
    expect(parseStockCode('sh600519')).toEqual({ market: 'sh', symbol: '600519' })
  })

  it('formats missing percent as placeholder', () => {
    expect(formatPercent(null)).toBe('--')
    expect(formatPercent(undefined)).toBe('--')
  })
})
