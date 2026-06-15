import { describe, expect, it } from 'vitest'
import { sortBoardRankings } from '@/services/rankings'
import type { Board } from '@/types'

function board(partial: Partial<Board>): Board {
  return {
    rank: partial.rank ?? 0,
    name: partial.name || '测试板块',
    code: partial.code || 'BK0001',
    changePercent: partial.changePercent ?? null,
    totalMarketCap: partial.totalMarketCap ?? null,
    turnoverRate: partial.turnoverRate ?? null,
  }
}

describe('ranking helpers', () => {
  it('sorts boards by rise fall amount and turnover', () => {
    const rows = [
      board({ code: 'BK1', changePercent: 3, totalMarketCap: 50, turnoverRate: 1 }),
      board({ code: 'BK2', changePercent: -2, totalMarketCap: 100, turnoverRate: 8 }),
      board({ code: 'BK3', changePercent: 1, totalMarketCap: null, turnoverRate: 3 }),
    ]

    expect(sortBoardRankings(rows, 'rise').map((item) => item.code)).toEqual(['BK1', 'BK3', 'BK2'])
    expect(sortBoardRankings(rows, 'fall').map((item) => item.code)).toEqual(['BK2', 'BK3', 'BK1'])
    expect(sortBoardRankings(rows, 'amount').map((item) => item.code)).toEqual(['BK2', 'BK1', 'BK3'])
    expect(sortBoardRankings(rows, 'turnover').map((item) => item.code)).toEqual(['BK2', 'BK3', 'BK1'])
  })

  it('limits board ranking rows when a limit is provided', () => {
    const rows = Array.from({ length: 60 }, (_, index) => board({ code: `BK${index}`, changePercent: index }))

    expect(sortBoardRankings(rows, 'rise', 50)).toHaveLength(50)
    expect(sortBoardRankings(rows, 'rise', 1)[0].code).toBe('BK59')
  })
})
