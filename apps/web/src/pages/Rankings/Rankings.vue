<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">榜单</h1>
        <div class="page-subtitle">行业与概念强弱、领涨股和资金容量对照</div>
      </div>
      <a-space>
        <a-segmented v-model:value="activeRankView" :options="rankViewOptions" />
        <a-segmented v-if="activeRankView !== 'limitUp'" v-model:value="rankType" :options="rankOptions" />
        <a-button :loading="activeLoading" @click="load">刷新</a-button>
      </a-space>
    </div>

    <a-alert v-if="errorMessage" type="warning" show-icon closable :message="errorMessage" @close="errorMessage = ''" />

    <a-card v-if="activeRankView === 'limitUp'" title="当日涨停板" size="small">
      <template #extra>
        <span class="muted">数据源：涨停池</span>
      </template>
      <div class="limit-up-summary">
        <a-statistic title="涨停家数" :value="limitUpPool.length" suffix="只" />
        <a-statistic title="最高连板" :value="maxContinuousBoards" suffix="板" />
        <a-statistic title="最早封板" :value="earliestFirstBoardTime" />
        <a-statistic title="封单金额" :value="formatYuanAmount(totalSealAmount)" />
      </div>
      <a-table
        :columns="limitUpColumns"
        :data-source="limitUpPool"
        :pagination="{ pageSize: 15 }"
        :loading="limitUpLoading && limitUpPool.length === 0"
        size="small"
        row-key="code"
        :custom-row="stockRowTo"
      >
        <template #bodyCell="{ column, record, index }">
          <template v-if="column.key === 'rank'">
            <span class="rank-num" :class="{ top: index < 3 }">{{ index + 1 }}</span>
          </template>
          <template v-else-if="column.key === 'name'">
            <div class="board-name">{{ record.name }}</div>
            <div class="muted">{{ record.code }}</div>
          </template>
          <template v-else-if="column.key === 'changePercent'">
            <span :class="getChangeColorClass(record.changePercent)">
              {{ formatPercent(record.changePercent) }}
            </span>
          </template>
          <template v-else-if="column.key === 'continuousBoardCount'">
            <a-tag v-if="record.continuousBoardCount" color="red">{{ formatBoardCount(record.continuousBoardCount) }}</a-tag>
            <span v-else>--</span>
          </template>
          <template v-else-if="column.key === 'boardTime'">
            <div>{{ record.firstBoardTime || '--' }}</div>
            <div class="muted">末封 {{ record.lastBoardTime || '--' }}</div>
          </template>
          <template v-else-if="column.key === 'industry'">
            {{ record.industry || '--' }}
          </template>
          <template v-else-if="column.key === 'ztStatistics'">
            {{ record.ztStatistics || '--' }}
          </template>
          <template v-else-if="column.key === 'amount'">
            {{ formatAmount(record.amount) }}
          </template>
          <template v-else-if="column.key === 'turnoverRate'">
            {{ formatTurnover(record.turnoverRate) }}
          </template>
        </template>
      </a-table>
    </a-card>

    <a-card v-else :title="activeRankView === 'industry' ? '行业板块' : '概念板块'" size="small">
      <a-table
        :columns="columns"
        :data-source="activeBoardRows"
        :pagination="{ pageSize: 20 }"
        :loading="boardLoading && activeBoardRows.length === 0"
        size="small"
        row-key="code"
        :custom-row="rowTo(activeRankView)"
      >
        <template #bodyCell="{ column, record, index }">
          <template v-if="column.key === 'rank'">
            <span class="rank-num" :class="{ top: index < 3 }">{{ index + 1 }}</span>
          </template>
          <template v-else-if="column.key === 'name'">
            <div class="board-name">{{ record.name }}</div>
            <div class="muted">{{ record.code }}</div>
          </template>
          <template v-else-if="column.key === 'changePercent'">
            <span :class="getChangeColorClass(record.changePercent)">
              {{ formatPercent(record.changePercent) }}
            </span>
          </template>
          <template v-else-if="column.key === 'leadingStock'">
            <div>{{ record.leadingStock || '--' }}</div>
            <div class="muted" :class="getChangeColorClass(record.leadingStockChangePercent)">
              {{ formatPercent(record.leadingStockChangePercent) }}
            </div>
          </template>
          <template v-else-if="column.key === 'breadth'">
            <span class="text-rise">{{ record.riseCount ?? '--' }}</span>
            <span class="muted"> / </span>
            <span class="text-fall">{{ record.fallCount ?? '--' }}</span>
          </template>
          <template v-else-if="column.key === 'totalMarketCap'">
            {{ formatMarketCap(record.totalMarketCap) }}
          </template>
          <template v-else-if="column.key === 'turnoverRate'">
            {{ formatTurnover(record.turnoverRate) }}
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getConceptList, getIndustryList, getZTPool } from '@/services/api'
import { sortBoardRankings, type BoardRankingType } from '@/services/rankings'
import type { Board, ZTPoolItem } from '@/types'
import {
  formatAmount,
  formatMarketCap,
  formatPercent,
  formatTurnover,
  formatYuanAmount,
  getChangeColorClass,
} from '@/utils/format'

type RankView = 'limitUp' | 'industry' | 'concept'

const router = useRouter()
const activeRankView = ref<RankView>('limitUp')
const rankType = ref<BoardRankingType>('rise')
const industry = ref<Board[]>([])
const concept = ref<Board[]>([])
const limitUpPool = ref<ZTPoolItem[]>([])
const boardLoading = ref(false)
const limitUpLoading = ref(false)
const errorMessage = ref('')
const rankViewOptions = [
  { label: '涨停板', value: 'limitUp' },
  { label: '行业板块', value: 'industry' },
  { label: '概念板块', value: 'concept' },
]
const rankOptions = [
  { label: '涨幅榜', value: 'rise' },
  { label: '跌幅榜', value: 'fall' },
  { label: '总市值', value: 'amount' },
  { label: '换手率', value: 'turnover' },
]

const columns = [
  { title: '排名', key: 'rank', width: 66 },
  { title: '名称', key: 'name', width: 140 },
  { title: '涨跌幅', key: 'changePercent', width: 90 },
  { title: '领涨股', key: 'leadingStock', width: 130 },
  { title: '涨/跌', key: 'breadth', width: 80 },
  { title: '总市值', key: 'totalMarketCap', width: 100 },
  { title: '换手', key: 'turnoverRate', width: 80 },
]

const limitUpColumns = [
  { title: '排名', key: 'rank', width: 66 },
  { title: '股票', key: 'name', width: 140, sorter: (a: ZTPoolItem, b: ZTPoolItem) => compareText(a.name, b.name) },
  { title: '涨幅', key: 'changePercent', width: 90, sorter: (a: ZTPoolItem, b: ZTPoolItem) => compareNumber(a.changePercent, b.changePercent) },
  { title: '连板', key: 'continuousBoardCount', width: 86, sorter: (a: ZTPoolItem, b: ZTPoolItem) => compareNumber(a.continuousBoardCount, b.continuousBoardCount) },
  { title: '封板时间', key: 'boardTime', width: 120, sorter: (a: ZTPoolItem, b: ZTPoolItem) => compareText(a.firstBoardTime, b.firstBoardTime) },
  { title: '行业', key: 'industry', width: 120, sorter: (a: ZTPoolItem, b: ZTPoolItem) => compareText(a.industry, b.industry) },
  { title: '统计', key: 'ztStatistics', width: 90, sorter: (a: ZTPoolItem, b: ZTPoolItem) => compareText(a.ztStatistics, b.ztStatistics) },
  { title: '成交额', key: 'amount', width: 100, sorter: (a: ZTPoolItem, b: ZTPoolItem) => compareNumber(a.amount, b.amount) },
  { title: '换手', key: 'turnoverRate', width: 80, sorter: (a: ZTPoolItem, b: ZTPoolItem) => compareNumber(a.turnoverRate, b.turnoverRate) },
]

const sortedIndustry = computed(() => sortBoardRankings(industry.value, rankType.value, 50))
const sortedConcept = computed(() => sortBoardRankings(concept.value, rankType.value, 50))
const activeBoardRows = computed(() => activeRankView.value === 'industry' ? sortedIndustry.value : sortedConcept.value)
const activeLoading = computed(() => activeRankView.value === 'limitUp' ? limitUpLoading.value : boardLoading.value)
const maxContinuousBoards = computed(() => {
  return Math.max(0, ...limitUpPool.value.map((item) => Number(item.continuousBoardCount || 0)))
})
const earliestFirstBoardTime = computed(() => {
  return limitUpPool.value
    .map((item) => item.firstBoardTime)
    .filter((time): time is string => Boolean(time))
    .sort()[0] || '--'
})
const totalSealAmount = computed(() => {
  const total = limitUpPool.value.reduce((sum, item) => sum + Number(item.sealAmount || item.boardAmount || 0), 0)
  return total > 0 ? total : null
})

async function load() {
  errorMessage.value = ''
  const results = await Promise.allSettled([loadBoards(), loadLimitUpPool()])

  const failed = results.filter((result) => result.status === 'rejected')
  if (failed.length > 0) {
    failed.forEach((result) => {
      if (result.status === 'rejected') console.warn('Rankings data loading failed', result.reason)
    })
    errorMessage.value = failed.length === results.length
      ? '榜单数据加载失败，请稍后刷新'
      : '榜单数据部分加载失败，请稍后刷新'
  }
}

async function loadBoards() {
  boardLoading.value = true
  try {
    const results = await Promise.allSettled([
      getIndustryList(),
      getConceptList(),
    ])
    const [industryResult, conceptResult] = results
    if (industryResult.status === 'fulfilled') industry.value = industryResult.value
    if (conceptResult.status === 'fulfilled') concept.value = conceptResult.value

    const failed = results.filter((result) => result.status === 'rejected')
    if (failed.length > 0) {
      throw new AggregateError(failed.map((result) => result.status === 'rejected' ? result.reason : undefined), 'board rankings failed')
    }
  } finally {
    boardLoading.value = false
  }
}

async function loadLimitUpPool() {
  limitUpLoading.value = true
  try {
    limitUpPool.value = await getZTPool()
  } finally {
    limitUpLoading.value = false
  }
}

function compareNumber(left: number | null | undefined, right: number | null | undefined) {
  return finiteNumber(left) - finiteNumber(right)
}

function compareText(left: string | null | undefined, right: string | null | undefined) {
  return String(left || '').localeCompare(String(right || ''), 'zh-Hans-CN')
}

function finiteNumber(value: number | null | undefined) {
  return typeof value === 'number' && Number.isFinite(value) ? value : 0
}

function rowTo(type: 'industry' | 'concept') {
  return (record: Board) => ({
    onClick: () => router.push(`/boards/${type}/${record.code}`),
  })
}

function stockRowTo(record: ZTPoolItem) {
  return {
    onClick: () => router.push(`/s/${record.code}`),
  }
}

function formatBoardCount(value: number | null | undefined) {
  if (!value || Number.isNaN(value)) return '--'
  return `${Math.floor(value)}连板`
}

onMounted(load)
</script>

<style scoped>
.page-subtitle,
.muted {
  color: var(--color-muted);
  font-size: 12px;
}

.board-name {
  font-weight: 600;
}

.limit-up-summary {
  display: grid;
  grid-template-columns: repeat(4, minmax(120px, 1fr));
  gap: 12px;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--color-border);
}

.rank-num {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 24px;
  color: var(--color-muted);
  font-variant-numeric: tabular-nums;
}

.rank-num.top {
  color: #d97706;
  font-weight: 700;
}

:deep(.ant-table-row) {
  cursor: pointer;
}

@media (max-width: 640px) {
  .page-header {
    align-items: flex-start;
    flex-direction: column;
  }

  .limit-up-summary {
    grid-template-columns: repeat(2, minmax(120px, 1fr));
  }
}
</style>
