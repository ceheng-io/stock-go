<template>
  <div class="page watchlist-page">
    <div class="page-header">
      <div>
        <h1 class="page-title">自选</h1>
        <div class="page-subtitle">{{ activeGroup?.name || '默认分组' }} · {{ activeCodes.length }} 只</div>
      </div>
      <a-space wrap>
        <a-input-search
          v-model:value="newCode"
          placeholder="输入代码加入当前分组"
          enter-button="添加"
          style="width: 260px"
          @search="addCode"
        />
        <a-button :loading="loading" @click="refreshQuotes">刷新</a-button>
        <a-button @click="showImport = true">导入</a-button>
        <a-button @click="exportCodes">导出</a-button>
        <a-button @click="showColumns = true">列设置</a-button>
      </a-space>
    </div>

    <a-alert v-if="error" type="error" show-icon :message="error" />

    <div class="watchlist-shell">
      <a-card class="group-panel" size="small" title="分组">
        <a-list :data-source="groups" size="small">
          <template #renderItem="{ item }">
            <a-list-item
              :class="{ active: item.id === activeGroupId }"
              class="group-row"
              @click="selectGroup(item.id)"
            >
              <template v-if="editingGroupId === item.id">
                <a-input
                  v-model:value="editingGroupName"
                  size="small"
                  @press-enter="commitRename"
                  @blur="commitRename"
                  @click.stop
                />
              </template>
              <template v-else>
                <span class="group-name">{{ item.name }}</span>
                <a-badge :count="item.codes.length" :number-style="{ backgroundColor: '#8c8c8c' }" />
                <a-dropdown v-if="item.id !== 'default'" trigger="click">
                  <a-button type="text" size="small" @click.stop>...</a-button>
                  <template #overlay>
                    <a-menu>
                      <a-menu-item @click="startRename(item)">重命名</a-menu-item>
                      <a-menu-item danger @click="removeGroup(item.id)">删除分组</a-menu-item>
                    </a-menu>
                  </template>
                </a-dropdown>
              </template>
            </a-list-item>
          </template>
        </a-list>
        <a-input-search
          v-model:value="newGroupName"
          placeholder="新建分组"
          enter-button="添加"
          size="small"
          class="new-group"
          @search="createGroup"
        />
      </a-card>

      <a-card class="table-panel" size="small">
        <template #title>
          <a-space wrap>
            <span>{{ activeGroup?.name || '自选股' }}</span>
            <a-tag v-if="triggeredAlerts.length" color="red">{{ triggeredAlerts.length }} 条预警</a-tag>
          </a-space>
        </template>
        <template #extra>
          <a-space wrap>
            <a-select v-model:value="sortField" size="small" style="width: 116px">
              <a-select-option value="default">默认排序</a-select-option>
              <a-select-option value="changePercent">涨跌幅</a-select-option>
              <a-select-option value="amount">成交额</a-select-option>
              <a-select-option value="turnoverRate">换手率</a-select-option>
              <a-select-option value="totalMarketCap">总市值</a-select-option>
            </a-select>
            <a-segmented v-model:value="sortOrder" :options="sortOrderOptions" size="small" />
            <a-button danger size="small" :disabled="selectedRowKeys.length === 0" @click="batchDelete">
              删除 {{ selectedRowKeys.length || '' }}
            </a-button>
            <a-button size="small" @click="showAlerts = true">告警中心</a-button>
          </a-space>
        </template>

        <a-empty v-if="activeCodes.length === 0" description="当前分组暂无自选股">
          <a-button type="primary" @click="showImport = true">导入股票</a-button>
        </a-empty>
        <a-table
          v-else
          :columns="tableColumns"
          :data-source="sortedQuotes"
          :loading="loading"
          :row-selection="rowSelection"
          row-key="code"
          size="small"
          :pagination="{ pageSize: 30 }"
          :custom-row="rowClick"
        />
      </a-card>
    </div>

    <a-drawer v-model:open="showImport" title="导入自选" width="420">
      <a-textarea v-model:value="importText" :rows="8" placeholder="支持逗号、空格或换行分隔，例如：600519 000001" />
      <template #footer>
        <a-space>
          <a-button @click="showImport = false">取消</a-button>
          <a-button type="primary" @click="importCodes">导入</a-button>
        </a-space>
      </template>
    </a-drawer>

    <a-drawer v-model:open="showColumns" title="列设置" width="360">
      <a-checkbox-group v-model:value="visibleColumnKeys" class="column-options">
        <a-checkbox v-for="column in defaultColumns" :key="column.key" :value="column.key">
          {{ column.label }}
        </a-checkbox>
      </a-checkbox-group>
      <template #footer>
        <a-space>
          <a-button @click="showColumns = false">取消</a-button>
          <a-button type="primary" @click="saveColumns">保存</a-button>
        </a-space>
      </template>
    </a-drawer>

    <a-drawer v-model:open="showAlerts" title="告警中心" width="480">
      <a-list :data-source="groupAlerts" size="small">
        <template #renderItem="{ item }">
          <a-list-item>
            <a-list-item-meta :title="item.name" :description="formatAlertRule(item)" />
            <a-space>
              <a-switch :checked="item.enabled" size="small" @change="toggleAlert(item.id, Boolean($event))" />
              <a-button danger type="link" size="small" @click="removeAlert(item.id)">删除</a-button>
            </a-space>
          </a-list-item>
        </template>
        <template #empty>
          <a-empty description="当前分组暂无预警" />
        </template>
      </a-list>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { computed, h, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import type { TableColumnsType, TableProps } from 'ant-design-vue'
import { usePolling } from '@/composables/usePolling'
import { getAllQuotesByCodes } from '@/services/api'
import {
  batchAddToWatchlist,
  batchRemoveFromWatchlist,
  createWatchlistGroup,
  deleteAlertRule,
  deleteWatchlistGroup,
  getAlertRules,
  getSettings,
  getTableColumns,
  getWatchlistGroups,
  renameWatchlistGroup,
  saveTableColumns,
  updateAlertRule,
  addToWatchlist,
  reorderWatchlist,
} from '@/services/storage'
import { moveWatchlistCodeBefore } from '@/services/watchlist'
import type { AlertRule, ColumnConfig, FullQuote, WatchlistGroup } from '@/types'
import {
  formatAmount,
  formatMarketCap,
  formatPercent,
  formatPrice,
  formatTurnover,
  normalizeStockCode,
} from '@/utils/format'

type SortField = 'default' | 'changePercent' | 'amount' | 'turnoverRate' | 'totalMarketCap'
type QuoteSortField = Exclude<SortField, 'default'>
type SortOrder = 'asc' | 'desc'

const defaultColumns: ColumnConfig[] = [
  { key: 'name', label: '名称/代码', visible: true },
  { key: 'price', label: '现价', visible: true },
  { key: 'change', label: '涨跌幅', visible: true },
  { key: 'amount', label: '成交额', visible: true },
  { key: 'turnover', label: '换手', visible: true },
  { key: 'marketCap', label: '总市值', visible: false },
]

const router = useRouter()
const groups = ref<WatchlistGroup[]>([])
const activeGroupId = ref('default')
const quotes = ref<FullQuote[]>([])
const loading = ref(false)
const error = ref('')
const newCode = ref('')
const newGroupName = ref('')
const importText = ref('')
const showImport = ref(false)
const showColumns = ref(false)
const showAlerts = ref(false)
const editingGroupId = ref('')
const editingGroupName = ref('')
const selectedRowKeys = ref<string[]>([])
const sortField = ref<SortField>('default')
const sortOrder = ref<SortOrder>('desc')
const visibleColumnKeys = ref<string[]>([])
const draggedCode = ref('')
const listRefreshInterval = getSettings().refreshInterval.list
const sortOrderOptions = [
  { label: '降序', value: 'desc' },
  { label: '升序', value: 'asc' },
]

const activeGroup = computed(() => groups.value.find((group) => group.id === activeGroupId.value) || groups.value[0])
const activeCodes = computed(() => activeGroup.value?.codes || [])
const activeCodeSet = computed(() => new Set(activeCodes.value.map(normalizeStockCode)))
const alertRules = computed(() => getAlertRules())
const groupAlerts = computed(() => alertRules.value.filter((rule) => activeCodeSet.value.has(normalizeStockCode(rule.code))))
const triggeredAlerts = computed(() => groupAlerts.value.filter((rule) => rule.enabled && quotes.value.some((quote) => quoteMatchesAlert(rule, quote))))

const sortedQuotes = computed(() => {
  const quoteMap = new Map(quotes.value.map((quote) => [normalizeStockCode(quote.code), quote]))
  const rows = activeCodes.value.map((code) => quoteMap.get(normalizeStockCode(code))).filter((item): item is FullQuote => Boolean(item))
  if (sortField.value === 'default') return rows
  const field = sortField.value as QuoteSortField
  return [...rows].sort((a, b) => {
    const left = Number(a[field] || 0)
    const right = Number(b[field] || 0)
    return sortOrder.value === 'desc' ? right - left : left - right
  })
})

const tableColumns = computed<TableColumnsType<FullQuote>>(() => {
  const columns: TableColumnsType<FullQuote> = []
  if (sortField.value === 'default') {
    columns.push({
      title: '',
      key: 'drag',
      width: 42,
      customRender: () => '⋮⋮',
    })
  }
  if (visibleColumnKeys.value.includes('name')) {
    columns.push({
      title: '名称/代码',
      dataIndex: 'name',
      width: 170,
      customRender: ({ record }) => (
        `${record.name}\n${record.code}`
      ),
    })
  }
  if (visibleColumnKeys.value.includes('price')) {
    columns.push({ title: '现价', customRender: ({ record }) => formatPrice(record.price), align: 'right' })
  }
  if (visibleColumnKeys.value.includes('change')) {
    columns.push({
      title: '涨跌幅',
      align: 'right',
      customRender: ({ record }) => formatPercent(record.changePercent),
    })
  }
  if (visibleColumnKeys.value.includes('amount')) {
    columns.push({ title: '成交额', align: 'right', customRender: ({ record }) => formatAmount(record.amount) })
  }
  if (visibleColumnKeys.value.includes('turnover')) {
    columns.push({ title: '换手', align: 'right', customRender: ({ record }) => formatTurnover(record.turnoverRate) })
  }
  if (visibleColumnKeys.value.includes('marketCap')) {
    columns.push({ title: '总市值', align: 'right', customRender: ({ record }) => formatMarketCap(record.totalMarketCap) })
  }
  columns.push({
    title: '状态',
    width: 96,
    customRender: ({ record }) => {
      const alerts = groupAlerts.value.filter((rule) => rule.enabled && quoteMatchesAlert(rule, record))
      return alerts.length > 0 ? '预警' : getChangeLabel(record.changePercent)
    },
  })
  columns.push({
    title: '操作',
    width: 96,
    customRender: ({ record }) => h(
      'button',
      {
        class: 'link-danger',
        onClick: (event: MouseEvent) => {
          event.stopPropagation()
          removeStock(record.code)
        },
      },
      '删除',
    ),
  })
  return columns
})

const rowSelection = computed<TableProps['rowSelection']>(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys) => {
    selectedRowKeys.value = keys.map(String)
  },
}))

function reloadGroups() {
  groups.value = getWatchlistGroups()
  if (!groups.value.some((group) => group.id === activeGroupId.value)) {
    activeGroupId.value = groups.value[0]?.id || 'default'
  }
}

function loadColumns() {
  const stored = getTableColumns('watchlist')
  const columns = stored || defaultColumns
  visibleColumnKeys.value = columns.filter((column) => column.visible).map((column) => column.key)
}

async function refreshQuotes() {
  if (activeCodes.value.length === 0) {
    quotes.value = []
    return
  }
  loading.value = true
  error.value = ''
  try {
    quotes.value = await getAllQuotesByCodes(activeCodes.value) as FullQuote[]
    announceTriggeredAlerts()
  } catch (err) {
    error.value = err instanceof Error ? err.message : '加载自选行情失败'
  } finally {
    loading.value = false
  }
}

function selectGroup(groupId: string) {
  activeGroupId.value = groupId
  selectedRowKeys.value = []
}

function createGroup() {
  if (!newGroupName.value.trim()) return
  const group = createWatchlistGroup(newGroupName.value)
  activeGroupId.value = group.id
  newGroupName.value = ''
  reloadGroups()
}

function startRename(group: WatchlistGroup) {
  editingGroupId.value = group.id
  editingGroupName.value = group.name
}

function commitRename() {
  if (!editingGroupId.value) return
  renameWatchlistGroup(editingGroupId.value, editingGroupName.value)
  editingGroupId.value = ''
  editingGroupName.value = ''
  reloadGroups()
}

function removeGroup(groupId: string) {
  Modal.confirm({
    title: '删除分组',
    content: '分组内的股票会从该分组移除。',
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    onOk: () => {
      deleteWatchlistGroup(groupId)
      reloadGroups()
      refreshQuotes()
    },
  })
}

async function addCode() {
  const code = normalizeStockCode(newCode.value)
  if (!/^(sh|sz|bj)\d{6}$/i.test(code)) {
    message.warning('请输入有效 A 股代码')
    return
  }
  addToWatchlist(code, activeGroupId.value)
  newCode.value = ''
  reloadGroups()
  await refreshQuotes()
}

function importCodes() {
  const codes = importText.value.split(/[\s,;，；]+/).map(normalizeStockCode).filter((code) => /^(sh|sz|bj)\d{6}$/i.test(code))
  if (codes.length === 0) {
    message.warning('未识别到有效 A 股代码')
    return
  }
  const before = activeCodes.value.length
  batchAddToWatchlist(codes, activeGroupId.value)
  reloadGroups()
  importText.value = ''
  showImport.value = false
  message.success(`已导入 ${Math.max(activeCodes.value.length - before, 0)} 只新股票`)
  refreshQuotes()
}

function exportCodes() {
  const text = activeCodes.value.join('\n')
  navigator.clipboard?.writeText(text)
    .then(() => message.success('已复制到剪贴板'))
    .catch(() => message.info(text || '当前分组为空'))
}

function batchDelete() {
  if (selectedRowKeys.value.length === 0) return
  batchRemoveFromWatchlist(selectedRowKeys.value, activeGroupId.value)
  selectedRowKeys.value = []
  reloadGroups()
  refreshQuotes()
}

function removeStock(code: string) {
  batchRemoveFromWatchlist([code], activeGroupId.value)
  selectedRowKeys.value = selectedRowKeys.value.filter((item) => item !== code)
  reloadGroups()
  refreshQuotes()
}

function saveColumns() {
  const columns = defaultColumns.map((column) => ({ ...column, visible: visibleColumnKeys.value.includes(column.key) }))
  saveTableColumns('watchlist', columns)
  showColumns.value = false
  message.success('列设置已保存')
}

function removeAlert(id: string) {
  deleteAlertRule(id)
  message.success('已删除预警')
}

function toggleAlert(id: string, enabled: boolean) {
  updateAlertRule(id, { enabled })
}

function rowClick(record: FullQuote) {
  return {
    draggable: sortField.value === 'default',
    onDragstart: () => {
      draggedCode.value = normalizeStockCode(record.code)
    },
    onDragover: (event: DragEvent) => {
      if (sortField.value !== 'default' || !draggedCode.value) return
      event.preventDefault()
    },
    onDrop: (event: DragEvent) => {
      event.preventDefault()
      if (sortField.value !== 'default' || !draggedCode.value) return
      const nextCodes = moveWatchlistCodeBefore(activeCodes.value, draggedCode.value, record.code)
      reorderWatchlist(activeGroupId.value, nextCodes)
      draggedCode.value = ''
      reloadGroups()
      refreshQuotes()
    },
    onDragend: () => {
      draggedCode.value = ''
    },
    onClick: (event: MouseEvent) => {
      const target = event.target as HTMLElement
      if (target.closest('button') || target.closest('.ant-checkbox-wrapper')) return
      router.push(`/s/${record.code}`)
    },
  }
}

function quoteMatchesAlert(rule: AlertRule, quote: FullQuote) {
  if (normalizeStockCode(rule.code) !== normalizeStockCode(quote.code)) return false
  switch (rule.type) {
    case 'price_gte':
      return quote.price >= rule.value
    case 'price_lte':
      return quote.price <= rule.value
    case 'change_percent_gte':
      return quote.changePercent >= rule.value
    case 'change_percent_lte':
      return quote.changePercent <= rule.value
    case 'amount_gte':
      return quote.amount >= rule.value
    case 'near_limit_up':
      return quote.limitUp !== null && quote.limitUp !== undefined && quote.price >= quote.limitUp * (1 - rule.value / 100)
    case 'near_limit_down':
      return quote.limitDown !== null && quote.limitDown !== undefined && quote.price <= quote.limitDown * (1 + rule.value / 100)
    default:
      return false
  }
}

function formatAlertRule(rule: AlertRule) {
  const labels: Record<string, string> = {
    price_gte: '价格 >=',
    price_lte: '价格 <=',
    change_percent_gte: '涨幅 >=',
    change_percent_lte: '涨幅 <=',
    amount_gte: '成交额 >=',
    near_limit_up: '接近涨停',
    near_limit_down: '接近跌停',
  }
  const suffix = rule.type.includes('percent') || rule.type.startsWith('near_') ? '%' : ''
  return `${labels[rule.type] || rule.type} ${rule.value}${suffix}`
}

function getChangeLabel(value?: number | null) {
  if (value === null || value === undefined || Number.isNaN(value) || value === 0) return '平'
  return value > 0 ? '上涨' : '下跌'
}

function announceTriggeredAlerts() {
  triggeredAlerts.value.slice(0, 3).forEach((rule) => {
    const now = Date.now()
    if (now - rule.lastTriggeredAt < rule.cooldownSec * 1000) return
    updateAlertRule(rule.id, { lastTriggeredAt: now })
    message.warning(`${rule.name} 触发预警：${formatAlertRule(rule)}`)
  })
}

watch(activeGroupId, refreshQuotes)
usePolling(refreshQuotes, {
  interval: listRefreshInterval,
  enabled: computed(() => activeCodes.value.length > 0 && !loading.value),
  immediate: false,
})

onMounted(async () => {
  reloadGroups()
  loadColumns()
  await refreshQuotes()
})
</script>

<style scoped>
.watchlist-page {
  min-width: 0;
}

.page-subtitle {
  margin-top: 4px;
  color: var(--color-text-secondary);
  font-size: 13px;
}

.watchlist-shell {
  display: grid;
  grid-template-columns: 260px minmax(0, 1fr);
  gap: 12px;
  align-items: start;
}

.group-panel,
.table-panel {
  min-width: 0;
}

.group-row {
  cursor: pointer;
  border-radius: 6px;
  padding-inline: 8px;
}

.group-row.active {
  background: #e6f4ff;
}

.group-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.new-group {
  margin-top: 12px;
}

.column-options {
  display: grid;
  gap: 10px;
}

:deep(.ant-table-row) {
  cursor: pointer;
}

:deep(.ant-table-row[draggable='true']) {
  cursor: grab;
}

:deep(.link-danger) {
  border: 0;
  padding: 0;
  color: #cf1322;
  background: transparent;
  cursor: pointer;
}

@media (max-width: 900px) {
  .watchlist-shell {
    grid-template-columns: 1fr;
  }
}
</style>
