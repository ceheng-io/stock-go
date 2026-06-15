<template>
  <div class="page eod-page">
    <div class="page-header">
      <div>
        <h1 class="page-title">尾盘选股</h1>
        <div class="page-subtitle">一日持股法筛选 · 分时强度与基础量价共振</div>
      </div>
      <a-space>
        <a-button @click="resetFilters">恢复默认</a-button>
        <a-button v-if="!isLoading" type="primary" @click="start">开始分析</a-button>
        <a-button v-else danger @click="cancel">取消</a-button>
      </a-space>
    </div>

    <a-row :gutter="[12, 12]">
      <a-col :xs="24" :xl="16">
        <a-card size="small" title="筛选条件">
          <template #extra>
            <a-space>
              <a-button size="small" @click="saveFiltersSnapshot">保存条件</a-button>
              <a-popover trigger="click" placement="bottomRight">
                <template #content>
                  <div class="scheme-popover">
                    <a-input-search
                      v-model:value="schemeName"
                      placeholder="方案名称"
                      enter-button="保存"
                      @search="saveScheme"
                    />
                    <a-divider />
                    <a-empty v-if="schemes.length === 0" description="暂无保存方案" />
                    <a-list v-else :data-source="schemes" size="small">
                      <template #renderItem="{ item }">
                        <a-list-item>
                          <a-list-item-meta :title="item.name" :description="formatDateTime(item.createdAt)" />
                          <a-space>
                            <a-button type="link" size="small" @click="applyFilters(item.filters)">应用</a-button>
                            <a-button type="link" danger size="small" @click="removeScheme(item.id)">删除</a-button>
                          </a-space>
                        </a-list-item>
                      </template>
                    </a-list>
                  </div>
                </template>
                <a-button size="small">方案</a-button>
              </a-popover>
              <a-popover trigger="click" placement="bottomRight">
                <template #content>
                  <div class="scheme-popover">
                    <a-empty v-if="recentUsage.length === 0" description="暂无最近使用" />
                    <a-list v-else :data-source="recentUsage" size="small">
                      <template #renderItem="{ item }">
                        <a-list-item class="clickable-row" @click="applyFilters(item.filters)">
                          <a-list-item-meta
                            :title="recentSummary(item.filters)"
                            :description="formatDateTime(item.usedAt)"
                          />
                        </a-list-item>
                      </template>
                    </a-list>
                  </div>
                </template>
                <a-button size="small">最近使用</a-button>
              </a-popover>
            </a-space>
          </template>

          <a-form layout="vertical">
            <a-row :gutter="12">
              <a-col :xs="24" :md="6">
                <a-form-item label="流通市值下限(亿)">
                  <a-input-number v-model:value="filters.marketCapMin" :min="0" />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :md="6">
                <a-form-item label="流通市值上限(亿)">
                  <a-input-number v-model:value="filters.marketCapMax" :min="0" />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :md="6">
                <a-form-item label="涨幅下限(%)">
                  <a-input-number v-model:value="filters.changePercentMin" :step="0.1" />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :md="6">
                <a-form-item label="涨幅上限(%)">
                  <a-input-number v-model:value="filters.changePercentMax" :step="0.1" />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :md="6">
                <a-form-item label="量比下限">
                  <a-input-number v-model:value="filters.volumeRatioMin" :min="0" :step="0.1" />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :md="6">
                <a-form-item label="换手率下限(%)">
                  <a-input-number v-model:value="filters.turnoverRateMin" :min="0" :step="0.1" />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :md="6">
                <a-form-item label="换手率上限(%)">
                  <a-input-number v-model:value="filters.turnoverRateMax" :min="0" :step="0.1" />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :md="6">
                <a-form-item label="分时强度(%)">
                  <a-input-number v-model:value="filters.timelineAboveAvgRatio" :min="0" :max="100" />
                </a-form-item>
              </a-col>
              <a-col :xs="24">
                <a-checkbox v-model:checked="filters.excludeST">排除 ST / *ST</a-checkbox>
              </a-col>
            </a-row>
          </a-form>

          <a-progress
            v-if="isLoading || progress.total > 0"
            :percent="percent"
            :status="isLoading ? 'active' : progress.stage === '已取消' ? 'exception' : 'normal'"
          />
          <p class="progress-text">{{ progress.stage }}</p>
        </a-card>
      </a-col>

      <a-col :xs="24" :xl="8">
        <a-card size="small" title="结果概览">
          <a-statistic title="入选股票" :value="stocks.length" suffix="只" />
          <a-descriptions :column="1" size="small" class="summary">
            <a-descriptions-item label="排序">{{ sortLabel }}</a-descriptions-item>
            <a-descriptions-item label="已选择">{{ selectedCodes.length }} 只</a-descriptions-item>
            <a-descriptions-item label="进度">{{ progress.completed }} / {{ progress.total }}</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>
    </a-row>

    <a-card size="small" title="结果">
      <template #extra>
        <a-space wrap>
          <a-select v-model:value="sortField" size="small" class="sort-select" :options="sortOptions" />
          <a-segmented v-model:value="sortOrder" size="small" :options="[{ label: '降序', value: 'desc' }, { label: '升序', value: 'asc' }]" />
          <a-button size="small" @click="toggleSelectAll">
            {{ selectedCodes.length === sortedStocks.length && sortedStocks.length > 0 ? '取消全选' : '全选' }}
          </a-button>
          <a-button size="small" type="primary" :disabled="selectedCodes.length === 0" @click="batchAdd">
            批量加自选
          </a-button>
        </a-space>
      </template>
      <a-table
        :columns="columns"
        :data-source="sortedStocks"
        row-key="routeCode"
        size="small"
        :pagination="{ pageSize: 20 }"
        :loading="isLoading"
        @row="rowToDetail"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'select'">
            <a-checkbox
              :checked="selectedCodes.includes(record.routeCode)"
              @click.stop
              @change="toggleRow(record.routeCode)"
            />
          </template>
          <template v-else-if="column.key === 'name'">
            <div class="stock-name">{{ record.name }}</div>
            <div class="muted">{{ record.code }}</div>
          </template>
          <template v-else-if="column.key === 'changePercent'">
            <span :class="getChangeColorClass(record.changePercent)">
              {{ formatPercent(record.changePercent) }}
            </span>
          </template>
          <template v-else-if="column.key === 'timelineAboveAvgRatio'">
            {{ formatPercent(record.timelineAboveAvgRatio, false) }}
          </template>
          <template v-else-if="column.key === 'turnoverRate'">
            {{ formatTurnover(record.turnoverRate) }}
          </template>
          <template v-else-if="column.key === 'volumeRatio'">
            {{ formatVolumeRatio(record.volumeRatio) }}
          </template>
          <template v-else-if="column.key === 'circulatingMarketCap'">
            {{ formatMarketCap(record.circulatingMarketCap) }}
          </template>
          <template v-else-if="column.key === 'amount'">
            {{ formatYuanAmount(record.amount) }}
          </template>
          <template v-else-if="column.key === 'actions'">
            <a-space>
              <a-button size="small" @click.stop="router.push(`/s/${record.routeCode}`)">详情</a-button>
              <a-button size="small" :disabled="isInWatchlist(record.routeCode)" @click.stop="add(record.routeCode)">
                {{ isInWatchlist(record.routeCode) ? '已自选' : '加自选' }}
              </a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import {
  type AnalysisProgress,
  type EndOfDayFilters,
  type EndOfDayStock,
  analyzeEndOfDayStocks,
  isAnalysisAborted,
} from '@/services/analysis'
import {
  DEFAULT_END_OF_DAY_FILTERS,
  addEndOfDayRecentUsage,
  deleteEndOfDayScheme,
  getBatchWatchlistCandidates,
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
  formatMarketCap,
  formatPercent,
  formatPrice,
  formatTurnover,
  formatVolumeRatio,
  formatYuanAmount,
  getChangeColorClass,
} from '@/utils/format'

const router = useRouter()
const filters = reactive<EndOfDayFilters>(getEndOfDayFilters())
const isLoading = ref(false)
const progress = ref<AnalysisProgress>({ completed: 0, total: 0, stage: '待开始' })
const stocks = ref<EndOfDayStock[]>([])
const abortController = ref<AbortController | null>(null)
const schemes = ref<EndOfDayScheme[]>([])
const recentUsage = ref<EndOfDayRecentUsage[]>([])
const schemeName = ref('')
const sortField = ref<EndOfDaySortField>('timelineAboveAvgRatio')
const sortOrder = ref<EndOfDaySortOrder>('desc')
const selectedCodes = ref<string[]>([])

const sortOptions: Array<{ label: string; value: EndOfDaySortField }> = [
  { label: '分时强度', value: 'timelineAboveAvgRatio' },
  { label: '涨幅', value: 'changePercent' },
  { label: '换手率', value: 'turnoverRate' },
  { label: '流通市值', value: 'circulatingMarketCap' },
  { label: '量比', value: 'volumeRatio' },
]

const columns = [
  { title: '', key: 'select', width: 48 },
  { title: '股票', key: 'name', fixed: 'left' },
  { title: '价格', dataIndex: 'price', customRender: ({ record }: { record: EndOfDayStock }) => formatPrice(record.price) },
  { title: '涨幅', key: 'changePercent' },
  { title: '分时强度', key: 'timelineAboveAvgRatio' },
  { title: '换手', key: 'turnoverRate' },
  { title: '量比', key: 'volumeRatio' },
  { title: '流通市值', key: 'circulatingMarketCap' },
  { title: '成交额', key: 'amount' },
  { title: '操作', key: 'actions', width: 150 },
]

const percent = computed(() => (progress.value.total > 0 ? Math.round((progress.value.completed / progress.value.total) * 100) : 0))
const sortedStocks = computed(() => sortEndOfDayStocks(stocks.value, sortField.value, sortOrder.value))
const sortLabel = computed(() => `${sortOptions.find((item) => item.value === sortField.value)?.label || '--'} · ${sortOrder.value === 'desc' ? '降序' : '升序'}`)

function reloadLocalState() {
  schemes.value = getEndOfDaySchemes()
  recentUsage.value = getEndOfDayRecentUsage()
}

function applyFilters(next: EndOfDayFilters) {
  Object.assign(filters, { ...DEFAULT_END_OF_DAY_FILTERS, ...next })
  saveFiltersSnapshot(false)
}

function resetFilters() {
  applyFilters(DEFAULT_END_OF_DAY_FILTERS)
  message.success('已恢复默认条件')
}

function saveFiltersSnapshot(showMessage = true) {
  saveEndOfDayFilters({ ...filters })
  if (showMessage) message.success('筛选条件已保存')
}

function saveScheme() {
  const name = schemeName.value.trim()
  if (!name) {
    message.warning('请输入方案名称')
    return
  }
  saveEndOfDayScheme(name, { ...filters })
  schemeName.value = ''
  reloadLocalState()
  message.success(`方案「${name}」已保存`)
}

function removeScheme(id: string) {
  deleteEndOfDayScheme(id)
  reloadLocalState()
  message.success('方案已删除')
}

function recentSummary(item: EndOfDayFilters) {
  return `市值 ${item.marketCapMin}-${item.marketCapMax} 亿 · 涨幅 ${item.changePercentMin}-${item.changePercentMax}% · 强度 ${item.timelineAboveAvgRatio}%`
}

function formatDateTime(value: number) {
  return new Date(value).toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

async function start() {
  saveFiltersSnapshot(false)
  abortController.value = new AbortController()
  isLoading.value = true
  progress.value = { completed: 0, total: 0, stage: '获取行情数据' }
  stocks.value = []
  selectedCodes.value = []
  addEndOfDayRecentUsage({ ...filters })
  reloadLocalState()

  try {
    const rows = await analyzeEndOfDayStocks({ ...filters }, {
      signal: abortController.value.signal,
      onProgress: (value) => (progress.value = value),
    })
    stocks.value = rows
    if (rows.length === 0) {
      message.info('没有符合条件的股票，请调整筛选条件')
    }
  } catch (error) {
    if (isAnalysisAborted(error)) {
      progress.value = { ...progress.value, stage: '已取消' }
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

function add(code: string) {
  addToWatchlist(code)
  message.success('已加入自选')
}

function toggleRow(routeCode: string) {
  selectedCodes.value = toggleSelectedCode(selectedCodes.value, routeCode)
}

function toggleSelectAll() {
  if (selectedCodes.value.length === sortedStocks.value.length && sortedStocks.value.length > 0) {
    selectedCodes.value = []
    return
  }
  selectedCodes.value = sortedStocks.value.map((stock) => stock.routeCode)
}

function batchAdd() {
  const candidates = getBatchWatchlistCandidates(stocks.value, selectedCodes.value, isInWatchlist)
  candidates.forEach((code) => addToWatchlist(code))
  selectedCodes.value = []
  if (candidates.length > 0) {
    message.success(`已将 ${candidates.length} 只股票加入自选`)
  } else {
    message.info('所选股票已在自选中')
  }
}

function rowToDetail(record: EndOfDayStock) {
  return {
    onClick: () => router.push(`/s/${record.routeCode}`),
  }
}

watch(filters, () => saveEndOfDayFilters({ ...filters }), { deep: true })
onMounted(reloadLocalState)
</script>

<style scoped>
.eod-page :deep(.ant-input-number) {
  width: 100%;
}

.page-subtitle,
.muted,
.progress-text {
  color: var(--color-muted);
  font-size: 12px;
}

.progress-text {
  margin: 8px 0 0;
}

.summary {
  margin-top: 12px;
}

.scheme-popover {
  width: min(520px, 80vw);
}

.sort-select {
  width: 120px;
}

.clickable-row {
  cursor: pointer;
}

.stock-name {
  font-weight: 600;
}

:deep(.ant-table-row) {
  cursor: pointer;
}

@media (max-width: 640px) {
  .page-header {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
