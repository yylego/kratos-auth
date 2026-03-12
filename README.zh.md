[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yylego/kratos-auth/release.yml?branch=main&label=BUILD)](https://github.com/yylego/kratos-auth/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yylego/kratos-auth)](https://pkg.go.dev/github.com/yylego/kratos-auth)
[![Coverage Status](https://img.shields.io/coveralls/github/yylego/kratos-auth/main.svg)](https://coveralls.io/github/yylego/kratos-auth?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25+-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/yylego/kratos-auth.svg)](https://github.com/yylego/kratos-auth/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yylego/kratos-auth)](https://goreportcard.com/report/github.com/yylego/kratos-auth)

# kratos-auth

Kratos 认证中间件的基座包，提供各中间件共用的类型和工具函数。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->

## 英文文档

[ENGLISH README](README.md)

<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

## 核心特性

🎯 **路由范围匹配**: 基于 INCLUDE/EXCLUDE 模式的路由匹配，提供 `RouteScope`、`SelectSide`、`Operation` 类型
🔌 **Span Hook 接口**: 可插拔的追踪集成，通过 `SpanHook` / `NewSpanHookFunc` / `RunSpanHooks` 实现
🛠️ **共享工具函数**: 如 `BooleanToNum` 等下游中间件包共用的公共函数

## 安装

```bash
go get github.com/yylego/kratos-auth/authkratos
```

## 使用方式

### 路由范围匹配

```go
import "github.com/yylego/kratos-auth/authkratos"

// INCLUDE 模式：仅匹配指定的操作
scope := authkratos.NewInclude(
    "/api.Service/CreateUser",
    "/api.Service/UpdateUser",
)

// EXCLUDE 模式：匹配除指定操作外的所有操作
scope := authkratos.NewExclude(
    "/api.Service/HealthCheck",
    "/api.Service/Metrics",
)

// 检查操作是否匹配
if scope.Match("/api.Service/CreateUser") {
    // 匹配成功
}

// 反转匹配模式
opposite := scope.Opposite()
```

### Span Hook — 可插拔追踪设计

下游中间件包需要 APM / 追踪支持，但基座包不应依赖具体的 APM 实现（如 Elastic APM、Jaeger、Datadog）。`SpanHook` 接口解决了这个问题：

- 每个中间件通过 `WithNewSpanHook(fn)` 接受一个 `NewSpanHookFunc` **回调列表**
- 处理请求时，中间件调用 `RunSpanHooks(ctx, hooks, spanName)` 启动追踪 span
- 返回的清理函数通过 `defer` 调用来关闭 span

这种设计使基座包与 APM 解耦。使用方通过实现 `SpanHook` 注入追踪：

```go
import "github.com/yylego/kratos-auth/authkratos"

// 实现 SpanHook 接口以集成具体的 APM / 追踪后端
type mySpanHook struct {
    span trace.Span
}

func (h *mySpanHook) Start(ctx context.Context, spanName string) {
    _, h.span = tracer.Start(ctx, spanName)
}

func (h *mySpanHook) Close() {
    h.span.End()
}

// 每次调用必须返回新实例（span 状态保存在内部）
newHook := func() authkratos.SpanHook {
    return &mySpanHook{}
}

// 在中间件代码中 — 启动钩子并在退出时清理
cleanup := authkratos.RunSpanHooks(ctx, []authkratos.NewSpanHookFunc{newHook}, "my-span")
defer cleanup()
```

下游中间件包使用此模式：

```go
// 通过配置接受钩子
cfg.WithNewSpanHook(func() authkratos.SpanHook { return &mySpanHook{} })

// 中间件内部 — 自动管理 span 生命周期
defer authkratos.RunSpanHooks(ctx, cfg.spanHooks, "rate-kratos-limits")()
```

### BooleanToNum

```go
import "github.com/yylego/kratos-auth/authkratos"

// 布尔值转整数（true=1, false=0），常用于调试日志
authkratos.BooleanToNum(true)  // 1
authkratos.BooleanToNum(false) // 0
```

## API 一览

| 类型 / 函数 | 说明 |
| --- | --- |
| `RouteScope` | 路由匹配范围，支持 INCLUDE/EXCLUDE 模式 |
| `NewInclude(ops...)` | 创建仅匹配指定操作的范围 |
| `NewExclude(ops...)` | 创建排除指定操作后匹配其他操作的范围 |
| `RouteScope.Match(op)` | 检查操作是否在范围内 |
| `RouteScope.Opposite()` | 返回反转模式的范围 |
| `SelectSide` | 匹配模式常量（`INCLUDE` / `EXCLUDE`） |
| `Operation` | 路由操作路径（string 别名） |
| `SpanHook` | 追踪 span 接口（`Start` + `Close`） |
| `NewSpanHookFunc` | 创建新 `SpanHook` 实例的函数类型 |
| `RunSpanHooks(ctx, hooks, name)` | 启动钩子并返回清理函数 |
| `BooleanToNum(b)` | 布尔值转整数（1 / 0） |

## 测试

```bash
go test -v ./...
```

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 许可证类型

MIT 许可证 - 详见 [LICENSE](LICENSE)。

---

## 💬 联系与反馈

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **问题报告？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **新颖思路？** 创建 issue 讨论
- 📖 **文档疑惑？** 报告问题，帮助我们完善文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，协助解决性能问题
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/kratos-auth.git`）
3. **导航**：进入克隆的项目（`cd kratos-auth`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：面向用户的更改需要更新文档
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来贡献此项目。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![标星点赞](https://starchart.cc/yylego/kratos-auth.svg?variant=adaptive)](https://starchart.cc/yylego/kratos-auth)
