package stock

import "github.com/ceheng.io/stock-go/timeutil"

type MarketTz = timeutil.MarketTz
type TimeMeta = timeutil.TimeMeta

var MarketTZ = timeutil.MarketTZ

// MARKET_TZ preserves the TypeScript SDK constant name.
var MARKET_TZ = MarketTZ

// ParseMarketTime 把市场本地时间字符串解析为 UTC unix 毫秒。
func ParseMarketTime(local string, tz MarketTz) (int64, bool) {
	return timeutil.ParseMarketTime(local, tz)
}

// BuildTimeMeta 构造时间元信息；解析失败时 Timestamp 为 nil。
func BuildTimeMeta(local string, tz MarketTz) TimeMeta {
	return timeutil.BuildTimeMeta(local, tz)
}

// BuildTimeMetaFromDateAndTime 将基础日期和 HH:mm/HH:mm:ss 组合后构造时间元信息。
func BuildTimeMetaFromDateAndTime(baseDate string, hhmm string, tz MarketTz) TimeMeta {
	return timeutil.BuildTimeMetaFromDateAndTime(baseDate, hhmm, tz)
}

// FormatInTz 把 UTC 毫秒时间戳格式化为指定市场时区的 YYYY-MM-DD HH:mm。
func FormatInTz(epoch *int64, tz MarketTz) string {
	return timeutil.FormatInTz(epoch, tz)
}
