import { flushPromises, mount } from '@vue/test-utils'
import { defineComponent, nextTick } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import Rankings from '@/pages/Rankings/Rankings.vue'
import { sortBoardRankings } from '@/services/rankings'
import type { Board, ZTPoolItem } from '@/types'

const push = vi.fn()
const apiMocks = vi.hoisted(() => ({
  getIndustryList: vi.fn(),
  getConceptList: vi.fn(),
  getTHSLimitUpPool: vi.fn(),
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({ push }),
}))

vi.mock('@/services/api', () => ({
  getIndustryList: apiMocks.getIndustryList,
  getConceptList: apiMocks.getConceptList,
  getTHSLimitUpPool: apiMocks.getTHSLimitUpPool,
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

function limitUpItem(partial: Partial<ZTPoolItem> = {}): ZTPoolItem {
  return {
    code: partial.code || '002190',
    name: partial.name || '成飞集成',
    changePercent: partial.changePercent ?? 10.01,
    price: partial.price ?? 22.45,
    amount: partial.amount ?? 125000,
    turnoverRate: partial.turnoverRate ?? 8.32,
    continuousBoardCount: partial.continuousBoardCount ?? 2,
    firstBoardTime: partial.firstBoardTime || '09:33:12',
    lastBoardTime: partial.lastBoardTime || '14:56:00',
    industry: partial.industry || '军工',
    ztStatistics: partial.ztStatistics || '2/3',
    reasonType: partial.reasonType || null,
    limitUpType: partial.limitUpType || null,
  }
}

function deferred<T>() {
  let resolve!: (value: T) => void
  let reject!: (reason?: unknown) => void
  const promise = new Promise<T>((innerResolve, innerReject) => {
    resolve = innerResolve
    reject = innerReject
  })
  return { promise, resolve, reject }
}

describe('ranking helpers', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    apiMocks.getIndustryList.mockResolvedValue([
      board({ rank: 1, code: 'BK0475', name: '酿酒行业', changePercent: 2.5 }),
    ])
    apiMocks.getConceptList.mockResolvedValue([
      board({ rank: 1, code: 'BK1234', name: '白酒概念', changePercent: 3.1 }),
    ])
    apiMocks.getTHSLimitUpPool.mockResolvedValue([limitUpItem()])
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
    const segmentedStub = defineComponent({
      name: 'ASegmented',
      props: {
        value: { type: String, default: '' },
        options: { type: Array, default: () => [] },
      },
      emits: ['update:value'],
      template: '<button v-for="option in options" :key="option.value" @click="$emit(\'update:value\', option.value)">{{ option.label }}</button>',
    })
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
          ASegmented: segmentedStub,
          ASpace: { template: '<div><slot /></div>' },
          AStatistic: true,
          ATable: tableStub,
          ATag: true,
        },
      },
    })
    await flushPromises()
    await nextTick()

    wrapper.findAllComponents(segmentedStub)[0].vm.$emit('update:value', 'industry')
    await nextTick()

    const customRow = wrapper.findAllComponents(tableStub)[0].props('customRow')
    expect(customRow).toEqual(expect.any(Function))

    customRow?.({ code: 'BK0475' }).onClick()
    expect(push).toHaveBeenCalledWith('/boards/industry/BK0475')
  })

  it('shows the daily limit-up pool in its own ranking tab', async () => {
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
          ASpace: { template: '<div><slot /></div>' },
          ASegmented: {
            props: ['options'],
            template: '<div><button v-for="option in options" :key="option.value">{{ option.label }}</button></div>',
          },
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
    expect(wrapper.text()).toContain('行业板块')
    expect(wrapper.text()).toContain('概念板块')
    expect(wrapper.text()).not.toContain('酿酒行业')
    expect(wrapper.text()).not.toContain('白酒概念')
  })

  it('shows the limit-up reason in the daily limit-up table', async () => {
    apiMocks.getTHSLimitUpPool.mockResolvedValue([
      limitUpItem({ reasonType: '低空经济', industry: '军工' }),
    ])

    const tableStub = defineComponent({
      name: 'ATable',
      props: {
        columns: { type: Array, default: () => [] },
        dataSource: { type: Array, default: () => [] },
      },
      template: `
        <table>
          <tbody>
            <tr v-for="(row, index) in dataSource" :key="row.code">
              <td v-for="column in columns" :key="column.key">
                <slot name="bodyCell" :column="column" :record="row" :index="index" />
              </td>
            </tr>
          </tbody>
        </table>
      `,
    })

    const wrapper = mount(Rankings, {
      global: {
        stubs: {
          AAlert: true,
          AButton: true,
          ACard: {
            props: { title: { type: String, default: '' } },
            template: '<section><h2>{{ title }}</h2><slot /></section>',
          },
          ASpace: { template: '<div><slot /></div>' },
          ASegmented: true,
          AStatistic: true,
          ATable: tableStub,
          ATag: { template: '<span><slot /></span>' },
        },
      },
    })
    await flushPromises()
    await nextTick()

    expect(wrapper.text()).toContain('低空经济')
  })

  it('keeps successful limit-up rows visible when board rankings fail', async () => {
    apiMocks.getIndustryList.mockRejectedValue(new Error('industry failed'))
    apiMocks.getConceptList.mockRejectedValue(new Error('concept failed'))

    const wrapper = mount(Rankings, {
      global: {
        stubs: {
          AAlert: {
            props: ['message'],
            template: '<div>{{ message }}</div>',
          },
          AButton: true,
          ACard: {
            props: { title: { type: String, default: '' } },
            template: '<section><h2>{{ title }}</h2><slot /></section>',
          },
          ACol: { template: '<div><slot /></div>' },
          ARow: { template: '<div><slot /></div>' },
          ASegmented: true,
          ASpace: { template: '<div><slot /></div>' },
          AStatistic: {
            props: ['title', 'value', 'suffix'],
            template: '<div>{{ title }} {{ value }}{{ suffix }}</div>',
          },
          ATable: {
            props: {
              dataSource: { type: Array, default: () => [] },
            },
            template: '<div><div v-for="row in dataSource" :key="row.code">{{ row.name }}</div></div>',
          },
          ATag: { template: '<span><slot /></span>' },
        },
      },
    })
    await flushPromises()
    await nextTick()

    expect(wrapper.text()).toContain('成飞集成')
    expect(wrapper.text()).toContain('榜单数据部分加载失败，请稍后刷新')
  })

  it('does not keep the limit-up table loading while board rankings are still pending', async () => {
    const industryRequest = deferred<Board[]>()
    const conceptRequest = deferred<Board[]>()
    apiMocks.getIndustryList.mockReturnValue(industryRequest.promise)
    apiMocks.getConceptList.mockReturnValue(conceptRequest.promise)
    apiMocks.getTHSLimitUpPool.mockResolvedValue([limitUpItem({ code: '000001', name: '平安银行' })])

    const tableStub = defineComponent({
      name: 'ATable',
      props: {
        dataSource: { type: Array, default: () => [] },
        loading: { type: Boolean, default: false },
      },
      template: '<div />',
    })

    const wrapper = mount(Rankings, {
      global: {
        stubs: {
          AAlert: true,
          AButton: true,
          ACard: {
            props: { title: { type: String, default: '' } },
            template: '<section><h2>{{ title }}</h2><slot /></section>',
          },
          ASpace: { template: '<div><slot /></div>' },
          ASegmented: true,
          AStatistic: true,
          ATable: tableStub,
          ATag: true,
        },
      },
    })
    await flushPromises()
    await nextTick()

    const table = wrapper.findComponent(tableStub)
    expect(table.props('dataSource')).toHaveLength(1)
    expect(table.props('loading')).toBe(false)

    industryRequest.resolve([])
    conceptRequest.resolve([])
    await flushPromises()
    await nextTick()
  })

  it('adds sortable columns to the limit-up table fields', async () => {
    const rows = [
      limitUpItem({ code: '000001', name: '平安银行', amount: 200, continuousBoardCount: 1, firstBoardTime: '09:45:00', industry: '银行', reasonType: '银行' }),
      limitUpItem({ code: '000002', name: '万科A', amount: 300, continuousBoardCount: 3, firstBoardTime: '09:31:00', industry: '地产', reasonType: '房地产' }),
      limitUpItem({ code: '000003', name: '中信证券', amount: 100, continuousBoardCount: 2, firstBoardTime: '10:02:00', industry: '证券', reasonType: '券商' }),
    ]
    apiMocks.getTHSLimitUpPool.mockResolvedValue(rows)

    const tableStub = defineComponent({
      name: 'ATable',
      props: {
        columns: { type: Array, default: () => [] },
      },
      template: '<div />',
    })

    const wrapper = mount(Rankings, {
      global: {
        stubs: {
          AAlert: true,
          AButton: true,
          ACard: {
            props: { title: { type: String, default: '' } },
            template: '<section><h2>{{ title }}</h2><slot /></section>',
          },
          ASpace: { template: '<div><slot /></div>' },
          ASegmented: true,
          AStatistic: true,
          ATable: tableStub,
          ATag: true,
        },
      },
    })
    await flushPromises()
    await nextTick()

    const columns = wrapper.findComponent(tableStub).props('columns') as Array<{ key: string; sorter?: (left: ZTPoolItem, right: ZTPoolItem) => number }>
    const amountSorter = columns.find((column) => column.key === 'amount')?.sorter
    const boardCountSorter = columns.find((column) => column.key === 'continuousBoardCount')?.sorter
    const boardTimeSorter = columns.find((column) => column.key === 'boardTime')?.sorter
    const industrySorter = columns.find((column) => column.key === 'industry')?.sorter
    const reasonSorter = columns.find((column) => column.key === 'reasonType')?.sorter

    expect(amountSorter).toEqual(expect.any(Function))
    expect(boardCountSorter).toEqual(expect.any(Function))
    expect(boardTimeSorter).toEqual(expect.any(Function))
    expect(industrySorter).toEqual(expect.any(Function))
    expect(reasonSorter).toEqual(expect.any(Function))
    expect([...rows].sort(amountSorter).map((item) => item.code)).toEqual(['000003', '000001', '000002'])
    expect([...rows].sort(boardCountSorter).map((item) => item.code)).toEqual(['000001', '000003', '000002'])
    expect([...rows].sort(boardTimeSorter).map((item) => item.code)).toEqual(['000002', '000001', '000003'])
    expect([...rows].sort(industrySorter).map((item) => item.code)).toEqual(['000002', '000001', '000003'])
    expect([...rows].sort(reasonSorter).map((item) => item.code)).toEqual(['000002', '000003', '000001'])
  })
})
