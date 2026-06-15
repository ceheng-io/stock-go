import { computed, onUnmounted, ref, toValue, watch, type MaybeRefOrGetter } from 'vue'

interface UsePollingOptions {
  interval: MaybeRefOrGetter<number>
  enabled?: MaybeRefOrGetter<boolean>
  pauseOnHidden?: boolean
  immediate?: boolean
}

export function usePolling(fetcher: () => Promise<unknown> | unknown, options: UsePollingOptions) {
  const isLoading = ref(false)
  const lastRefresh = ref<number | null>(null)
  const isPaused = ref(false)
  let timer: number | undefined
  let stopped = false

  const enabled = computed(() => options.enabled === undefined ? true : toValue(options.enabled))
  const interval = computed(() => Math.max(0, Number(toValue(options.interval)) || 0))

  function clearTimer() {
    if (timer !== undefined) {
      window.clearTimeout(timer)
      timer = undefined
    }
  }

  async function refresh() {
    if (stopped || !enabled.value || isPaused.value || interval.value <= 0) return
    isLoading.value = true
    try {
      await fetcher()
      lastRefresh.value = Date.now()
    } catch (error) {
      console.warn('[usePolling] refresh failed', error)
    } finally {
      isLoading.value = false
    }
  }

  function scheduleNext() {
    clearTimer()
    if (stopped || !enabled.value || isPaused.value || interval.value <= 0) return
    if (options.pauseOnHidden !== false && document.hidden) return
    timer = window.setTimeout(async () => {
      await refresh()
      scheduleNext()
    }, interval.value)
  }

  function pause() {
    isPaused.value = true
    clearTimer()
  }

  function resume() {
    isPaused.value = false
    if (!enabled.value || interval.value <= 0) return
    refresh().then(scheduleNext)
  }

  function handleVisibilityChange() {
    if (options.pauseOnHidden === false) return
    if (document.hidden) {
      clearTimer()
      return
    }
    if (!isPaused.value && enabled.value && interval.value > 0) {
      refresh().then(scheduleNext)
    }
  }

  watch([enabled, interval, isPaused], () => {
    clearTimer()
    if (!enabled.value || isPaused.value || interval.value <= 0) return
    if (options.immediate === false) {
      scheduleNext()
    } else {
      refresh().then(scheduleNext)
    }
  }, { immediate: true })

  if (options.pauseOnHidden !== false) {
    document.addEventListener('visibilitychange', handleVisibilityChange)
  }

  onUnmounted(() => {
    stopped = true
    clearTimer()
    document.removeEventListener('visibilitychange', handleVisibilityChange)
  })

  return {
    isLoading,
    lastRefresh,
    refresh,
    pause,
    resume,
    isPaused,
  }
}
