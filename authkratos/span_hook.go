package authkratos

import "context"

// SpanHook provides tracing span lifecycle management
// Implement Start to create and store span, implement Close to finish it
//
// SpanHook 提供追踪 span 的生命周期管理
// 实现 Start 创建并保存 span，实现 Close 结束 span
type SpanHook interface {
	Start(ctx context.Context, spanName string)
	Close()
}

// NewSpanHookFunc creates a fresh SpanHook instance for each span
// Each call must return a new instance since span state is stored internally
//
// NewSpanHookFunc 为每个 span 创建新的 SpanHook 实例
// 每次调用必须返回新实例，因为 span 状态保存在内部
type NewSpanHookFunc func() SpanHook

// RunSpanHooks starts all span hooks and returns a cleanup function
// Caller should defer the returned function to close all hooks
//
// RunSpanHooks 启动所有 span hooks 并返回清理函数
// 调用方应 defer 返回的函数来关闭所有 hooks
func RunSpanHooks(ctx context.Context, spanHooks []NewSpanHookFunc, spanName string) func() {
	hooks := make([]SpanHook, 0, len(spanHooks))
	for _, newHook := range spanHooks {
		hook := newHook()
		hook.Start(ctx, spanName)
		hooks = append(hooks, hook)
	}
	return func() {
		for _, hook := range hooks {
			hook.Close()
		}
	}
}
