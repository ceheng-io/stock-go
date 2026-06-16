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
        <a-space>
          <a-date-picker
            v-model:value="limitUpDate"
            value-format="YYYY-MM-DD"
            allow-clear
            placeholder="当前交易日"
            @change="loadLimitUpPool"
          />
          <span class="muted">数据源：涨停池</span>
        </a-space>
      </template>
      <div class="limit-up-summary">
        <a-statistic title="涨停家数" :value="limitUpPool.length" suffix="只" />
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
          <template v-else-if="column.key === 'ztStatistics'">
            {{ record.ztStatistics || '--' }}
          </template>
          <template v-else-if="column.key === 'boardTime'">
            <div>{{ record.firstBoardTime || '--' }}</div>
            <div class="muted">末封 {{ record.lastBoardTime || '--' }}</div>
          </template>
          <template v-else-if="column.key === 'industry'">
            {{ record.industry || '--' }}
          </template>
          <template v-else-if="column.key === 'reasonType'">
            <div class="limit-up-reason">{{ record.reasonType || '--' }}</div>
            <div v-if="record.limitUpType" class="muted">{{ record.limitUpType }}</div>
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
import { getConceptList, getIndustryList, getTHSLimitUpPool } from '@/services/api'
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
const limitUpDate = ref('')
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
  { title: '统计', key: 'ztStatistics', width: 90, sorter: (a: ZTPoolItem, b: ZTPoolItem) => compareText(a.ztStatistics, b.ztStatistics) },
  { title: '封板时间', key: 'boardTime', width: 120, sorter: (a: ZTPoolItem, b: ZTPoolItem) => compareText(a.firstBoardTime, b.firstBoardTime) },
  { title: '行业', key: 'industry', width: 120, sorter: (a: ZTPoolItem, b: ZTPoolItem) => compareText(a.industry, b.industry) },
  { title: '涨停原因', key: 'reasonType', width: 150, sorter: (a: ZTPoolItem, b: ZTPoolItem) => compareText(a.reasonType, b.reasonType) },
  { title: '成交额', key: 'amount', width: 100, sorter: (a: ZTPoolItem, b: ZTPoolItem) => compareNumber(a.amount, b.amount) },
  { title: '换手', key: 'turnoverRate', width: 80, sorter: (a: ZTPoolItem, b: ZTPoolItem) => compareNumber(a.turnoverRate, b.turnoverRate) },
]

const sortedIndustry = computed(() => sortBoardRankings(industry.value, rankType.value, 50))
const sortedConcept = computed(() => sortBoardRankings(concept.value, rankType.value, 50))
const activeBoardRows = computed(() => activeRankView.value === 'industry' ? sortedIndustry.value : sortedConcept.value)
const activeLoading = computed(() => activeRankView.value === 'limitUp' ? limitUpLoading.value : boardLoading.value)
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
  const [boardResults, limitUpResults] = await Promise.all([loadBoards(), loadLimitUpPool()])
  const results = [...boardResults, ...limitUpResults]

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
    return await Promise.allSettled([
      getIndustryList().then((value) => {
        industry.value = value
      }),
      getConceptList().then((value) => {
        concept.value = value
      }),
    ])
  } finally {
    boardLoading.value = false
  }
}

async function loadLimitUpPool() {
  limitUpLoading.value = true
  try {
    return await Promise.allSettled([
      getTHSLimitUpPool({ date: limitUpDate.value, limit: 100 }).then((value) => {
        limitUpPool.value = value
      }),
    ])
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

.limit-up-reason {
  max-width: 150px;
  overflow-wrap: anywhere;
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
