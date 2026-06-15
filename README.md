# 策衡 stock-go

`stock-go` 是“策衡”的 Go 版本股票行情 SDK，迁移自 TypeScript 项目 `/Users/xingyys/project/html/stock-sdk`。

当前仓库的 Go module：

```go
module github.com/ceheng.io/stock-go
```

## 当前阶段

第一阶段先建设 SDK 库底座：

- 根包 `stock`：`Client`、`New`、请求超时、重试、限流、熔断、User-Agent、HTTP client 注入和结构化统一错误类型。
- `cache` 与根包缓存工具：`MemoryCache`、`MemoryCacheStore`、`NewMemoryCache`、`GetSharedCache`、`CreateCacheKey`、`CacheThrough`，提供内存 TTL/LRU 缓存、共享缓存、缓存 key 生成和 single-flight 缓存穿透 helper。
- `symbols` 与根包符号工具：`NormalizeSymbol`、`ToTencentSymbol`、`ToEastmoneySecID`、`ToEastmoneySecid`、`ToPlainCode`、`InferAShareExchange`、`ExtractVariety`、`FuturesExchanges`，以及股票/基金/期货/板块代码归一化和数据源代码适配。
- `indicators` 与根包指标工具：MA/SMA/EMA/WMA、MACD、BOLL、KDJ、RSI、WR、BIAS、CCI、ATR、OBV、ROC、DMI、SAR、KC、`AddIndicators`、`BuildIndicatorContext`、`GetEnabledIndicatorKeys`、`EstimateIndicatorLookback` 等技术指标和指标注册辅助能力。
- `signals` 与根包信号工具：`CalcSignals`，基于已计算指标的金叉/死叉、超买/超卖、布林突破和 SAR 反转信号识别。
- `screener` 与根包选股/回测工具：`Screen`、`Backtest`，提供本地链式选股筛选器和单标的全仓多头轻量回测。
- `types`：行情、K 线、板块、资金流、北向、龙虎榜、大宗交易、两融、分红、基金、期货、期权、事件等公开数据结构，按领域拆分并限制单文件不超过 1000 行。
- `Quotes.CN`、`Quotes.SimpleCN`、`Quotes.HK`、`Quotes.US`、`Quotes.Fund`、`Quotes.FundFlow`、`Quotes.PanelLargeOrder`、`Quotes.TodayTimeline`：腾讯行情、资金流、盘口大单占比和当日分时 provider/service/root 闭环。
- `Quotes.Search`、`Quotes.TradingCalendar`：腾讯 Smartbox 搜索和 A 股交易日历基础能力。
- `GenerateSearchExternalLinks`：按搜索结果生成东方财富、雪球外部财经链接。
- `timeutil` 与根包时间工具：`MarketTZ`、`ParseMarketTime`、`BuildTimeMeta`、`BuildTimeMetaFromDateAndTime`、`FormatInTz`。
- `utils` 与根包工具：`ChunkArray`、`AsyncPool`、周期/复权参数校验、`PeriodCode`/`GetPeriodCode`、`AdjustCode`/`GetAdjustCode`，对应 TS 版工具函数基础能力。
- `parser` 与根包解析工具：`DecodeGBK`、`ParseResponse`、`SafeNumber`、`SafeNumberOrNil`、`SafeNumberOrNull`、`ToNumber`、`ToNumberSafe`。
- 根包脚本响应解析工具：`ExtractJSONP`、`ExtractJsonFromJsonp`、`JSONPRequest`、`ExtractJSVar`、`ParseJSVars`、`FetchJSVars`，对应 TS 版 `core/jsonp.ts` 和 `core/jsVars.ts` 的 JSONP 请求、文本解包、JS 变量声明请求与解析能力。
- `constants` 与根包常量：腾讯、东方财富、新浪 API URL、token、默认请求参数、期货/期权映射表，对应 TS 版 `core/constants.ts`。
- `useragent` 与根包 User-Agent 池：`AllUserAgents`、`NextUserAgent`、`RandomUserAgent`、`WithNextUserAgent`、`WithRandomUserAgent`。
- `RateLimiter` 与根包限流工具：`NewRateLimiter`、`RateLimiterOptions`，对应 TS 版 `core/rateLimiter.ts` 的令牌桶能力。
- `CircuitBreaker` 与根包熔断工具：`NewCircuitBreaker`、`CircuitBreakerOptions`、`CircuitBreakerStats` 和状态常量，对应 TS 版 `core/circuitBreaker.ts` 的状态机能力。
- `HostFallbackManager` 与根包 host fallback 工具：`NewHostFallbackManager`、`HostFallbackOptions`、`HostHealthStats`，对应 TS 版 `core/fallback.ts` 的备用 host 治理能力。
- Provider 策略工具：`MergeProviderPolicy`、`ResolveProviderPolicy`、`InferProviderFromURL`，用于合并、解析和推断数据源级请求治理配置。
- `Quotes.DividendDetail`：东方财富个股分红派送详情的兼容入口。
- `Calendar.IsTradingDay`、`Calendar.NextTradingDay`、`Calendar.PrevTradingDay`、`Calendar.MarketStatus`：A 股交易日判断、前后交易日和市场时段状态辅助能力。
- `Quotes.CodesCN`、`Quotes.CodesUS`、`Quotes.CodesHK`、`Quotes.CodesFund`：代码列表基础能力。
- `Quotes.BatchCN`、`Quotes.BatchHK`、`Quotes.BatchUS`、`Quotes.AllCN`、`Quotes.AllHK`、`Quotes.AllUS`：按代码批量行情和全市场代码列表 + 批量行情组合入口，支持分块、并发和进度回调。
- `Quotes.BatchRaw`：透传腾讯行情 raw query，并返回拆分后的原始字段。
- `Kline.CN`、`Kline.CNMinute`：东方财富 A 股历史日/周/月 K 线，以及 1/5/15/30/60 分钟线基础能力。
- `Kline.HK`、`Kline.US`、`Kline.HKMinute`、`Kline.USMinute`：东方财富港股、美股历史日/周/月 K 线，以及 1/5/15/30/60 分钟线基础能力。
- K 线与板块选项别名：`HKKlineOptions`、`USKlineOptions`、`HKMinuteKlineOptions`、`USMinuteKlineOptions`、`IndustryBoardKlineOptions`、`ConceptBoardKlineOptions` 等，对齐 TS 顶层公开类型名。
- `Indicator.KlineWithIndicators`：拉取足够历史 K 线并按需附加 MA、MACD、BOLL、RSI 等技术指标。
- `Board.IndustryList`、`Board.ConceptList`、`Board.IndustrySpot`、`Board.ConceptSpot`、`Board.IndustryConstituents`、`Board.ConceptConstituents`、`Board.IndustryKline`、`Board.ConceptKline`、`Board.IndustryMinute`、`Board.ConceptMinute`：东方财富行业/概念板块列表、盘口指标、成分股、历史 K 线和分钟线基础能力。
- `FundFlow.Individual`、`FundFlow.Market`、`FundFlow.Rank`、`FundFlow.SectorRank`、`FundFlow.SectorHistory`：东方财富个股、大盘和板块资金流基础能力。
- `Northbound.Minute`、`Northbound.Summary`、`Northbound.HoldingRank`、`Northbound.History`、`Northbound.Individual`：东方财富沪深港通/北向资金分时、汇总、持股排行、历史和个股持仓基础能力。
- `DragonTiger.Detail`、`DragonTiger.StockStats`、`DragonTiger.Institution`、`DragonTiger.BranchRank`、`DragonTiger.SeatDetail`：东方财富龙虎榜详情、个股统计、机构买卖、营业部排行和席位明细基础能力。
- `BlockTrade.MarketStat`、`BlockTrade.Detail`、`BlockTrade.DailyStat`：东方财富大宗交易市场统计、成交明细和每日个股统计基础能力。
- `Margin.AccountInfo`、`Margin.TargetList`：东方财富融资融券账户统计和标的明细基础能力。
- `Dividend.Detail`：东方财富个股分红派送详情基础能力。
- `Data.Search`、`Data.CodesCN`、`Data.CodesUS`、`Data.CodesHK`、`Data.CodesFund`、`Data.BlockTradeDetail`、`Data.MarginTargetList`、`Data.DividendDetail`：搜索、代码列表、大宗交易、融资融券和分红等数据类聚合入口。
- `DatacenterQuery`、`DatacenterResult`、`ParseDCDate`：东方财富 datacenter-web 通用分页查询参数、结果模型和日期解析辅助能力。
- `MarketEvent.ZTPool`、`MarketEvent.StockChanges`、`MarketEvent.BoardChanges`：东方财富涨停股池、个股盘口异动和板块异动基础能力。
- `Fund.Estimate`、`Fund.NavHistory`、`Fund.RankHistory`、`Fund.DividendList`：天天基金/东方财富基金当日实时估值、历史净值、同类排名走势和分红列表基础能力。
- `Futures.Kline`、`Futures.GlobalSpot`、`Futures.GlobalKline`、`Futures.InventorySymbols`、`Futures.Inventory`、`Futures.ComexInventory`：东方财富国内期货历史 K 线、全球期货实时行情、全球期货历史 K 线、国内期货库存和 COMEX 黄金/白银库存基础能力。
- 期货工具：`ExtractFuturesVariety`、`FuturesMarketCode`，对应 TS 版 Eastmoney 期货 K 线 provider 的品种提取和 market code 查找能力。
- `Options.CFFEXQuotes`、`Options.LHB`、`Options.IndexOptionSpot`、`Options.IndexOptionKline`、`Options.CommodityOptionSpot`、`Options.CommodityOptionKline`、`Options.ETFOptionMonths`、`Options.ETFOptionExpireDay`、`Options.ETFOptionMinute`、`Options.ETFOptionDailyKline`、`Options.ETFOption5DayMinute`：东方财富中金所期权实时行情、期权龙虎榜和新浪股指/商品/ETF 期权基础能力。
- `internal/`：请求治理、数据源 provider 和 service 编排目录。
- `cmd/ceheng`、`apps/api`、`apps/web`：预留 CLI、后端 API 和 Web 前端入口。

当前阶段聚焦 Go SDK 库；CLI、MCP、后端 API 和 Web 前端属于后续应用层。SDK v0.1 的收口边界见 `docs/sdk-v0.1-acceptance.md`。

## 使用示例

```go
package main

import (
	"context"
	"fmt"
	"time"

	stock "github.com/ceheng.io/stock-go"
	"github.com/ceheng.io/stock-go/symbols"
)

func main() {
	sdk := stock.New(stock.WithTimeout(10 * time.Second))
	_ = sdk

	symbol, err := symbols.Normalize("600519", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(symbols.ToTencent(symbol))

	quotes, err := sdk.Quotes.SimpleCN(context.Background(), []string{"sh600519"})
	if err != nil {
		panic(err)
	}
	fmt.Println(quotes[0].Name, quotes[0].Price)
}
```

## 开发验证

```bash
go test ./...
```

`apps/web` 作为后续前端应用目录有独立 `go.mod` 边界，用于避免根模块的 `go list ./...` 穿透 `node_modules` 中的第三方 Go 包。

写任何前端代码前，必须先执行本地 `skills/frontend-design` 并确认页面结构、交互状态、数据依赖和验收方式。
