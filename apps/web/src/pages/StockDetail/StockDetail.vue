<template>
  <div class="page stock-detail">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ quote?.name || code }} 个股详情</h1>
        <div class="page-subtitle">{{ code }} · {{ quote?.time || '实时行情' }}</div>
      </div>
      <a-space wrap>
        <a-button :loading="loading" @click="load">刷新</a-button>
        <a-button :type="inWatchlist ? 'default' : 'primary'" @click="toggleWatchlist">
          {{ inWatchlist ? '移出自选' : '加入自选' }}
        </a-button>
        <a-button @click="router.back()">返回</a-button>
      </a-space>
    </div>

    <a-alert v-if="error" type="error" show-icon :message="error" />

    <a-row :gutter="[12, 12]">
      <a-col :xs="24" :lg="8">
        <a-card :loading="loading" size="small">
          <template v-if="quote">
            <div class="quote-card">
              <div>
                <div class="quote-price">{{ formatPrice(quote.price) }}</div>
                <div :class="['quote-change', getChangeColorClass(quote.changePercent)]">
                  {{ formatChange(quote.change) }} · {{ formatPercent(quote.changePercent) }}
                </div>
              </div>
              <a-tag :color="quote.changePercent > 0 ? 'red' : quote.changePercent < 0 ? 'green' : 'default'">
                {{ quote.source || quote.assetType || 'A 股' }}
              </a-tag>
            </div>
            <a-descriptions :column="2" size="small" class="quote-metrics">
              <a-descriptions-item label="今开">{{ formatPrice(quote.open) }}</a-descriptions-item>
              <a-descriptions-item label="昨收">{{ formatPrice(quote.prevClose) }}</a-descriptions-item>
              <a-descriptions-item label="最高">{{ formatPrice(quote.high) }}</a-descriptions-item>
              <a-descriptions-item label="最低">{{ formatPrice(quote.low) }}</a-descriptions-item>
              <a-descriptions-item label="成交量">{{ formatVolume(quote.volume) }}</a-descriptions-item>
              <a-descriptions-item label="成交额">{{ formatAmount(quote.amount) }}</a-descriptions-item>
              <a-descriptions-item label="换手">{{ formatTurnover(quote.turnoverRate) }}</a-descriptions-item>
              <a-descriptions-item label="量比">{{ formatVolumeRatio(quote.volumeRatio) }}</a-descriptions-item>
              <a-descriptions-item label="PE">{{ formatRatio(quote.pe) }}</a-descriptions-item>
              <a-descriptions-item label="PB">{{ formatRatio(quote.pb) }}</a-descriptions-item>
            </a-descriptions>
          </template>
          <a-empty v-else-if="!loading" description="暂无行情" />
        </a-card>
      </a-col>

      <a-col :xs="24" :lg="8">
        <a-card title="资金与大单" :loading="sideLoading" size="small">
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="主力净流入">{{ formatYuanAmount(fundFlow?.mainNet) }}</a-descriptions-item>
            <a-descriptions-item label="主力净占比">{{ formatPercent(fundFlow?.mainNetRatio) }}</a-descriptions-item>
            <a-descriptions-item label="散户净流入">{{ formatYuanAmount(fundFlow?.retailNet) }}</a-descriptions-item>
            <a-descriptions-item label="大单买入占比">{{ formatPercent(panelOrder?.buyLargeRatio, false) }}</a-descriptions-item>
            <a-descriptions-item label="大单卖出占比">{{ formatPercent(panelOrder?.sellLargeRatio, false) }}</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>

      <a-col :xs="24" :lg="8">
        <a-card title="本地预警" size="small">
          <a-form layout="inline" class="alert-form" @finish="createAlert">
            <a-form-item>
              <a-select v-model:value="alertType" style="width: 132px">
                <a-select-option value="price_gte">价格 >=</a-select-option>
                <a-select-option value="price_lte">价格 <=</a-select-option>
                <a-select-option value="change_percent_gte">涨幅 >=</a-select-option>
                <a-select-option value="change_percent_lte">涨幅 <=</a-select-option>
                <a-select-option value="amount_gte">成交额 >=</a-select-option>
              </a-select>
            </a-form-item>
            <a-form-item>
              <a-input-number v-model:value="alertValue" :precision="2" style="width: 112px" />
            </a-form-item>
            <a-form-item>
              <a-button html-type="submit">添加</a-button>
            </a-form-item>
          </a-form>
          <a-list :data-source="alerts" size="small">
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta :title="alertLabel(item.type)" :description="String(item.value)" />
                <a-button danger size="small" type="link" @click="removeAlert(item.id)">删除</a-button>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </a-col>
    </a-row>

    <a-card size="small">
      <a-tabs v-model:active-key="activeTab">
        <a-tab-pane key="timeline" tab="分时">
          <div class="chart-toolbar">
            <a-segmented v-model:value="minutePeriod" :options="minutePeriodOptions" />
          </div>
          <v-chart class="chart" :option="timelineOption" autoresize :not-merge="true" />
        </a-tab-pane>
        <a-tab-pane key="kline" tab="K 线">
          <div class="chart-toolbar">
            <a-segmented v-model:value="klinePeriod" :options="klinePeriodOptions" />
            <a-checkbox-group v-model:value="selectedOverlays" :options="overlayOptions" />
            <a-segmented v-model:value="selectedOscillator" :options="oscillatorOptions" />
          </div>
          <v-chart class="chart" :option="klineOption" autoresize :not-merge="true" />
        </a-tab-pane>
        <a-tab-pane key="fund" tab="历史资金">
          <a-table
            :columns="fundColumns"
            :data-source="fundHistory"
            row-key="date"
            size="small"
            :pagination="{ pageSize: 10 }"
          />
        </a-tab-pane>
        <a-tab-pane key="northbound" tab="北向持仓">
          <a-table
            :columns="northboundColumns"
            :data-source="northboundHistory"
            row-key="date"
            size="small"
            :pagination="{ pageSize: 10 }"
          />
        </a-tab-pane>
        <a-tab-pane key="dividend" tab="分红">
          <a-table
            :columns="dividendColumns"
            :data-source="dividends"
            :row-key="dividendRowKey"
            size="small"
            :pagination="{ pageSize: 10 }"
          />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { CandlestickChart, LineChart, BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, DataZoomComponent, LegendComponent, TitleComponent } from 'echarts/components'
import type { AlertRule, AlertType, FullQuote } from '@/types'
import { usePolling } from '@/composables/usePolling'
import {
  getDividendDetail,
  getFullQuotes,
  getIndividualFundFlow,
  getKlineWithIndicators,
  getMinuteKline,
  getNorthboundIndividual,
  getPanelLargeOrder,
  getQuoteFundFlow,
  getTodayTimeline,
} from '@/services/api'
import {
  buildIndicatorKlineOption,
  buildMinuteChartOption,
  type IndicatorKlineRow,
  type MinuteKlineRow,
  type OscillatorIndicatorKey,
  type OverlayIndicatorKey,
  type TimelineChartPayload,
} from '@/services/charts'
import {
  addAlertRule,
  addToWatchlist,
  deleteAlertRule,
  getAlertsByCode,
  isInWatchlist,
  getIndicatorConfig,
  getSettings,
  removeFromWatchlist,
} from '@/services/storage'
import {
  formatChange,
  formatAmount,
  formatPercent,
  formatPrice,
  formatRatio,
  formatTurnover,
  formatVolume,
  formatVolumeRatio,
  formatYuanAmount,
  getChangeColorClass,
  normalizeStockCode,
} from '@/utils/format'

use([CanvasRenderer, CandlestickChart, LineChart, BarChart, GridComponent, TooltipComponent, DataZoomComponent, LegendComponent, TitleComponent])

type NumberRecord = Record<string, number | string | null | undefined>
type KlineRow = NumberRecord & IndicatorKlineRow
type TimelineResponse = TimelineChartPayload
type FundFlow = { mainNet?: number | null; mainNetRatio?: number | null; retailNet?: number | null }
type PanelOrder = { buyLargeRatio?: number | null; sellLargeRatio?: number | null }
type FundHistoryRow = { date: string; close?: number | null; changePercent?: number | null; mainNetInflow?: number | null; mainNetInflowPercent?: number | null }
type NorthboundRow = { date: string; holdShares?: number | null; holdMarketValue?: number | null; holdRatioFloat?: number | null; close?: number | null; changePercent?: number | null }
type DividendRow = { reportDate?: string | null; noticeDate?: string | null; dividendDesc?: string | null; dividendYield?: number | null; exDividendDate?: string | null; payDate?: string | null }

const route = useRoute()
const router = useRouter()
const code = computed(() => normalizeStockCode(String(route.params.code || '')))
const quote = ref<FullQuote | null>(null)
const timeline = ref<TimelineResponse | null>(null)
const minuteKlines = ref<MinuteKlineRow[]>([])
const klines = ref<KlineRow[]>([])
const fundFlow = ref<FundFlow | null>(null)
const panelOrder = ref<PanelOrder | null>(null)
const fundHistory = ref<FundHistoryRow[]>([])
const northboundHistory = ref<NorthboundRow[]>([])
const dividends = ref<DividendRow[]>([])
const alerts = ref<AlertRule[]>([])
const loading = ref(false)
const sideLoading = ref(false)
const error = ref('')
const activeTab = ref('timeline')
const minutePeriod = ref('1')
const klinePeriod = ref('daily')
const selectedOverlays = ref<OverlayIndicatorKey[]>(['ma'])
const selectedOscillator = ref<OscillatorIndicatorKey>('macd')
const inWatchlist = ref(false)
const alertType = ref<AlertType>('price_gte')
const alertValue = ref<number | null>(null)
const detailRefreshInterval = getSettings().refreshInterval.detail
const fundRefreshInterval = Math.max(detailRefreshInterval * 6, 30000)
const klinePeriodOptions = [
  { label: '日K', value: 'daily' },
  { label: '周K', value: 'weekly' },
  { label: '月K', value: 'monthly' },
]
const minutePeriodOptions = [
  { label: '分时', value: '1' },
  { label: '5分', value: '5' },
  { label: '15分', value: '15' },
  { label: '30分', value: '30' },
  { label: '60分', value: '60' },
]
const overlayOptions: Array<{ label: string; value: OverlayIndicatorKey }> = [
  { label: 'MA', value: 'ma' },
  { label: 'BOLL', value: 'boll' },
  { label: 'SAR', value: 'sar' },
  { label: 'KC', value: 'kc' },
]
const oscillatorOptions: Array<{ label: string; value: OscillatorIndicatorKey }> = [
  { label: 'MACD', value: 'macd' },
  { label: 'KDJ', value: 'kdj' },
  { label: 'RSI', value: 'rsi' },
  { label: 'OBV', value: 'obv' },
  { label: 'ROC', value: 'roc' },
  { label: 'DMI', value: 'dmi' },
]
const indicatorConfig = getIndicatorConfig()

const fundColumns = [
  { title: '日期', dataIndex: 'date' },
  { title: '收盘', customRender: ({ record }: { record: FundHistoryRow }) => formatPrice(record.close) },
  { title: '涨跌幅', customRender: ({ record }: { record: FundHistoryRow }) => formatPercent(record.changePercent) },
  { title: '主力净流入', customRender: ({ record }: { record: FundHistoryRow }) => formatYuanAmount(record.mainNetInflow) },
  { title: '主力净占比', customRender: ({ record }: { record: FundHistoryRow }) => formatPercent(record.mainNetInflowPercent) },
]

const northboundColumns = [
  { title: '日期', dataIndex: 'date' },
  { title: '收盘', customRender: ({ record }: { record: NorthboundRow }) => formatPrice(record.close) },
  { title: '涨跌幅', customRender: ({ record }: { record: NorthboundRow }) => formatPercent(record.changePercent) },
  { title: '持股数', customRender: ({ record }: { record: NorthboundRow }) => formatVolume(record.holdShares) },
  { title: '持股市值', customRender: ({ record }: { record: NorthboundRow }) => formatYuanAmount(record.holdMarketValue) },
  { title: '流通占比', customRender: ({ record }: { record: NorthboundRow }) => formatPercent(record.holdRatioFloat) },
]

const dividendColumns = [
  { title: '报告期', dataIndex: 'reportDate' },
  { title: '公告日', dataIndex: 'noticeDate' },
  { title: '方案', dataIndex: 'dividendDesc' },
  { title: '股息率', customRender: ({ record }: { record: DividendRow }) => formatPercent(record.dividendYield, false) },
  { title: '除权除息日', dataIndex: 'exDividendDate' },
  { title: '派息日', dataIndex: 'payDate' },
]

const timelineOption = computed(() => {
  return buildMinuteChartOption({
    period: minutePeriod.value,
    timeline: timeline.value,
    minuteKline: minuteKlines.value,
    emptyText: minutePeriod.value === '1' ? '暂无分时数据' : '暂无分钟 K 数据',
  })
})

const klineOption = computed(() => {
  return buildIndicatorKlineOption(klines.value, {
    emptyText: '暂无 K 线数据',
    overlays: selectedOverlays.value,
    oscillator: selectedOscillator.value,
    indicatorConfig,
  })
})

async function load() {
  if (!code.value) return
  loading.value = true
  sideLoading.value = true
  error.value = ''
  try {
    const [quoteRows, timelineRow, klineRows] = await Promise.all([
      getFullQuotes([code.value]) as Promise<FullQuote[]>,
      getTodayTimeline(code.value) as Promise<TimelineResponse>,
      getKlineWithIndicators(code.value, { period: klinePeriod.value, adjust: 'qfq' }) as Promise<KlineRow[]>,
    ])
    quote.value = quoteRows[0] || null
    timeline.value = timelineRow
    klines.value = klineRows
    if (minutePeriod.value !== '1') {
      minuteKlines.value = await getMinuteKline(code.value, { period: minutePeriod.value }) as MinuteKlineRow[]
    } else {
      minuteKlines.value = []
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : '加载个股数据失败'
  } finally {
    loading.value = false
  }

  try {
    const [fundRows, panelRows, fundHistoryRows, northboundRows, dividendRows] = await Promise.allSettled([
      getQuoteFundFlow([code.value]) as Promise<FundFlow[]>,
      getPanelLargeOrder([code.value]) as Promise<PanelOrder[]>,
      getIndividualFundFlow(code.value, { period: 'daily' }) as Promise<FundHistoryRow[]>,
      getNorthboundIndividual(code.value) as Promise<NorthboundRow[]>,
      getDividendDetail(code.value) as Promise<DividendRow[]>,
    ])
    fundFlow.value = fundRows.status === 'fulfilled' ? fundRows.value[0] || null : null
    panelOrder.value = panelRows.status === 'fulfilled' ? panelRows.value[0] || null : null
    fundHistory.value = fundHistoryRows.status === 'fulfilled' ? fundHistoryRows.value : []
    northboundHistory.value = northboundRows.status === 'fulfilled' ? northboundRows.value : []
    dividends.value = dividendRows.status === 'fulfilled' ? dividendRows.value : []
  } finally {
    sideLoading.value = false
    refreshLocalState()
  }
}

async function refreshRealtimeData() {
  if (!code.value) return
  try {
    const [quoteRows, timelineRow, klineRows] = await Promise.all([
      getFullQuotes([code.value]) as Promise<FullQuote[]>,
      minutePeriod.value === '1'
        ? getTodayTimeline(code.value) as Promise<TimelineResponse>
        : getMinuteKline(code.value, { period: minutePeriod.value }) as Promise<MinuteKlineRow[]>,
      getKlineWithIndicators(code.value, { period: klinePeriod.value, adjust: 'qfq' }) as Promise<KlineRow[]>,
    ])
    quote.value = quoteRows[0] || null
    if (minutePeriod.value === '1') {
      timeline.value = timelineRow as TimelineResponse
      minuteKlines.value = []
    } else {
      timeline.value = null
      minuteKlines.value = timelineRow as MinuteKlineRow[]
    }
    klines.value = klineRows
  } catch (err) {
    console.warn('Stock detail realtime refresh failed', err)
  }
}

async function refreshSideData() {
  if (!code.value) return
  try {
    const [fundRows, panelRows, fundHistoryRows, northboundRows, dividendRows] = await Promise.allSettled([
      getQuoteFundFlow([code.value]) as Promise<FundFlow[]>,
      getPanelLargeOrder([code.value]) as Promise<PanelOrder[]>,
      getIndividualFundFlow(code.value, { period: 'daily' }) as Promise<FundHistoryRow[]>,
      getNorthboundIndividual(code.value) as Promise<NorthboundRow[]>,
      getDividendDetail(code.value) as Promise<DividendRow[]>,
    ])
    fundFlow.value = fundRows.status === 'fulfilled' ? fundRows.value[0] || null : fundFlow.value
    panelOrder.value = panelRows.status === 'fulfilled' ? panelRows.value[0] || null : panelOrder.value
    fundHistory.value = fundHistoryRows.status === 'fulfilled' ? fundHistoryRows.value : fundHistory.value
    northboundHistory.value = northboundRows.status === 'fulfilled' ? northboundRows.value : northboundHistory.value
    dividends.value = dividendRows.status === 'fulfilled' ? dividendRows.value : dividends.value
    refreshLocalState()
  } catch (err) {
    console.warn('Stock detail side refresh failed', err)
  }
}

function refreshLocalState() {
  inWatchlist.value = isInWatchlist(code.value)
  alerts.value = getAlertsByCode(code.value)
}

function toggleWatchlist() {
  if (inWatchlist.value) {
    removeFromWatchlist(code.value)
    message.success('已移出自选')
  } else {
    addToWatchlist(code.value)
    message.success('已加入自选')
  }
  refreshLocalState()
}

function createAlert() {
  if (alertValue.value === null) {
    message.warning('请输入预警值')
    return
  }
  addAlertRule({
    code: code.value,
    name: quote.value?.name || code.value,
    type: alertType.value,
    value: alertValue.value,
    cooldownSec: 300,
    enabled: true,
  })
  alertValue.value = null
  refreshLocalState()
}

function removeAlert(id: string) {
  deleteAlertRule(id)
  refreshLocalState()
}

function dividendRowKey(record: DividendRow) {
  return `${record.reportDate || ''}-${record.noticeDate || ''}-${record.exDividendDate || ''}`
}

function alertLabel(type: AlertType) {
  const labels: Record<AlertType, string> = {
    price_gte: '价格大于等于',
    price_lte: '价格小于等于',
    change_percent_gte: '涨幅大于等于',
    change_percent_lte: '涨幅小于等于',
    amount_gte: '成交额大于等于',
    near_limit_up: '接近涨停',
    near_limit_down: '接近跌停',
  }
  return labels[type] || type
}

watch(klinePeriod, async () => {
  if (!code.value) return
  klines.value = await getKlineWithIndicators(code.value, { period: klinePeriod.value, adjust: 'qfq' }) as KlineRow[]
})

watch(minutePeriod, async () => {
  if (!code.value) return
  if (minutePeriod.value === '1') {
    minuteKlines.value = []
    timeline.value = await getTodayTimeline(code.value) as TimelineResponse
    return
  }
  minuteKlines.value = await getMinuteKline(code.value, { period: minutePeriod.value }) as MinuteKlineRow[]
})

watch(code, load)
onMounted(load)
usePolling(refreshRealtimeData, {
  interval: detailRefreshInterval,
  enabled: computed(() => Boolean(code.value) && !loading.value),
  immediate: false,
})
usePolling(refreshSideData, {
  interval: fundRefreshInterval,
  enabled: computed(() => Boolean(code.value) && !sideLoading.value),
  immediate: false,
})
</script>

<style scoped>
.stock-detail {
  min-width: 0;
}

.page-subtitle {
  margin-top: 4px;
  color: var(--color-text-secondary);
  font-size: 13px;
}

.quote-card {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.quote-price {
  font-size: 34px;
  line-height: 1;
  font-weight: 700;
}

.quote-change {
  margin-top: 8px;
  font-size: 15px;
}

.quote-metrics {
  margin-top: 16px;
}

.alert-form {
  row-gap: 8px;
  margin-bottom: 8px;
}

.chart-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px 14px;
  margin-bottom: 10px;
}

.chart {
  width: 100%;
  height: 430px;
}

@media (max-width: 768px) {
  .chart {
    height: 360px;
  }
}
</style>
