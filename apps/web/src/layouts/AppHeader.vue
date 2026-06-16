<template>
  <a-layout-header class="header">
    <a-button class="menu-button" type="text" @click="$emit('toggle-menu')">
      <MenuOutlined />
    </a-button>
    <a-auto-complete
      v-model:value="keyword"
      class="search"
      :options="options"
      allow-clear
      placeholder="搜索股票、板块..."
      @focus="showHistory"
      @search="handleSearch"
      @select="handleAutoCompleteSelect"
    >
      <template #option="slotProps">
        <div
          v-if="slotOption(slotProps).kind === 'history-title'"
          class="history-title"
          @mousedown.prevent
        >
          <span>最近搜索</span>
          <button class="clear-history" type="button" @click.stop="clearHistory">清空</button>
        </div>
        <div
          v-else
          class="search-option"
          @mousedown.prevent
          @click="handleSelectOption(slotOption(slotProps))"
        >
          <span class="option-label">{{ slotOption(slotProps).label }}</span>
          <span class="option-code">{{ slotOption(slotProps).value }}</span>
          <button
            v-if="slotOption(slotProps).entityType === 'stock'"
            class="quick-add"
            type="button"
            :disabled="isStockInWatchlist(slotOption(slotProps).value)"
            @click.stop="quickAdd(slotOption(slotProps))"
          >
            <CheckOutlined v-if="isStockInWatchlist(slotOption(slotProps).value)" />
            <StarOutlined v-else />
          </button>
        </div>
      </template>
    </a-auto-complete>
    <a-space class="header-actions" :size="6">
      <a-button class="theme-toggle" type="text" :title="themeMode === 'dark' ? '切换到浅色主题' : '切换到深色主题'" @click="toggleTheme">
        <BulbOutlined v-if="themeMode === 'dark'" />
        <EyeInvisibleOutlined v-else />
      </a-button>
      <a href="https://github.com/ceheng-io/stock-go" target="_blank" rel="noopener noreferrer" title="stock-go">
        <DatabaseOutlined />
      </a>
      <a href="https://github.com/ceheng-io/stock-go" target="_blank" rel="noopener noreferrer" title="GitHub">
        <GithubOutlined />
      </a>
    </a-space>
  </a-layout-header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  CheckOutlined,
  BulbOutlined,
  DatabaseOutlined,
  EyeInvisibleOutlined,
  GithubOutlined,
  MenuOutlined,
  StarOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { search } from '@/services/api'
import { addSearchHistory, addToWatchlist, clearSearchHistory, getSearchHistory, isInWatchlist } from '@/services/storage'
import { getThemeMode, toggleThemeMode, type ThemeMode } from '@/services/theme'
import type { SearchHistoryItem, SearchResult } from '@/types'
import { normalizeStockCode } from '@/utils/format'

defineEmits<{ 'toggle-menu': [] }>()

const router = useRouter()
const keyword = ref('')
const themeMode = ref<ThemeMode>(getThemeMode())
type HeaderOption = {
  value: string
  label: string
  route?: string
  market?: string
  type?: string
  entityType?: string
  kind?: 'result' | 'history' | 'history-title'
}
type HeaderOptionSlot = HeaderOption | { option?: HeaderOption }

const options = ref<HeaderOption[]>([])
const addedCodes = ref<Set<string>>(new Set())
let timer: number | undefined

function routeFor(item: { code: string; market: string; type: string }) {
  if (item.type === '行业板块') return `/boards/industry/${item.code}`
  if (item.type === '概念板块') return `/boards/concept/${item.code}`
  const code = normalizeStockCode(item.code)
  if (/^(sh|sz|bj)\d{6}$/i.test(code)) return `/s/${code}`
  return ''
}

function toResultOption(item: SearchResult): HeaderOption | null {
  const route = routeFor(item)
  if (!route) return null
  return {
    value: item.code,
    label: item.name,
    route,
    market: item.market,
    type: item.type,
    entityType: item.type === '股票' ? 'stock' : item.category,
    kind: 'result',
  }
}

function toHistoryOption(item: SearchHistoryItem): HeaderOption | null {
  const route = routeFor(item)
  if (!route) return null
  return {
    value: item.code,
    label: item.name,
    route,
    market: item.market,
    type: item.type,
    entityType: item.type === '股票' ? 'stock' : undefined,
    kind: 'history',
  }
}

function slotOption(slotProps: HeaderOptionSlot) {
  if ('option' in slotProps && slotProps.option) return slotProps.option
  return slotProps as HeaderOption
}

function selectedOption(value: string, option?: HeaderOptionSlot) {
  if (option) return slotOption(option)
  return options.value.find((item) => item.value === value)
}

function historyOptions() {
  const history = getSearchHistory().map(toHistoryOption).filter((item): item is HeaderOption => Boolean(item))
  if (history.length === 0) return []
  return [
    { value: '__history_title__', label: '最近搜索', kind: 'history-title' as const },
    ...history,
  ]
}

function showHistory() {
  if (keyword.value.trim()) return
  options.value = historyOptions()
}

function handleSearch(value: string) {
  window.clearTimeout(timer)
  if (!value.trim()) {
    options.value = historyOptions()
    return
  }
  timer = window.setTimeout(async () => {
    const results = await search(value)
    options.value = results.map(toResultOption).filter((item): item is HeaderOption => Boolean(item))
  }, 250)
}

function handleSelectOption(option: HeaderOption) {
  if (option.kind === 'history-title') return
  if (option.route) {
    if (option.kind !== 'history' && option.market && option.type) {
      addSearchHistory({
        code: option.value,
        name: option.label,
        market: option.market,
        type: option.type,
      })
    }
    keyword.value = ''
    options.value = []
    router.push(option.route)
  }
}

function handleAutoCompleteSelect(value: string, option?: HeaderOptionSlot) {
  const optionFromSelect = selectedOption(value, option)
  if (optionFromSelect) handleSelectOption(optionFromSelect)
}

function isStockInWatchlist(code: string) {
  return addedCodes.value.has(normalizeStockCode(code)) || isInWatchlist(code)
}

function quickAdd(option: HeaderOption) {
  const code = normalizeStockCode(option.value)
  if (!code || isStockInWatchlist(code)) return
  addToWatchlist(code)
  addedCodes.value = new Set([...addedCodes.value, code])
  message.success(`已将 ${option.label} 加入自选`)
}

function clearHistory() {
  clearSearchHistory()
  options.value = []
}

function toggleTheme() {
  themeMode.value = toggleThemeMode()
}
</script>

<style scoped>
.header {
  display: flex;
  align-items: center;
  gap: 12px;
  height: 64px;
  padding: 0 18px;
  background: var(--color-header-bg);
  border-bottom: 1px solid var(--color-border);
}

.search {
  max-width: 520px;
  flex: 1;
}

.search-option,
.history-title {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 8px;
}

.search-option {
  cursor: pointer;
}

.option-label {
  min-width: 0;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.option-code {
  color: var(--color-muted);
}

.quick-add,
.clear-history {
  border: 0;
  padding: 0;
  color: var(--color-link);
  background: transparent;
  cursor: pointer;
}

.quick-add:disabled {
  color: var(--color-muted);
  cursor: not-allowed;
}

.header-actions {
  margin-left: auto;
  flex: 0 0 auto;
}

.header-actions a {
  display: inline-flex;
  width: 32px;
  height: 32px;
  align-items: center;
  justify-content: center;
  color: var(--color-muted);
  border-radius: 6px;
}

.header-actions a:hover {
  color: var(--color-link);
  background: var(--color-hover);
}

.history-title {
  justify-content: space-between;
  color: var(--color-muted);
  font-size: 12px;
}

.menu-button {
  display: none;
}

@media (max-width: 768px) {
  .menu-button {
    display: inline-flex;
  }
}
</style>
