<template>
  <a-layout-sider class="sidebar" :width="232">
    <div class="brand">
      <div class="brand-mark">策</div>
      <div>
        <div class="brand-name">策衡</div>
        <div class="brand-sub">A 股工作台</div>
      </div>
    </div>
    <a-menu mode="inline" :selected-keys="[activeKey]" class="nav-menu" @click="handleClick">
      <a-menu-item v-for="item in navItems" :key="item.path">
        <component :is="item.icon" />
        <span>{{ item.label }}</span>
      </a-menu-item>
    </a-menu>
  </a-layout-sider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AppstoreOutlined,
  BarChartOutlined,
  FireOutlined,
  FundProjectionScreenOutlined,
  SearchOutlined,
  SettingOutlined,
  StarOutlined,
  TableOutlined,
} from '@ant-design/icons-vue'

const emit = defineEmits<{ navigate: [] }>()
const router = useRouter()
const route = useRoute()

const navItems = [
  { path: '/', label: '总览', icon: AppstoreOutlined },
  { path: '/heatmap', label: '热力图', icon: FireOutlined },
  { path: '/rankings', label: '榜单', icon: BarChartOutlined },
  { path: '/boards', label: '板块', icon: TableOutlined },
  { path: '/watchlist', label: '自选', icon: StarOutlined },
  { path: '/scanner', label: '扫描', icon: SearchOutlined },
  { path: '/eod-picker', label: '尾盘选股', icon: FundProjectionScreenOutlined },
  { path: '/settings', label: '设置', icon: SettingOutlined },
]

const activeKey = computed(() => {
  const match = [...navItems].reverse().find((item) =>
    item.path === '/' ? route.path === '/' : route.path.startsWith(item.path),
  )
  return match?.path || '/'
})

function handleClick(event: { key: string }) {
  router.push(event.key)
  emit('navigate')
}
</script>

<style scoped>
.sidebar {
  min-height: 100vh;
  background: var(--color-sidebar-bg);
  border-inline-end: 1px solid var(--color-border);
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 64px;
  padding: 0 18px;
  color: var(--color-text-strong);
  border-bottom: 1px solid var(--color-border);
}

.brand-mark {
  display: grid;
  width: 36px;
  height: 36px;
  place-items: center;
  border-radius: 8px;
  color: #fff;
  background: var(--color-link);
  font-weight: 700;
}

.brand-name {
  font-size: 16px;
  font-weight: 700;
}

.brand-sub {
  color: var(--color-muted);
  font-size: 12px;
}

.nav-menu {
  padding: 8px;
  border-inline-end: 0;
  background: transparent;
  color: var(--color-menu-text);
}

.nav-menu :deep(.ant-menu-item) {
  width: 100%;
  margin-inline: 0;
  color: var(--color-menu-text);
  border-radius: 6px;
}

.nav-menu :deep(.ant-menu-item:hover) {
  color: var(--color-text-strong);
  background: var(--color-hover);
}

.nav-menu :deep(.ant-menu-item-selected) {
  color: var(--color-menu-active-text);
  background: var(--color-menu-active-bg);
}

.nav-menu :deep(.ant-menu-item-selected::after) {
  display: none;
}
</style>
