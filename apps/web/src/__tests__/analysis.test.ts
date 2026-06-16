import { beforeEach, describe, expect, it, vi } from 'vitest'
import { analyzeEndOfDayStocks, calculateTimelineStrength, scanSignalPool } from '@/services/analysis'
import { getAllAShareQuotes, getTodayTimeline } from '@/services/api'

vi.mock('@/services/api', () => ({
  getAllAShareQuotes: vi.fn(),
  getKlineWithIndicators: vi.fn(),
  getTodayTimeline: vi.fn(),
}))

const mockGetAllAShareQuotes = vi.mocked(getAllAShareQuotes)
const mockGetTodayTimeline = vi.mocked(getTodayTimeline)

describe('analysis service', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('calculates percent of timeline points above average price', () => {
    const result = calculateTimelineStrength({
      code: 'sh600519',
      date: '2026-06-15',
      data: [
        { time: '09:30', price: 10, avgPrice: 10, volume: 100, amount: 1000 },
        { time: '09:31', price: 11, avgPrice: 10.5, volume: 100, amount: 1100 },
        { time: '09:32', price: 9, avgPrice: 10, volume: 100, amount: 900 },
      ],
    })

    expect(result.ratio).toBeCloseTo(66.666, 2)
    expect(result.points).toHaveLength(3)
  })

  it('matches TypeScript end-of-day filtering, progress, field mapping, and candidate order', async () => {
    const progress: Array<{ completed: number; total: number; stage: string }> = []
    mockGetAllAShareQuotes.mockImplementation(async (options) => {
      options?.onProgress?.(250, 500)
      return [
        fullQuote({ code: '600001', name: '低涨幅', changePercent: 3.2, pe: 12, pb: 1.4, high: 11, low: 9, open: 9.5, prevClose: 9.7 }),
        fullQuote({ code: '600002', name: '高涨幅', changePercent: 4.8, pe: 21, pb: 2.2, high: 12, low: 10, open: 10.5, prevClose: 10.2 }),
        fullQuote({ code: '600003', name: 'ST过滤', changePercent: 4.9, nameOverride: '*ST测试' }),
      ]
    })
    mockGetTodayTimeline.mockResolvedValue({
      code: 'sh600001',
      date: '2026-06-15',
      data: [
        { time: '14:57', price: 10.1, avgPrice: 10, volume: 100, amount: 1010 },
        { time: '14:58', price: 10.2, avgPrice: 10, volume: 100, amount: 1020 },
      ],
    })

    const rows = await analyzeEndOfDayStocks(
      {
        marketCapMin: 50,
        marketCapMax: 200,
        volumeRatioMin: 1.2,
        changePercentMin: 3,
        changePercentMax: 5,
        turnoverRateMin: 5,
        turnoverRateMax: 10,
        excludeST: true,
        timelineAboveAvgRatio: 80,
      },
      {
        timelineConcurrency: 1,
        onProgress: (item) => progress.push(item),
      },
    )

    expect(mockGetAllAShareQuotes).toHaveBeenCalledWith(expect.objectContaining({
      batchSize: 500,
      concurrency: 4,
      onProgress: expect.any(Function),
    }))
    expect(progress).toContainEqual({ completed: 250, total: 500, stage: '获取行情数据' })
    expect(mockGetTodayTimeline).toHaveBeenNthCalledWith(1, 'sh600002')
    expect(mockGetTodayTimeline).toHaveBeenNthCalledWith(2, 'sh600001')
    expect(rows.map((row) => row.code)).toEqual(['600002', '600001'])
    expect(rows[0]).toMatchObject({
      code: '600002',
      routeCode: 'sh600002',
      name: '高涨幅',
      pe: 21,
      pb: 2.2,
      high: 12,
      low: 10,
      open: 10.5,
      prevClose: 10.2,
      timelineAboveAvgRatio: 100,
    })
  })

  it('returns early without timeline requests when basic end-of-day filters have no candidates', async () => {
    mockGetAllAShareQuotes.mockResolvedValue([
      fullQuote({ code: '600001', name: '低量比', volumeRatio: 0.8 }),
    ])

    const rows = await analyzeEndOfDayStocks({
      marketCapMin: 50,
      marketCapMax: 200,
      volumeRatioMin: 1.2,
      changePercentMin: 3,
      changePercentMax: 5,
      turnoverRateMin: 5,
      turnoverRateMax: 10,
      excludeST: true,
      timelineAboveAvgRatio: 80,
    })

    expect(rows).toEqual([])
    expect(mockGetTodayTimeline).not.toHaveBeenCalled()
  })

  it('matches only requested technical scanner signals from indicator kline data', async () => {
    const rows = [
      {
        date: '2026-06-11',
        close: 10,
        ma: { 5: 9.8, 10: 10.2 },
        macd: { dif: -0.2, dea: -0.1 },
        rsi: { 6: 45 },
        boll: { upper: 11, lower: 9 },
      },
      {
        date: '2026-06-12',
        close: 8.8,
        ma: { 5: 10.5, 10: 10.1 },
        macd: { dif: 0.1, dea: -0.05 },
        rsi: { 6: 24 },
        boll: { upper: 11, lower: 9 },
      },
    ]

    const results = await scanSignalPool(
      [{ code: '600519', routeCode: 'sh600519', name: '贵州茅台' }],
      ['ma_golden', 'rsi_oversold', 'boll_upper'],
      {
        loadKline: async () => rows,
      },
    )

    expect(results).toEqual([
      {
        code: '600519',
        routeCode: 'sh600519',
        name: '贵州茅台',
        matchedSignals: ['MA金叉', 'RSI超卖'],
      },
    ])
  })

  it('omits stocks when selected scanner signals are not matched', async () => {
    const results = await scanSignalPool(
      [{ code: '000001', routeCode: 'sz000001', name: '平安银行' }],
      ['ma_death', 'macd_death'],
      {
        loadKline: async () => [
          { date: '2026-06-11', close: 10, ma: { 5: 10, 10: 9 }, macd: { dif: 0.2, dea: 0.1 } },
          { date: '2026-06-12', close: 10.2, ma: { 5: 10.3, 10: 9.5 }, macd: { dif: 0.3, dea: 0.15 } },
        ],
      },
    )

    expect(results).toEqual([])
  })
})

function fullQuote(partial: {
  code: string
  name?: string
  nameOverride?: string
  changePercent?: number
  volumeRatio?: number | null
  turnoverRate?: number | null
  circulatingMarketCap?: number | null
  totalMarketCap?: number | null
  pe?: number | null
  pb?: number | null
  high?: number
  low?: number
  open?: number
  prevClose?: number
}) {
  return {
    code: partial.code,
    name: partial.nameOverride ?? partial.name ?? '测试股票',
    price: 10,
    changePercent: partial.changePercent ?? 4,
    change: 0.4,
    volume: 1000,
    amount: 100000,
    turnoverRate: partial.turnoverRate ?? 6,
    volumeRatio: partial.volumeRatio ?? 1.5,
    circulatingMarketCap: partial.circulatingMarketCap ?? 100,
    totalMarketCap: partial.totalMarketCap ?? 150,
    pe: partial.pe ?? 10,
    pb: partial.pb ?? 1.2,
    high: partial.high ?? 11,
    low: partial.low ?? 9,
    open: partial.open ?? 9.8,
    prevClose: partial.prevClose ?? 9.6,
  }
}
