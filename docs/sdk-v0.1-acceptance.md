# SDK v0.1 验收清单

本文用于记录 `github.com/ceheng-io/stock-go` SDK v0.1 的完成范围。v0.1 的交付物是 Go SDK 库，不包含 CLI、MCP、后端 API 和 Web 前端。

项目仓库：[ceheng-io/stock-go](https://github.com/ceheng-io/stock-go)。

## 验收结论

截至 2026-06-15，Go SDK v0.1 可以按当前库能力范围收口。

完成口径：

- Go module、根包 `stock`、公开子包、`internal` 分层和预留应用目录已按目录设计落地。
- 根包、服务字段和公开子包已覆盖主要行情、K 线、板块、资金、事件、基金、期货、期权和纯计算能力。
- 腾讯、东方财富、新浪、同花顺 provider 的 SDK 库能力已形成 provider/service/root 闭环。
- 错误体系、请求治理、缓存、符号、指标、信号、筛选、公开类型和常量已覆盖 v0.1 需要的公共能力。
- `types/` 已按领域拆分，并通过测试约束单个非测试 Go 文件不超过 1000 行。
- CLI、MCP、后端 API、Web 前端和浏览器专用 script mutex 已明确列为后续应用层或不适用能力。

## 已验收范围

### 工程结构

- `go.mod` 使用 `module github.com/ceheng-io/stock-go`。
- 根包名为 `stock`，提供 `stock.New()` 和 `stock.StockSDK` 入口。
- 公开纯能力包包括 `types`、`cache`、`symbols`、`indicators`、`signals`、`screener`、`parser`、`utils`、`errors`、`constants`、`timeutil`、`useragent`。
- 数据源适配和服务编排位于 `internal/providers`、`internal/services`、`internal/core`。
- `cmd/ceheng`、`apps/api`、`apps/web`、`testdata` 已作为后续扩展位置保留。
- `.gitignore` 已忽略 `docs/superpowers/`、构建产物、依赖目录、编辑器文件和环境变量文件。
- `AGENTS.md` 已记录中文协作约定、项目标识、Go 代码约定和前端设计前置要求。

### 公开 API

- 根包 `Client`/`StockSDK`、配置选项、服务字段和 `Client.Get*` 薄委托已覆盖主要 SDK 调用路径。
- `docs/api-matrix.md` 已记录根包、服务字段和子包能力索引。
- 根包保留常用兼容命名别名，例如 `JsonpRequest`、`FetchJsVars`、`SdkError`、`InferProviderFromUrl`、`SafeNumberOrNull`、`INDICATOR_REGISTRY`、`EM_PUSH_TOKEN`、`SINA_OPTION_API_URL` 等。
- 公开领域类型由 `types` 子包承载，根包为常用类型提供别名。

### 核心能力

- 请求治理覆盖 timeout、retry、指数退避、可重试状态码、网络错误/超时重试、User-Agent 轮换、限流、熔断、host fallback 和 provider policy。
- 错误码覆盖 `INVALID_ARGUMENT`、`INVALID_SYMBOL`、`PARSE_ERROR`、`UPSTREAM_ERROR`、`HTTP_ERROR`、`RATE_LIMITED`、`NETWORK_ERROR`、`TIMEOUT`、`ABORTED`、`NOT_FOUND` 等场景。
- 缓存覆盖内存 TTL/LRU、共享缓存、缓存 key、single-flight 和常用兼容命名。
- JSONP、JS 变量、GBK、数值解析、市场时间工具和请求常量已覆盖。

### 领域能力

- 腾讯：A 股、港股、美股、基金行情、批量行情、资金流、盘口大单、搜索、代码列表、交易日历、当日分时。
- 东方财富：A/HK/US K 线、板块、资金流、北向、龙虎榜、大宗交易、融资融券、分红、datacenter、市场事件、基金、期货、期权龙虎榜、中金所期权。
- 新浪：ETF、股指、商品期权现货、K 线、分钟线和 ETF 期权月份/到期日/分钟线/日 K/五日分钟线。
- 同花顺：涨停池和相关市场事件数据。
- 纯计算：技术指标、信号识别、本地筛选和轻量回测。

## 后续范围

以下内容不作为 SDK v0.1 完成条件：

- `cmd/ceheng`：后续作为 CLI 独立目标实现。
- MCP server/tools：后续作为独立目标实现。
- 后端 API：后续基于当前 SDK 在 `apps/api` 规划。
- Web 前端：后续基于当前 SDK/API 在 `apps/web` 规划；写前端前必须先执行本地 `skills/frontend-design`。
- `core/scriptMutex.ts`：浏览器 `<script>` 注入互斥锁不适用于 Go SDK 服务端实现；当前只保留 `BROWSER_JSVARS_MUTEX_KEY` 常量。

## v0.1 后续 Backlog

这些事项可以持续改进，但不阻塞当前 v0.1 收口：

- 继续抽样审计 provider 字段兼容、`null` 与零值/指针、nil slice 与空 slice、非数组 payload 等细节。
- 为高频真实接口补充更多 fixture，降低上游字段漂移带来的解析风险。
- 如新增便捷 wrappers 或应用层入口，同步维护 `docs/api-matrix.md`。
- 启动 CLI、MCP、API 或 Web 时，分别建立独立设计文档、验收清单和测试策略。

## 收口验证

收口时至少执行：

```bash
go test -count=1 ./...
go test -count=1 ./types -run 'Test(TypesFilesStaySmall|DomainTypesStayInDomainFiles)'
git diff --check
```
