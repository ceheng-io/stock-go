import type { FullQuote } from '@/types'

export type DashboardRankingTab = 'rise' | 'fall' | 'amount' | 'turnover'

export interface MarketBreadthSummary {
  riseCount: number
  fallCount: number
  flatCount: number
  limitUpCount: number
  limitDownCount: number
  totalAmount: number
}

export interface DashboardSection<T> {
  name: string
  load: () => Promise<T>
  commit: (value: T) => void
}

function finiteNumber(value: number | null | undefined, fallback = 0): number {
  return typeof value === 'number' && Number.isFinite(value) ? value : fallback
}

export function summarizeMarketBreadth(quotes: FullQuote[]): MarketBreadthSummary {
  return quotes.reduce<MarketBreadthSummary>(
    (summary, quote) => {
      const changePercent = finiteNumber(quote.changePercent)
      if (changePercent > 0) summary.riseCount += 1
      else if (changePercent < 0) summary.fallCount += 1
      else summary.flatCount += 1

      if (changePercent >= 9.8) summary.limitUpCount += 1
      if (changePercent <= -9.8) summary.limitDownCount += 1
      summary.totalAmount += finiteNumber(quote.amount)
      return summary
    },
    {
      riseCount: 0,
      fallCount: 0,
      flatCount: 0,
      limitUpCount: 0,
      limitDownCount: 0,
      totalAmount: 0,
    },
  )
}

export function rankDashboardQuotes(quotes: FullQuote[], tab: DashboardRankingTab): FullQuote[] {
  const sorted = [...quotes]
  switch (tab) {
    case 'fall':
      sorted.sort((a, b) => finiteNumber(a.changePercent) - finiteNumber(b.changePercent))
      break
    case 'amount':
      sorted.sort((a, b) => finiteNumber(b.amount) - finiteNumber(a.amount))
      break
    case 'turnover':
      sorted.sort((a, b) => finiteNumber(b.turnoverRate, Number.NEGATIVE_INFINITY) - finiteNumber(a.turnoverRate, Number.NEGATIVE_INFINITY))
      break
    case 'rise':
    default:
      sorted.sort((a, b) => finiteNumber(b.changePercent) - finiteNumber(a.changePercent))
      break
  }
  return sorted
}

export function pickNorthboundSnapshot<T extends { direction?: string; boardName?: string }>(rows: T[]): T | null {
  return (
    rows.find((item) => (item.direction || '').includes('北向') || (item.boardName || '').includes('北向')) ||
    rows.find((item) => (item.direction || '').includes('沪深港通') || (item.boardName || '').includes('沪深港通')) ||
    rows[0] ||
    null
  )
}

export function getLatestMarketFundFlow<T>(rows: T[]): T | null {
  return rows.at(-1) ?? null
}

export async function loadDashboardSections(
  sections: DashboardSection<unknown>[],
  onError?: (name: string, error: unknown) => void,
): Promise<void> {
  await Promise.all(
    sections.map(async (section) => {
      try {
        section.commit(await section.load())
      } catch (error) {
        onError?.(section.name, error)
      }
    }),
  )
}

export function buildDashboardFailureMessage(failures: string[]): string {
  if (failures.length === 0) return ''
  return `${failures.join('、')}加载失败，其他区域已保留可用数据`
}
