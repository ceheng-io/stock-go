import { defineComponent, nextTick, ref } from 'vue'
import { mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { usePolling } from '@/composables/usePolling'

function setDocumentHidden(hidden: boolean) {
  Object.defineProperty(document, 'hidden', {
    configurable: true,
    value: hidden,
  })
  document.dispatchEvent(new Event('visibilitychange'))
}

describe('usePolling', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    setDocumentHidden(false)
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('runs immediately and then on the configured interval', async () => {
    const fetcher = vi.fn(async () => undefined)

    mount(defineComponent({
      setup() {
        usePolling(fetcher, { interval: 1000 })
        return () => null
      },
    }))

    await nextTick()
    await Promise.resolve()
    expect(fetcher).toHaveBeenCalledTimes(1)

    await vi.advanceTimersByTimeAsync(1000)
    expect(fetcher).toHaveBeenCalledTimes(2)
  })

  it('does not start when disabled or interval is zero', async () => {
    const disabledFetcher = vi.fn(async () => undefined)
    const zeroIntervalFetcher = vi.fn(async () => undefined)

    mount(defineComponent({
      setup() {
        usePolling(disabledFetcher, { interval: 1000, enabled: false })
        usePolling(zeroIntervalFetcher, { interval: 0 })
        return () => null
      },
    }))

    await vi.advanceTimersByTimeAsync(3000)

    expect(disabledFetcher).not.toHaveBeenCalled()
    expect(zeroIntervalFetcher).not.toHaveBeenCalled()
  })

  it('pauses while the page is hidden and refreshes when visible again', async () => {
    const fetcher = vi.fn(async () => undefined)

    mount(defineComponent({
      setup() {
        usePolling(fetcher, { interval: 1000, immediate: false })
        return () => null
      },
    }))

    setDocumentHidden(true)
    await vi.advanceTimersByTimeAsync(1000)
    expect(fetcher).not.toHaveBeenCalled()

    setDocumentHidden(false)
    await nextTick()
    await Promise.resolve()
    expect(fetcher).toHaveBeenCalledTimes(1)
  })

  it('reacts to enabled changes', async () => {
    const fetcher = vi.fn(async () => undefined)
    const enabled = ref(false)

    mount(defineComponent({
      setup() {
        usePolling(fetcher, { interval: 1000, enabled })
        return () => null
      },
    }))

    await vi.advanceTimersByTimeAsync(1000)
    expect(fetcher).not.toHaveBeenCalled()

    enabled.value = true
    await nextTick()
    await Promise.resolve()
    expect(fetcher).toHaveBeenCalledTimes(1)
  })
})
