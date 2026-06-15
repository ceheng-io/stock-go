<template>
  <div class="page board-detail">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ boardName }} 板块详情</h1>
        <div class="page-subtitle">{{ boardTypeLabel }} · {{ boardCode }}</div>
      </div>
      <a-space wrap>
        <a-button :loading="loading" @click="load">刷新</a-button>
        <a-button @click="router.back()">返回</a-button>
      </a-space>
    </div>

    <a-alert v-if="error" type="error" show-icon :message="error" />

    <a-row :gutter="[12, 12]">
      <a-col v-for="item in spotCards" :key="item.item" :xs="12" :md="8" :xl="4">
        <a-card size="small">
          <a-statistic :title="item.item" :value="formatSpotValue(item.value)" />
        </a-card>
      </a-col>
    </a-row>

    <a-row v-if="latestFundFlow" :gutter="[12, 12]">
      <a-col :xs="12" :md="6">
        <a-card size="small">
          <a-statistic title="主力净流入" :value="formatYuanAmount(latestFundFlow.mainNetInflow)" />
        </a-card>
      </a-col>
      <a-col :xs="12" :md="6">
        <a-card size="small">
          <a-statistic title="净占比" :value="formatPercent(latestFundFlow.mainNetInflowPercent)" />
        </a-card>
      </a-col>
      <a-col :xs="12" :md="6">
        <a-card size="small">
          <a-statistic title="超大单净流入" :value="formatYuanAmount(latestFundFlow.superLargeNetInflow)" />
        </a-card>
      </a-col>
      <a-col :xs="12" :md="6">
        <a-card size="small">
          <a-statistic title="小单净流入" :value="formatYuanAmount(latestFundFlow.smallNetInflow)" />
        </a-card>
      </a-col>
    </a-row>

    <a-card size="small">
      <a-tabs v-model:active-key="activeTab">
        <a-tab-pane key="trend" tab="趋势">
          <a-segmented v-model:value="period" :options="periodOptions" class="toolbar-control" />
          <KLineChart :rows="klines" empty-text="暂无板块 K 线" />
        </a-tab-pane>
        <a-tab-pane key="minute" tab="分时">
          <v-chart class="chart" :option="minuteOption" autoresize :not-merge="true" />
        </a-tab-pane>
        <a-tab-pane key="fund" tab="资金流">
          <v-chart v-if="fundFlowHistory.length > 0" class="chart" :option="fundFlowOption" autoresize :not-merge="true" />
          <a-empty v-else class="chart-empty" description="暂无板块资金流" />
        </a-tab-pane>
        <a-tab-pane key="constituents" tab="成分股">
          <a-table
            :columns="columns"
            :data-source="rows"
            row-key="code"
            size="small"
            :loading="loading"
            :pagination="{ pageSize: 30 }"
            :custom-row="rowClick"
          />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import VChart from 'vue-echarts'
import { h } from 'vue'
import { use } from 'echarts/core'
import { message } from 'ant-design-vue'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, DataZoomComponent, LegendComponent, TitleComponent } from 'echarts/components'
import KLineChart from '@/components/charts/KLineChart.vue'
import {
  getBoardKline,
  getBoardMinuteKline,
  getBoardSpot,
  getConceptList,
  getConceptConstituents,
  getSectorFundFlowHistory,
  getIndustryList,
  getIndustryConstituents,
} from '@/services/api'
import { buildEmptyChartOption, buildFundFlowOption, normalizeBoardSpotRows, type BoardSpotRow, type FundFlowRow } from '@/services/charts'
import { addToWatchlist, isInWatchlist } from '@/services/storage'
import { formatPercent, formatPrice, formatTurnover, formatVolume, formatYuanAmount } from '@/utils/format'
import type { Board } from '@/types'

use([CanvasRenderer, LineChart, BarChart, GridComponent, TooltipComponent, DataZoomComponent, LegendComponent, TitleComponent])

type BoardType = 'industry' | 'concept'
type Spot = BoardSpotRow
type Constituent = {
  code: string
  name: string
  price?: number | null
  changePercent?: number | null
  amount?: number | null
  turnoverRate?: number | null
  pe?: number | null
  pb?: number | null
}
type BoardKline = { date: string; open?: number | null; close?: number | null; low?: number | null; high?: number | null; volume?: number | null; changePercent?: number | null }
type MinuteRow = { time: string; price?: number | null; close?: number | null }
type MinuteResult = { timeline?: MinuteRow[]; klines?: MinuteRow[] }

const route = useRoute()
const router = useRouter()
const rows = ref<Constituent[]>([])
const spot = ref<Spot[]>([])
const boardTitle = ref('')
const klines = ref<BoardKline[]>([])
const minute = ref<MinuteResult | null>(null)
const fundFlowHistory = ref<FundFlowRow[]>([])
const addedCodes = ref<Set<string>>(new Set())
const loading = ref(false)
const error = ref('')
const activeTab = ref('trend')
const period = ref('daily')
const periodOptions = [
  { label: '日K', value: 'daily' },
  { label: '周K', value: 'weekly' },
  { label: '月K', value: 'monthly' },
]

const boardType = computed(() => String(route.params.type) === 'concept' ? 'concept' : 'industry')
const boardCode = computed(() => String(route.params.code || ''))
const boardTypeLabel = computed(() => boardType.value === 'industry' ? '行业' : '概念')
const boardName = computed(() => boardTitle.value || boardCode.value)
const spotCards = computed(() => spot.value.slice(0, 6))
const latestFundFlow = computed(() => fundFlowHistory.value.at(-1) || null)

const columns = [
  { title: '代码', dataIndex: 'code', width: 110 },
  { title: '名称', dataIndex: 'name', width: 120 },
  { title: '价格', customRender: ({ record }: { record: Constituent }) => formatPrice(record.price) },
  {
    title: '涨跌幅',
    customRender: ({ record }: { record: Constituent }) => formatPercent(record.changePercent),
    sorter: (a: Constituent, b: Constituent) => (a.changePercent || 0) - (b.changePercent || 0),
  },
  { title: '成交额', customRender: ({ record }: { record: Constituent }) => formatYuanAmount(record.amount) },
  { title: '换手', customRender: ({ record }: { record: Constituent }) => formatTurnover(record.turnoverRate) },
  { title: 'PE', customRender: ({ record }: { record: Constituent }) => record.pe?.toFixed(2) || '--' },
  { title: 'PB', customRender: ({ record }: { record: Constituent }) => record.pb?.toFixed(2) || '--' },
  {
    title: '操作',
    width: 108,
    customRender: ({ record }: { record: Constituent }) => h(
      'button',
      {
        class: ['link-action', isStockAdded(record.code) ? 'is-added' : ''],
        disabled: isStockAdded(record.code),
        onClick: (event: MouseEvent) => addConstituentToWatchlist(record, event),
      },
      isStockAdded(record.code) ? '已自选' : '加自选',
    ),
  },
]

const minuteOption = computed(() => {
  const rows = minute.value?.timeline?.length ? minute.value.timeline : minute.value?.klines || []
  if (rows.length === 0) return buildEmptyChartOption('暂无板块分时')
  return {
    animation: false,
    tooltip: { trigger: 'axis' },
    grid: { left: 52, right: 18, top: 24, bottom: 34 },
    xAxis: [{ type: 'category', data: rows.map((item) => item.time), boundaryGap: false }],
    yAxis: [{ type: 'value', scale: true }],
    series: [{ name: '价格', type: 'line', symbol: 'none', data: rows.map((item) => item.price ?? item.close), lineStyle: { color: '#1677ff', width: 1.5 } }],
  }
})

const fundFlowOption = computed(() => buildFundFlowOption(fundFlowHistory.value, { emptyText: '暂无板块资金流' }))

function formatSpotValue(value?: number | string | null) {
  if (value === null || value === undefined || Number.isNaN(value)) return '--'
  if (typeof value === 'string') return value
  if (Math.abs(value) >= 10000) return formatVolume(value)
  return Number(value).toFixed(2)
}

function rowClick(record: Constituent) {
  return { onClick: () => router.push(`/s/${record.code}`) }
}

function isStockAdded(code: string) {
  return addedCodes.value.has(code) || isInWatchlist(code)
}

function addConstituentToWatchlist(record: Constituent, event: MouseEvent) {
  event.stopPropagation()
  if (isStockAdded(record.code)) return
  addToWatchlist(record.code)
  addedCodes.value = new Set([...addedCodes.value, record.code])
  message.success(`已将 ${record.name || record.code} 加入自选`)
}

async function load() {
  if (!boardCode.value) return
  loading.value = true
  error.value = ''
  const type = boardType.value as BoardType
  const [boards, constituents, spotRows, klineRows, minuteRows, fundFlowRows] = await Promise.allSettled([
    type === 'industry' ? getIndustryList() as Promise<Board[]> : getConceptList() as Promise<Board[]>,
    type === 'industry' ? getIndustryConstituents(boardCode.value) as Promise<Constituent[]> : getConceptConstituents(boardCode.value) as Promise<Constituent[]>,
    getBoardSpot(type, boardCode.value) as Promise<Spot[]>,
    getBoardKline(type, boardCode.value, { period: period.value }) as Promise<BoardKline[]>,
    getBoardMinuteKline(type, boardCode.value, { period: '1' }) as Promise<MinuteResult>,
    getSectorFundFlowHistory(boardCode.value, { period: 'daily' }) as Promise<FundFlowRow[]>,
  ])

  if (boards.status === 'fulfilled') {
    boardTitle.value = boards.value.find((item) => item.code === boardCode.value)?.name || boardCode.value
  } else {
    boardTitle.value = boardCode.value
  }
  if (constituents.status === 'fulfilled') {
    rows.value = constituents.value
  } else {
    rows.value = []
  }
  if (spotRows.status === 'fulfilled') {
    spot.value = normalizeBoardSpotRows(spotRows.value)
  } else {
    spot.value = []
  }
  if (klineRows.status === 'fulfilled') {
    klines.value = klineRows.value
  } else {
    klines.value = []
  }
  if (minuteRows.status === 'fulfilled') {
    minute.value = minuteRows.value
  } else {
    minute.value = null
  }
  if (fundFlowRows.status === 'fulfilled') {
    fundFlowHistory.value = fundFlowRows.value.slice(-30)
  } else {
    fundFlowHistory.value = []
  }

  const failed = [boards, constituents, spotRows, klineRows, minuteRows, fundFlowRows].filter((item) => item.status === 'rejected')
  if (failed.length > 0) {
    error.value = `部分板块数据暂不可用，已显示 ${6 - failed.length} 项可用数据`
  }
  loading.value = false
}

watch(period, async () => {
  const type = boardType.value as BoardType
  klines.value = await getBoardKline(type, boardCode.value, { period: period.value }) as BoardKline[]
})

watch(() => route.fullPath, load)
onMounted(load)
</script>

<style scoped>
.board-detail {
  min-width: 0;
}

.page-subtitle {
  margin-top: 4px;
  color: var(--color-text-secondary);
  font-size: 13px;
}

.toolbar-control {
  margin-bottom: 10px;
}

.chart {
  width: 100%;
  height: 430px;
}

.chart-empty {
  display: flex;
  min-height: 360px;
  align-items: center;
  justify-content: center;
}

.link-action {
  border: 0;
  padding: 0;
  color: #1677ff;
  background: transparent;
  cursor: pointer;
}

.link-action.is-added,
.link-action:disabled {
  color: var(--color-text-tertiary);
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .chart {
    height: 360px;
  }
}
</style>
