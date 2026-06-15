# apps/web

策衡 Web 前端，使用 Vite + TypeScript + Vue 3.x + Ant Design Vue。

## 开发

```bash
npm install
npm run dev
```

默认通过 Vite proxy 访问 `apps/server` 的 `/api`。可用环境变量覆盖：

```bash
VITE_DEV_API_TARGET=http://127.0.0.1:8080 npm run dev
```

## 构建与测试

```bash
npm run test
npm run typecheck
npm run build
```

## 数据来源

前端不直连外部行情源；所有行情、板块、K 线、资金流和搜索请求都通过 `apps/server`，由后端调用当前 Go SDK。
