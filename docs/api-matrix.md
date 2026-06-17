# 公开 API 速查

本文记录 `github.com/ceheng-io/stock-go` 当前公开入口，方便接入方按领域查找根包、服务字段和公开子包能力。

项目仓库：[ceheng-io/stock-go](https://github.com/ceheng-io/stock-go)。

## 根入口

| 能力 | 入口 | 说明 |
| --- | --- | --- |
| 创建客户端 | `stock.New(options...)` | 使用函数式选项构造 `*stock.Client`。 |
| 根入口类型 | `stock.Client`、`stock.StockSDK` | `StockSDK` 为根入口类型别名。 |
| 服务字段 | `client.Quotes`、`client.Kline`、`client.Board` 等 | 推荐按领域浏览和调用能力。 |
| 便捷方法 | `client.GetSimpleQuotes`、`client.GetHistoryKline`、`client.GetTHSLimitUpPool` 等 | 适合偏好扁平方法名的调用方。 |
| 服务类型 | `stock.QuoteService`、`stock.TradingCalendarService`、`stock.FundService` 等 | 外部无需 import `internal/services`。 |

## 服务字段

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

## 常用便捷方法

| 能力 | 方法 |
| --- | --- |
| A 股行情 | `GetSimpleQuotes`、`GetFullQuotes`、`GetAllAShareQuotes`、`GetAllQuotesByCodes` |
| 港股/美股/基金行情 | `GetHKQuotes`、`GetUSQuotes`、`GetFundQuotes`、`GetAllHKShareQuotes`、`GetAllUSShareQuotes` |
| 代码列表 | `GetAShareCodeList`、`GetHKCodeList`、`GetUSCodeList`、`GetFundCodeList` |
| K 线 | `GetHistoryKline`、`GetMinuteKline`、`GetHKHistoryKline`、`GetUSHistoryKline`、`GetKlineWithIndicators` |
| 板块 | `GetIndustryList`、`GetConceptList`、`GetIndustryConstituents`、`GetConceptConstituents` |
| 资金与事件 | `GetIndividualFundFlow`、`GetNorthboundHoldingRank`、`GetZTPool`、`GetStockChanges`、`GetBoardChanges` |
| 龙虎榜/大宗/两融 | `GetDragonTigerDetail`、`GetBlockTradeDeals`、`GetMarginSummary` |
| 基金/期货/期权 | `GetFundEstimate`、`GetFuturesKline`、`GetGlobalFuturesSpot`、`GetETFOptionMinute`、`GetOptionLHB` |
| 工具能力 | `Search`、`BatchRaw`、`GetTradingCalendar`、`IsTradingDay`、`NextTradingDay`、`PrevTradingDay`、`GetMarketStatus` |

## 公开子包

| 子包 | 能力 |
| --- | --- |
| `types` | 行情、K 线、板块、资金流、北向、龙虎榜、大宗交易、两融、分红、基金、期货、期权、事件等领域类型 |
| `indicators` | MA、MACD、BOLL、KDJ、RSI、WR、BIAS、CCI、ATR、OBV、ROC、DMI、SAR、KC 等技术指标 |
| `signals` | 金叉/死叉、超买/超卖、BOLL 突破、SAR 反转等信号识别 |
| `symbols` | 统一符号模型、市场推断、腾讯/东方财富代码转换、期货品种映射 |
| `screener` | 链式本地选股和轻量回测 |
| `cache` | 内存 TTL/LRU 缓存、共享缓存、缓存 key、single-flight |
| `errors` | 错误类型、错误码提取、结构化错误判断和错误元数据工具 |
| `parser` | 腾讯行情文本解析 |
| `utils` | 分片、并发池、周期/复权参数工具 |
| `constants` | 数据源 URL、请求默认值、期货交易所映射等常量 |
| `timeutil` | 市场时区、时间解析和时间元信息工具 |
| `useragent` | 内置浏览器 User-Agent 池 |

## 根包工具

| 类型 | 代表入口 |
| --- | --- |
| 指标与信号 | `CalcMA`、`CalcMACD`、`CalcSignals`、`Screen`、`Backtest` |
| 缓存 | `NewMemoryCache`、`CreateCacheKey`、`CacheThrough`、`ClearSharedCaches` |
| 符号 | `NormalizeSymbol`、`ToTencentSymbol`、`ToEastmoneySecID`、`ToEastmoneySecid` |
| 响应解析 | `DecodeGBK`、`ParseResponse`、`ExtractJSONP`、`ExtractJSVar`、`ParseJSVars`、`FetchJSVars` |
| 请求治理 | `WithTimeout`、`WithRetry`、`WithRateLimit`、`WithCircuitBreaker`、`WithProviderPolicy`、`WithRequestHooks` |
| 错误处理 | `GetErrorCode`、`GetSDKErrorCode`、`GetSdkErrorCode`、`IsSDKError`、`IsSdkError` |
| 时间 | `ParseMarketTime`、`BuildTimeMeta`、`FormatInTz`、`MARKET_TZ`、`MarketTZ` |

## 维护约定

- 新增公开服务字段、便捷方法、公开子包或根包工具时，同步更新本文。
- 根包保持组合与薄委托职责，provider 解析逻辑放在 `internal/`。
- 新增公开领域类型优先放入 `types/` 对应领域文件。
