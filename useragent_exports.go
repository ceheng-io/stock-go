package stock

import "github.com/ceheng-io/stock-go/useragent"

// AllUserAgents 返回所有可用 User-Agent 的副本。
func AllUserAgents() []string {
	return useragent.All()
}

// GetAllUserAgents preserves the TypeScript SDK naming style.
func GetAllUserAgents() []string {
	return AllUserAgents()
}

// NextUserAgent 以轮询方式返回下一个 User-Agent。
func NextUserAgent() string {
	return useragent.Next()
}

// GetNextUserAgent preserves the TypeScript SDK naming style.
func GetNextUserAgent() string {
	return NextUserAgent()
}

// RandomUserAgent 随机返回一个 User-Agent。
func RandomUserAgent() string {
	return useragent.Random()
}

// GetRandomUserAgent preserves the TypeScript SDK naming style.
func GetRandomUserAgent() string {
	return RandomUserAgent()
}

// ResetUserAgents 重置 User-Agent 轮询位置。
func ResetUserAgents() {
	useragent.Reset()
}
