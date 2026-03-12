[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yylego/kratos-auth/release.yml?branch=main&label=BUILD)](https://github.com/yylego/kratos-auth/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yylego/kratos-auth)](https://pkg.go.dev/github.com/yylego/kratos-auth)
[![Coverage Status](https://img.shields.io/coveralls/github/yylego/kratos-auth/main.svg)](https://coveralls.io/github/yylego/kratos-auth?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25+-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/yylego/kratos-auth.svg)](https://github.com/yylego/kratos-auth/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yylego/kratos-auth)](https://goreportcard.com/report/github.com/yylego/kratos-auth)

# kratos-auth

Foundation package providing shared types and utilities consumed across Kratos authentication middleware packages.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->

## CHINESE README

[中文说明](README.zh.md)

<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Main Features

🎯 **Route Scope Matching**: INCLUDE/EXCLUDE mode route matching with `RouteScope`, `SelectSide`, and `Operation` types
🔌 **Span Hook Interface**: Pluggable tracing integration via `SpanHook` / `NewSpanHookFunc` / `RunSpanHooks`
🛠️ **Shared Utilities**: Common functions such as `BooleanToNum` consumed across downstream packages

## Installation

```bash
go get github.com/yylego/kratos-auth/authkratos
```

## Usage

### Route Scope Matching

```go
import "github.com/yylego/kratos-auth/authkratos"

// INCLUDE mode: match specified operations
scope := authkratos.NewInclude(
    "/api.Service/CreateUser",
    "/api.Service/UpdateUser",
)

// EXCLUDE mode: match everything except specified operations
scope := authkratos.NewExclude(
    "/api.Service/HealthCheck",
    "/api.Service/Metrics",
)

// Check if an operation matches
if scope.Match("/api.Service/CreateUser") {
    // matched
}

// Invert matching mode
opposite := scope.Opposite()
```

### Span Hook — Pluggable Tracing Design

Downstream middleware packages need APM / tracing support, but this base package must not depend on specific APM implementations (e.g., Elastic APM, Datadog, Zipkin). The `SpanHook` interface solves this:

- Each middleware accepts a **list** of `NewSpanHookFunc` callbacks via `WithNewSpanHook(fn)`
- When processing a request, the middleware invokes `RunSpanHooks(ctx, hooks, spanName)` to start tracing spans
- The returned cleanup function closes spans when done

This design keeps the base package APM-agnostic. Consumers inject tracing through implementing `SpanHook`:

```go
import "github.com/yylego/kratos-auth/authkratos"

// Implement SpanHook to integrate with specific APM / tracing backends
type mySpanHook struct {
    span trace.Span
}

func (h *mySpanHook) Start(ctx context.Context, spanName string) {
    _, h.span = tracer.Start(ctx, spanName)
}

func (h *mySpanHook) Close() {
    h.span.End()
}

// Each invocation must produce a fresh instance (span state is kept inside)
newHook := func() authkratos.SpanHook {
    return &mySpanHook{}
}

// In middleware code — start hooks and cleanup on exit
cleanup := authkratos.RunSpanHooks(ctx, []authkratos.NewSpanHookFunc{newHook}, "my-span")
defer cleanup()
```

Downstream middleware packages use this pattern:

```go
// Accept hooks via config
cfg.WithNewSpanHook(func() authkratos.SpanHook { return &mySpanHook{} })

// Inside middleware — automatic span management
defer authkratos.RunSpanHooks(ctx, cfg.spanHooks, "rate-kratos-limits")()
```

### BooleanToNum

```go
import "github.com/yylego/kratos-auth/authkratos"

// Converts boolean to int (true=1, false=0), common in debug logging
authkratos.BooleanToNum(true)  // 1
authkratos.BooleanToNum(false) // 0
```

## API

| Type / Function | Description |
| --- | --- |
| `RouteScope` | Route matching scope with INCLUDE/EXCLUDE mode |
| `NewInclude(ops...)` | Create scope that matches specified operations |
| `NewExclude(ops...)` | Create scope that matches everything except specified operations |
| `RouteScope.Match(op)` | Check if operation is within scope |
| `RouteScope.Opposite()` | Produce scope with inverted mode |
| `SelectSide` | Matching mode constant (`INCLUDE` / `EXCLUDE`) |
| `Operation` | Route operation path (string alias) |
| `SpanHook` | Tracing span interface (`Start` + `Close`) |
| `NewSpanHookFunc` | Function that creates a fresh `SpanHook` instance |
| `RunSpanHooks(ctx, hooks, name)` | Start hooks and produce cleanup function |
| `BooleanToNum(b)` | Converts boolean to int (1 / 0) |

## Testing

```bash
go test -v ./...
```

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 License

MIT License - see [LICENSE](LICENSE).

---

## 💬 Contact & Feedback

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Mistake reports?** Open an issue on GitHub with reproduction steps
- 💡 **Fresh ideas?** Create an issue to discuss
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package improved the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/kratos-auth.git`).
3. **Navigate**: Navigate to the cloned project (`cd kratos-auth`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/yylego/kratos-auth.svg?variant=adaptive)](https://starchart.cc/yylego/kratos-auth)
