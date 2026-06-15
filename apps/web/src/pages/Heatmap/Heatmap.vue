<template>
  <div class="page heatmap-page">
    <div class="page-header">
      <div>
        <h1 class="page-title">热力图</h1>
        <div class="page-subtitle">{{ dimensionLabel }} · {{ treemapData.length }} 项</div>
      </div>
      <a-button :loading="loading" @click="load">刷新</a-button>
    </div>

    <a-alert v-if="error" type="warning" show-icon :message="error" />

    <a-card size="small" class="control-card">
      <a-space wrap>
        <a-segmented v-model:value="config.dimension" :options="dimensionOptions" />
        <a-select v-model:value="config.colorField" size="small" style="width: 112px">
          <a-select-option value="changePercent">涨跌幅</a-select-option>
          <a-select-option value="turnoverRate">换手率</a-select-option>
          <a-select-option value="volumeRatio">量比</a-select-option>
        </a-select>
        <a-select v-model:value="config.sizeField" size="small" style="width: 112px">
          <a-select-option value="totalMarketCap">总市值</a-select-option>
          <a-select-option value="amount">成交额</a-select-option>
          <a-select-option value="volume">成交量</a-select-option>
        </a-select>
        <a-segmented v-model:value="config.topK" :options="topKOptions" />
        <a-segmented v-model:value="config.colorMode" :options="colorModeOptions" />
      </a-space>
    </a-card>

    <a-card size="small" class="chart-card" :loading="loading">
      <v-chart v-if="treemapData.length > 0" class="chart" :option="chartOption" autoresize @click="handleChartClick" />
      <a-empty v-else class="empty" description="暂无热力图数据" />
    </a-card>

    <div class="legend">
      <span>{{ config.colorMode === 'red-rise' ? '跌' : '涨' }}</span>
      <div class="legend-bar" />
      <span>{{ config.colorMode === 'red-rise' ? '涨' : '跌' }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { TreemapChart } from 'echarts/charts'
import { TooltipComponent } from 'echarts/components'
import type { ECElementEvent } from 'echarts/core'
import type { ECBasicOption } from 'echarts/types/dist/shared'
import { usePolling } from '@/composables/usePolling'
import { getAllQuotesByCodes, getConceptList, getIndustryList } from '@/services/api'
import { getAllWatchlistCodes, getHeatmapConfig, getSettings, updateSettings } from '@/services/storage'
import type { Board, FullQuote, HeatmapConfig } from '@/types'
import { formatAmount, formatMarketCap, formatPercent, formatVolume } from '@/utils/format'

use([CanvasRenderer, TreemapChart, TooltipComponent])

type HeatmapItem = {
  name: string
  code: string
  value: number
  changePercent?: number | null
  turnoverRate?: number | null
  volumeRatio?: number | null
  totalMarketCap?: number | null
  amount?: number | null
  volume?: number | null
  price?: number | null
  riseCount?: number | null
  fallCount?: number | null
  leadingStock?: string | null
  leadingStockChangePercent?: number | null
  itemStyle: { color: string }
}

const router = useRouter()
const config = reactive<HeatmapConfig>({ ...getHeatmapConfig() })
const loading = ref(false)
const error = ref('')
const boards = ref<Board[]>([])
const stockQuotes = ref<FullQuote[]>([])
const heatmapRefreshInterval = getSettings().refreshInterval.heatmap
const dimensionOptions = [
  { label: '行业', value: 'industry' },
  { label: '概念', value: 'concept' },
  { label: '自选', value: 'watchlist' },
]
const topKOptions = [
  { label: 'Top 50', value: 50 },
  { label: 'Top 100', value: 100 },
  { label: 'Top 200', value: 200 },
]
const colorModeOptions = [
  { label: '红涨绿跌', value: 'red-rise' },
  { label: '绿涨红跌', value: 'green-rise' },
]

const dimensionLabel = computed(() => dimensionOptions.find((item) => item.value === config.dimension)?.label || '热力')

const treemapData = computed<HeatmapItem[]>(() => {
  const source = config.dimension === 'watchlist' ? stockQuotes.value : boards.value
  return source.slice(0, config.topK).map((item) => {
    const colorValue = getColorValue(item)
    const value = getSizeValue(item)
    return {
      name: item.name || item.code,
      code: item.code,
      value,
      changePercent: item.changePercent,
      turnoverRate: item.turnoverRate,
      volumeRatio: 'volumeRatio' in item ? item.volumeRatio : undefined,
      totalMarketCap: item.totalMarketCap,
      amount: 'amount' in item ? item.amount : undefined,
      volume: 'volume' in item ? item.volume : undefined,
      price: 'price' in item ? item.price : undefined,
      riseCount: 'riseCount' in item ? item.riseCount : undefined,
      fallCount: 'fallCount' in item ? item.fallCount : undefined,
      leadingStock: 'leadingStock' in item ? item.leadingStock : undefined,
      leadingStockChangePercent: 'leadingStockChangePercent' in item ? item.leadingStockChangePercent : undefined,
      itemStyle: { color: getColor(colorValue, config.colorField) },
    }
  })
})

const chartOption = computed<ECBasicOption>(() => ({
  tooltip: {
    trigger: 'item',
    formatter: (params: { data?: HeatmapItem }) => formatTooltip(params.data),
  },
  series: [
    {
      type: 'treemap',
      roam: false,
      nodeClick: false,
      breadcrumb: { show: false },
      left: 0,
      right: 0,
      top: 0,
      bottom: 0,
      itemStyle: {
        borderColor: '#ffffff',
        borderWidth: 1,
        gapWidth: 1,
      },
      label: {
        show: true,
        formatter: (params: { data?: HeatmapItem }) => {
          const item = params.data
          if (!item) return ''
          return `${item.name}\n${formatPercent(item.changePercent)}`
        },
        color: '#fff',
        fontSize: 12,
        textShadowBlur: 2,
        textShadowColor: 'rgba(0,0,0,.45)',
      },
      emphasis: {
        itemStyle: { borderColor: '#1677ff', borderWidth: 2 },
      },
      data: treemapData.value,
    },
  ],
}))

function getColorValue(item: Board | FullQuote) {
  if (config.colorField === 'turnoverRate') return item.turnoverRate ?? 0
  if (config.colorField === 'volumeRatio') return 'volumeRatio' in item ? item.volumeRatio ?? 1 : 1
  return item.changePercent ?? 0
}

function getSizeValue(item: Board | FullQuote) {
  if (config.sizeField === 'amount') return 'amount' in item ? Math.max(item.amount || 1, 1) : 1
  if (config.sizeField === 'volume') return 'volume' in item ? Math.max(item.volume || 1, 1) : 1
  return Math.max(item.totalMarketCap || 1, 1)
}

function getColor(value: number, field: HeatmapConfig['colorField']) {
  const riseRed = config.colorMode === 'red-rise'
  if (field === 'changePercent') {
    if (!value) return '#8c8c8c'
    const alpha = 0.35 + Math.min(Math.abs(value) / 10, 1) * 0.6
    const red = `rgba(207, 19, 34, ${alpha})`
    const green = `rgba(56, 158, 13, ${alpha})`
    return value > 0 ? (riseRed ? red : green) : (riseRed ? green : red)
  }
  const alpha = 0.3 + Math.min(value / (field === 'turnoverRate' ? 20 : 5), 1) * 0.65
  return `rgba(22, 119, 255, ${alpha})`
}

function formatTooltip(item?: HeatmapItem) {
  if (!item) return ''
  const lines = [
    `<strong>${item.name}</strong>`,
    `涨跌幅：${formatPercent(item.changePercent)}`,
    `换手率：${formatPercent(item.turnoverRate, false)}`,
  ]
  if (item.price !== undefined) lines.push(`现价：${item.price?.toFixed(2) ?? '--'}`)
  if (item.amount !== undefined) lines.push(`成交额：${formatAmount(item.amount)}`)
  if (item.volume !== undefined) lines.push(`成交量：${formatVolume(item.volume)}`)
  if (item.totalMarketCap !== undefined) lines.push(`总市值：${formatMarketCap(item.totalMarketCap)}`)
  if (item.leadingStock) lines.push(`领涨：${item.leadingStock} ${formatPercent(item.leadingStockChangePercent)}`)
  if (item.riseCount !== undefined) lines.push(`涨跌家数：${item.riseCount ?? 0} / ${item.fallCount ?? 0}`)
  return lines.join('<br/>')
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    if (config.dimension === 'industry') {
      boards.value = await getIndustryList()
      stockQuotes.value = []
    } else if (config.dimension === 'concept') {
      boards.value = await getConceptList()
      stockQuotes.value = []
    } else {
      const codes = getAllWatchlistCodes().slice(0, config.topK)
      boards.value = []
      stockQuotes.value = codes.length > 0 ? await getAllQuotesByCodes(codes) as FullQuote[] : []
    }
  } catch (err) {
    console.warn('Heatmap data loading failed', err)
    error.value = '热力图数据加载失败，请稍后刷新'
    boards.value = []
    stockQuotes.value = []
  } finally {
    loading.value = false
  }
}

function persistConfig() {
  updateSettings({ heatmapConfig: { ...config } })
}

function handleChartClick(params: ECElementEvent) {
  const item = params.data as HeatmapItem | undefined
  if (!item?.code) return
  if (config.dimension === 'industry') {
    router.push(`/boards/industry/${item.code}`)
  } else if (config.dimension === 'concept') {
    router.push(`/boards/concept/${item.code}`)
  } else {
    router.push(`/s/${item.code}`)
  }
}

watch(config, () => {
  persistConfig()
  load()
})

onMounted(load)
usePolling(load, { interval: heatmapRefreshInterval, enabled: computed(() => !loading.value), immediate: false })
</script>

<style scoped>
.heatmap-page {
  min-width: 0;
}

.page-subtitle {
  margin-top: 4px;
  color: var(--color-text-secondary);
  font-size: 13px;
}

.control-card :deep(.ant-card-body) {
  padding: 10px 12px;
}

.chart-card {
  min-height: 560px;
}

.chart-card :deep(.ant-card-body) {
  height: min(66vh, 680px);
  min-height: 520px;
  padding: 0;
}

.chart {
  width: 100%;
  height: 100%;
}

.empty {
  padding-top: 160px;
}

.legend {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: var(--color-text-secondary);
  font-size: 12px;
}

.legend-bar {
  width: 220px;
  height: 10px;
  border-radius: 999px;
  background: linear-gradient(90deg, #389e0d, #8c8c8c, #cf1322);
}

@media (max-width: 768px) {
  .chart-card :deep(.ant-card-body) {
    height: 62vh;
    min-height: 420px;
  }
}
</style>
