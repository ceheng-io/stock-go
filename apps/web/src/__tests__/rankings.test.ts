import { flushPromises, mount } from '@vue/test-utils'
import { defineComponent, nextTick } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import Rankings from '@/pages/Rankings/Rankings.vue'
import { sortBoardRankings } from '@/services/rankings'
import type { Board, ZTPoolItem } from '@/types'

const push = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({ push }),
}))

vi.mock('@/services/api', () => ({
  getIndustryList: vi.fn(async (): Promise<Board[]> => [
    board({ rank: 1, code: 'BK0475', name: '酿酒行业', changePercent: 2.5 }),
  ]),
  getConceptList: vi.fn(async (): Promise<Board[]> => [
    board({ rank: 1, code: 'BK1234', name: '白酒概念', changePercent: 3.1 }),
  ]),
  getZTPool: vi.fn(async (): Promise<ZTPoolItem[]> => [
    {
      code: '002190',
      name: '成飞集成',
      changePercent: 10.01,
      price: 22.45,
      amount: 125000,
      turnoverRate: 8.32,
      continuousBoardCount: 2,
      firstBoardTime: '09:33:12',
      lastBoardTime: '14:56:00',
      industry: '军工',
      ztStatistics: '2/3',
    },
  ]),
}))

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
  beforeEach(() => {
    vi.clearAllMocks()
  })

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

  it('passes ranking row navigation through Ant Design table customRow', async () => {
    const tableStub = defineComponent({
      name: 'ATable',
      props: {
        customRow: { type: Function, default: undefined },
        dataSource: { type: Array, default: () => [] },
      },
      template: '<div data-testid="ranking-table" />',
    })

    const wrapper = mount(Rankings, {
      global: {
        stubs: {
          AAlert: true,
          AButton: true,
          ACard: { template: '<section><slot /><slot name="title" /><slot name="extra" /></section>' },
          ACol: { template: '<div><slot /></div>' },
          ARow: { template: '<div><slot /></div>' },
          ASegmented: true,
          ASpace: { template: '<div><slot /></div>' },
          AStatistic: true,
          ATable: tableStub,
          ATag: true,
        },
      },
    })
    await flushPromises()
    await nextTick()

    const customRow = wrapper.findAllComponents(tableStub)[1].props('customRow')
    expect(customRow).toEqual(expect.any(Function))

    customRow?.({ code: 'BK0475' }).onClick()
    expect(push).toHaveBeenCalledWith('/boards/industry/BK0475')
  })

  it('shows the daily limit-up pool on the rankings page', async () => {
    const cardStub = defineComponent({
      name: 'ACard',
      props: { title: { type: String, default: '' } },
      template: '<section><h2>{{ title }}</h2><slot /></section>',
    })
    const tableStub = defineComponent({
      name: 'ATable',
      props: {
        dataSource: { type: Array, default: () => [] },
      },
      template: `
        <div>
          <div v-for="row in dataSource" :key="row.code">
            {{ row.name }} {{ row.code }} {{ row.industry }} {{ row.continuousBoardCount }}连板
          </div>
        </div>
      `,
    })

    const wrapper = mount(Rankings, {
      global: {
        stubs: {
          AAlert: true,
          AButton: true,
          ACard: cardStub,
          ACol: { template: '<div><slot /></div>' },
          ARow: { template: '<div><slot /></div>' },
          ASegmented: true,
          ASpace: { template: '<div><slot /></div>' },
          AStatistic: {
            props: ['title', 'value', 'suffix'],
            template: '<div>{{ title }} {{ value }}{{ suffix }}</div>',
          },
          ATable: tableStub,
          ATag: { template: '<span><slot /></span>' },
        },
      },
    })
    await flushPromises()
    await nextTick()

    expect(wrapper.text()).toContain('当日涨停板')
    expect(wrapper.text()).toContain('成飞集成')
    expect(wrapper.text()).toContain('2连板')
    expect(wrapper.text()).toContain('军工')
    expect(wrapper.text()).toContain('涨停家数 1')
  })
})
