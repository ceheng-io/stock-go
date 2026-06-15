import type { KLineData } from 'klinecharts'
import type { IndicatorConfig } from '@/types'
import type { OscillatorIndicatorKey, OverlayIndicatorKey } from '@/services/charts'

export interface KLineChartSourceRow {
  date?: string | null
  time?: string | null
  open?: number | null
  close?: number | null
  low?: number | null
  high?: number | null
  volume?: number | null
}

const OSCILLATOR_INDICATOR_MAP: Record<OscillatorIndicatorKey, string> = {
  macd: 'MACD',
  kdj: 'KDJ',
  rsi: 'RSI',
  obv: 'OBV',
  roc: 'ROC',
  dmi: 'DMI',
}

const OVERLAY_INDICATOR_MAP: Record<OverlayIndicatorKey, string> = {
  ma: 'MA',
  boll: 'BOLL',
  sar: 'SAR',
  kc: 'KC',
}

function toTimestamp(value?: string | null): number | null {
  if (!value) return null
  const normalized = value.trim()
  if (!normalized) return null
  const withTime = normalized.includes(' ') ? normalized.replace(' ', 'T') : `${normalized}T00:00:00`
  const timestamp = new Date(`${withTime}+08:00`).getTime()
  return Number.isFinite(timestamp) ? timestamp : null
}

function finiteNumber(value?: number | null): value is number {
  return typeof value === 'number' && Number.isFinite(value)
}

export function toKLineChartData(rows: KLineChartSourceRow[]): KLineData[] {
  return rows
    .map((row): KLineData | null => {
      const timestamp = toTimestamp(row.date || row.time)
      if (
        timestamp === null ||
        !finiteNumber(row.open) ||
        !finiteNumber(row.close) ||
        !finiteNumber(row.low) ||
        !finiteNumber(row.high)
      ) {
        return null
      }
      return {
        timestamp,
        open: row.open,
        close: row.close,
        low: row.low,
        high: row.high,
        volume: finiteNumber(row.volume) ? row.volume : 0,
      }
    })
    .filter((row): row is KLineData => row !== null)
}

export function buildKLineIndicatorNames(options: {
  overlays: OverlayIndicatorKey[]
  oscillator: OscillatorIndicatorKey
  indicatorConfig: IndicatorConfig
}): { overlays: string[]; panes: string[] } {
  return {
    overlays: options.overlays.map((item) => OVERLAY_INDICATOR_MAP[item]),
    panes: ['VOL', OSCILLATOR_INDICATOR_MAP[options.oscillator]],
  }
}
