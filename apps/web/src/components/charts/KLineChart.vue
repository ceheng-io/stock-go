<template>
  <div class="kline-chart-shell">
    <div ref="container" class="kline-chart" />
    <div v-if="chartRows.length === 0" class="kline-empty">{{ emptyText }}</div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { dispose, init, type Chart, type DeepPartial, type Styles } from 'klinecharts'
import type { IndicatorConfig } from '@/types'
import type { OscillatorIndicatorKey, OverlayIndicatorKey } from '@/services/charts'
import { buildKLineIndicatorNames, toKLineChartData, type KLineChartSourceRow } from '@/services/klineChart'

const props = withDefaults(defineProps<{
  rows: KLineChartSourceRow[]
  overlays?: OverlayIndicatorKey[]
  oscillator?: OscillatorIndicatorKey
  indicatorConfig?: IndicatorConfig | null
  emptyText?: string
}>(), {
  overlays: () => [],
  oscillator: 'macd',
  indicatorConfig: null,
  emptyText: '暂无 K 线数据',
})

const container = ref<HTMLElement | null>(null)
let chart: Chart | null = null
let resizeObserver: ResizeObserver | null = null

const chartRows = computed(() => toKLineChartData(props.rows))

const defaultIndicatorConfig: IndicatorConfig = {
  ma: [5, 10, 20, 60],
  macd: { short: 12, long: 26, signal: 9 },
  boll: { period: 20, stdDev: 2 },
  kdj: { period: 9, kPeriod: 3, dPeriod: 3 },
  rsi: [6, 12, 24],
  dmi: { period: 14, adxPeriod: 14 },
  sar: { afStart: 0.02, afIncrement: 0.02, afMax: 0.2 },
  kc: { emaPeriod: 20, atrPeriod: 10, multiplier: 2 },
}

const chartStyles: DeepPartial<Styles> = {
  candle: {
    bar: {
      upColor: '#cf1322',
      downColor: '#389e0d',
      noChangeColor: '#8c8c8c',
      upBorderColor: '#cf1322',
      downBorderColor: '#389e0d',
      noChangeBorderColor: '#8c8c8c',
      upWickColor: '#cf1322',
      downWickColor: '#389e0d',
      noChangeWickColor: '#8c8c8c',
    },
  },
  grid: {
    horizontal: { color: 'rgba(148, 163, 184, 0.22)' },
    vertical: { color: 'rgba(148, 163, 184, 0.16)' },
  },
}

function ensureChart() {
  if (chart || !container.value) return
  chart = init(container.value, {
    locale: 'zh-CN',
    timezone: 'Asia/Shanghai',
    styles: chartStyles,
  })
  chart?.setPriceVolumePrecision(2, 0)
  resizeObserver = new ResizeObserver(() => chart?.resize())
  resizeObserver.observe(container.value)
}

function applyData() {
  ensureChart()
  chart?.applyNewData(chartRows.value)
}

function applyIndicators() {
  ensureChart()
  if (!chart) return
  chart.removeIndicator('candle_pane')
  chart.removeIndicator('volume_pane')
  chart.removeIndicator('indicator_pane')
  const names = buildKLineIndicatorNames({
    overlays: props.overlays,
    oscillator: props.oscillator,
    indicatorConfig: props.indicatorConfig || defaultIndicatorConfig,
  })
  names.overlays.forEach((name) => chart?.createIndicator(name, true, { id: 'candle_pane' }))
  chart.createIndicator(names.panes[0], false, { id: 'volume_pane', height: 72 })
  chart.createIndicator(names.panes[1], false, { id: 'indicator_pane', height: 92 })
}

onMounted(async () => {
  await nextTick()
  ensureChart()
  applyData()
  applyIndicators()
})

watch(chartRows, applyData)
watch(
  () => [props.overlays.join(','), props.oscillator, props.indicatorConfig],
  applyIndicators,
)

onBeforeUnmount(() => {
  resizeObserver?.disconnect()
  resizeObserver = null
  if (chart) {
    dispose(chart)
    chart = null
  }
})
</script>

<style scoped>
.kline-chart-shell {
  position: relative;
  min-width: 0;
}

.kline-chart {
  width: 100%;
  height: 520px;
}

.kline-empty {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-muted);
  font-size: 14px;
  pointer-events: none;
}
</style>
