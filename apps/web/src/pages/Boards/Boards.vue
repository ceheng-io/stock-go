<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">板块</h1>
      <a-space>
        <a-segmented v-model:value="boardType" :options="[{ label: '行业', value: 'industry' }, { label: '概念', value: 'concept' }]" />
        <a-input-search v-model:value="keyword" placeholder="搜索板块" allow-clear />
      </a-space>
    </div>
    <a-card size="small">
      <a-table :columns="columns" :data-source="filteredRows" row-key="code" size="small" :pagination="{ pageSize: 30 }" :custom-row="rowClick">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'changePercent'">
            <span :class="getChangeColorClass(record.changePercent)">
              {{ formatPercent(record.changePercent) }}
            </span>
          </template>
          <template v-else-if="column.key === 'breadth'">
            <span class="text-rise">{{ record.riseCount ?? '--' }}</span>
            <span class="muted"> / </span>
            <span class="text-fall">{{ record.fallCount ?? '--' }}</span>
          </template>
          <template v-else-if="column.key === 'leader'">
            <div>{{ record.leadingStock || '--' }}</div>
            <div class="muted" :class="getChangeColorClass(record.leadingStockChangePercent)">
              {{ formatPercent(record.leadingStockChangePercent) }}
            </div>
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
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getConceptList, getIndustryList } from '@/services/api'
import type { Board } from '@/types'
import { formatMarketCap, formatPercent, formatTurnover, getChangeColorClass } from '@/utils/format'

const router = useRouter()
const boardType = ref<'industry' | 'concept'>('industry')
const keyword = ref('')
const rows = ref<Board[]>([])

const columns = [
  { title: '排名', dataIndex: 'rank', width: 80 },
  { title: '名称', dataIndex: 'name' },
  { title: '涨跌幅', key: 'changePercent' },
  { title: '涨/跌', key: 'breadth' },
  { title: '总市值', key: 'totalMarketCap' },
  { title: '换手', key: 'turnoverRate' },
  { title: '领涨股', key: 'leader' },
]

async function load() {
  rows.value = boardType.value === 'industry' ? await getIndustryList() : await getConceptList()
}

const filteredRows = computed(() => {
  const value = keyword.value.trim().toLowerCase()
  if (!value) return rows.value
  return rows.value.filter((item) => item.name.toLowerCase().includes(value) || item.leadingStock?.toLowerCase().includes(value))
})

function rowClick(record: Board) {
  return {
    onClick: () => router.push(`/boards/${boardType.value}/${record.code}`),
  }
}

watch(boardType, load)
onMounted(load)
</script>

<style scoped>
.muted {
  color: var(--color-muted);
  font-size: 12px;
}

:deep(.ant-table-row) {
  cursor: pointer;
}
</style>
