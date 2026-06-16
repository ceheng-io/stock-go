import { afterEach, describe, expect, it, vi } from 'vitest'
import { apiRequest, getBoardSpot, getFullQuotes, getKlineWithIndicators, getTHSLimitUpPool, getZTPool } from '@/services/api'

describe('api client', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('joins quote codes in the query string', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => [{ Code: 'sh600519', Name: '贵州茅台' }],
    })
    vi.stubGlobal('fetch', fetchMock)

    const result = await getFullQuotes(['sh600519', 'sz000001'])

    expect(result).toEqual([{ code: 'sh600519', name: '贵州茅台' }])
    expect(fetchMock).toHaveBeenCalledWith('/api/quotes/full?codes=sh600519%2Csz000001', {
      headers: { Accept: 'application/json' },
    })
  })

  it('keeps the requested exchange prefix when full quote payloads return bare codes', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => [
          { Code: '000001', Name: '上证指数' },
          { Code: '000001', Name: '平安银行' },
        ],
      }),
    )

    await expect(getFullQuotes(['sh000001', 'sz000001'])).resolves.toEqual([
      { code: 'sh000001', name: '上证指数' },
      { code: 'sz000001', name: '平安银行' },
    ])
  })

  it('normalizes Go json field names to camel case', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => [
          {
            Code: 'sh600519',
            ChangePercent: 1.23,
            TotalMarketCap: 25000,
            PEStatic: 28.5,
            TZ: 'Asia/Shanghai',
            Bid: [{ Price: 1688, Volume: 100 }],
          },
        ],
      }),
    )

    await expect(getFullQuotes(['sh600519'])).resolves.toEqual([
      {
        code: 'sh600519',
        changePercent: 1.23,
        totalMarketCap: 25000,
        peStatic: 28.5,
        tz: 'Asia/Shanghai',
        bid: [{ price: 1688, volume: 100 }],
      },
    ])
  })

  it('throws structured api errors', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: false,
        status: 400,
        json: async () => ({ error: { message: 'keyword is required' } }),
      }),
    )

    await expect(apiRequest('/search')).rejects.toMatchObject({
      message: 'keyword is required',
      status: 400,
    })
  })

  it('deduplicates concurrent matching K-line requests', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => [{ Date: '2026-06-12', Close: 1500 }],
    })
    vi.stubGlobal('fetch', fetchMock)

    const [left, right] = await Promise.all([
      getKlineWithIndicators('sh600519', { period: 'daily', adjust: 'qfq' }),
      getKlineWithIndicators('sh600519', { period: 'daily', adjust: 'qfq' }),
    ])

    expect(left).toEqual([{ date: '2026-06-12', close: 1500 }])
    expect(right).toEqual(left)
    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(fetchMock).toHaveBeenCalledWith('/api/kline/indicators?symbol=sh600519&period=daily&adjust=qfq', {
      headers: { Accept: 'application/json' },
    })
  })

  it('releases failed K-line requests so users can retry', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce({
        ok: false,
        status: 502,
        statusText: 'Bad Gateway',
        json: async () => ({ error: { message: 'eastmoney EOF' } }),
      })
      .mockResolvedValueOnce({
        ok: true,
        json: async () => [{ Date: '2026-06-15', Close: 1510 }],
      })
    vi.stubGlobal('fetch', fetchMock)

    await expect(Promise.all([
      getKlineWithIndicators('sh600519', { period: 'daily', adjust: 'qfq' }),
      getKlineWithIndicators('sh600519', { period: 'daily', adjust: 'qfq' }),
    ])).rejects.toMatchObject({ message: 'eastmoney EOF', status: 502 })

    await expect(getKlineWithIndicators('sh600519', { period: 'daily', adjust: 'qfq' })).resolves.toEqual([
      { date: '2026-06-15', close: 1510 },
    ])
    expect(fetchMock).toHaveBeenCalledTimes(2)
  })

  it('normalizes object-like board spot payloads into metric rows', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => ({
          Price: 1234.56,
          ChangePercent: 1.23,
          TotalMarketCap: 987654321,
          Name: '酿酒行业',
        }),
      }),
    )

    await expect(getBoardSpot('industry', 'BK0475')).resolves.toEqual([
      { item: 'price', value: 1234.56 },
      { item: 'changePercent', value: 1.23 },
      { item: 'totalMarketCap', value: 987654321 },
    ])
  })

  it('normalizes wrapped limit-up pool payloads into rows', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => ({
          Items: [
            {
              Code: '002190',
              Name: '成飞集成',
              Latest: 22.45,
              ChangeRate: 10.01,
              FirstLimitUpTimeText: '09:33:12',
              LastLimitUpTimeText: '14:56:00',
              OrderAmount: 420000000,
              TurnoverRate: 8.32,
              HighDays: '2天2板',
              HighDaysValue: 2,
              ReasonType: '低空经济',
              LimitUpType: '换手板',
            },
          ],
        }),
      }),
    )

    await expect(getZTPool()).resolves.toEqual([
      expect.objectContaining({
        code: '002190',
        name: '成飞集成',
        price: 22.45,
        changePercent: 10.01,
        firstBoardTime: '09:33:12',
        lastBoardTime: '14:56:00',
        sealAmount: 420000000,
        turnoverRate: 8.32,
        continuousBoardCount: 2,
        industry: '',
        ztStatistics: '2天2板',
        limitUpType: '换手板',
        reasonType: '低空经济',
      }),
    ])
  })

  it('normalizes Tonghuashun limit-up pool payloads into rows with reasons and board counts', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        Items: [
          {
            Code: '002190',
            Name: '成飞集成',
            Latest: 22.45,
            ChangeRate: 10.01,
            FirstLimitUpTimeText: '09:33:12',
            LastLimitUpTimeText: '14:56:00',
            OrderAmount: 420000000,
            TurnoverRate: 8.32,
            HighDays: '首板',
            HighDaysValue: 65537,
            ReasonType: '低空经济',
            LimitUpType: '换手板',
          },
        ],
      }),
    })
    vi.stubGlobal('fetch', fetchMock)

    await expect(getTHSLimitUpPool({ date: '2026-06-15', limit: 20 })).resolves.toEqual([
      expect.objectContaining({
        code: '002190',
        name: '成飞集成',
        industry: '',
        reasonType: '低空经济',
        continuousBoardCount: 1,
        ztStatistics: '首板',
        limitUpType: '换手板',
      }),
    ])
    expect(fetchMock).toHaveBeenCalledWith('/api/market-event/ths-limit-up-pool?date=2026-06-15&limit=20', {
      headers: { Accept: 'application/json' },
    })
  })
})
