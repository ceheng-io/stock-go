# 公开 API 迁移矩阵

本文按 TypeScript 版 `/Users/xingyys/project/html/stock-sdk/src` 的顶层导出和 `StockSDK` 扁平方法，记录 Go 版 `github.com/ceheng.io/stock-go` 当前映射。当前阶段只建设 Go SDK 库，不迁移 CLI、MCP、后端 API 和 Web 前端。

## 根入口

| TS 入口 | Go 入口 | 状态 | 说明 |
| --- | --- | --- | --- |
| `new StockSDK(options)` | `stock.New(options...)`、`stock.StockSDK` | 已覆盖 | Go 使用函数式选项构造 `*stock.Client`；`StockSDK` 保留为根入口类型别名。 |
| `sdk.quotes.*` / `getFullQuotes` 等 | `client.Quotes.*`、`client.GetFullQuotes` 等 | 已覆盖 | 推荐服务字段；`Client.Get*` 作为 TS 扁平方法迁移便利入口，并支持省略常见可选参数。 |
| `sdk.codes.*` / 代码列表方法 | `client.Quotes.Codes*`、`client.Data.Codes*`、`client.Get*CodeList` | 已覆盖 | 代码列表同时在 Quotes 与 Data 聚合入口保留。 |
| `sdk.batch.*` / 全市场批量行情 | `client.Quotes.Batch*`、`client.Quotes.All*`、`client.GetAll*` | 已覆盖 | Go service 层仍将 code-list options 与 batch options 分开；`Client.GetAll*` 兼容 TS 聚合选项。 |
| `sdk.kline.*` / K 线方法 | `client.Kline.*`、`client.Indicator.KlineWithIndicators`、`client.Get*Kline` | 已覆盖 | A/HK/US 历史和分钟 K 线均已迁移。 |
| `sdk.board.*` / 板块方法 | `client.Board.*`、`client.GetIndustry*`、`client.GetConcept*` | 已覆盖 | 行业/概念列表、盘口、成分股、K 线、分钟线。 |
| `sdk.calendar.*` / 日历方法 | `client.Calendar.*`、`client.Quotes.TradingCalendar`、`client.GetTradingCalendar` 等 | 已覆盖 | 市场状态为本地同步判断，交易日依赖交易日历。 |
| `sdk.reference.*` | `client.Quotes.DividendDetail`、`client.Quotes.TradingCalendar`、`client.GetDividendDetail` | 已覆盖 | Go 未单独建 `Reference` 字段，能力保留在 Quotes/Data 和根薄委托。 |
| 顶层 service 类型 | `stock.QuoteService`、`TradingCalendarService`、`FundService` 等 | 已覆盖 | `Client` 上所有公开服务字段均提供根包类型别名，外部无需 import `internal/services`。 |

## `StockSDK` 扁平方法

| TS 方法 | Go 映射 | 状态 |
| --- | --- | --- |
| `getFullQuotes` | `client.Quotes.CN`、`client.GetFullQuotes` | 已覆盖 |
| `getSimpleQuotes` | `client.Quotes.SimpleCN`、`client.GetSimpleQuotes` | 已覆盖 |
| `getHKQuotes` | `client.Quotes.HK`、`client.GetHKQuotes` | 已覆盖 |
| `getUSQuotes` | `client.Quotes.US`、`client.GetUSQuotes` | 已覆盖 |
| `getFundQuotes` | `client.Quotes.Fund`、`client.GetFundQuotes` | 已覆盖 |
| `getFundFlow` | `client.Quotes.FundFlow`、`client.GetFundFlow` | 已覆盖 |
| `getPanelLargeOrder` | `client.Quotes.PanelLargeOrder`、`client.GetPanelLargeOrder` | 已覆盖 |
| `getTodayTimeline` | `client.Quotes.TodayTimeline`、`client.GetTodayTimeline` | 已覆盖 |
| `getIndustryList` / `getConceptList` | `client.Board.IndustryList` / `ConceptList`、`client.GetIndustryList` / `GetConceptList` | 已覆盖 |
| `getIndustrySpot` / `getConceptSpot` | `client.Board.IndustrySpot` / `ConceptSpot`、`client.GetIndustrySpot` / `GetConceptSpot` | 已覆盖 |
| `getIndustryConstituents` / `getConceptConstituents` | `client.Board.IndustryConstituents` / `ConceptConstituents`、`client.GetIndustryConstituents` / `GetConceptConstituents` | 已覆盖 |
| `getIndustryKline` / `getConceptKline` | `client.Board.IndustryKline` / `ConceptKline`、`client.GetIndustryKline` / `GetConceptKline` | 已覆盖 |
| `getIndustryMinuteKline` / `getConceptMinuteKline` | `client.Board.IndustryMinute` / `ConceptMinute`、`client.GetIndustryMinuteKline` / `GetConceptMinuteKline` | 已覆盖 |
| `getHistoryKline` / `getMinuteKline` | `client.Kline.CN` / `CNMinute`、`client.GetHistoryKline` / `GetMinuteKline` | 已覆盖 |
| `getHKHistoryKline` / `getHKMinuteKline` | `client.Kline.HK` / `HKMinute`、`client.GetHKHistoryKline` / `GetHKMinuteKline` | 已覆盖 |
| `getUSHistoryKline` / `getUSMinuteKline` | `client.Kline.US` / `USMinute`、`client.GetUSHistoryKline` / `GetUSMinuteKline` | 已覆盖 |
| `search` | `client.Quotes.Search`、`client.Data.Search`、`client.Search` | 已覆盖 |
| `getAShareCodeList` / `getUSCodeList` | `client.Quotes.CodesCN` / `CodesUS`、`client.GetAShareCodeList` / `GetUSCodeList` | 已覆盖 |
| `getHKCodeList` / `getFundCodeList` | `client.Quotes.CodesHK` / `CodesFund`、`client.GetHKCodeList` / `GetFundCodeList` | 已覆盖 |
| `getAllAShareQuotes` / `getAllHKShareQuotes` / `getAllUSShareQuotes` | `client.Quotes.AllCN` / `AllHK` / `AllUS`、`client.GetAllAShareQuotes` / `GetAllHKShareQuotes` / `GetAllUSShareQuotes` | 已覆盖 |
| `getAllQuotesByCodes` | `client.Quotes.BatchCN`、`client.GetAllQuotesByCodes` | 已覆盖 |
| `batchRaw` | `client.Quotes.BatchRaw`、`client.BatchRaw` | 已覆盖 |
| `getTradingCalendar` | `client.Quotes.TradingCalendar`、`client.GetTradingCalendar` | 已覆盖 |
| `isTradingDay` / `nextTradingDay` / `prevTradingDay` | `client.Calendar.IsTradingDay` / `NextTradingDay` / `PrevTradingDay`、`client.IsTradingDay` / `NextTradingDay` / `PrevTradingDay` | 已覆盖 |
| `getMarketStatus` | `client.Calendar.MarketStatus`、`client.GetMarketStatus` | 已覆盖 |
| `getDividendDetail` | `client.Quotes.DividendDetail`、`client.Dividend.Detail`、`client.Data.DividendDetail`、`client.GetDividendDetail` | 已覆盖 |
| `getFuturesKline` / `getGlobalFuturesSpot` / `getGlobalFuturesKline` | `client.Futures.Kline` / `GlobalSpot` / `GlobalKline`、`client.GetFuturesKline` / `GetGlobalFuturesSpot` / `GetGlobalFuturesKline` | 已覆盖 |
| `getFuturesInventorySymbols` / `getFuturesInventory` / `getComexInventory` | `client.Futures.InventorySymbols` / `Inventory` / `ComexInventory`、`client.GetFuturesInventorySymbols` / `GetFuturesInventory` / `GetComexInventory` | 已覆盖 |
| 期权方法 | `client.Options.*`、`client.Get*Option*` / `GetOptionLHB` | 已覆盖 |
| `getKlineWithIndicators` | `client.Indicator.KlineWithIndicators`、`client.GetKlineWithIndicators` | 已覆盖 |
| 深度资金流方法 | `client.FundFlow.*`、`client.GetIndividualFundFlow` 等 | 已覆盖 |
| 北向资金方法 | `client.Northbound.*`、`client.GetNorthbound*` | 已覆盖 |
| 涨停/异动方法 | `client.MarketEvent.*`、`client.GetZTPool` / `GetStockChanges` / `GetBoardChanges` / `GetTHSLimitUpPool` | 已覆盖 |
| 龙虎榜方法 | `client.DragonTiger.*`、`client.GetDragonTiger*` | 已覆盖 |
| 大宗交易方法 | `client.BlockTrade.*`、`client.Data.BlockTrade*`、`client.GetBlockTrade*` | 已覆盖 |
| 融资融券方法 | `client.Margin.*`、`client.Data.Margin*`、`client.GetMargin*` | 已覆盖 |
| 基金扩展方法 | `client.Fund.Estimate`、`NavHistory`、`RankHistory`、`DividendList`、`client.GetFund*` | 已覆盖 |

## 顶层工具导出

| TS 导出 | Go 映射 | 状态 |
| --- | --- | --- |
| 指标函数与类型 | 根包薄转发 + `indicators` 子包 | 已覆盖 | 根包和 `indicators` 子包均保留 `BuildIndicatorContext`、`GetEnabledIndicatorKeys`、`EstimateIndicatorLookback`、`IndicatorRegistry` / `INDICATOR_REGISTRY` 等 TS 风格命名。 |
| 公开领域类型 | 根包类型别名 + `types` 子包 | 已覆盖 | 行业/概念板块类型、`StockFundFlowDaily`、`AnyHistoryKline`、`ZTPoolType`、`StockChangeType`、`THSLimitUpPoolOptions`、`THSLimitUpPoolResult`、`NorthboundDirection`、`NorthboundMarket`、`NorthboundRankPeriod`、`DragonTigerPeriod`、`DragonTigerDateOptions`、`BlockTradeDateOptions`、`FundDividendRank`、`ETFOptionCate`、`IndexOptionProduct`、`types.FuturesExchange` 等 TS 兼容别名、联合类型、字符串枚举和领域选项已保留。 |
| 公开选项类型 | 根包类型别名/兼容结构 | 已覆盖 | `GetAShareCodeListOptions`、`GetAllAShareQuotesOptions`、`GetAllHKQuotesOptions`、`GetUSCodeListOptions`、`GetAllUSQuotesOptions` 等 TS 顶层名称已保留；Go service 方法仍按 code-list options 与 batch options 拆分入参。 |
| `calcSignals` | `stock.CalcSignals` + `signals` 子包 | 已覆盖 |
| `screen` / `backtest` | `stock.Screen` / `stock.Backtest` + `screener` 子包 | 已覆盖 |
| `MemoryCacheStore` / `MemoryCache.getOrFetch` / `cacheThrough` | `stock.MemoryCacheStore`、`stock.NewMemoryCache().GetOrFetch` / `stock.CacheThrough` + `cache` 子包 | 已覆盖 | 根包和 `cache` 子包均保留 `MemoryCache`、`MemoryCacheStore`、`CacheOptions`、`CacheStore`、`NewMemoryCache`、`CreateCacheKey`、`CacheThrough` 等 TS 风格命名；强类型场景推荐 `CacheThrough[T]`。 |
| `normalizeSymbol` / symbol adapters | `stock.NormalizeSymbol`、`ToTencentSymbol`、`ToEastmoneySecID`/`ToEastmoneySecid`、`SymbolInput` 等 + `symbols` 子包 | 已覆盖 | 根包和 `symbols` 子包均保留 TS 风格函数名；`symbols.FUTURES_EXCHANGES` 对齐 TS `symbols` 子入口常量。 |
| `decodeGBK` / `parseResponse` / number helpers | `stock.DecodeGBK` / `ParseResponse` / `SafeNumber` / `SafeNumberOrNull` 等 | 已覆盖 |
| JSONP / JS 变量解析 | `stock.ExtractJSONP`、`ExtractJsonFromJsonp`、`JSONPRequest` / `JsonpRequest`、`ExtractJSVar`、`ParseJSVars` / `ParseJsVars`、`FetchJSVars` / `FetchJsVars` | 已覆盖 | `JsonpOptions`、`JSONPOptions`、`FetchJsVarsOptions`、`FetchJSVarsOptions`、`BROWSER_JSVARS_MUTEX_KEY` 等 TS/Go 命名均已保留；Go 请求方法仍推荐使用函数式选项。 |
| `chunkArray` / `asyncPool` / 周期工具 | `stock.ChunkArray`、`AsyncPool`、`GetPeriodCode`、`GetAdjustCode` 等 | 已覆盖 | `GetAdjustCode` 同时接受 TS 的 `""` 与 Go 的 `AdjustNone`/`"none"`，均映射为东方财富 `fqt=0`。 |
| 错误体系 | `stock.Error`、`SdkError`、`HttpError`、`RequestError`、`GetErrorCode`、`GetSdkErrorCode`、`IsSDKError`、`IsSdkError`、`stock-go/errors` | 已覆盖 | 保留 TS 风格名称，同时提供 Go initialism 风格名称；`stock-go/errors` 对应 TS `stock-sdk/errors` 子入口，并提供 `AttachErrorMetadata` / `NormalizeRequestError`。 |
| `MARKET_TZ` / time helpers | `stock.MARKET_TZ` / `stock.MarketTZ`、`ParseMarketTime`、`BuildTimeMeta`、`FormatInTz` | 已覆盖 |
| User-Agent、限流、熔断、host fallback、provider policy、核心常量 | 根包对应工具 | 已覆盖 | `GetNextUserAgent` / `GetRandomUserAgent` / `GetAllUserAgents` 已保留 TS 风格命名；`RetryOptions.onRetry`、`rotateUserAgent` 已支持全局配置和 provider policy 覆盖；`InferProviderFromURL` / `InferProviderFromUrl` 均保留；`TENCENT_BASE_URL`、`EM_KLINE_URL`、`FUTURES_EXCHANGE_MAP`、`DEFAULT_RETRYABLE_STATUS_CODES` 等 TS 风格大写下划线常量/映射已在根包和 `constants` 子包保留。 |
| `withScriptMutex` / `BROWSER_JSVARS_MUTEX_KEY` | `BROWSER_JSVARS_MUTEX_KEY` | 部分覆盖 | 浏览器 `<script>` 注入互斥锁第一阶段 Go SDK 服务端不需要；固定 key 常量保留为迁移兼容。 |

## v0.1 后续审计项

- 继续抽样审计 provider 字段、空值语义和时间字段，尤其是 TS 中 `null` 与 Go 零值/指针的对应关系。
- SDK v0.1 收口边界见 `docs/sdk-v0.1-acceptance.md`；CLI、MCP、后端 API 和 Web 前端不属于本矩阵的第一阶段完成条件。
