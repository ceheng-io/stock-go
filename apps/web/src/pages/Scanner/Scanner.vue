<template>
  <div class="page scanner-page">
    <div class="page-header">
      <div>
        <h1 class="page-title">信号扫描</h1>
        <div class="page-subtitle">{{ progress.stage }} · {{ results.length }} 个触发</div>
      </div>
      <a-space>
        <a-button v-if="!isScanning" type="primary" :disabled="selectedSignals.length === 0" @click="startScan">开始扫描</a-button>
        <a-button v-else danger @click="cancelScan">取消扫描</a-button>
      </a-space>
    </div>

    <a-alert v-if="error" type="error" show-icon :message="error" />

    <div class="scanner-shell">
      <a-card title="股票池来源" size="small">
        <a-form layout="vertical">
          <a-form-item label="来源">
            <a-segmented v-model:value="poolSource" :options="poolSourceOptions" />
          </a-form-item>

          <template v-if="poolSource === 'board'">
            <a-form-item label="板块类型">
              <a-segmented v-model:value="boardType" :options="boardTypeOptions" />
            </a-form-item>
            <a-form-item label="板块">
              <a-select
                v-model:value="selectedBoardCode"
                show-search
                allow-clear
                :loading="boardsLoading"
                :options="boardOptions"
                :filter-option="filterBoardOption"
                placeholder="选择行业或概念"
              />
            </a-form-item>
            <a-form-item label="成分数量">
              <a-segmented v-model:value="boardLimit" :options="[30, 50, 80]" />
            </a-form-item>
          </template>

          <template v-if="poolSource === 'ranking'">
            <a-form-item label="榜单字段">
              <a-segmented v-model:value="rankingField" :options="rankingOptions" />
            </a-form-item>
            <a-form-item label="范围">
              <a-segmented v-model:value="poolTopN" :options="topNOptions" />
            </a-form-item>
          </template>

          <template v-if="poolSource === 'zt_pool'">
            <a-form-item label="涨停/强势池">
              <a-select v-model:value="ztPoolType" :options="ztPoolOptions" />
            </a-form-item>
            <a-form-item label="范围">
              <a-segmented v-model:value="poolTopN" :options="topNOptions" />
            </a-form-item>
          </template>

          <template v-if="poolSource === 'stock_changes'">
            <a-form-item label="异动类型">
              <a-select v-model:value="stockChangeType" :options="stockChangeOptions" />
            </a-form-item>
            <a-form-item label="范围">
              <a-segmented v-model:value="poolTopN" :options="topNOptions" />
            </a-form-item>
          </template>

          <a-form-item label="信号模板">
            <div class="signal-grid">
              <a-checkable-tag
                v-for="signal in signalTemplates"
                :key="signal.value"
                :checked="selectedSignals.includes(signal.value)"
                @change="toggleSignal(signal.value)"
              >
                <strong>{{ signal.label }}</strong>
                <span>{{ signal.desc }}</span>
              </a-checkable-tag>
            </div>
          </a-form-item>

          <a-progress
            v-if="isScanning || progress.total > 0"
            :percent="percent"
            :status="isScanning ? 'active' : 'normal'"
          />
        </a-form>
      </a-card>

      <a-card size="small">
        <template #title>
          <a-space>
            <span>扫描结果</span>
            <a-tag>{{ results.length }}</a-tag>
          </a-space>
        </template>
        <a-empty v-if="results.length === 0" :description="isScanning ? '正在扫描...' : '暂无扫描结果'" />
        <a-table
          v-else
          :columns="columns"
          :data-source="results"
          row-key="routeCode"
          size="small"
          :pagination="{ pageSize: 20 }"
          :custom-row="rowClick"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'matchedSignals'">
              <a-space wrap>
                <a-tag v-for="signal in record.matchedSignals" :key="signal" color="blue">{{ signal }}</a-tag>
              </a-space>
            </template>
            <template v-if="column.key === 'actions'">
              <a-space>
                <a-button size="small" @click.stop="router.push(`/s/${record.routeCode}`)">详情</a-button>
                <a-button size="small" :disabled="record.added" @click.stop="saveResultToWatchlist(record)">
                  {{ record.added ? '已加入' : '加自选' }}
                </a-button>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import {
  getAllAShareQuotes,
  getConceptConstituents,
  getConceptList,
  getIndustryConstituents,
  getIndustryList,
  getStockChanges,
  getZTPool,
} from '@/services/api'
import { addToWatchlist as saveToWatchlist, getAllWatchlistCodes, isInWatchlist } from '@/services/storage'
import {
  isAnalysisAborted,
  type AnalysisProgress,
  type ScannerSignalKey,
  type ScannerSignalResult,
  type ScannerStockPoolItem,
  scanSignalPool,
} from '@/services/analysis'
import type { Board, FullQuote } from '@/types'
import { formatAmount, formatPercent, normalizeStockCode, parseStockCode } from '@/utils/format'

type PoolSource = 'watchlist' | 'board' | 'ranking' | 'zt_pool' | 'stock_changes'
type BoardType = 'industry' | 'concept'
type RankingField = 'amount' | 'changePercent' | 'turnoverRate'

interface ScanResultRow extends ScannerSignalResult {
  added: boolean
  time: string
}

const router = useRouter()
const poolSource = ref<PoolSource>('watchlist')
const boardType = ref<BoardType>('industry')
const selectedBoardCode = ref<string>()
const boardLimit = ref(50)
const rankingField = ref<RankingField>('amount')
const poolTopN = ref(20)
const ztPoolType = ref('strong')
const stockChangeType = ref('rocket_launch')
const selectedSignals = ref<ScannerSignalKey[]>(['ma_golden'])
const isScanning = ref(false)
const results = ref<ScanResultRow[]>([])
const progress = ref<AnalysisProgress>({ completed: 0, total: 0, stage: '待开始' })
const abortController = ref<AbortController | null>(null)
const boards = ref<Board[]>([])
const boardsLoading = ref(false)
const error = ref('')

const poolSourceOptions = [
  { label: '自选', value: 'watchlist' },
  { label: '板块', value: 'board' },
  { label: '榜单', value: 'ranking' },
  { label: '涨停池', value: 'zt_pool' },
  { label: '盘口异动池', value: 'stock_changes' },
]
const boardTypeOptions = [{ label: '行业', value: 'industry' }, { label: '概念', value: 'concept' }]
const rankingOptions = [{ label: '成交额', value: 'amount' }, { label: '涨幅', value: 'changePercent' }, { label: '换手率', value: 'turnoverRate' }]
const topNOptions = [20, 50, 100]
const ztPoolOptions = [
  { label: '涨停池', value: 'zt' },
  { label: '强势股', value: 'strong' },
  { label: '昨日涨停', value: 'yesterday' },
  { label: '次新股', value: 'sub_new' },
  { label: '炸板池', value: 'broken' },
  { label: '跌停池', value: 'dt' },
]
const stockChangeOptions = [
  { label: '火箭发射', value: 'rocket_launch' },
  { label: '大笔买入', value: 'large_buy' },
  { label: '大单扫货', value: 'big_buy_order' },
  { label: '封涨停板', value: 'limit_up_seal' },
  { label: '向上缺口', value: 'gap_up' },
  { label: '60日新高', value: 'high_60d' },
]
const signalTemplates: Array<{ value: ScannerSignalKey; label: string; desc: string }> = [
  { value: 'ma_golden', label: 'MA金叉', desc: '短期均线上穿长期均线' },
  { value: 'ma_death', label: 'MA死叉', desc: '短期均线下穿长期均线' },
  { value: 'macd_golden', label: 'MACD金叉', desc: 'DIF 上穿 DEA' },
  { value: 'macd_death', label: 'MACD死叉', desc: 'DIF 下穿 DEA' },
  { value: 'rsi_oversold', label: 'RSI超卖', desc: 'RSI 低于 30' },
  { value: 'rsi_overbought', label: 'RSI超买', desc: 'RSI 高于 70' },
  { value: 'boll_upper', label: 'BOLL上轨', desc: '收盘价突破上轨' },
  { value: 'boll_lower', label: 'BOLL下轨', desc: '收盘价跌破下轨' },
]

const columns = [
  { title: '代码', dataIndex: 'code', width: 110 },
  { title: '名称', dataIndex: 'name', width: 140 },
  { title: '命中信号', key: 'matchedSignals' },
  { title: '时间', dataIndex: 'time', width: 180 },
  { title: '操作', key: 'actions', width: 150 },
]

const percent = computed(() => {
  if (progress.value.total <= 0) return 0
  return Math.round((progress.value.completed / progress.value.total) * 100)
})

const boardOptions = computed(() => boards.value.map((item) => ({ label: `${item.name} ${formatPercent(item.changePercent)}`, value: item.code })))

function filterBoardOption(input: string, option?: { label?: string; value?: string }) {
  return String(option?.label || '').toLowerCase().includes(input.toLowerCase())
}

async function loadBoards() {
  if (poolSource.value !== 'board') return
  boardsLoading.value = true
  try {
    boards.value = boardType.value === 'industry' ? await getIndustryList() : await getConceptList()
  } finally {
    boardsLoading.value = false
  }
}

function toggleSignal(signal: ScannerSignalKey) {
  selectedSignals.value = selectedSignals.value.includes(signal)
    ? selectedSignals.value.filter((item) => item !== signal)
    : [...selectedSignals.value, signal]
}

async function resolvePool(): Promise<ScannerStockPoolItem[]> {
  if (poolSource.value === 'watchlist') return resolveWatchlistPool()
  if (poolSource.value === 'board') return resolveBoardPool()
  if (poolSource.value === 'zt_pool') return resolveZTPool()
  if (poolSource.value === 'stock_changes') return resolveStockChanges()
  return resolveRankingPool()
}

function resolveWatchlistPool() {
  return getAllWatchlistCodes().map((code) => {
    const routeCode = normalizeStockCode(code)
    return { code: parseStockCode(routeCode).symbol || routeCode, routeCode, name: routeCode }
  })
}

async function resolveBoardPool() {
  if (!selectedBoardCode.value) {
    message.info('请先选择一个板块')
    return []
  }
  const rows = boardType.value === 'industry'
    ? await getIndustryConstituents(selectedBoardCode.value)
    : await getConceptConstituents(selectedBoardCode.value)
  return (rows as Array<{ code: string; name: string }>).slice(0, boardLimit.value).map(toPoolItem)
}

async function resolveRankingPool() {
  const quotes = await getAllAShareQuotes({ batchSize: 500, concurrency: 4 }) as FullQuote[]
  return [...quotes]
    .sort((left, right) => Number(right[rankingField.value] || 0) - Number(left[rankingField.value] || 0))
    .slice(0, poolTopN.value)
    .map((quote) => toPoolItem({ code: quote.code, name: `${quote.name} · ${rankingField.value === 'amount' ? formatAmount(quote.amount) : formatPercent(Number(quote[rankingField.value] || 0))}` }))
}

async function resolveZTPool() {
  const rows = await getZTPool(ztPoolType.value) as Array<{ code: string; name: string }>
  return rows.slice(0, poolTopN.value).map(toPoolItem)
}

async function resolveStockChanges() {
  const rows = await getStockChanges(stockChangeType.value) as Array<{ code: string; name: string }>
  const deduped = new Map<string, ScannerStockPoolItem>()
  rows.forEach((row) => {
    const item = toPoolItem(row)
    if (item.routeCode && !deduped.has(item.routeCode)) deduped.set(item.routeCode, item)
  })
  return Array.from(deduped.values()).slice(0, poolTopN.value)
}

function toPoolItem(item: { code: string; name: string }): ScannerStockPoolItem {
  const routeCode = normalizeStockCode(item.code)
  return {
    code: parseStockCode(routeCode).symbol || item.code,
    routeCode,
    name: item.name,
  }
}

async function startScan() {
  if (selectedSignals.value.length === 0) {
    message.info('请至少选择一个信号模板')
    return
  }
  error.value = ''
  results.value = []
  progress.value = { completed: 0, total: 0, stage: '准备股票池' }
  const pool = await resolvePool()
  if (pool.length === 0) {
    progress.value = { completed: 0, total: 0, stage: '股票池为空' }
    return
  }
  abortController.value = new AbortController()
  isScanning.value = true
  const startedAt = new Date().toLocaleString('zh-CN', { hour12: false })
  try {
    const rows = await scanSignalPool(pool, selectedSignals.value, {
      signal: abortController.value.signal,
      concurrency: 4,
      onProgress: (value) => (progress.value = value),
      onResult: (row) => {
        results.value.push({ ...row, time: startedAt, added: isInWatchlist(row.routeCode) })
      },
    })
    results.value = rows.map((row) => ({ ...row, time: startedAt, added: isInWatchlist(row.routeCode) }))
    if (rows.length === 0) message.info('没有股票命中当前信号')
  } catch (err) {
    if (isAnalysisAborted(err)) {
      progress.value = { ...progress.value, stage: '已取消' }
      message.info('已取消扫描')
    } else {
      error.value = err instanceof Error ? err.message : '扫描失败'
    }
  } finally {
    abortController.value = null
    isScanning.value = false
  }
}

function cancelScan() {
  abortController.value?.abort()
}

function rowClick(record: ScanResultRow) {
  return { onClick: () => router.push(`/s/${record.routeCode}`) }
}

function saveResultToWatchlist(record: ScanResultRow) {
  saveToWatchlist(record.routeCode)
  record.added = true
  message.success(`已将 ${record.name} 加入自选`)
}

watch([poolSource, boardType], () => {
  selectedBoardCode.value = undefined
  loadBoards()
})

onMounted(loadBoards)
</script>

<style scoped>
.scanner-page {
  min-width: 0;
}

.page-subtitle {
  margin-top: 4px;
  color: var(--color-text-secondary);
  font-size: 13px;
}

.scanner-shell {
  display: grid;
  grid-template-columns: 360px minmax(0, 1fr);
  gap: 12px;
  align-items: start;
}

.signal-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}

.signal-grid :deep(.ant-tag) {
  display: flex;
  min-height: 64px;
  flex-direction: column;
  justify-content: center;
  gap: 4px;
  margin: 0;
  white-space: normal;
}

.signal-grid span {
  color: var(--color-text-secondary);
  font-size: 12px;
}

@media (max-width: 980px) {
  .scanner-shell {
    grid-template-columns: 1fr;
  }
}
</style>
