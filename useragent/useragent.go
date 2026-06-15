package useragent

import (
	"math/rand"
	"sync"
	"time"
)

var values = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0",
}

var (
	mu      sync.Mutex
	current int
	rng     = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// All 返回所有可用 User-Agent 的副本。
func All() []string {
	copied := make([]string, len(values))
	copy(copied, values)
	return copied
}

// Next 以轮询方式返回下一个 User-Agent。
func Next() string {
	mu.Lock()
	defer mu.Unlock()
	if len(values) == 0 {
		return ""
	}
	value := values[current]
	current = (current + 1) % len(values)
	return value
}

// Random 随机返回一个 User-Agent。
func Random() string {
	mu.Lock()
	defer mu.Unlock()
	if len(values) == 0 {
		return ""
	}
	return values[rng.Intn(len(values))]
}

// Reset 重置轮询位置，主要用于测试和可重复调用场景。
func Reset() {
	mu.Lock()
	current = 0
	mu.Unlock()
}
