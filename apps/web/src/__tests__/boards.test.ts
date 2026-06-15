import { mount } from '@vue/test-utils'
import { defineComponent, nextTick } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import Boards from '@/pages/Boards/Boards.vue'
import BoardDetail from '@/pages/Boards/BoardDetail.vue'
import type { Board } from '@/types'

const push = vi.fn()
const back = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({ push, back }),
  useRoute: () => ({ params: { type: 'industry', code: 'BK0475' }, fullPath: '/boards/industry/BK0475' }),
}))

vi.mock('@/services/api', () => ({
  getBoardKline: vi.fn(async () => []),
  getBoardMinuteKline: vi.fn(async () => ({ timeline: [] })),
  getBoardSpot: vi.fn(async () => []),
  getConceptList: vi.fn(async () => []),
  getConceptConstituents: vi.fn(async () => []),
  getSectorFundFlowHistory: vi.fn(async () => []),
  getIndustryList: vi.fn(async (): Promise<Board[]> => [
    { rank: 1, code: 'BK0475', name: '酿酒行业', changePercent: 2.5 },
  ]),
  getIndustryConstituents: vi.fn(async () => [
    { code: 'sh600519', name: '贵州茅台', price: 1271.1, changePercent: 1.2 },
  ]),
}))

const tableStub = defineComponent({
  name: 'ATable',
  props: {
    customRow: { type: Function, default: undefined },
    dataSource: { type: Array, default: () => [] },
  },
  template: '<div data-testid="table" />',
})

function mountPage(component: typeof Boards | typeof BoardDetail) {
  return mount(component, {
    global: {
      stubs: {
        AAlert: true,
        AButton: true,
        ACard: { template: '<section><slot /><slot name="title" /><slot name="extra" /></section>' },
        ACol: { template: '<div><slot /></div>' },
        AEmpty: true,
        AInputSearch: true,
        ARow: { template: '<div><slot /></div>' },
        ASegmented: true,
        ASpace: { template: '<div><slot /></div>' },
        AStatistic: true,
        ATable: tableStub,
        ATabs: { template: '<div><slot /></div>' },
        ATabPane: { template: '<div><slot /></div>' },
        Echarts: true,
        KLineChart: true,
        VChart: true,
      },
    },
  })
}

describe('Boards pages', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('passes board row navigation through Ant Design table customRow', async () => {
    const wrapper = mountPage(Boards)
    await nextTick()
    await nextTick()

    const customRow = wrapper.getComponent(tableStub).props('customRow')
    expect(customRow).toEqual(expect.any(Function))

    customRow?.({ code: 'BK0475' }).onClick()
    expect(push).toHaveBeenCalledWith('/boards/industry/BK0475')
  })

  it('passes board constituent row navigation through Ant Design table customRow', async () => {
    const wrapper = mountPage(BoardDetail)
    await nextTick()
    await nextTick()

    const customRow = wrapper.getComponent(tableStub).props('customRow')
    expect(customRow).toEqual(expect.any(Function))

    customRow?.({ code: 'sh600519' }).onClick()
    expect(push).toHaveBeenCalledWith('/s/sh600519')
  })
})
