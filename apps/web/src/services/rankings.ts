import type { Board } from '@/types'

export type BoardRankingType = 'rise' | 'fall' | 'amount' | 'turnover'

function finiteNumber(value: number | null | undefined, fallback = 0): number {
  return typeof value === 'number' && Number.isFinite(value) ? value : fallback
}

export function sortBoardRankings(rows: Board[], rankType: BoardRankingType, limit = 50): Board[] {
  const sorted = [...rows]
  switch (rankType) {
    case 'fall':
      sorted.sort((a, b) => finiteNumber(a.changePercent) - finiteNumber(b.changePercent))
      break
    case 'amount':
      sorted.sort((a, b) => finiteNumber(b.totalMarketCap) - finiteNumber(a.totalMarketCap))
      break
    case 'turnover':
      sorted.sort((a, b) => finiteNumber(b.turnoverRate) - finiteNumber(a.turnoverRate))
      break
    case 'rise':
    default:
      sorted.sort((a, b) => finiteNumber(b.changePercent) - finiteNumber(a.changePercent))
      break
  }
  return sorted.slice(0, limit)
}
