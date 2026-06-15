<template>
  <div class="page dashboard-page">
    <div class="page-header">
      <div>
        <h1 class="page-title">总览</h1>
        <div class="page-subtitle">指数、宽度、资金与热点的一屏快照</div>
      </div>
      <a-space>
        <a-button @click="router.push('/watchlist')">自选管理</a-button>
        <a-button type="primary" :loading="loading" @click="load">
          刷新
        </a-button>
      </a-space>
    </div>

    <a-alert
      v-if="errorMessage"
      type="warning"
      show-icon
      closable
      :message="errorMessage"
      @close="errorMessage = ''"
    />

    <a-row :gutter="[12, 12]">
      <a-col v-for="item in indices" :key="item.code" :xs="24" :sm="12" :lg="8" :xl="4">
        <a-card size="small" hoverable class="index-card" @click="goStock(item.code)">
          <div class="index-name">{{ item.name }}</div>
          <div class="index-price" :class="getChangeColorClass(item.changePercent)">
            {{ formatPrice(item.price) }}
          </div>
          <div class="index-meta">
            <span :class="getChangeColorClass(item.changePercent)">
              {{ formatPercent(item.changePercent) }}
            </span>
            <span :class="getChangeColorClass(item.change)">
              {{ formatChange(item.change) }}
            </span>
          </div>
          <div class="muted">成交 {{ formatAmount(item.amount) }}</div>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[12, 12]">
      <a-col :xs="24" :sm="12" :xl="4">
        <a-card title="市场涨跌" size="small">
          <div class="stat-main">
            <span class="text-rise">{{ marketSummary.riseCount }}</span>
            <span class="stat-divider">/</span>
            <span class="text-fall">{{ marketSummary.fallCount }}</span>
          </div>
          <div class="stat-meta">上涨 / 下跌 · {{ marketSummary.flatCount }} 平</div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :xl="4">
        <a-card title="涨跌停" size="small">
          <div class="stat-main">
            <span class="text-rise">{{ marketSummary.limitUpCount }}</span>
            <span class="stat-divider">/</span>
            <span class="text-fall">{{ marketSummary.limitDownCount }}</span>
          </div>
          <div class="stat-meta">涨停 / 跌停</div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :xl="4">
        <a-card title="全市场成交额" size="small">
          <div class="stat-value">{{ formatAmount(marketSummary.totalAmount) }}</div>
          <div class="stat-meta">A 股实时快照</div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :xl="4">
        <a-card title="北向资金" size="small">
          <div class="stat-value" :class="getChangeColorClass(northboundNetInflow)">
            {{ formatYuanAmount(northboundNetInflow) }}
          </div>
          <div class="stat-meta">{{ northboundSnapshot?.boardName || northboundSnapshot?.direction || '北向汇总' }}</div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :xl="4">
        <a-card title="大盘主力" size="small">
          <div class="stat-value" :class="getChangeColorClass(latestMarketFundFlow?.mainNetInflow)">
            {{ formatYuanAmount(latestMarketFundFlow?.mainNetInflow) }}
          </div>
          <div class="stat-meta">占比 {{ formatPercent(latestMarketFundFlow?.mainNetInflowPercent) }}</div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :xl="4">
        <a-card title="最强板块" size="small">
          <div class="stat-value text-clip">{{ strongestBoard?.name || '--' }}</div>
          <div class="stat-meta" :class="getChangeColorClass(strongestBoard?.changePercent)">
            {{ formatPercent(strongestBoard?.changePercent) }}
          </div>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[12, 12]">
      <a-col :xs="24" :xl="10">
        <a-card size="small" title="自选股">
          <template #extra>
            <a-button type="link" size="small" @click="router.push('/watchlist')">管理</a-button>
          </template>
          <a-empty v-if="watchlistCodes.length === 0" description="暂无自选股">
            <a-button type="primary" @click="router.push('/watchlist')">添加自选</a-button>
          </a-empty>
          <a-list v-else :data-source="watchlistQuotes.slice(0, 10)" size="small" :loading="loading && watchlistQuotes.length === 0">
            <template #renderItem="{ item }">
              <a-list-item class="clickable-row" @click="goStock(item.code)">
                <a-list-item-meta :title="item.name" :description="item.code" />
                <div class="quote-side">
                  <div>{{ formatPrice(item.price) }}</div>
                  <div :class="getChangeColorClass(item.changePercent)">
                    {{ formatPercent(item.changePercent) }}
                  </div>
                </div>
              </a-list-item>
            </template>
          </a-list>
        </a-card>

        <a-card size="small" title="市场榜单">
          <template #extra>
            <a-segmented v-model:value="rankingTab" :options="rankingOptions" size="small" />
          </template>
          <a-list :data-source="rankingItems" size="small" :loading="loading && marketQuotes.length === 0">
            <template #renderItem="{ item, index }">
              <a-list-item class="clickable-row" @click="goStock(item.code)">
                <span class="rank-num">{{ index + 1 }}</span>
                <a-list-item-meta :title="item.name" :description="item.code" />
                <div class="quote-side">
                  <div>{{ formatPrice(item.price) }}</div>
                  <div :class="rankingValueClass(item)">
                    {{ rankingValue(item) }}
                  </div>
                </div>
              </a-list-item>
            </template>
            <template #empty>
              <a-empty description="暂无全市场快照" />
            </template>
          </a-list>
        </a-card>

        <a-card size="small" title="主力净流入榜">
          <a-list :data-source="fundFlowRanks.slice(0, 8)" size="small" :loading="loading && fundFlowRanks.length === 0">
            <template #renderItem="{ item, index }">
              <a-list-item class="clickable-row" @click="goStock(item.code)">
                <span class="rank-num">{{ index + 1 }}</span>
                <a-list-item-meta :title="item.name" :description="item.code" />
                <div class="quote-side">
                  <div>{{ formatPrice(item.price) }}</div>
                  <div :class="getChangeColorClass(item.mainNetInflow)">
                    {{ formatYuanAmount(item.mainNetInflow) }}
                  </div>
                </div>
              </a-list-item>
            </template>
            <template #empty>
              <a-empty description="暂无资金榜单" />
            </template>
          </a-list>
        </a-card>
      </a-col>

      <a-col :xs="24" :xl="14">
        <a-card size="small" title="热点板块">
          <template #extra>
            <a-segmented v-model:value="boardTab" :options="boardOptions" size="small" />
          </template>
          <a-list :data-source="currentBoards.slice(0, 15)" size="small" :loading="loading && currentBoards.length === 0">
            <template #renderItem="{ item, index }">
              <a-list-item class="clickable-row board-row" @click="goBoard(item.code)">
                <span class="rank-num">{{ item.rank || index + 1 }}</span>
                <a-list-item-meta>
                  <template #title>{{ item.name }}</template>
                  <template #description>
                    领涨：{{ item.leadingStock || '--' }}
                    <span :class="getChangeColorClass(item.leadingStockChangePercent)">
                      {{ formatPercent(item.leadingStockChangePercent) }}
                    </span>
                  </template>
                </a-list-item-meta>
                <div class="board-side">
                  <div :class="getChangeColorClass(item.changePercent)">
                    {{ formatPercent(item.changePercent) }}
                  </div>
                  <div class="muted">
                    <span class="text-rise">{{ item.riseCount ?? '--' }}↑</span>
                    <span class="text-fall">{{ item.fallCount ?? '--' }}↓</span>
                  </div>
                </div>
              </a-list-item>
            </template>
            <template #empty>
              <a-empty description="暂无板块数据" />
            </template>
          </a-list>
        </a-card>

        <a-card size="small" title="板块资金流">
          <template #extra>
            <a-segmented v-model:value="boardTab" :options="boardOptions" size="small" />
          </template>
          <a-list :data-source="currentFundFlowBoards.slice(0, 12)" size="small" :loading="loading && currentFundFlowBoards.length === 0">
            <template #renderItem="{ item, index }">
              <a-list-item class="clickable-row board-row" @click="goBoard(item.code)">
                <span class="rank-num">{{ index + 1 }}</span>
                <a-list-item-meta>
                  <template #title>{{ item.name }}</template>
                  <template #description>
                    领流：{{ item.topStockName || '--' }}
                    <span v-if="item.topStockCode"> · {{ item.topStockCode }}</span>
                  </template>
                </a-list-item-meta>
                <div class="board-side">
                  <div :class="getChangeColorClass(item.mainNetInflow)">
                    {{ formatYuanAmount(item.mainNetInflow) }}
                  </div>
                  <div class="muted">
                    <span :class="getChangeColorClass(item.changePercent)">
                      {{ formatPercent(item.changePercent) }}
                    </span>
                    <span :class="getChangeColorClass(item.mainNetInflowPercent)">
                      {{ formatPercent(item.mainNetInflowPercent) }}
                    </span>
                  </div>
                </div>
              </a-list-item>
            </template>
            <template #empty>
              <a-empty description="暂无板块资金流" />
            </template>
          </a-list>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { usePolling } from '@/composables/usePolling'
import {
  getAllAShareQuotes,
  getConceptList,
  getFullQuotes,
  getFundFlowRank,
  getIndustryList,
  getMarketFundFlow,
  getNorthboundFlowSummary,
  getSectorFundFlowRank,
} from '@/services/api'
import {
  buildDashboardFailureMessage,
  getLatestMarketFundFlow,
  loadDashboardSections,
  pickNorthboundSnapshot,
  rankDashboardQuotes,
  summarizeMarketBreadth,
  type DashboardRankingTab,
} from '@/services/dashboard'
import { getAllWatchlistCodes, getRefreshInterval } from '@/services/storage'
import type { Board, FullQuote } from '@/types'
import {
  formatAmount,
  formatChange,
  formatPercent,
  formatPrice,
  formatTurnover,
  formatYuanAmount,
  getChangeColorClass,
} from '@/utils/format'

const MAIN_INDICES = ['sh000001', 'sz399001', 'sz399006', 'sh000688', 'sz399300', 'sh000016']

interface MarketFundFlowRow {
  date?: string
  mainNetInflow?: number | null
  mainNetInflowPercent?: number | null
}

interface NorthboundSummaryRow {
  direction?: string
  boardName?: string
  netInflow?: number | null
  netBuyAmount?: number | null
  upCount?: number | null
  downCount?: number | null
}

interface SectorFundFlowRankRow {
  code: string
  name: string
  price?: number | null
  changePercent?: number | null
  mainNetInflow?: number | null
  mainNetInflowPercent?: number | null
  topStockCode?: string | null
  topStockName?: string | null
}

interface FundFlowRankRow {
  code: string
  name: string
  price?: number | null
  mainNetInflow?: number | null
}

const router = useRouter()
const loading = ref(false)
const errorMessage = ref('')
const rankingTab = ref<DashboardRankingTab>('rise')
const boardTab = ref<'industry' | 'concept'>('industry')
const indices = ref<FullQuote[]>([])
const watchlistCodes = ref<string[]>([])
const watchlistQuotes = ref<FullQuote[]>([])
const marketQuotes = ref<FullQuote[]>([])
const industryList = ref<Board[]>([])
const conceptList = ref<Board[]>([])
const marketFundFlowHistory = ref<MarketFundFlowRow[]>([])
const northboundSummary = ref<NorthboundSummaryRow[]>([])
const industryFundFlowRanks = ref<SectorFundFlowRankRow[]>([])
const conceptFundFlowRanks = ref<SectorFundFlowRankRow[]>([])
const fundFlowRanks = ref<FundFlowRankRow[]>([])
const listRefreshInterval = getRefreshInterval('list')

const rankingOptions = [
  { label: '涨幅榜', value: 'rise' },
  { label: '跌幅榜', value: 'fall' },
  { label: '成交额', value: 'amount' },
  { label: '换手率', value: 'turnover' },
]
const boardOptions = [
  { label: '行业', value: 'industry' },
  { label: '概念', value: 'concept' },
]

const marketSummary = computed(() => summarizeMarketBreadth(marketQuotes.value))
const latestMarketFundFlow = computed(() => getLatestMarketFundFlow(marketFundFlowHistory.value))
const northboundSnapshot = computed(() => pickNorthboundSnapshot(northboundSummary.value))
const northboundNetInflow = computed(
  () => northboundSnapshot.value?.netInflow ?? northboundSnapshot.value?.netBuyAmount ?? null,
)
const currentBoards = computed(() => (boardTab.value === 'industry' ? industryList.value : conceptList.value))
const currentFundFlowBoards = computed(() =>
  boardTab.value === 'industry' ? industryFundFlowRanks.value : conceptFundFlowRanks.value,
)
const strongestBoard = computed(() => currentBoards.value[0] || null)
const rankingItems = computed(() => rankDashboardQuotes(marketQuotes.value, rankingTab.value).slice(0, 10))

function normalizeList<T>(value: unknown): T[] {
  return Array.isArray(value) ? (value as T[]) : []
}

async function load() {
  loading.value = true
  errorMessage.value = ''
  watchlistCodes.value = getAllWatchlistCodes()

  const failures: string[] = []
  await loadDashboardSections(
    [
      {
        name: '指数行情',
        load: () => getFullQuotes(MAIN_INDICES),
        commit: (value) => {
          indices.value = value as FullQuote[]
        },
      },
      {
        name: '自选行情',
        load: () => (watchlistCodes.value.length > 0 ? getFullQuotes(watchlistCodes.value.slice(0, 50)) : Promise.resolve([])),
        commit: (value) => {
          watchlistQuotes.value = value as FullQuote[]
        },
      },
      {
        name: '全市场行情',
        load: () => getAllAShareQuotes({ batchSize: 500, concurrency: 4 }),
        commit: (value) => {
          marketQuotes.value = value as FullQuote[]
        },
      },
      {
        name: '板块行情',
        load: () => Promise.all([getIndustryList(), getConceptList()]),
        commit: (value) => {
          const [industryRows, conceptRows] = value as [Board[], Board[]]
          industryList.value = industryRows
          conceptList.value = conceptRows
        },
      },
      {
        name: '大盘资金',
        load: () => getMarketFundFlow(),
        commit: (value) => {
          marketFundFlowHistory.value = normalizeList<MarketFundFlowRow>(value)
        },
      },
      {
        name: '北向资金',
        load: () => getNorthboundFlowSummary(),
        commit: (value) => {
          northboundSummary.value = normalizeList<NorthboundSummaryRow>(value)
        },
      },
      {
        name: '板块资金',
        load: () =>
          Promise.all([
            getSectorFundFlowRank({ indicator: 'today', sectorType: 'industry' }),
            getSectorFundFlowRank({ indicator: 'today', sectorType: 'concept' }),
          ]),
        commit: (value) => {
          const [industryRows, conceptRows] = value as [unknown, unknown]
          industryFundFlowRanks.value = normalizeList<SectorFundFlowRankRow>(industryRows)
          conceptFundFlowRanks.value = normalizeList<SectorFundFlowRankRow>(conceptRows)
        },
      },
      {
        name: '个股资金榜',
        load: () => getFundFlowRank({ indicator: 'today' }),
        commit: (value) => {
          fundFlowRanks.value = normalizeList<FundFlowRankRow>(value)
        },
      },
    ],
    (name, error) => {
      failures.push(name)
      console.warn(`Dashboard section failed: ${name}`, error)
    },
  )
  errorMessage.value = buildDashboardFailureMessage(failures)
  loading.value = false
}

function goStock(code: string) {
  router.push(`/s/${code}`)
}

function goBoard(code: string) {
  router.push(`/boards/${boardTab.value}/${code}`)
}

function rankingValue(item: FullQuote) {
  if (rankingTab.value === 'amount') return formatAmount(item.amount)
  if (rankingTab.value === 'turnover') return formatTurnover(item.turnoverRate)
  return formatPercent(item.changePercent)
}

function rankingValueClass(item: FullQuote) {
  if (rankingTab.value === 'amount' || rankingTab.value === 'turnover') return ''
  return getChangeColorClass(item.changePercent)
}

onMounted(load)
usePolling(load, { interval: listRefreshInterval, enabled: computed(() => !loading.value), immediate: false })
</script>

<style scoped>
.dashboard-page {
  gap: 12px;
}

.page-subtitle,
.muted,
.stat-meta {
  color: var(--color-muted);
  font-size: 12px;
}

.index-card {
  cursor: pointer;
}

.index-name {
  color: var(--color-muted);
  font-size: 13px;
}

.index-price {
  margin-top: 4px;
  font-size: 24px;
  font-weight: 700;
  line-height: 1.2;
}

.index-meta {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  margin: 6px 0 2px;
  font-size: 13px;
}

.stat-main {
  display: flex;
  align-items: baseline;
  gap: 8px;
  font-size: 26px;
  font-weight: 700;
  line-height: 1.25;
}

.stat-divider {
  color: var(--color-muted);
  font-weight: 500;
}

.stat-value {
  min-width: 0;
  font-size: 21px;
  font-weight: 700;
  line-height: 1.3;
}

.text-clip {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.clickable-row {
  cursor: pointer;
  background: transparent;
}

.clickable-row:hover {
  background: var(--color-hover);
}

.quote-side,
.board-side {
  min-width: 92px;
  text-align: right;
  font-variant-numeric: tabular-nums;
}

.board-side {
  min-width: 132px;
}

.board-side .muted {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.rank-num {
  display: inline-flex;
  flex: 0 0 28px;
  align-items: center;
  justify-content: center;
  width: 28px;
  color: var(--color-muted);
  font-variant-numeric: tabular-nums;
}

:deep(.ant-card) {
  border-radius: 6px;
}

:deep(.ant-card-body) {
  padding: 12px;
}

:deep(.ant-list-item) {
  padding: 8px 4px;
}

:deep(.ant-list-item-meta-title) {
  margin-bottom: 0;
}

@media (max-width: 640px) {
  .page-header {
    align-items: flex-start;
    flex-direction: column;
  }

  .quote-side,
  .board-side {
    min-width: 82px;
  }

  .board-side .muted {
    flex-direction: column;
    gap: 0;
  }
}
</style>
