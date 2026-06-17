# 策衡 stock-go

`stock-go` 是“策衡”的 Go 版本股票行情 SDK。代码仓库：[ceheng-io/stock-go](https://github.com/ceheng-io/stock-go)。

它面向 Go 服务端、命令行工具、数据任务和后续 Web/API 应用，提供 A 股、港股、美股、公募基金、期货、期权、资金流、北向资金、龙虎榜、大宗交易、融资融券、涨停池等公开行情数据能力。

```go
module github.com/ceheng-io/stock-go
```

当前阶段已经按 SDK v0.1 收口；CLI、MCP、后端 API 和 Web 前端属于后续应用层。收口边界见 [docs/sdk-v0.1-acceptance.md](docs/sdk-v0.1-acceptance.md)。

## 特性

- 统一根入口：`stock.New()` 创建 `*stock.Client`，服务字段按领域组织。
- Go 友好的命名空间 API：`client.Quotes`、`client.Kline`、`client.Board`、`client.MarketEvent` 等。
- 便捷入口：保留大量 `Client.Get*` 薄委托和常用常量/类型别名。
- 统一符号模型：A 股、港股、美股、基金、期货、板块代码归一化和数据源代码转换。
- 行情数据：A 股、港股、美股、公募基金实时行情、代码列表、批量行情、分时、K 线。
- 扩展数据：板块、资金流、北向资金、龙虎榜、大宗交易、融资融券、分红、基金、期货、期权、涨停池。
- 技术分析：MA、MACD、BOLL、KDJ、RSI、WR、BIAS、CCI、ATR、OBV、ROC、DMI、SAR、KC。
- 信号与策略：金叉/死叉、超买/超卖、BOLL 突破、SAR 反转、本地选股器和轻量回测。
- 请求治理：timeout、retry、限流、熔断、host fallback、User-Agent 轮换、provider policy、请求 hooks。
- 统一错误码：`INVALID_ARGUMENT`、`INVALID_SYMBOL`、`HTTP_ERROR`、`RATE_LIMITED`、`NETWORK_ERROR`、`TIMEOUT`、`ABORTED`、`PARSE_ERROR`、`UPSTREAM_ERROR`、`NOT_FOUND`。
- 公开纯能力子包：`indicators`、`signals`、`symbols`、`screener`、`cache`、`errors`、`parser`、`utils` 等可独立使用。

## 安装

```bash
go get github.com/ceheng-io/stock-go
```

```go
import stock "github.com/ceheng-io/stock-go"
```

更完整的接入说明、场景示例、错误处理、请求治理和生产注意事项见 [docs/usage.md](docs/usage.md)。

## 快速开始

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
		fmt.Printf("%s: %.2f\n", quote.Name, quote.Price)
	}
}
```

## 常用示例

### 符号归一化

```go
symbol, err := stock.NormalizeSymbol("600519", nil)
if err != nil {
	panic(err)
}

fmt.Println(stock.ToTencentSymbol(symbol))    // sh600519
fmt.Println(stock.ToEastmoneySecID(symbol))   // 1.600519
fmt.Println(stock.ToEastmoneySecid(symbol))   // 兼容别名
```

### 历史 K 线与技术指标

```go
rows, err := sdk.Kline.CN(ctx, "600519", stock.HistoryKlineOptions{
	Period:    stock.KlinePeriodDaily,
	Adjust:    stock.AdjustQFQ,
	StartDate: "20250101",
	EndDate:   "20250613",
})
if err != nil {
	panic(err)
}

withIndicators, err := sdk.Indicator.KlineWithIndicators(ctx, "600519", stock.KlineWithIndicatorsOptions{
	Period: stock.KlinePeriodDaily,
	Indicators: stock.IndicatorOptions{
		MA:   &stock.MAOptions{Periods: []int{5, 10, 20}},
		MACD: &stock.MACDOptions{},
		RSI:  &stock.RSIOptions{Periods: []int{6, 12, 24}},
	},
})
if err != nil {
	panic(err)
}

fmt.Println(len(rows), len(withIndicators))
```

### 全市场批量行情

```go
all, err := sdk.Quotes.AllCN(ctx, stock.CodeListOptions{}, stock.BatchOptions{
	Concurrency: 5,
	OnProgress: func(done, total int) {
		fmt.Printf("%d/%d\n", done, total)
	},
})
if err != nil {
	panic(err)
}

fmt.Println("A 股数量:", len(all))
```

也可以使用根包兼容入口：

```go
all, err := sdk.GetAllAShareQuotes(ctx, stock.GetAllAShareQuotesOptions{
	Concurrency: 5,
})
```

### 同花顺涨停池

```go
// 当日涨停池：Date 留空，按同花顺默认交易日查询。
today, err := sdk.MarketEvent.THSLimitUpPool(ctx, stock.THSLimitUpPoolOptions{
	Limit: 50,
})
if err != nil {
	panic(err)
}
fmt.Println(today.Date, today.Page.Total)

// 历史涨停池：支持 YYYY-MM-DD 或 YYYYMMDD。
history, err := sdk.GetTHSLimitUpPool(ctx, stock.THSLimitUpPoolOptions{
	Date:       "2025-06-13",
	Page:       1,
	Limit:      20,
	OrderField: stock.THSLimitUpOrderLastLimitUpTime,
	OrderType:  stock.THSLimitUpOrderDesc,
})
if err != nil {
	panic(err)
}

for _, item := range history.Items {
	fmt.Printf("%s %s %s %s\n", item.Code, item.Name, item.LimitUpType, item.LastLimitUpTimeText)
}
```

默认参数：

- `Date=""`：查询同花顺默认交易日。
- `Page=1`、`Limit=50`。
- `Filter="HS,GEM2STAR"`。
- `OrderField=THSLimitUpOrderLastLimitUpTime`。
- `OrderType=THSLimitUpOrderDesc`。

同花顺接口有反爬限制，SDK 已为 `ProviderTHS` 配置浏览器式默认 `User-Agent`、`Referer` 和 `Cookie`；如上游策略变化，可通过 `WithProviderPolicy(stock.ProviderTHS, ...)` 覆盖。

### 资金流、北向、龙虎榜

```go
fundFlow, err := sdk.FundFlow.Individual(ctx, "600519", stock.FundFlowOptions{})
northbound, err := sdk.Northbound.HoldingRank(ctx, stock.NorthboundHoldingRankOptions{
	Market: stock.NorthboundMarketAll,
	Period: stock.NorthboundRankToday,
})
dragonTiger, err := sdk.DragonTiger.Detail(ctx, stock.DragonTigerDateOptions{
	StartDate: "2025-06-13",
	EndDate:   "2025-06-13",
})

_, _, _, _ = fundFlow, northbound, dragonTiger, err
```

### 期货和期权

```go
futures, err := sdk.Futures.Kline(ctx, "rb2605", stock.FuturesKlineOptions{
	Period:    stock.KlinePeriodDaily,
	StartDate: "20250101",
	EndDate:   "20250613",
})
if err != nil {
	panic(err)
}

etfMinute, err := sdk.Options.ETFOptionMinute(ctx, "10009633")
if err != nil {
	panic(err)
}

fmt.Println(len(futures), len(etfMinute))
```

### 纯计算子包

```go
import (
	stock "github.com/ceheng-io/stock-go"
	"github.com/ceheng-io/stock-go/indicators"
	"github.com/ceheng-io/stock-go/screener"
	"github.com/ceheng-io/stock-go/signals"
)

macd := indicators.CalcMACD(closes, indicators.MACDOptions{})
signalRows, err := signals.CalcSignals(signalKlines, signals.SignalOptions{
	MA: &signals.MAOptions{Fast: 5, Slow: 20},
})
picks, err := screener.Screen(quotes).
	Where(func(q stock.FullQuote) bool { return q.ChangePercent > 3 }).
	Top(20)

_, _, _, _ = macd, signalRows, picks, err
```

## 请求治理与错误处理

```go
sdk := stock.New(
	stock.WithTimeout(12*time.Second),
	stock.WithRetry(stock.RetryOptions{
		MaxRetries: 2,
		BaseDelay: 500 * time.Millisecond,
	}),
	stock.WithProviderPolicy(stock.ProviderEastmoney, stock.ProviderPolicy{
		Timeout: 15 * time.Second,
		RateLimit: &stock.RateLimitOptions{
			RequestsPerSecond: 3,
			MaxBurst:          3,
		},
	}),
)

_, err := sdk.Quotes.SimpleCN(context.Background(), []string{"bad"})
if err != nil {
	fmt.Println(stock.GetErrorCode(err))
	if stock.IsSdkError(err) {
		fmt.Println("structured sdk error")
	}
}
```

## API 概览

| 服务字段 | 代表能力 |
| --- | --- |
| `client.Quotes` | A/HK/US/基金行情、资金流、盘口大单、搜索、交易日历、代码列表、批量行情 |
| `client.Kline` | A/HK/US 历史 K 线和分钟 K 线 |
| `client.Indicator` | K 线 + 技术指标 |
| `client.Board` | 行业/概念板块列表、盘口、成分股、K 线、分钟线 |
| `client.FundFlow` | 个股、大盘、排行、板块资金流 |
| `client.Northbound` | 北向/南向分时、汇总、持股排行、历史和个股持仓 |
| `client.MarketEvent` | 东方财富涨停池、盘口异动、板块异动、同花顺涨停池 |
| `client.DragonTiger` | 龙虎榜详情、个股统计、机构买卖、营业部排行、席位明细 |
| `client.BlockTrade` | 大宗交易市场统计、成交明细、每日个股统计 |
| `client.Margin` | 融资融券账户统计和标的明细 |
| `client.Dividend` | 个股分红派送详情 |
| `client.Fund` | 基金估值、历史净值、排名走势、基金分红 |
| `client.Futures` | 国内/全球期货 K 线、全球期货现货、库存、COMEX 库存 |
| `client.Options` | 中金所期权、ETF 期权、股指期权、商品期权、期权龙虎榜 |
| `client.Calendar` | 交易日判断、前后交易日、市场状态 |
| `client.Data` | 搜索、代码列表、大宗交易、融资融券、分红等聚合入口 |

常用根包便捷入口包括：

```go
client.GetSimpleQuotes(ctx, []string{"sh600519"})
client.GetHistoryKline(ctx, "600519", stock.HistoryKlineOptions{})
client.GetKlineWithIndicators(ctx, "600519", stock.KlineWithIndicatorsOptions{})
client.GetZTPool(ctx, stock.ZTPoolZT, "2025-06-13")
client.GetTHSLimitUpPool(ctx, stock.THSLimitUpPoolOptions{Date: "2025-06-13"})
client.GetDragonTigerDetail(ctx, stock.DragonTigerDateOptions{StartDate: "20250613", EndDate: "20250613"})
```

完整公开 API 速查见 [docs/api-matrix.md](docs/api-matrix.md)。

## 市场支持矩阵

| 能力 | A 股 | 港股 | 美股 | 公募基金 | 期货 | 期权 |
| --- | :---: | :---: | :---: | :---: | :---: | :---: |
| 实时行情 | 支持 | 支持 | 支持 | 支持 | 全球期货 | ETF / 中金所 / 商品 |
| 历史 K 线（日/周/月） | 支持 | 支持 | 支持 | 场内 ETF/LOF | 国内 + 全球 | 支持 |
| 分钟 K 线 | 1/5/15/30/60 | 1/5/15/30/60 | 1/5/15/30/60 | 场内 ETF/LOF | 暂无 | ETF 期权 |
| 当日分时 | 支持 | period=1 | period=1 | 场内 ETF/LOF | 暂无 | ETF 期权 |
| 资金流向 | 个股/大盘/排行/板块 | 暂无 | 暂无 | 不适用 | 不适用 | 不适用 |
| 板块（行业/概念） | 支持 | 暂无 | 暂无 | 暂无 | 不适用 | 不适用 |
| 龙虎榜 | 支持 | 不适用 | 不适用 | 不适用 | 不适用 | 期权龙虎榜 |
| 北向/南向资金 | 北向 | 南向 | 不适用 | 不适用 | 不适用 | 不适用 |
| 大宗交易 / 融资融券 | 支持 | 暂无 | 暂无 | 不适用 | 不适用 | 不适用 |
| 涨停池 / 盘口异动 | 东方财富 + 同花顺 | 不适用 | 不适用 | 不适用 | 不适用 | 不适用 |
| 全市场代码列表 / 批量行情 | 支持 | 支持 | 支持 | 支持 | 暂无 | 暂无 |
| 库存数据 | 不适用 | 不适用 | 不适用 | 不适用 | 国内 + COMEX | 不适用 |
| 交易日历 | 支持 | 市场状态 | 市场状态 | 不适用 | 不适用 | 不适用 |

> 数据来自腾讯财经、东方财富、新浪财经、同花顺等公开接口，通常有延迟，不适合高频交易决策。

## 目录结构

```text
stock-go/
├── stock.go              # 根 Client 与服务组合
├── options.go            # 请求、治理和 provider 配置
├── types/                # 公开领域类型
├── indicators/           # 技术指标
├── symbols/              # 符号归一化
├── signals/              # 指标信号
├── screener/             # 本地选股和回测
├── cache/                # 缓存能力
├── internal/
│   ├── core/             # HTTP、重试、限流、熔断、fallback、解析
│   ├── providers/        # tencent / eastmoney / sina / ths
│   └── services/         # 业务领域编排
├── cmd/ceheng/           # 后续 CLI 入口
├── apps/api/             # 后续 API 服务
├── apps/web/             # 后续 Web 前端
└── docs/
```

## 开发验证

```bash
go test ./...
go test ./types -run 'Test(TypesFilesStaySmall|DomainTypesStayInDomainFiles)' -count=1
```

真实网络集成测试默认跳过，可按需开启：

```bash
CEHENG_INTEGRATION=1 go test ./internal/providers/ths -run TestGetLimitUpPoolIntegration -count=1
```

`apps/web` 作为后续前端应用目录有独立 `go.mod` 边界，用于避免根模块的 `go list ./...` 穿透 `node_modules` 中的第三方 Go 包。

写任何前端代码前，必须先执行本地 `skills/frontend-design` 并确认页面结构、交互状态、数据依赖和验收方式。
