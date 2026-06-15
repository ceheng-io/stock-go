# apps/server

策衡 Web 后端 API，使用 Go 标准库 `net/http` 暴露 HTTP 接口，并调用当前仓库的 Go SDK。

## 开发

```bash
go run ./apps/server
```

默认监听 `127.0.0.1:8080`。可用环境变量覆盖：

```bash
CEHENG_SERVER_ADDR=127.0.0.1:8081 go run ./apps/server
```

## 测试

```bash
go test ./apps/server
```

## API 说明

当前 API 前缀为 `/api`，包括健康检查、搜索、行情、板块、K 线、资金流、北向资金、市场事件、分红和交易日历等接口。后端保持薄封装：解析 HTTP 参数、调用 Go SDK、返回 JSON。
