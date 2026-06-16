<template>
  <div class="page eod-page">
    <header class="page-header eod-header">
      <div class="title-section">
        <h1 class="page-title eod-title">
          <TrendingUp :size="24" />
          尾盘选股法
        </h1>
        <p class="page-subtitle">一日持股法分析工具 · 筛选分时强势股票</p>
      </div>
    </header>

    <main class="eod-main">
      <section v-if="!hasAnalyzed" class="start-screen">
        <div class="start-button-container">
          <button class="start-button" data-testid="eod-start" :disabled="isLoading" @click="start">
            <span class="button-glow" />
            <span class="button-content">
              <Zap :size="28" />
              <span>开始分析</span>
            </span>
          </button>
        </div>

        <section class="filter-card" :class="{ editing: isEditing }" @click="activateEditing">
          <div class="filter-header">
            <div class="filter-title">
              <SlidersHorizontal :size="20" />
              <span>筛选条件</span>
            </div>
            <div class="filter-actions">
              <template v-if="isEditing">
                <button class="reset-btn" title="恢复默认" @click.stop="resetFilters">
                  <RotateCcw :size="14" />
                  默认
                </button>
                <button class="scheme-btn" title="管理方案" @click.stop="toggleSchemePanel">
                  <FolderOpen :size="14" />
                  方案
                </button>
                <button class="recent-btn" title="最近使用" @click.stop="toggleRecentPanel">
                  <Clock :size="14" />
                  历史
                </button>
                <button class="save-btn" @click.stop="finishEditing">
                  <Check :size="14" />
                  完成
                </button>
              </template>
              <span v-else class="edit-hint">点击编辑</span>
            </div>
          </div>

          <section v-if="showSchemePanel" class="scheme-panel" @click.stop>
            <div class="scheme-panel-header">
              <span>保存的方案</span>
              <button class="add-scheme-btn" @click="showSaveInput = !showSaveInput">
                <Plus :size="14" />
                保存当前
              </button>
            </div>
            <div v-if="showSaveInput" class="save-input-row">
              <input
                v-model="schemeName"
                class="scheme-name-input"
                type="text"
                placeholder="输入方案名称"
                @keydown.enter="saveScheme"
              >
              <button class="confirm-save-btn" @click="saveScheme">
                <Save :size="14" />
              </button>
            </div>
            <div v-if="schemes.length > 0" class="scheme-list">
              <div v-for="scheme in schemes" :key="scheme.id" class="scheme-item">
                <button class="scheme-load-btn" @click="loadScheme(scheme)">{{ scheme.name }}</button>
                <button class="scheme-delete-btn" @click="removeScheme(scheme.id)">
                  <Trash2 :size="12" />
                </button>
              </div>
            </div>
            <p v-else class="empty-hint">暂无保存的方案</p>
          </section>

          <section v-if="showRecentPanel" class="recent-panel" @click.stop>
            <div class="scheme-panel-header">
              <span>最近使用</span>
            </div>
            <div v-if="recentUsage.length > 0" class="scheme-list">
              <button v-for="item in recentUsage" :key="item.usedAt" class="recent-item" @click="loadRecent(item)">
                <span class="recent-summary">{{ recentSummary(item.filters) }}</span>
                <span class="recent-time">{{ formatDateTime(item.usedAt) }}</span>
              </button>
            </div>
            <p v-else class="empty-hint">暂无使用记录</p>
          </section>

          <div class="filter-grid">
            <div class="filter-item">
              <span class="filter-label">流通市值</span>
              <div class="filter-value">
                <template v-if="isEditing">
                  <input v-model.number="filters.marketCapMin" class="filter-input" type="number" @click.stop>
                  <span class="filter-separator">~</span>
                  <input v-model.number="filters.marketCapMax" class="filter-input" type="number" @click.stop>
                  <span class="filter-unit">亿</span>
                </template>
                <span v-else class="filter-display">{{ filters.marketCapMin }} ~ {{ filters.marketCapMax }}<span class="filter-unit">亿</span></span>
              </div>
            </div>

            <div class="filter-item">
              <span class="filter-label">量比</span>
              <div class="filter-value">
                <input v-if="isEditing" v-model.number="filters.volumeRatioMin" class="filter-input" type="number" step="0.1" @click.stop>
                <span v-else class="filter-display">≥ {{ filters.volumeRatioMin }}</span>
              </div>
            </div>

            <div class="filter-item">
              <span class="filter-label">当日涨幅</span>
              <div class="filter-value">
                <template v-if="isEditing">
                  <input v-model.number="filters.changePercentMin" class="filter-input" type="number" step="0.5" @click.stop>
                  <span class="filter-separator">~</span>
                  <input v-model.number="filters.changePercentMax" class="filter-input" type="number" step="0.5" @click.stop>
                  <span class="filter-unit">%</span>
                </template>
                <span v-else class="filter-display">{{ filters.changePercentMin }} ~ {{ filters.changePercentMax }}<span class="filter-unit">%</span></span>
              </div>
            </div>

            <div class="filter-item">
              <span class="filter-label">换手率</span>
              <div class="filter-value">
                <template v-if="isEditing">
                  <input v-model.number="filters.turnoverRateMin" class="filter-input" type="number" step="0.5" @click.stop>
                  <span class="filter-separator">~</span>
                  <input v-model.number="filters.turnoverRateMax" class="filter-input" type="number" step="0.5" @click.stop>
                  <span class="filter-unit">%</span>
                </template>
                <span v-else class="filter-display">{{ filters.turnoverRateMin }} ~ {{ filters.turnoverRateMax }}<span class="filter-unit">%</span></span>
              </div>
            </div>

            <div class="filter-item">
              <span class="filter-label">过滤ST股票</span>
              <div class="filter-value">
                <button v-if="isEditing" class="toggle-btn" :class="{ active: filters.excludeST }" @click.stop="filters.excludeST = !filters.excludeST">
                  <component :is="filters.excludeST ? ToggleRight : ToggleLeft" :size="22" />
                  <span>{{ filters.excludeST ? '开启' : '关闭' }}</span>
                </button>
                <span v-else class="filter-display">{{ filters.excludeST ? '开启' : '关闭' }}</span>
              </div>
            </div>

            <div class="filter-item">
              <span class="filter-label">分时强度</span>
              <div class="filter-value">
                <template v-if="isEditing">
                  <input v-model.number="filters.timelineAboveAvgRatio" class="filter-input" type="number" min="0" max="100" step="5" @click.stop>
                  <span class="filter-unit">%</span>
                </template>
                <span v-else class="filter-display">≥ {{ filters.timelineAboveAvgRatio }}<span class="filter-unit">%</span></span>
              </div>
            </div>
          </div>
        </section>
      </section>

      <section v-else class="results-screen">
        <div class="results-header">
          <div class="results-summary">
            <Target :size="20" />
            <span class="summary-text">共筛选出 <strong>{{ stocks.length }}</strong> 只符合条件的股票</span>
          </div>
          <div class="results-actions">
            <button v-if="stocks.length > 0" class="select-mode-btn" data-testid="eod-select-mode" :class="{ active: showSelectMode }" @click="toggleSelectMode">
              <component :is="showSelectMode ? X : CheckSquare" :size="16" />
              {{ showSelectMode ? '取消' : '批量选' }}
            </button>
            <button class="back-button" data-testid="eod-reset" @click="returnToStart">
              <ChevronLeft :size="18" />
              重新筛选
            </button>
          </div>
        </div>

        <div v-if="stocks.length > 0" class="sort-bar">
          <div class="sort-section">
            <ArrowUpDown :size="14" />
            <span class="sort-label">排序：</span>
            <button
              v-for="option in sortOptions"
              :key="option.value"
              class="sort-option"
              :class="{ active: sortField === option.value }"
              @click="handleSortChange(option.value)"
            >
              {{ option.label }}
              <component v-if="sortField === option.value" :is="sortOrder === 'desc' ? ArrowDown : ArrowUp" :size="12" />
            </button>
          </div>
          <div v-if="showSelectMode" class="batch-section">
            <button class="select-all-btn" @click="toggleSelectAll">
              {{ selectedCodes.length === sortedStocks.length ? '取消全选' : '全选' }}
            </button>
            <button class="batch-add-btn" data-testid="eod-batch-add" :disabled="selectedCodes.length === 0" @click="batchAdd">
              <Plus :size="14" />
              加入自选 ({{ selectedCodes.length }})
            </button>
          </div>
        </div>

        <div v-if="sortedStocks.length > 0" class="stock-grid">
          <article
            v-for="(stock, index) in sortedStocks"
            :key="stock.code"
            class="stock-card"
            :class="{ positive: stock.changePercent >= 0, negative: stock.changePercent < 0, selected: selectedCodes.includes(stock.code) }"
            :style="{ '--card-delay': `${index * 45}ms` }"
            @click="router.push(`/s/${stock.routeCode}`)"
          >
            <div class="stock-header">
              <div class="stock-info">
                <div class="stock-name-row">
                  <button v-if="showSelectMode" class="select-btn" @click.stop="toggleRow(stock.code)">
                    <component :is="selectedCodes.includes(stock.code) ? CheckSquare : Square" :size="16" />
                  </button>
                  <h3 class="stock-name">{{ stock.name }}</h3>
                </div>
                <span class="stock-code">{{ stock.code }}</span>
              </div>
              <div class="change-badge">
                <span class="change-icon">{{ stock.changePercent >= 0 ? '▲' : '▼' }}</span>
                <span class="change-percent">{{ formatPercent(stock.changePercent, false) }}</span>
              </div>
            </div>

            <div class="price-section">
              <div class="current-price">
                <span class="price-label">现价</span>
                <span class="price-value">{{ formatPrice(stock.price) }}</span>
              </div>
              <div class="price-change">
                <span class="change-value">{{ formatChange(stock.change) }}</span>
                <span v-if="stock.timelineAboveAvgRatio !== undefined" class="timeline-ratio">强度 {{ stock.timelineAboveAvgRatio.toFixed(0) }}%</span>
              </div>
            </div>

            <div v-if="stock.timeline?.length" class="timeline-section">
              <TimelineChart :points="stock.timeline" :prev-close="stock.prevClose" />
            </div>

            <div class="data-grid">
              <div class="data-item">
                <span class="data-label">流通市值</span>
                <span class="data-value">{{ formatMarketCap(stock.circulatingMarketCap) }}</span>
              </div>
              <div class="data-item">
                <span class="data-label">量比</span>
                <span class="data-value">{{ formatVolumeRatio(stock.volumeRatio) }}</span>
              </div>
              <div class="data-item">
                <span class="data-label">换手率</span>
                <span class="data-value">{{ formatTurnover(stock.turnoverRate) }}</span>
              </div>
              <div class="data-item">
                <span class="data-label">成交额</span>
                <span class="data-value">{{ formatYuanAmount(stock.amount) }}</span>
              </div>
            </div>

            <button class="add-watchlist-btn" :class="{ added: isInWatchlist(stock.routeCode) }" :disabled="isInWatchlist(stock.routeCode)" @click.stop="add(stock.routeCode, stock.name)">
              <component :is="isInWatchlist(stock.routeCode) ? Check : Plus" :size="14" />
              {{ isInWatchlist(stock.routeCode) ? '已自选' : '加自选' }}
            </button>
          </article>
        </div>

        <div v-else class="no-results">
          <SearchX :size="64" :stroke-width="1" />
          <p class="no-results-title">没有找到符合条件的股票</p>
          <p class="no-results-hint">请尝试调整筛选条件后重新分析</p>
        </div>
      </section>
    </main>

    <div v-if="isLoading" class="loading-overlay">
      <div class="loading-content">
        <div class="loading-spinner">
          <div class="spinner-ring" />
          <div class="spinner-center">
            <span class="loading-percentage">{{ percent }}%</span>
          </div>
        </div>
        <div class="loading-progress">
          <div class="loading-progress-fill" :style="{ width: `${percent}%` }" />
        </div>
        <div class="loading-status">
          <p class="loading-text">{{ progress.stage || '正在扫描全市场股票数据...' }}</p>
          <p class="loading-detail">{{ progress.total > 0 ? `已处理 ${progress.completed} / ${progress.total}` : '正在初始化连接...' }}</p>
        </div>
        <button class="cancel-analysis-btn" @click="cancel">
          <CircleStop :size="16" />
          取消分析
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, onMounted, reactive, ref, watch, type PropType } from 'vue'
import {
  ArrowDown,
  ArrowUp,
  ArrowUpDown,
  Check,
  CheckSquare,
  ChevronLeft,
  CircleStop,
  Clock,
  FolderOpen,
  Plus,
  RotateCcw,
  Save,
  SearchX,
  SlidersHorizontal,
  Square,
  Target,
  ToggleLeft,
  ToggleRight,
  Trash2,
  TrendingUp,
  X,
  Zap,
} from '@lucide/vue'
import { message } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import {
  type AnalysisProgress,
  type EndOfDayFilters,
  type EndOfDayStock,
  type TimelinePoint,
  analyzeEndOfDayStocks,
  isAnalysisAborted,
} from '@/services/analysis'
import {
  DEFAULT_END_OF_DAY_FILTERS,
  addEndOfDayRecentUsage,
  deleteEndOfDayScheme,
  getEndOfDayFilters,
  getEndOfDayRecentUsage,
  getEndOfDaySchemes,
  saveEndOfDayFilters,
  saveEndOfDayScheme,
  sortEndOfDayStocks,
  toggleSelectedCode,
  type EndOfDayRecentUsage,
  type EndOfDayScheme,
  type EndOfDaySortField,
  type EndOfDaySortOrder,
} from '@/services/endOfDayPicker'
import { addToWatchlist, isInWatchlist } from '@/services/storage'
import {
  formatChange,
  formatMarketCap,
  formatPercent,
  formatPrice,
  formatTurnover,
  formatVolumeRatio,
  formatYuanAmount,
} from '@/utils/format'

const TimelineChart = defineComponent({
  name: 'TimelineChart',
  props: {
    points: { type: Array as PropType<TimelinePoint[]>, required: true },
    prevClose: { type: Number, required: true },
  },
  setup(props) {
    const chart = computed(() => {
      if (props.points.length === 0) return null
      const width = 320
      const height = 100
      const padding = { top: 6, right: 6, bottom: 6, left: 6 }
      const prices = props.points.map((item) => item.price)
      const avgPrices = props.points.map((item) => item.avgPrice)
      const allValues = [...prices, ...avgPrices, props.prevClose]
      const minValue = Math.min(...allValues)
      const maxValue = Math.max(...allValues)
      const range = maxValue - minValue || 1
      const chartWidth = width - padding.left - padding.right
      const chartHeight = height - padding.top - padding.bottom
      const x = (index: number) => padding.left + (index / Math.max(1, props.points.length - 1)) * chartWidth
      const y = (value: number) => padding.top + ((maxValue - value) / range) * chartHeight
      const pricePath = props.points.map((item, index) => `${index === 0 ? 'M' : 'L'} ${x(index)} ${y(item.price)}`).join(' ')
      const avgPath = props.points.map((item, index) => `${index === 0 ? 'M' : 'L'} ${x(index)} ${y(item.avgPrice)}`).join(' ')
      const fillPath = `${pricePath} L ${x(props.points.length - 1)} ${height - padding.bottom} L ${padding.left} ${height - padding.bottom} Z`
      const last = props.points[props.points.length - 1]
      return {
        width,
        height,
        pricePath,
        avgPath,
        fillPath,
        prevCloseY: y(props.prevClose),
        lastX: x(props.points.length - 1),
        lastY: y(last.price),
        positive: last.price >= props.prevClose,
      }
    })

    return () => {
      const value = chart.value
      if (!value) return h('div', { class: 'chart-empty' }, '暂无分时数据')
      return h('div', { class: 'timeline-chart' }, [
        h('svg', { width: '100%', height: '100%', viewBox: `0 0 ${value.width} ${value.height}`, preserveAspectRatio: 'none' }, [
          h('defs', [
            h('linearGradient', { id: 'priceGradientEOD', x1: '0%', y1: '0%', x2: '0%', y2: '100%' }, [
              h('stop', { offset: '0%', stopColor: value.positive ? 'rgba(220, 38, 38, 0.24)' : 'rgba(22, 163, 74, 0.24)' }),
              h('stop', { offset: '100%', stopColor: 'rgba(0, 0, 0, 0)' }),
            ]),
          ]),
          h('path', { d: value.fillPath, fill: 'url(#priceGradientEOD)' }),
          h('line', { x1: 5, y1: value.prevCloseY, x2: value.width - 5, y2: value.prevCloseY, stroke: 'var(--color-muted)', strokeWidth: 1, strokeDasharray: '4 2', opacity: 0.6 }),
          h('path', { d: value.avgPath, fill: 'none', stroke: '#f59e0b', strokeWidth: 1.5, opacity: 0.85 }),
          h('path', { d: value.pricePath, fill: 'none', stroke: value.positive ? 'var(--color-rise)' : 'var(--color-fall)', strokeWidth: 1.5 }),
          h('circle', { cx: value.lastX, cy: value.lastY, r: 3, fill: value.positive ? 'var(--color-rise)' : 'var(--color-fall)' }),
        ]),
        h('div', { class: 'chart-legend' }, [
          h('span', { class: 'legend-item' }, [h('span', { class: 'legend-line price' }), '价格']),
          h('span', { class: 'legend-item' }, [h('span', { class: 'legend-line avg' }), '均价']),
        ]),
      ])
    }
  },
})

const router = useRouter()
const filters = reactive<EndOfDayFilters>(getEndOfDayFilters())
const isEditing = ref(false)
const isLoading = ref(false)
const progress = ref<AnalysisProgress>({ completed: 0, total: 0, stage: '' })
const stocks = ref<EndOfDayStock[]>([])
const hasAnalyzed = ref(false)
const abortController = ref<AbortController | null>(null)
const schemes = ref<EndOfDayScheme[]>([])
const recentUsage = ref<EndOfDayRecentUsage[]>([])
const showSchemePanel = ref(false)
const showRecentPanel = ref(false)
const showSaveInput = ref(false)
const schemeName = ref('')
const sortField = ref<EndOfDaySortField>('timelineAboveAvgRatio')
const sortOrder = ref<EndOfDaySortOrder>('desc')
const selectedCodes = ref<string[]>([])
const showSelectMode = ref(false)

const sortOptions: Array<{ label: string; value: EndOfDaySortField }> = [
  { label: '涨幅', value: 'changePercent' },
  { label: '分时强度', value: 'timelineAboveAvgRatio' },
  { label: '换手率', value: 'turnoverRate' },
  { label: '流通市值', value: 'circulatingMarketCap' },
  { label: '量比', value: 'volumeRatio' },
]

const percent = computed(() => (progress.value.total > 0 ? Math.round((progress.value.completed / progress.value.total) * 100) : 0))
const sortedStocks = computed(() => sortEndOfDayStocks(stocks.value, sortField.value, sortOrder.value))

function reloadLocalState() {
  schemes.value = getEndOfDaySchemes()
  recentUsage.value = getEndOfDayRecentUsage()
}

function activateEditing() {
  if (!isEditing.value) isEditing.value = true
}

function finishEditing() {
  isEditing.value = false
  showSchemePanel.value = false
  showRecentPanel.value = false
}

function toggleSchemePanel() {
  showSchemePanel.value = !showSchemePanel.value
  showRecentPanel.value = false
}

function toggleRecentPanel() {
  showRecentPanel.value = !showRecentPanel.value
  showSchemePanel.value = false
}

function applyFilters(next: EndOfDayFilters) {
  Object.assign(filters, { ...DEFAULT_END_OF_DAY_FILTERS, ...next })
  saveEndOfDayFilters({ ...filters })
}

function resetFilters() {
  applyFilters(DEFAULT_END_OF_DAY_FILTERS)
}

function saveScheme() {
  const name = schemeName.value.trim()
  if (!name) {
    message.warning('请输入方案名称')
    return
  }
  saveEndOfDayScheme(name, { ...filters })
  schemeName.value = ''
  showSaveInput.value = false
  reloadLocalState()
  message.success(`方案「${name}」已保存`)
}

function loadScheme(scheme: EndOfDayScheme) {
  applyFilters(scheme.filters)
  showSchemePanel.value = false
  message.success(`已加载方案「${scheme.name}」`)
}

function removeScheme(id: string) {
  deleteEndOfDayScheme(id)
  reloadLocalState()
  message.success('方案已删除')
}

function loadRecent(item: EndOfDayRecentUsage) {
  applyFilters(item.filters)
  showRecentPanel.value = false
  message.success('已加载历史配置')
}

function recentSummary(item: EndOfDayFilters) {
  return `市值 ${item.marketCapMin}-${item.marketCapMax}亿 · 涨幅 ${item.changePercentMin}-${item.changePercentMax}%`
}

function formatDateTime(value: number) {
  return new Date(value).toLocaleString('zh-CN', {
    month: 'numeric',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

async function start() {
  saveEndOfDayFilters({ ...filters })
  abortController.value = new AbortController()
  isLoading.value = true
  progress.value = { completed: 0, total: 0, stage: '获取行情数据' }
  stocks.value = []
  selectedCodes.value = []
  showSelectMode.value = false
  addEndOfDayRecentUsage({ ...filters })
  reloadLocalState()

  try {
    const rows = await analyzeEndOfDayStocks({ ...filters }, {
      signal: abortController.value.signal,
      onProgress: (value) => (progress.value = value),
    })
    if (rows.length === 0) {
      message.info('没有符合基础条件的股票，请尝试调整筛选条件')
      hasAnalyzed.value = true
      return
    }
    if (!abortController.value.signal.aborted) {
      stocks.value = rows
    }
    hasAnalyzed.value = true
  } catch (error) {
    if (isAnalysisAborted(error)) {
      message.info('已取消分析')
    } else {
      message.error('获取股票数据失败，请重试')
      console.error('End-of-day analysis failed', error)
    }
  } finally {
    abortController.value = null
    isLoading.value = false
  }
}

function cancel() {
  abortController.value?.abort()
}

function handleSortChange(field: EndOfDaySortField) {
  if (field === sortField.value) {
    sortOrder.value = sortOrder.value === 'desc' ? 'asc' : 'desc'
    return
  }
  sortField.value = field
  sortOrder.value = 'desc'
}

function toggleRow(code: string) {
  selectedCodes.value = toggleSelectedCode(selectedCodes.value, code)
}

function toggleSelectAll() {
  if (selectedCodes.value.length === sortedStocks.value.length) {
    selectedCodes.value = []
    return
  }
  selectedCodes.value = sortedStocks.value.map((stock) => stock.code)
}

function toggleSelectMode() {
  showSelectMode.value = !showSelectMode.value
  if (!showSelectMode.value) selectedCodes.value = []
}

function batchAdd() {
  const selected = new Set(selectedCodes.value)
  const candidates = stocks.value
    .filter((stock) => selected.has(stock.code) && !isInWatchlist(stock.routeCode))
    .map((stock) => stock.routeCode)
  candidates.forEach((code) => addToWatchlist(code))
  selectedCodes.value = []
  showSelectMode.value = false
  if (candidates.length > 0) {
    message.success(`已将 ${candidates.length} 只股票加入自选`)
  } else {
    message.info('所选股票已在自选中')
  }
}

function add(code: string, name: string) {
  addToWatchlist(code)
  message.success(`已将 ${name} 加入自选`)
  stocks.value = [...stocks.value]
}

function returnToStart() {
  hasAnalyzed.value = false
  stocks.value = []
  showSelectMode.value = false
  selectedCodes.value = []
}

watch(filters, () => saveEndOfDayFilters({ ...filters }), { deep: true })
onMounted(reloadLocalState)
</script>

<style scoped>
.eod-page {
  --eod-surface: var(--color-surface);
  --eod-surface-elevated: var(--color-surface-elevated);
  --eod-border: var(--color-split);
  --eod-accent: var(--color-link);
  min-height: 100%;
}

.eod-header {
  border-bottom: 1px solid var(--eod-border);
  padding-bottom: 16px;
}

.title-section {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.eod-title {
  align-items: center;
  display: flex;
  gap: 8px;
}

.page-subtitle {
  color: var(--color-muted);
  font-size: 13px;
  margin: 0;
}

.eod-main,
.start-screen,
.results-screen {
  display: flex;
  flex: 1;
  flex-direction: column;
}

.start-screen {
  align-items: center;
  gap: 28px;
  padding: 28px 0;
}

.start-button-container {
  position: relative;
}

.start-button {
  align-items: center;
  background: linear-gradient(135deg, #2563eb 0%, #0f766e 100%);
  border: 0;
  border-radius: 999px;
  box-shadow: 0 12px 36px rgba(37, 99, 235, 0.28), 0 0 0 5px rgba(15, 118, 110, 0.12);
  color: #fff;
  cursor: pointer;
  display: flex;
  font-size: 18px;
  font-weight: 700;
  height: 180px;
  justify-content: center;
  overflow: hidden;
  position: relative;
  transition: box-shadow 0.2s ease, transform 0.2s ease;
  width: 180px;
}

.start-button:hover {
  box-shadow: 0 16px 48px rgba(37, 99, 235, 0.34), 0 0 0 7px rgba(15, 118, 110, 0.14);
  transform: scale(1.02);
}

.start-button:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.button-glow {
  animation: shimmer 2.4s infinite;
  background: linear-gradient(45deg, transparent 28%, rgba(255, 255, 255, 0.32) 50%, transparent 72%);
  background-size: 220% 220%;
  border-radius: inherit;
  inset: -2px;
  position: absolute;
}

.button-content {
  align-items: center;
  display: flex;
  flex-direction: column;
  gap: 10px;
  position: relative;
  z-index: 1;
}

@keyframes shimmer {
  0% { background-position: 220% 0; }
  100% { background-position: -220% 0; }
}

.filter-card {
  background: var(--eod-surface);
  border: 1px solid var(--eod-border);
  border-radius: 8px;
  cursor: pointer;
  max-width: 760px;
  padding: 22px;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
  width: 100%;
}

.filter-card:hover:not(.editing) {
  border-color: rgba(37, 99, 235, 0.35);
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.08);
}

.filter-card.editing {
  border-color: var(--eod-accent);
  cursor: default;
}

.filter-header,
.results-header,
.stock-header,
.scheme-panel-header,
.results-actions,
.sort-bar,
.sort-section,
.batch-section {
  align-items: center;
  display: flex;
}

.filter-header,
.results-header,
.sort-bar,
.scheme-panel-header {
  justify-content: space-between;
}

.filter-header {
  margin-bottom: 18px;
}

.filter-title,
.filter-actions,
.edit-actions,
.results-summary,
.stock-name-row {
  align-items: center;
  display: flex;
  gap: 8px;
}

.filter-title {
  color: var(--color-text-strong);
  font-size: 16px;
  font-weight: 650;
}

.edit-hint,
.data-label,
.price-label,
.stock-code,
.chart-legend,
.empty-hint,
.sort-label,
.recent-time,
.no-results-hint {
  color: var(--color-muted);
  font-size: 12px;
}

.reset-btn,
.save-btn,
.scheme-btn,
.recent-btn,
.add-scheme-btn,
.confirm-save-btn,
.scheme-load-btn,
.scheme-delete-btn,
.recent-item,
.back-button,
.select-mode-btn,
.sort-option,
.select-all-btn,
.batch-add-btn,
.select-btn,
.add-watchlist-btn,
.cancel-analysis-btn,
.toggle-btn {
  align-items: center;
  border: 0;
  cursor: pointer;
  display: inline-flex;
  font-family: inherit;
  justify-content: center;
  transition: background 0.16s ease, border-color 0.16s ease, color 0.16s ease, transform 0.16s ease;
}

.reset-btn,
.scheme-btn,
.recent-btn,
.back-button,
.select-mode-btn,
.select-all-btn {
  background: var(--color-hover);
  border-radius: 6px;
  color: var(--color-muted);
  gap: 4px;
  padding: 7px 12px;
}

.reset-btn,
.scheme-btn,
.recent-btn,
.save-btn {
  font-size: 12px;
  font-weight: 600;
}

.save-btn,
.batch-add-btn,
.confirm-save-btn {
  background: var(--eod-accent);
  color: #fff;
}

.save-btn {
  border-radius: 6px;
  gap: 4px;
  padding: 7px 12px;
}

.reset-btn:hover,
.scheme-btn:hover,
.recent-btn:hover,
.back-button:hover,
.select-mode-btn:hover,
.select-all-btn:hover {
  color: var(--color-text-strong);
  background: var(--eod-border);
}

.scheme-panel,
.recent-panel {
  background: var(--color-hover);
  border-radius: 8px;
  margin: 0 0 16px;
  overflow: hidden;
  padding: 14px;
}

.scheme-panel-header {
  color: var(--color-muted);
  font-size: 13px;
  font-weight: 650;
  margin-bottom: 10px;
}

.add-scheme-btn {
  background: rgba(37, 99, 235, 0.1);
  border-radius: 6px;
  color: var(--eod-accent);
  gap: 4px;
  padding: 5px 9px;
}

.save-input-row,
.scheme-item {
  display: flex;
  gap: 8px;
}

.save-input-row {
  margin-bottom: 10px;
}

.scheme-name-input {
  background: var(--eod-surface);
  border: 1px solid var(--eod-border);
  border-radius: 6px;
  color: var(--color-text);
  flex: 1;
  min-width: 0;
  padding: 7px 10px;
}

.scheme-name-input:focus,
.filter-input:focus {
  border-color: var(--eod-accent);
  outline: none;
}

.confirm-save-btn,
.scheme-delete-btn {
  border-radius: 6px;
  height: 32px;
  width: 32px;
}

.scheme-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.scheme-load-btn,
.recent-item {
  background: var(--eod-surface);
  border-radius: 6px;
  color: var(--color-text-strong);
  flex: 1;
  justify-content: flex-start;
  padding: 8px 10px;
  text-align: left;
}

.scheme-delete-btn {
  background: transparent;
  color: var(--color-muted);
}

.scheme-delete-btn:hover {
  background: rgba(220, 38, 38, 0.1);
  color: #dc2626;
}

.recent-item {
  justify-content: space-between;
  width: 100%;
}

.recent-summary,
.filter-display,
.filter-input,
.stock-code,
.change-percent,
.price-value,
.change-value,
.data-value,
.loading-percentage,
.loading-detail {
  font-variant-numeric: tabular-nums;
}

.filter-grid {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.filter-item {
  align-items: center;
  background: var(--color-hover);
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  min-height: 50px;
  padding: 12px;
}

.filter-label {
  color: var(--color-muted);
  font-size: 13px;
}

.filter-value {
  align-items: center;
  display: flex;
  gap: 6px;
}

.filter-input {
  background: var(--eod-surface);
  border: 1px solid var(--eod-border);
  border-radius: 5px;
  color: var(--color-text-strong);
  font-size: 13px;
  padding: 6px 8px;
  text-align: center;
  width: 64px;
}

.filter-separator,
.filter-unit {
  color: var(--color-muted);
  font-size: 12px;
}

.filter-display {
  color: var(--color-text-strong);
  font-size: 13px;
}

.toggle-btn {
  background: transparent;
  color: var(--color-muted);
  gap: 5px;
  padding: 0;
}

.toggle-btn.active {
  color: var(--eod-accent);
}

.results-screen {
  gap: 18px;
}

.results-header {
  background: var(--eod-surface);
  border: 1px solid var(--eod-border);
  border-radius: 8px;
  padding: 14px 18px;
}

.results-summary {
  color: var(--color-muted);
}

.summary-text {
  font-size: 14px;
}

.summary-text strong {
  color: var(--eod-accent);
}

.results-actions {
  gap: 8px;
}

.select-mode-btn.active {
  background: rgba(37, 99, 235, 0.1);
  color: var(--eod-accent);
}

.sort-bar {
  background: var(--eod-surface);
  border: 1px solid var(--eod-border);
  border-radius: 8px;
  flex-wrap: wrap;
  gap: 12px;
  padding: 10px 12px;
}

.sort-section,
.batch-section {
  flex-wrap: wrap;
  gap: 8px;
}

.sort-option {
  background: var(--color-hover);
  border-radius: 6px;
  color: var(--color-muted);
  gap: 2px;
  padding: 5px 9px;
}

.sort-option.active {
  background: rgba(37, 99, 235, 0.1);
  color: var(--eod-accent);
}

.batch-add-btn {
  border-radius: 6px;
  gap: 4px;
  padding: 6px 12px;
}

.batch-add-btn:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.stock-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
}

.stock-card {
  animation: card-in 0.32s ease both;
  animation-delay: var(--card-delay);
  background: var(--eod-surface);
  border: 1px solid var(--eod-border);
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
  padding: 18px;
  position: relative;
  transition: border-color 0.18s ease, box-shadow 0.18s ease, transform 0.18s ease;
}

.stock-card::before {
  background: var(--color-flat);
  content: '';
  height: 3px;
  inset: 0 0 auto;
  position: absolute;
}

.stock-card.positive::before {
  background: var(--color-rise);
}

.stock-card.negative::before {
  background: var(--color-fall);
}

.stock-card:hover {
  box-shadow: 0 10px 26px rgba(15, 23, 42, 0.1);
  transform: translateY(-2px);
}

.stock-card.selected {
  border-color: var(--eod-accent);
  box-shadow: 0 0 0 1px var(--eod-accent);
}

@keyframes card-in {
  from {
    opacity: 0;
    transform: translateY(18px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.stock-header {
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 14px;
}

.stock-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.stock-name {
  color: var(--color-text-strong);
  font-size: 16px;
  font-weight: 700;
  margin: 0;
}

.stock-code {
  color: var(--color-muted);
  font-size: 12px;
}

.select-btn {
  background: transparent;
  color: var(--color-muted);
  height: 22px;
  padding: 0;
  width: 22px;
}

.select-btn:hover,
.stock-card.selected .select-btn {
  color: var(--eod-accent);
}

.change-badge {
  align-items: center;
  border-radius: 6px;
  display: flex;
  font-weight: 700;
  gap: 4px;
  padding: 5px 8px;
}

.stock-card.positive .change-badge {
  background: rgba(220, 38, 38, 0.1);
  color: var(--color-rise);
}

.stock-card.negative .change-badge {
  background: rgba(22, 163, 74, 0.1);
  color: var(--color-fall);
}

.change-icon {
  font-size: 10px;
}

.price-section {
  align-items: flex-end;
  display: flex;
  justify-content: space-between;
  margin-bottom: 14px;
}

.current-price,
.price-change {
  display: flex;
  flex-direction: column;
}

.current-price {
  gap: 2px;
}

.price-change {
  align-items: flex-end;
  gap: 4px;
}

.price-value {
  color: var(--color-text-strong);
  font-size: 25px;
  font-weight: 750;
}

.stock-card.positive .change-value {
  color: var(--color-rise);
}

.stock-card.negative .change-value {
  color: var(--color-fall);
}

.timeline-ratio {
  background: rgba(37, 99, 235, 0.1);
  border-radius: 5px;
  color: var(--eod-accent);
  font-size: 12px;
  padding: 2px 6px;
}

.timeline-section {
  background: var(--color-hover);
  border-radius: 8px;
  margin-bottom: 14px;
  padding: 6px;
}

:deep(.timeline-chart) {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

:deep(.timeline-chart svg) {
  display: block;
  height: 100px;
  width: 100%;
}

:deep(.chart-legend) {
  color: var(--color-muted);
  display: flex;
  font-size: 12px;
  gap: 14px;
}

:deep(.legend-item) {
  align-items: center;
  display: flex;
  gap: 4px;
}

:deep(.legend-line) {
  border-radius: 1px;
  height: 2px;
  width: 12px;
}

:deep(.legend-line.price) {
  background: var(--color-rise);
}

:deep(.legend-line.avg) {
  background: #f59e0b;
}

:deep(.chart-empty) {
  color: var(--color-muted);
  font-size: 12px;
  padding: 18px;
  text-align: center;
}

.data-grid {
  display: grid;
  gap: 8px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  margin-bottom: 14px;
}

.data-item {
  background: var(--color-hover);
  border-radius: 6px;
  display: flex;
  justify-content: space-between;
  min-width: 0;
  padding: 8px;
}

.data-value {
  color: var(--color-text-strong);
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.add-watchlist-btn {
  background: var(--color-hover);
  border: 1px solid var(--eod-border);
  border-radius: 7px;
  color: var(--color-muted);
  gap: 6px;
  padding: 9px;
  width: 100%;
}

.add-watchlist-btn:hover:not(:disabled) {
  background: rgba(37, 99, 235, 0.1);
  border-color: var(--eod-accent);
  color: var(--eod-accent);
}

.add-watchlist-btn.added {
  background: rgba(22, 163, 74, 0.1);
  border-color: #16a34a;
  color: #16a34a;
  cursor: default;
}

.no-results {
  align-items: center;
  color: var(--color-muted);
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.no-results-title {
  color: var(--color-text);
  font-size: 18px;
  font-weight: 650;
  margin: 16px 0 0;
}

.no-results-hint {
  margin: 8px 0 0;
}

.loading-overlay {
  align-items: center;
  background: rgba(15, 23, 42, 0.76);
  backdrop-filter: blur(4px);
  display: flex;
  inset: 0;
  justify-content: center;
  position: fixed;
  z-index: 1000;
}

.loading-content {
  align-items: center;
  color: #fff;
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 32px;
}

.loading-spinner {
  height: 120px;
  position: relative;
  width: 120px;
}

.spinner-ring {
  animation: spin 1.4s linear infinite;
  border: 3px solid rgba(255, 255, 255, 0.24);
  border-radius: 50%;
  border-top-color: #60a5fa;
  inset: 0;
  position: absolute;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.spinner-center {
  align-items: center;
  display: flex;
  inset: 0;
  justify-content: center;
  position: absolute;
}

.loading-percentage {
  font-size: 28px;
  font-weight: 800;
}

.loading-progress {
  background: rgba(255, 255, 255, 0.16);
  border-radius: 99px;
  height: 4px;
  overflow: hidden;
  width: 280px;
}

.loading-progress-fill {
  background: linear-gradient(90deg, #60a5fa, #2dd4bf);
  border-radius: inherit;
  height: 100%;
  transition: width 0.25s ease;
}

.loading-status {
  text-align: center;
}

.loading-text {
  font-size: 16px;
  margin: 0 0 6px;
}

.loading-detail {
  color: rgba(255, 255, 255, 0.7);
  font-size: 13px;
  margin: 0;
}

.cancel-analysis-btn {
  background: #dc2626;
  border-radius: 7px;
  color: #fff;
  gap: 6px;
  padding: 9px 14px;
}

.cancel-analysis-btn:hover {
  background: #b91c1c;
}

@media (max-width: 768px) {
  .filter-grid,
  .stock-grid {
    grid-template-columns: 1fr;
  }

  .results-header,
  .sort-bar {
    align-items: flex-start;
    flex-direction: column;
  }

  .results-actions {
    justify-content: space-between;
    width: 100%;
  }

  .filter-card {
    padding: 16px;
  }

  .edit-actions,
  .filter-actions {
    flex-wrap: wrap;
    justify-content: flex-end;
  }
}

@media (max-width: 520px) {
  .start-button {
    height: 148px;
    width: 148px;
  }

  .filter-item,
  .price-section {
    align-items: flex-start;
    flex-direction: column;
    gap: 8px;
  }

  .filter-value {
    flex-wrap: wrap;
  }

  .data-grid {
    grid-template-columns: 1fr;
  }
}
</style>
