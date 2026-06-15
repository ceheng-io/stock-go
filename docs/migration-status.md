# stock-go 迁移状态

本文记录 `/Users/xingyys/project/html/stock-sdk` 到 `github.com/ceheng.io/stock-go` 的当前迁移状态。当前阶段只建设 Go SDK 库；CLI、MCP、后端 API 和 Web 前端属于后续应用层。

## SDK v0.1 收口结论

截至 2026-06-15，第一阶段 Go SDK 库迁移可以按 v0.1 收口。完成边界和延期项见 `docs/sdk-v0.1-acceptance.md`；后续 CLI、MCP、后端 API 和 Web 前端应作为独立目标重新规划。

## 当前已覆盖

- 根包 `stock`：`Client`/`StockSDK`、`New`、配置选项、错误类型、请求治理配置、主要服务字段、TS 扁平方法对应的 `Client.Get*` 薄委托、根包 service 类型别名、纯工具薄导出，以及大部分公开领域类型别名；公开 API 映射见 `docs/api-matrix.md`。
- `cache`：内存 TTL/LRU 缓存、共享缓存、缓存 key、single-flight `CacheThrough`，并在根包和 `cache` 子包同时保留 TS 风格缓存命名；默认 `maxSize=1000`、`CreateCacheKey` 跳过 `nil`、`MemoryCache.Has`、`MemoryCache.Delete` 返回是否删除、`MemoryCache.GetOrFetch` 实例方法、`ClearSharedCaches` 清空内容但保留 namespace 实例、`MemoryCache.Clear` 同步清理 in-flight 表等行为已对齐 TS 版。
- `errors`：`stock-go/errors` 子包兼容 TS `stock-sdk/errors` 子入口，复用根包错误类型、构造函数和 `GetSdkErrorCode` / `IsSdkError` 等工具。
- `utils`：`ChunkArray`、`AsyncPool`、K 线周期/分钟周期/复权参数校验、`PeriodCode`/`GetPeriodCode`、`AdjustCode`/`GetAdjustCode` 已覆盖；工具参数非法时按 TS `InvalidArgumentError` 语义映射为 `INVALID_ARGUMENT` 错误码。
- 脚本响应解析工具：`ExtractJSONP`/`ExtractJsonFromJsonp`、`JSONPRequest`/`JsonpRequest`、`ExtractJSVar`、`ParseJSVars`/`ParseJsVars`、`FetchJSVars`/`FetchJsVars` 已覆盖；JSONP 文本解析失败按 TS `PARSE_ERROR` 语义映射错误码。
- 请求治理：默认 timeout、retry 次数和 retry base delay 已对齐 TS 默认值（30s、3 次、1s），`RetryOptions` 已覆盖 TS 的 `maxDelay`、`backoffMultiplier`、`retryableStatusCodes`、`retryOnNetworkError`、`retryOnTimeout`、`onRetry` 语义，`rotateUserAgent` 已在全局配置和 provider policy 中对齐 TS 可选布尔覆盖行为，并保留 Go 侧 provider policy、限流、熔断和 host fallback 能力；HTTP 状态、429 限流、网络错误、请求超时/取消和 JSON 解析失败已按 TS `HttpError`/`SdkError` 语义映射为 `HTTP_ERROR`、`RATE_LIMITED`、`NETWORK_ERROR`、`TIMEOUT`、`ABORTED`、`PARSE_ERROR` 错误码，同时保留重试和 host fallback 识别；`TENCENT_BASE_URL`、`EM_KLINE_URL`、`FUTURES_EXCHANGE_MAP`、`DEFAULT_RETRYABLE_STATUS_CODES` 等 TS 风格核心常量/映射已在根包和 `constants` 子包保留；交易日历 `NextTradingDay`/`PrevTradingDay` 越界按 TS `InvalidArgumentError` 语义映射为 `INVALID_ARGUMENT` 错误码。
- `symbols`：统一符号模型、A 股交易所推断、腾讯/东方财富代码转换、期货品种和交易所映射，并保留 `FUTURES_EXCHANGES` 子包常量名；`Normalize` 非法标的按 TS `InvalidSymbolError` 语义映射为 `INVALID_SYMBOL` 错误码，腾讯/东方财富符号适配器遇到不支持的 market/exchange/assetType 时按 TS `InvalidArgumentError` 语义映射为 `INVALID_ARGUMENT` 错误码。
- `indicators`：MA/SMA/EMA/WMA、MACD、BOLL、KDJ、RSI、WR、BIAS、CCI、ATR、OBV、ROC、DMI、SAR、KC、指标上下文、启用 key 顺序、lookback 估算和 TS 风格注册表命名。
- `signals`：基于已计算指标的金叉/死叉、超买/超卖、BOLL 突破和 SAR 反转；`CalcSignals` 缺失 MA/RSI 指标 key 时按 TS `InvalidArgumentError` 语义映射为 `INVALID_ARGUMENT` 错误码。
- `screener`：链式本地筛选和轻量全仓多头回测；`Top(n)` 负数参数按 TS `InvalidArgumentError` 语义映射为 `INVALID_ARGUMENT` 错误码。
- `types`：行情、K 线、板块、资金流、北向、龙虎榜、大宗交易、两融、分红、基金、期货、期权、事件等公开类型，按领域拆分并限制单文件不超过 1000 行。
- `internal/core`：腾讯行情文本解析、HTTP client、重试、限流、熔断、host fallback、JSONP、JS 变量解析、provider policy。
- `internal/providers/tencent`：A/HK/US/基金行情、资金流、盘口大单、搜索、交易日历、代码列表、分时；行情、基金和资金流解析已补齐 `timestamp`/`tz` 时间元信息，腾讯当日分时 `preClose` 已按 TS 可选字段语义改为 Go 指针字段，且缺少股票数据或昨收字段时按 TS `parseFloat(...) || 0` 语义返回 `0`，接口返回 `code != 0` 时按 TS `UpstreamError` 语义映射为 `UPSTREAM_ERROR` 错误码；Smartbox 搜索的空关键词、`v_hint` 空/N、Unicode 解码、代码拼接和资产类型归一化已按 TS 语义覆盖；批量行情负数 `batchSize`/`concurrency` 按 TS `InvalidArgumentError` 语义映射为 `INVALID_ARGUMENT` 错误码；交易日历在 service 层按 TS `getSharedCache('tencent:trade-calendar').getOrFetch(...)` 语义复用按 `CalendarURL` 分桶的模块级共享结果、合并并发冷启动请求并返回副本；股票/基金代码列表默认 URL 已对齐 TS 的 linkdiary 端点，其中基金代码列表按 TS 纯文本逗号格式解析。
- `internal/providers/eastmoney`：A/HK/US K 线、板块、资金流、北向、龙虎榜、大宗交易、融资融券、分红、数据中心、市场事件、基金、期货、期权龙虎榜和中金所期权；A 股历史/分钟 K 线已补齐 `timestamp`/`tz` 时间元信息，A 股历史 K 线请求字段已对齐 TS 的 `f116`，Eastmoney push token 已对齐 TS `EM_PUSH_TOKEN`，HK/US `period=1` 分时默认 `ndays=1` 且保留显式传入值，A/HK/US 分钟线、板块分钟线和资金流历史遇到非数组 `trends`/`klines` 时按 TS `Array.isArray` 语义返回空结果，涨停池、个股异动和板块异动遇到非数组 `pool`/`allstock`/`allbk` 时按 TS 语义返回空结果，全球期货现货遇到非数组 `list` 时按 TS 语义停止翻页并返回已累计数据，A/HK/US 历史/分钟 K 线、板块历史/分钟 K 线、市场事件枚举、资金流 period/indicator/sectorType、北向持股排行 period/market、龙虎榜统计周期、国内/全球期货 K 线和 COMEX 库存的非法入参按 TS `InvalidArgumentError` 语义映射为 `INVALID_ARGUMENT` 错误码，US K 线空标的按 `InvalidSymbolError` 语义映射为 `INVALID_SYMBOL` 错误码，行业/概念板块服务支持按名称解析为 BK 代码且未命中时按 TS `NotFoundError` 语义映射为 `NOT_FOUND` 错误码，板块列表/成分股和资金流排名已对齐 TS 的自动翻页行为，板块列表/成分股及资金流 clist 翻页遇到非数组 `diff` 时按 TS 语义停止并返回已累计数据，板块 spot 缺失 `data` 时按 TS 语义返回空结果，通用 datacenter 翻页器遇到缺失 `result` 或非数组 `result.data` 时按 TS `Array.isArray` 语义停止翻页并返回已累计数据，龙虎榜营业部/席位金额字段，以及大宗交易、两融账户、涨停池、分红详情、北向持股排行名称和期货库存代码中的兼容字段回退均已按 TS `??` 的 nullish 语义处理，分红详情文本字段会按 TS 保留空字符串，基金估值空/坏 JSONP 响应会按 TS 返回 fallback code 与空字段，基金历史名称/时间戳、基金分红日期、股票分红详情、涨停池封板时间、涨停池空字符串封板时间、北向历史领涨股和板块资金流领涨股已按 TS 可空/可选语义改为 Go 指针字段，中金所期权和期权龙虎榜遇到非数组 payload 时按 TS `Array.isArray` 语义返回空结果。
- `internal/providers/sina`：ETF、股指、商品期权现货、K 线和分钟线；新浪 ETF 期权月份缺失数组字段时按 TS `?? []` 语义返回非 nil 空切片，新浪期权日 K 线遇到非数组 JSONP payload 时按 TS `Array.isArray` 语义返回空结果，ETF 期权分钟线和 5 日分钟线也按 TS 跳过非数组 payload / 非数组分组，商品期权未知品种按 TS `InvalidArgumentError` 语义映射为 `INVALID_ARGUMENT` 错误码。
- `internal/services`：Quotes、Calendar、Kline、Indicator、Board、FundFlow、Northbound、DragonTiger、BlockTrade、Margin、Dividend、Data、MarketEvent、Fund、Futures、Options。
- 工程卫生：`apps/web` 使用独立 `go.mod` 作为前端应用边界，根模块 `go list ./...` 不会穿透 `node_modules` 中的第三方 Go 包；根目录已按设计预留 `testdata/` fixture 目录；根包测试已覆盖这些结构约束。

## 等价但不逐字照搬

- TS `StockSDK` 的大量扁平方法在 Go 中优先推荐使用 `Client` 上的服务字段，例如 `client.Quotes.CN`、`client.Kline.CN`、`client.Board.IndustryList`；同时根 `Client.Get*` 保留薄委托作为迁移便利入口，并对 TS 中常见的 `options?`、`date?`、`period?`、`direction?` 形态提供可省略参数。
- TS batch options 中 `batchSize`/`concurrency` 传入 `0` 会触发参数错误；Go 中 `BatchOptions{}` 和零值字段按惯例表示使用默认值，负数仍视为非法参数。
- TS 顶层 service class 在 Go 中以 `Client` 服务字段和根包类型别名表达，例如 `stock.TradingCalendarService`、`stock.FundService`、`stock.QuoteService`。
- TS v2 namespace getter 在 Go 中用显式服务字段表达，不额外构造懒加载 namespace。
- TS 顶层的纯工具 export 在 Go 中以根包薄转发和公开子包并存，避免根包承载实现逻辑；`SafeNumberOrNull`、`GetPeriodCode`、`GetAdjustCode`、`ExtractJsonFromJsonp`、`JsonpRequest`、`JsonpOptions`、`FetchJsVars`、`FetchJsVarsOptions`、`SdkError`、`HttpError`、`MemoryCacheStore`、`NewMemoryCache`、`CreateCacheKey`、`CacheThrough`、`BuildIndicatorContext`、`GetEnabledIndicatorKeys`、`EstimateIndicatorLookback`、`INDICATOR_REGISTRY`、`RequestClientOptions`、`ProviderRequestPolicy`、`InferProviderFromUrl`、`ParseDcDate`、`MARKET_TZ`、`BROWSER_JSVARS_MUTEX_KEY`、`GetNextUserAgent`、`GetRandomUserAgent`、`GetAllUserAgents`、`SymbolInput`、`ToEastmoneySecid`、`TENCENT_BASE_URL`、`EM_PUSH_TOKEN`、`SINA_OPTION_API_URL`、`COMMODITY_OPTION_MAP`、`symbols.FUTURES_EXCHANGES` 等高频 TS 命名已保留兼容别名；`GetAdjustCode` 同时接受 TS 的 `""` 和 Go 的 `AdjustNone`/`"none"`。
- TS 顶层 `types` 大部分已在 Go 根包做类型别名；行业/概念板块类型名和 `StockFundFlowDaily` 等 TS 兼容别名已保留；`types.Quote`、`types.AnyHistoryKline` 联合类型在 Go 中用封闭接口表达；`ZTPoolType`、`StockChangeType` 等市场事件字符串枚举，`NorthboundDirection`、`NorthboundMarket`、`NorthboundRankPeriod` 等北向资金字符串枚举，`DragonTigerPeriod`、`DragonTigerDateOptions`、`BlockTradeDateOptions` 等龙虎榜/大宗交易领域选项，`FundDividendRank`、`FundSortDirection`、`FundDividendListOptions` 等基金分红领域选项，`ETFOptionCate`、`IndexOptionProduct` 等期权字符串枚举，以及 `types.FuturesExchange` 期货交易所字符串枚举由 `types` 子包承载；`types.Market` 与根包符号模型 `Market`、`types.FuturesExchange` 与根包符号元数据 `FuturesExchange` 存在命名冲突，保留在 `types` 子包，根包分别用 `SymbolMarketCN`、`SymbolMarketHK`、`SymbolMarketUS`、`SymbolMarketGlobal` 和 `FuturesExchangeCode*` 表达领域枚举常量。
- `types.TencentQuoteItem` 与 parser 原始项同名，根包保留 parser 版本；领域类型版本可通过 `types.TencentQuoteItem` 使用。

## 明确不属于第一阶段

- `src/cli/*`：Go 仓库仅预留 `cmd/ceheng/`，第一阶段不实现完整 CLI。
- `src/mcp/*`：MCP server/tools 属于后续应用层，第一阶段不迁移。
- Web 前端和后端 API：仅预留 `apps/web/`、`apps/api/`。写前端前必须先执行 `skills/frontend-design`。
- `core/scriptMutex.ts`：这是浏览器 `<script>` 注入场景的全局名互斥锁，Go SDK 服务端实现不需要等价能力；根包仅保留 `BROWSER_JSVARS_MUTEX_KEY` 常量，方便迁移代码引用固定 key。

## v0.1 后续 Backlog

- 持续维护 `docs/api-matrix.md`，后续若新增 convenience wrappers 或应用层入口，需同步映射关系。
- 继续抽样检查 provider 解析字段细节，尤其是接口返回兼容字段和空值语义；这些属于兼容性加固，不阻塞 SDK v0.1 收口。
- 后续如果启动 CLI/MCP/API/Web，需要基于当前 SDK 重新规划应用层目录和验收方式。

## 当前验证命令

```bash
go test ./...
go test ./types -run 'Test(TypesFilesStaySmall|DomainTypesStayInDomainFiles)' -count=1
```
