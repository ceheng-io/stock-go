<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">榜单</h1>
        <div class="page-subtitle">行业与概念强弱、领涨股和资金容量对照</div>
      </div>
      <a-space>
        <a-segmented v-model:value="rankType" :options="rankOptions" />
        <a-button :loading="loading" @click="load">刷新</a-button>
      </a-space>
    </div>

    <a-alert v-if="errorMessage" type="warning" show-icon closable :message="errorMessage" @close="errorMessage = ''" />

    <a-row :gutter="[12, 12]">
      <a-col :xs="24" :lg="12">
        <a-card title="行业板块" size="small">
          <a-table
            :columns="columns"
            :data-source="sortedIndustry"
            :pagination="{ pageSize: 20 }"
            :loading="loading && industry.length === 0"
            size="small"
            row-key="code"
            @row="rowTo('industry')"
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
      </a-col>
      <a-col :xs="24" :lg="12">
        <a-card title="概念板块" size="small">
          <a-table
            :columns="columns"
            :data-source="sortedConcept"
            :pagination="{ pageSize: 20 }"
            :loading="loading && concept.length === 0"
            size="small"
            row-key="code"
            @row="rowTo('concept')"
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
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getConceptList, getIndustryList } from '@/services/api'
import { sortBoardRankings, type BoardRankingType } from '@/services/rankings'
import type { Board } from '@/types'
import { formatMarketCap, formatPercent, formatTurnover, getChangeColorClass } from '@/utils/format'

const router = useRouter()
const rankType = ref<BoardRankingType>('rise')
const industry = ref<Board[]>([])
const concept = ref<Board[]>([])
const loading = ref(false)
const errorMessage = ref('')
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

const sortedIndustry = computed(() => sortBoardRankings(industry.value, rankType.value, 50))
const sortedConcept = computed(() => sortBoardRankings(concept.value, rankType.value, 50))

async function load() {
  loading.value = true
  errorMessage.value = ''
  try {
    ;[industry.value, concept.value] = await Promise.all([getIndustryList(), getConceptList()])
  } catch (error) {
    console.warn('Rankings data loading failed', error)
    errorMessage.value = '榜单数据加载失败，请稍后刷新'
  } finally {
    loading.value = false
  }
}

function rowTo(type: 'industry' | 'concept') {
  return (record: Board) => ({
    onClick: () => router.push(`/boards/${type}/${record.code}`),
  })
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
}
</style>
