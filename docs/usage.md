# 策衡 stock-go 使用手册

本文面向第一次接入 `github.com/ceheng-io/stock-go` 的 Go 开发者，按真实使用场景说明如何安装、配置、拉取行情、计算指标、处理错误和治理请求。

`stock-go` 当前定位是“策衡”的 Go 版本 SDK 库。根包提供稳定入口和薄委托，具体数据源适配、请求治理和服务编排放在 `internal/` 下；业务代码通常只需要导入根包，以及按需导入 `indicators`、`symbols`、`signals`、`screener`、`cache` 等纯能力子包。

## 安装

```bash
go get github.com/ceheng-io/stock-go
```

```go
import stock "github.com/ceheng-io/stock-go"
```

最低接入方式是创建一个 `*stock.Client`：

```go
package main

import (
	"context"
	"fmt"
	"time"

	stock "github.com/ceheng-io/stock-go"
)

func main() {
	ctx := context.Background()
	sdk := stock.New(stock.WithTimeout(10 * time.Second))

	quotes, err := sdk.Quotes.SimpleCN(ctx, []string{"sh000001", "sz000858", "sh600519"})
	if err != nil {
		panic(err)
	}

	for _, quote := range quotes {
		fmt.Printf("%s %s %.2f %.2f%%\n", quote.Code, quote.Name, quote.Price, quote.ChangePercent)
	}
}
```

建议在业务服务里复用同一个 `Client`。`Client` 内部没有保存单次请求状态，适合在多个请求或任务之间共享；具体并发度由调用方、批量参数和限流配置共同控制。

## 入口选择

SDK 同时提供两类入口：

| 入口 | 示例 | 推荐场景 |
| --- | --- | --- |
| 服务字段 | `sdk.Quotes.SimpleCN(ctx, codes)` | 新 Go 项目，命名清晰，按领域浏览能力 |
| 便捷薄委托 | `sdk.GetSimpleQuotes(ctx, codes)` | 希望沿用扁平方法名或减少服务字段层级 |

两类入口最终调用同一套服务实现。新代码优先使用服务字段；偏好扁平方法名时可直接使用 `Get*` 方法。

常用服务字段：

| 字段 | 能力 |
| --- | --- |
| `Quotes` | A 股、港股、美股、基金实时行情、搜索、代码列表、交易日历、批量行情 |
| `Kline` | A 股、港股、美股历史 K 线和分钟 K 线 |
| `Indicator` | 拉取 K 线并附加技术指标 |
| `Board` | 行业、概念板块列表、盘口、成分股、K 线和分钟线 |
| `FundFlow` | 个股、大盘、排行、行业/概念/地域资金流 |
| `Northbound` | 北向/南向分时、汇总、持股排行、历史和个股持仓 |
| `MarketEvent` | 东方财富涨停池、盘口异动、板块异动、同花顺涨停池 |
| `DragonTiger` | 龙虎榜详情、个股统计、机构买卖、营业部排行、席位明细 |
| `BlockTrade` | 大宗交易市场统计、成交明细、每日个股统计 |
| `Margin` | 融资融券账户统计和标的明细 |
| `Dividend` | 个股分红派送详情 |
| `Fund` | 基金估值、历史净值、排名走势、基金分红 |
| `Futures` | 国内/全球期货 K 线、全球期货现货、库存、COMEX 库存 |
| `Options` | 中金所期权、ETF 期权、股指期权、商品期权、期权龙虎榜 |
| `Calendar` | 交易日判断、前后交易日、市场状态 |
| `Data` | 搜索、代码列表、大宗交易、融资融券、分红等聚合入口 |

完整公开 API 速查见 [api-matrix.md](api-matrix.md)。

## Client 配置

`stock.New()` 接受函数式选项：

```go
sdk := stock.New(
	stock.WithTimeout(12*time.Second),
	stock.WithUserAgent("my-service/1.0"),
	stock.WithRetry(stock.RetryOptions{
		MaxRetries:        2,
		BaseDelay:         300 * time.Millisecond,
		MaxDelay:          2 * time.Second,
		BackoffMultiplier: 2,
	}),
	stock.WithRateLimit(stock.RateLimitOptions{
		RequestsPerSecond: 5,
		MaxBurst:          10,
	}),
)
```

常用配置说明：

| 配置 | 作用 |
| --- | --- |
| `WithHTTPClient` | 注入自定义 `*http.Client`，适合代理、观测、测试替身 |
| `WithTimeout` | 设置默认请求超时 |
| `WithUserAgent` | 设置全局 `User-Agent` |
| `WithNextUserAgent` / `WithRandomUserAgent` / `WithRotateUserAgent` | 使用内置浏览器 UA 池 |
| `WithHeaders` | 设置全局请求头 |
| `WithRetry` | 配置重试次数、退避、可重试状态码和重试回调 |
| `WithRateLimit` | 配置全局 token bucket 限流 |
| `WithCircuitBreaker` | 配置连续失败后的熔断保护 |
| `WithProviderPolicy` | 针对单个数据源覆盖超时、请求头、重试、限流和熔断 |
| `WithRequestHooks` | 观察 request、response、error、retry、fallback 生命周期 |

当某个上游更敏感时，可以只对该数据源加治理：

```go
sdk := stock.New(
	stock.WithProviderPolicy(stock.ProviderEastmoney, stock.ProviderPolicy{
		Timeout: 15 * time.Second,
		RateLimit: &stock.RateLimitOptions{
			RequestsPerSecond: 3,
			MaxBurst:          3,
		},
	}),
	stock.WithProviderPolicy(stock.ProviderTHS, stock.ProviderPolicy{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
		Headers: map[string]string{
			"Referer": "https://data.10jqka.com.cn/market/ztStock/",
		},
	}),
)
```

同花顺接口默认已经内置浏览器式 `User-Agent`、`Referer` 和 `Cookie`，一般不需要手动配置；如果上游策略变化，可使用 `WithProviderPolicy(stock.ProviderTHS, ...)` 覆盖。

## 符号输入与归一化

行情方法通常接受常见代码写法，例如：

| 市场 | 可用输入示例 | 常用归一化结果 |
| --- | --- | --- |
| A 股 | `600519`、`sh600519`、`000858`、`sz000858` | 腾讯 `sh600519`，东方财富 `1.600519` |
| 港股 | `00700`、`hk00700` | 腾讯 `hk00700` |
| 美股 | `AAPL`、`usAAPL` | 美股符号 |
| 基金 | `510300`、`sh510300` | 场内基金代码 |
| 期货 | `rb2605`、`sc2605` | 按品种推断交易所 |

需要在业务侧显式处理代码时，可使用根包或 `symbols` 子包：

```go
symbol, err := stock.NormalizeSymbol("600519", nil)
if err != nil {
	return err
}

fmt.Println(stock.ToTencentSymbol(symbol))  // sh600519
fmt.Println(stock.ToEastmoneySecID(symbol)) // 1.600519
```

如果你的上游数据混合了市场前缀和纯数字，建议先统一归一化，再进入缓存键、日志和业务表。

## 实时行情

### A 股行情

简版行情适合列表和看板：

```go
quotes, err := sdk.Quotes.SimpleCN(ctx, []string{"sh000001", "sz000858", "sh600519"})
if err != nil {
	return err
}

for _, quote := range quotes {
	fmt.Printf("%s %s %.2f %.2f%%\n", quote.Code, quote.Name, quote.Price, quote.ChangePercent)
}
```

详细行情适合个股页或策略计算：

```go
quotes, err := sdk.Quotes.CN(ctx, []string{"600519", "000858"})
if err != nil {
	return err
}

for _, quote := range quotes {
	fmt.Printf("%s open=%.2f high=%.2f low=%.2f volume=%d\n",
		quote.Name, quote.Open, quote.High, quote.Low, quote.Volume)
}
```

兼容入口：

```go
quotes, err := sdk.GetFullQuotes(ctx, []string{"600519", "000858"})
```

### 港股、美股、基金

```go
hkQuotes, err := sdk.Quotes.HK(ctx, []string{"hk00700", "hk09988"})
if err != nil {
	return err
}

usQuotes, err := sdk.Quotes.US(ctx, []string{"AAPL", "MSFT"})
if err != nil {
	return err
}

fundQuotes, err := sdk.Quotes.Fund(ctx, []string{"510300", "159919"})
if err != nil {
	return err
}

fmt.Println(len(hkQuotes), len(usQuotes), len(fundQuotes))
```

## 搜索与代码列表

搜索适合把用户输入转换成可交易标的：

```go
results, err := sdk.Quotes.Search(ctx, "茅台")
if err != nil {
	return err
}

for _, result := range results {
	fmt.Printf("%s %s %s\n", result.Code, result.Name, result.Type)
}
```

获取市场代码列表：

```go
aShareCodes, err := sdk.Quotes.CodesCN(ctx, stock.CodeListOptions{
	Market: stock.AShareMarketSH,
})
if err != nil {
	return err
}

hkCodes, err := sdk.Quotes.CodesHK(ctx)
if err != nil {
	return err
}

usCodes, err := sdk.Quotes.CodesUS(ctx, stock.USCodeListOptions{
	Market: stock.USMarketNASDAQ,
})
if err != nil {
	return err
}

fmt.Println(len(aShareCodes), len(hkCodes), len(usCodes))
```

## 全市场批量行情

全市场行情会先拉取代码列表，再分批请求行情。建议显式设置并发和进度回调，避免对上游造成压力：

```go
quotes, err := sdk.Quotes.AllCN(ctx, stock.CodeListOptions{
	Market: stock.AShareMarketSH,
}, stock.BatchOptions{
	BatchSize:   80,
	Concurrency: 4,
	OnProgress: func(done, total int) {
		fmt.Printf("quotes progress: %d/%d\n", done, total)
	},
})
if err != nil {
	return err
}

fmt.Println("A 股行情数量:", len(quotes))
```

只对已有代码批量取行情：

```go
quotes, err := sdk.Quotes.BatchCN(ctx, []string{"600519", "000858", "300750"}, stock.BatchOptions{
	BatchSize:   50,
	Concurrency: 3,
})
```

根包便捷入口会把代码列表选项和批量选项合并到一个结构体里：

```go
quotes, err := sdk.GetAllAShareQuotes(ctx, stock.GetAllAShareQuotesOptions{
	Market:      stock.AShareMarketSH,
	BatchSize:   80,
	Concurrency: 4,
})
```

## K 线与分时

### A 股历史 K 线

```go
rows, err := sdk.Kline.CN(ctx, "600519", stock.HistoryKlineOptions{
	Period:    stock.KlinePeriodDaily,
	Adjust:    stock.AdjustQFQ,
	StartDate: "20250101",
	EndDate:   "20250613",
})
if err != nil {
	return err
}

for _, row := range rows {
	if row.Close != nil {
		fmt.Printf("%s close=%.2f\n", row.Date, *row.Close)
	}
}
```

`Period` 常用值：

| 常量 | 含义 |
| --- | --- |
| `KlinePeriodDaily` | 日线 |
| `KlinePeriodWeekly` | 周线 |
| `KlinePeriodMonthly` | 月线 |

`Adjust` 常用值：

| 常量 | 含义 |
| --- | --- |
| `AdjustNone` | 不复权 |
| `AdjustQFQ` | 前复权 |
| `AdjustHFQ` | 后复权 |

### 分钟线和当日分时

```go
minute, err := sdk.Kline.CNMinute(ctx, "600519", stock.MinuteKlineOptions{
	Period: stock.MinutePeriodFive,
	NDays:  5,
})
if err != nil {
	return err
}

fmt.Println(len(minute.Klines), len(minute.Timeline))
```

分钟周期：

| 常量 | 含义 |
| --- | --- |
| `MinutePeriodOne` | 1 分钟或当日分时 |
| `MinutePeriodFive` | 5 分钟 |
| `MinutePeriodFifteen` | 15 分钟 |
| `MinutePeriodThirty` | 30 分钟 |
| `MinutePeriodSixty` | 60 分钟 |

港股和美股使用对应服务：

```go
hkRows, err := sdk.Kline.HK(ctx, "00700", stock.HKKlineOptions{
	Period:    stock.KlinePeriodDaily,
	Adjust:    stock.AdjustQFQ,
	StartDate: "20250101",
	EndDate:   "20250613",
})
if err != nil {
	return err
}

usRows, err := sdk.Kline.US(ctx, "AAPL", stock.USKlineOptions{
	Period:    stock.KlinePeriodDaily,
	Adjust:    stock.AdjustNone,
	StartDate: "20250101",
	EndDate:   "20250613",
})
if err != nil {
	return err
}

fmt.Println(len(hkRows), len(usRows))
```

## 技术指标

如果希望“拉取 K 线 + 附加指标”一步完成，使用 `sdk.Indicator`：

```go
rows, err := sdk.Indicator.KlineWithIndicators(ctx, "600519", stock.KlineWithIndicatorsOptions{
	Period:    stock.KlinePeriodDaily,
	Adjust:    stock.AdjustQFQ,
	StartDate: "20250101",
	EndDate:   "20250613",
	Indicators: stock.IndicatorOptions{
		MA:   &stock.MAOptions{Periods: []int{5, 10, 20}},
		MACD: &stock.MACDOptions{},
		RSI:  &stock.RSIOptions{Periods: []int{6, 12, 24}},
		BOLL: &stock.BOLLOptions{Period: 20, StdDev: 2},
	},
})
if err != nil {
	return err
}

last := rows[len(rows)-1]
fmt.Println(last.Date, last.Close, last.MA, last.MACD)
```

也可以完全离线使用 `indicators` 子包：

```go
import "github.com/ceheng-io/stock-go/indicators"

closes := indicators.Values(10.2, 10.5, 10.8, 10.4, 11.0)
ma := indicators.CalcMA(closes, indicators.MAOptions{
	Periods: []int{3, 5},
	Type:    indicators.MATypeSMA,
})

fmt.Println(ma)
```

指标输出使用 `*float64` 表示可空值。窗口不足、输入为空或无法计算时返回 `nil`，业务层展示前需要判空。

## 板块数据

行业和概念板块共用相似 API：

```go
industries, err := sdk.Board.IndustryList(ctx)
if err != nil {
	return err
}

if len(industries) == 0 {
	return nil
}

boardCode := industries[0].Code
spot, err := sdk.Board.IndustrySpot(ctx, boardCode)
if err != nil {
	return err
}

constituents, err := sdk.Board.IndustryConstituents(ctx, boardCode)
if err != nil {
	return err
}

fmt.Println(industries[0].Name, len(spot), len(constituents))
```

概念板块把 `Industry*` 换成 `Concept*`：

```go
concepts, err := sdk.Board.ConceptList(ctx)
conceptRows, err := sdk.Board.ConceptConstituents(ctx, "BK0815")
_, _ = concepts, conceptRows
```

板块 K 线和分钟线：

```go
rows, err := sdk.Board.IndustryKline(ctx, "BK0475", stock.IndustryBoardKlineOptions{
	Period:    stock.KlinePeriodDaily,
	StartDate: "20250101",
	EndDate:   "20250613",
})
```

## 资金流

腾讯简版资金流可从 `Quotes` 获取，东方财富深度资金流从 `FundFlow` 获取。

个股资金流：

```go
rows, err := sdk.FundFlow.Individual(ctx, "600519", stock.FundFlowOptions{
	Period: stock.FundFlowPeriodDaily,
})
if err != nil {
	return err
}

fmt.Println(len(rows))
```

市场资金流：

```go
market, err := sdk.FundFlow.Market(ctx)
if err != nil {
	return err
}

fmt.Println(len(market))
```

资金流排行：

```go
rank, err := sdk.FundFlow.Rank(ctx, stock.FundFlowRankOptions{
	Indicator: stock.FundFlowRankToday,
})
if err != nil {
	return err
}

sectorRank, err := sdk.FundFlow.SectorRank(ctx, stock.FundFlowRankOptions{
	SectorType: stock.FundFlowSectorIndustry,
	Indicator: stock.FundFlowRankFiveDay,
})
if err != nil {
	return err
}

fmt.Println(len(rank), len(sectorRank))
```

## 北向资金

```go
summary, err := sdk.Northbound.Summary(ctx)
if err != nil {
	return err
}

rank, err := sdk.Northbound.HoldingRank(ctx, stock.NorthboundHoldingRankOptions{
	Market: stock.NorthboundMarketAll,
	Period: stock.NorthboundRankToday,
})
if err != nil {
	return err
}

history, err := sdk.Northbound.Individual(ctx, "600519", stock.NorthboundHistoryOptions{
	StartDate: "2025-01-01",
	EndDate:   "2025-01-31",
})
if err != nil {
	return err
}

fmt.Println(len(summary), len(rank), len(history))
```

日期参数通常支持 `YYYY-MM-DD` 或 `YYYYMMDD`，具体选项以对应类型字段为准。

## 龙虎榜、大宗交易、融资融券

龙虎榜：

```go
detail, err := sdk.DragonTiger.Detail(ctx, stock.DragonTigerDateOptions{
	StartDate: "2025-06-13",
	EndDate:   "2025-06-13",
})
if err != nil {
	return err
}

stats, err := sdk.DragonTiger.StockStats(ctx, stock.DragonTigerPeriodOneMonth)
if err != nil {
	return err
}

fmt.Println(len(detail), len(stats))
```

大宗交易：

```go
marketStat, err := sdk.BlockTrade.MarketStat(ctx)
if err != nil {
	return err
}

detail, err := sdk.BlockTrade.Detail(ctx, stock.BlockTradeDateOptions{
	StartDate: "2025-06-13",
	EndDate:   "2025-06-13",
})
if err != nil {
	return err
}

fmt.Println(len(marketStat), len(detail))
```

融资融券：

```go
accounts, err := sdk.Margin.AccountInfo(ctx)
if err != nil {
	return err
}

targets, err := sdk.Margin.TargetList(ctx, "2025-06-13")
if err != nil {
	return err
}

fmt.Println(len(accounts), len(targets))
```

## 涨停池和市场异动

东方财富涨停池：

```go
items, err := sdk.MarketEvent.ZTPool(ctx, stock.ZTPoolZT, "2025-06-13")
if err != nil {
	return err
}

fmt.Println(len(items))
```

盘口异动：

```go
changes, err := sdk.MarketEvent.StockChanges(ctx, stock.StockChangeLargeBuy)
if err != nil {
	return err
}

boardChanges, err := sdk.MarketEvent.BoardChanges(ctx)
if err != nil {
	return err
}

fmt.Println(len(changes), len(boardChanges))
```

同花顺涨停池：

```go
today, err := sdk.MarketEvent.THSLimitUpPool(ctx, stock.THSLimitUpPoolOptions{
	Limit: 50,
})
if err != nil {
	return err
}

history, err := sdk.GetTHSLimitUpPool(ctx, stock.THSLimitUpPoolOptions{
	Date:       "2025-06-13",
	Page:       1,
	Limit:      20,
	OrderField: stock.THSLimitUpOrderLastLimitUpTime,
	OrderType:  stock.THSLimitUpOrderDesc,
})
if err != nil {
	return err
}

fmt.Println(today.Date, today.Page.Total, len(history.Items))
```

`THSLimitUpPoolOptions` 默认值：

| 字段 | 默认值 |
| --- | --- |
| `Date` | 空字符串，使用同花顺默认交易日 |
| `Page` | `1` |
| `Limit` | `50` |
| `Filter` | `HS,GEM2STAR` |
| `OrderField` | `THSLimitUpOrderLastLimitUpTime` |
| `OrderType` | `THSLimitUpOrderDesc` |

## 基金、期货和期权

基金：

```go
estimate, err := sdk.Fund.Estimate(ctx, "000001")
if err != nil {
	return err
}

nav, err := sdk.Fund.NavHistory(ctx, "000001")
if err != nil {
	return err
}

dividends, err := sdk.Fund.DividendList(ctx, stock.FundDividendListOptions{
	Page: 1,
})
if err != nil {
	return err
}

fmt.Println(estimate.Code, len(nav.Items), len(dividends.Items))
```

期货：

```go
kline, err := sdk.Futures.Kline(ctx, "rb2605", stock.FuturesKlineOptions{
	Period:    stock.KlinePeriodDaily,
	StartDate: "20250101",
	EndDate:   "20250613",
})
if err != nil {
	return err
}

global, err := sdk.Futures.GlobalSpot(ctx, stock.GlobalFuturesSpotOptions{})
if err != nil {
	return err
}

inventory, err := sdk.Futures.Inventory(ctx, "螺纹钢", stock.FuturesInventoryOptions{})
if err != nil {
	return err
}

fmt.Println(len(kline), len(global), len(inventory))
```

期权：

```go
months, err := sdk.Options.ETFOptionMonths(ctx, stock.ETFOptionCate50ETF)
if err != nil {
	return err
}

minute, err := sdk.Options.ETFOptionMinute(ctx, "10009633")
if err != nil {
	return err
}

cffex, err := sdk.Options.CFFEXQuotes(ctx, stock.CFFEXOptionQuotesOptions{})
if err != nil {
	return err
}

fmt.Println(months, len(minute), len(cffex))
```

## 交易日历和市场状态

交易日历来自上游，适合判断 A 股交易日：

```go
ok, err := sdk.Calendar.IsTradingDay(ctx, "2025-06-13")
if err != nil {
	return err
}

next, err := sdk.Calendar.NextTradingDay(ctx, "2025-06-13")
if err != nil {
	return err
}

prev, err := sdk.Calendar.PrevTradingDay(ctx, "2025-06-13")
if err != nil {
	return err
}

fmt.Println(ok, prev, next)
```

市场状态是本地同步判断，不发起网络请求：

```go
status := sdk.Calendar.MarketStatus(stock.MarketA, time.Now())
fmt.Println(status)
```

## 本地选股和信号

`screener` 子包用于对已获取的数据做本地过滤、排序和截取：

```go
import "github.com/ceheng-io/stock-go/screener"

picks, err := screener.Screen(quotes).
	Where(func(q stock.FullQuote) bool {
		return q.ChangePercent > 3 && q.Volume > 0
	}).
	SortBy(func(q stock.FullQuote) *float64 {
		return &q.ChangePercent
	}, screener.Desc).
	Top(20)
if err != nil {
	return err
}

fmt.Println(len(picks))
```

根包也保留了 `stock.Screen` 兼容入口。

`signals` 子包用于根据 K 线和指标计算金叉/死叉、超买/超卖、BOLL 突破、SAR 反转等信号：

```go
import "github.com/ceheng-io/stock-go/signals"

signalsRows, err := signals.CalcSignals(signalKlines, signals.SignalOptions{
	MA: &signals.MAOptions{Fast: 5, Slow: 20},
	RSI: &signals.RSIOptions{
		Period:     14,
		Overbought: 70,
		Oversold:   30,
	},
})
if err != nil {
	return err
}

fmt.Println(len(signalsRows))
```

## 缓存

SDK 提供轻量内存缓存能力，适合缓存代码列表、交易日历、低频榜单等数据。

```go
store := stock.NewMemoryCache(stock.CacheOptions{
	DefaultTTL: 5 * time.Minute,
	MaxSize:    1000,
})

key := stock.CreateCacheKey("quotes", []string{"600519", "000858"})

quotes, err := stock.CacheThrough(ctx, store, key, func(ctx context.Context) ([]stock.FullQuote, error) {
	return sdk.Quotes.CN(ctx, []string{"600519", "000858"})
}, 30*time.Second)
if err != nil {
	return err
}

fmt.Println(len(quotes))
```

`MemoryCache` 带 single-flight 保护：多个 goroutine 同时请求同一个 key 时，只会有一个 goroutine 执行 fetcher，其余等待同一个结果。

## 错误处理

所有上游请求和参数错误都会尽量归一成 SDK 错误码。常见错误码：

| 错误码 | 含义 |
| --- | --- |
| `INVALID_ARGUMENT` | 参数不合法 |
| `INVALID_SYMBOL` | 标的代码无法识别 |
| `HTTP_ERROR` | 上游 HTTP 状态码错误 |
| `RATE_LIMITED` | 被限流 |
| `NETWORK_ERROR` | 网络错误 |
| `TIMEOUT` | 超时 |
| `ABORTED` | 请求被取消 |
| `PARSE_ERROR` | 响应解析失败 |
| `UPSTREAM_ERROR` | 上游返回异常数据 |
| `NOT_FOUND` | 上游未找到目标数据 |

推荐统一按错误码分支：

```go
quotes, err := sdk.Quotes.SimpleCN(ctx, []string{"bad"})
if err != nil {
	switch stock.GetErrorCode(err) {
	case stock.CodeInvalidSymbol, stock.CodeInvalidArgument:
		return fmt.Errorf("请求参数错误: %w", err)
	case stock.CodeRateLimited, stock.CodeTimeout, stock.CodeNetwork:
		return fmt.Errorf("上游暂时不可用: %w", err)
	default:
		return err
	}
}

fmt.Println(quotes)
```

判断是否为结构化 SDK 错误：

```go
if stock.IsSdkError(err) {
	fmt.Println(stock.GetErrorCode(err))
}
```

根包同时保留 Go initialism 风格和常用兼容命名，例如 `IsSDKError` / `IsSdkError`、`GetSDKErrorCode` / `GetSdkErrorCode`。

## 请求观测

使用 `WithRequestHooks` 可以接入日志、指标或 tracing：

```go
sdk := stock.New(stock.WithRequestHooks(stock.RequestHooks{
	OnRequest: func(ctx stock.RequestContext) {
		fmt.Println("request", ctx.Provider, ctx.URL)
	},
	OnResponse: func(ctx stock.RequestContext, meta stock.ResponseMeta) {
		fmt.Println("response", ctx.Provider, ctx.URL, meta.StatusCode, meta.Duration)
	},
	OnError: func(ctx stock.RequestContext, err error) {
		fmt.Println("error", ctx.Provider, ctx.URL, err)
	},
	OnRetry: func(ctx stock.RequestContext, err error, delay time.Duration) {
		fmt.Println("retry", ctx.Provider, ctx.URL, delay, err)
	},
	Trace: func(event stock.RequestTraceEvent, ctx stock.RequestContext) {
		fmt.Println("trace", event, ctx.Provider, ctx.URL)
	},
}))
```

请求 hook 不应执行耗时阻塞逻辑。需要写入外部日志或指标系统时，建议在 hook 中做轻量封装或投递到异步队列。

## 测试建议

SDK 单元测试不依赖真实网络：

```bash
go test ./...
```

真实网络集成测试默认跳过，按需打开：

```bash
CEHENG_INTEGRATION=1 go test ./internal/providers/ths -run TestGetLimitUpPoolIntegration -count=1
```

业务项目中建议按两层测试：

| 层级 | 建议 |
| --- | --- |
| 单元测试 | 注入 `WithHTTPClient` 的测试替身，或包装 SDK 调用接口后 mock |
| 集成测试 | 少量覆盖关键数据源，并用环境变量开关控制 |

## 生产使用注意事项

- 上游来自腾讯财经、东方财富、新浪财经、同花顺等公开接口，通常有延迟，不适合高频交易决策。
- 批量和轮询任务应配置 `WithRateLimit`、`WithRetry` 和合理并发。
- 同一业务进程内建议复用 `Client`，不要为每次请求重复创建。
- 外部输入的代码先归一化，再进入缓存键、日志和持久化。
- 指标和信号输出可能存在窗口不足导致的 `nil` 值，展示和计算前需要判空。
- 日期参数优先使用 `YYYY-MM-DD`；兼容方法通常也接受 `YYYYMMDD`。
- 网络接口返回字段可能随上游变化，关键链路应保存原始错误和请求元信息，便于排查。

## 进一步阅读

- [README](../README.md)：项目概览、能力矩阵和快速示例。
- [公开 API 速查](api-matrix.md)：根包、服务字段和子包能力索引。
- [SDK v0.1 验收文档](sdk-v0.1-acceptance.md)：当前阶段收口边界和验收项。
