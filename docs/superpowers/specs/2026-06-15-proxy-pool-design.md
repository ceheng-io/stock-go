# 代理池轮换设计

## 背景

策衡 Go SDK 当前已支持统一 `HTTPClient` 注入、超时、重试、限速、熔断、数据源策略、fallback host 和请求 hooks。用户可以手动通过 `WithHTTPClient` 注入带代理的 `http.Transport`，但这要求了解 Go 网络栈，也不能直接表达“多个代理按请求轮换”。

本设计新增一层 SDK 内建代理池能力，让高频请求场景可以通过多个代理分散出口 IP。第一版聚焦全局代理池，每次 HTTP attempt 使用下一个代理。

## 目标

- 提供简单公开 API 配置多个代理 URL。
- 首次请求、重试请求、fallback host 请求都参与轮换。
- 不配置代理池时保持现有行为不变。
- 保持 `WithHTTPClient` 的高级扩展能力。
- 实现线程安全，支持 SDK 并发请求。

## 非目标

- 第一版不做代理健康检查、失败剔除、权重、地区选择。
- 第一版不做 provider 级代理池。
- 第一版不改变限速、熔断、重试和 fallback 的语义。
- 第一版不绕过或规避上游服务规则；代理仅作为用户可控网络出口配置。

## 公开 API

新增 `ProxyPoolOptions`：

```go
type ProxyPoolOptions struct {
    URLs []string
}
```

新增快捷选项：

```go
func WithProxyPool(urls []string) Option
func WithProxyPoolOptions(options ProxyPoolOptions) Option
```

示例：

```go
client := stock.New(
    stock.WithProxyPool([]string{
        "http://user:pass@1.2.3.4:8080",
        "http://user:pass@5.6.7.8:8080",
    }),
)
```

`ProxyPoolOptions` 第一版只包含 `URLs`，保留 options 入口是为了后续兼容扩展随机轮换、失败冷却、provider 级策略等能力。

## 配置规则

- 空字符串会被忽略。
- URL 解析失败会被忽略。
- URL scheme 为空或 host 为空会被忽略。
- 有效 URL 按输入顺序保存。
- 所有 URL 都无效时等同于未配置代理池。
- 支持 Go 标准库 `http.Transport.Proxy` 支持的代理 URL 形式，包括用户名密码。

## 内部结构

在根 `Config` 增加：

```go
ProxyPool ProxyPoolOptions
```

在 `internal/core.Config` 增加：

```go
ProxyPool ProxyPoolConfig
```

在 `internal/core` 新增线程安全 round-robin 选择器：

```go
type ProxyPool struct {
    urls []*url.URL
    next atomic.Uint64
}

func NewProxyPool(rawURLs []string) *ProxyPool
func (p *ProxyPool) Proxy(req *http.Request) (*url.URL, error)
```

`Proxy(req)` 每次被调用时返回下一个代理 URL。由于 Go 的 `http.Transport.Proxy` 会在每次 request round trip 时调用代理函数，而当前 core 每次 retry 都会重新执行 `httpClient.Do(req)`，所以可以自然覆盖首次请求、重试和 fallback host 尝试。

## HTTP client 组装

`stock.New` 将 `Config.ProxyPool` 传给 `core.NewClient`。

`core.NewClient` 在创建 client 时处理代理池：

- 如果没有有效代理 URL，沿用传入的 `HTTPClient`。
- 如果传入 client 的 `Transport` 为 `nil`，克隆 `http.DefaultTransport.(*http.Transport)` 并设置 `Proxy`。
- 如果传入 client 的 `Transport` 是 `*http.Transport`，克隆该 transport 并设置 `Proxy`，保留用户已有超时、TLS、Dialer 等设置。
- 如果传入 client 的 `Transport` 不是 `*http.Transport`，不强行改写 transport，代理池不生效。这样避免破坏用户自定义 round tripper。
- 复制 `http.Client` 的 `Timeout`、`Jar`、`CheckRedirect` 等字段，只替换克隆后的 transport。

如果用户需要完全自定义代理行为，仍可继续使用 `WithHTTPClient`。

## 轮换语义

轮换单位为 HTTP attempt，不是业务方法调用。

示例：代理池 `[p1, p2]`，一次请求遇到 429 并重试一次：

1. 首次 attempt 使用 `p1`。
2. 重试 attempt 使用 `p2`。

示例：代理池 `[p1, p2, p3]`，原始 host 失败后走 fallback host：

1. 原始 host 首次 attempt 使用 `p1`。
2. 原始 host retry 使用 `p2`。
3. fallback host attempt 使用 `p3`。

并发请求共享同一个代理池计数器，整体按原子计数递增选择代理。

## 错误处理

- 配置阶段不会因为非法代理 URL panic。
- 无有效代理时行为与未配置代理一致。
- 代理连接失败会以现有 `NETWORK_ERROR` 包装，并继续受已有 retry 策略控制。
- 请求 hooks 不新增代理字段，避免把含账号密码的代理 URL 暴露给日志。后续如需可增加脱敏后的 proxy id。

## 测试计划

- `WithProxyPool` 会把有效代理 URL 写入根配置。
- `stock.New` 会把代理池传入 core，并且不会影响默认配置。
- core 层代理池 round-robin 顺序正确，且空值和非法 URL 被忽略。
- 配置代理池后，连续请求或 retry 会触发不同代理。
- 使用 `WithHTTPClient` 且 transport 为 `*http.Transport` 时，会保留 client 字段并替换为带代理函数的克隆 transport。
- 使用自定义非 `*http.Transport` round tripper 时，SDK 不改写它。
- 执行 `go test ./...`。

