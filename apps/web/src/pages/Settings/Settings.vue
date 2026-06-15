<template>
  <div class="page settings-page">
    <div class="page-header">
      <h1 class="page-title">设置</h1>
      <a-space>
        <a-button @click="reset">恢复默认</a-button>
        <a-button type="primary" @click="save">保存</a-button>
      </a-space>
    </div>

    <a-row :gutter="[12, 12]">
      <a-col :xs="24" :lg="12">
        <a-card title="刷新频率" size="small">
          <a-form layout="vertical">
            <a-form-item label="列表 / 自选">
              <a-select v-model:value="settings.refreshInterval.list">
                <a-select-option :value="0">默认</a-select-option>
                <a-select-option :value="5000">5 秒</a-select-option>
                <a-select-option :value="10000">10 秒</a-select-option>
                <a-select-option :value="15000">15 秒</a-select-option>
                <a-select-option :value="30000">30 秒</a-select-option>
              </a-select>
            </a-form-item>
            <a-form-item label="个股详情">
              <a-select v-model:value="settings.refreshInterval.detail">
                <a-select-option :value="5000">5 秒</a-select-option>
                <a-select-option :value="10000">10 秒</a-select-option>
                <a-select-option :value="15000">15 秒</a-select-option>
                <a-select-option :value="30000">30 秒</a-select-option>
              </a-select>
            </a-form-item>
            <a-form-item label="热力图">
              <a-select v-model:value="settings.refreshInterval.heatmap">
                <a-select-option :value="5000">5 秒</a-select-option>
                <a-select-option :value="10000">10 秒</a-select-option>
                <a-select-option :value="15000">15 秒</a-select-option>
                <a-select-option :value="30000">30 秒</a-select-option>
              </a-select>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>

      <a-col :xs="24" :lg="12">
        <a-card title="显示偏好" size="small">
          <a-form layout="vertical">
            <a-form-item label="涨跌颜色">
              <a-radio-group v-model:value="settings.colorMode" button-style="solid">
                <a-radio-button value="red-rise">红涨绿跌</a-radio-button>
                <a-radio-button value="green-rise">绿涨红跌</a-radio-button>
              </a-radio-group>
            </a-form-item>
            <a-form-item label="热力图维度">
              <a-select v-model:value="settings.heatmapConfig.dimension">
                <a-select-option value="industry">行业</a-select-option>
                <a-select-option value="concept">概念</a-select-option>
                <a-select-option value="watchlist">自选</a-select-option>
              </a-select>
            </a-form-item>
            <a-form-item label="热力图数量">
              <a-input-number v-model:value="settings.heatmapConfig.topK" :min="20" :max="500" :step="20" />
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>

      <a-col :xs="24">
        <a-card title="指标参数" size="small">
          <a-row :gutter="[12, 12]">
            <a-col :xs="24" :md="8">
              <a-form-item label="MA 周期">
                <a-input v-model:value="maDraft" placeholder="5, 10, 20, 60" />
              </a-form-item>
            </a-col>
            <a-col :xs="24" :md="8">
              <a-form-item label="RSI 周期">
                <a-input v-model:value="rsiDraft" placeholder="6, 12, 24" />
              </a-form-item>
            </a-col>
            <a-col :xs="24" :md="8">
              <a-form-item label="BOLL">
                <a-space>
                  <a-input-number v-model:value="settings.indicatorConfig.boll.period" :min="2" />
                  <a-input-number v-model:value="settings.indicatorConfig.boll.stdDev" :min="0.1" :step="0.1" />
                </a-space>
              </a-form-item>
            </a-col>
            <a-col :xs="24" :md="8">
              <a-form-item label="MACD">
                <a-space>
                  <a-input-number v-model:value="settings.indicatorConfig.macd.short" :min="1" />
                  <a-input-number v-model:value="settings.indicatorConfig.macd.long" :min="1" />
                  <a-input-number v-model:value="settings.indicatorConfig.macd.signal" :min="1" />
                </a-space>
              </a-form-item>
            </a-col>
            <a-col :xs="24" :md="8">
              <a-form-item label="KDJ">
                <a-space>
                  <a-input-number v-model:value="settings.indicatorConfig.kdj.period" :min="1" />
                  <a-input-number v-model:value="settings.indicatorConfig.kdj.kPeriod" :min="1" />
                  <a-input-number v-model:value="settings.indicatorConfig.kdj.dPeriod" :min="1" />
                </a-space>
              </a-form-item>
            </a-col>
            <a-col :xs="24" :md="8">
              <a-form-item label="DMI / ADX">
                <a-space>
                  <a-input-number v-model:value="settings.indicatorConfig.dmi.period" :min="1" />
                  <a-input-number v-model:value="settings.indicatorConfig.dmi.adxPeriod" :min="1" />
                </a-space>
              </a-form-item>
            </a-col>
            <a-col :xs="24" :md="8">
              <a-form-item label="SAR">
                <a-space>
                  <a-input-number v-model:value="settings.indicatorConfig.sar.afStart" :min="0.01" :step="0.01" />
                  <a-input-number v-model:value="settings.indicatorConfig.sar.afIncrement" :min="0.01" :step="0.01" />
                  <a-input-number v-model:value="settings.indicatorConfig.sar.afMax" :min="0.01" :step="0.01" />
                </a-space>
              </a-form-item>
            </a-col>
            <a-col :xs="24" :md="8">
              <a-form-item label="KC">
                <a-space>
                  <a-input-number v-model:value="settings.indicatorConfig.kc.emaPeriod" :min="1" />
                  <a-input-number v-model:value="settings.indicatorConfig.kc.atrPeriod" :min="1" />
                  <a-input-number v-model:value="settings.indicatorConfig.kc.multiplier" :min="0.1" :step="0.1" />
                </a-space>
              </a-form-item>
            </a-col>
          </a-row>
        </a-card>
      </a-col>

      <a-col :xs="24">
        <a-card title="数据源" size="small">
          <a-descriptions :column="{ xs: 1, md: 2 }" size="small">
            <a-descriptions-item label="应用">策衡 A 股看板</a-descriptions-item>
            <a-descriptions-item label="前端">Vue 3 + Ant Design Vue + ECharts</a-descriptions-item>
            <a-descriptions-item label="后端">apps/server 调用当前 Go SDK</a-descriptions-item>
            <a-descriptions-item label="API 前缀">/api</a-descriptions-item>
            <a-descriptions-item label="本地数据">自选、预警和显示设置保存在浏览器 localStorage</a-descriptions-item>
          </a-descriptions>
          <a-divider />
          <a-typography-text strong>数据说明：</a-typography-text>
          <ul class="note-list">
            <li>成交量单位：手（1 手 = 100 股）</li>
            <li>成交额单位：万元</li>
            <li>资金流、北向、龙虎榜等数据默认使用元级展示</li>
            <li>市值单位：亿元</li>
            <li>详情页以 A 股为主；港股、美股、基金搜索结果不跳转详情页</li>
          </ul>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import type { AppSettings } from '@/types'
import { getSettings, saveSettings } from '@/services/storage'
import { applyColorMode } from '@/services/theme'

function cloneSettings(value: AppSettings): AppSettings {
  return JSON.parse(JSON.stringify(value)) as AppSettings
}

function parsePeriods(value: string, fallback: number[]) {
  const periods = value
    .split(/[,\s，]+/)
    .map((item) => Number(item.trim()))
    .filter((item) => Number.isFinite(item) && item > 0)
  return periods.length > 0 ? periods : fallback
}

const defaultSettings = getSettings()
const settings = reactive<AppSettings>(cloneSettings(defaultSettings))
const maDraft = ref(settings.indicatorConfig.ma.join(', '))
const rsiDraft = ref(settings.indicatorConfig.rsi.join(', '))

function applySettings(next: AppSettings) {
  Object.assign(settings, cloneSettings(next))
  maDraft.value = settings.indicatorConfig.ma.join(', ')
  rsiDraft.value = settings.indicatorConfig.rsi.join(', ')
}

function save() {
  settings.indicatorConfig.ma = parsePeriods(maDraft.value, settings.indicatorConfig.ma)
  settings.indicatorConfig.rsi = parsePeriods(rsiDraft.value, settings.indicatorConfig.rsi)
  settings.heatmapConfig.colorMode = settings.colorMode
  saveSettings(cloneSettings(settings))
  applyColorMode(settings.colorMode)
  message.success('设置已保存')
}

function reset() {
  const cleared: AppSettings = {
    refreshInterval: { list: 0, detail: 5000, heatmap: 10000 },
    colorMode: 'red-rise',
    heatmapConfig: {
      dimension: 'industry',
      colorField: 'changePercent',
      sizeField: 'totalMarketCap',
      colorMode: 'red-rise',
      topK: 200,
    },
    indicatorConfig: {
      ma: [5, 10, 20, 60],
      macd: { short: 12, long: 26, signal: 9 },
      boll: { period: 20, stdDev: 2 },
      kdj: { period: 9, kPeriod: 3, dPeriod: 3 },
      rsi: [6, 12, 24],
      dmi: { period: 14, adxPeriod: 14 },
      sar: { afStart: 0.02, afIncrement: 0.02, afMax: 0.2 },
      kc: { emaPeriod: 20, atrPeriod: 10, multiplier: 2 },
    },
  }
  applySettings(cleared)
  saveSettings(cleared)
  applyColorMode(cleared.colorMode)
  message.success('已恢复默认设置')
}
</script>

<style scoped>
.settings-page :deep(.ant-input-number) {
  width: 92px;
}

.note-list {
  margin: 8px 0 0;
  padding-left: 18px;
  color: var(--color-muted);
  line-height: 1.8;
}
</style>
