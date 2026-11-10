package authkratos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockSpanHook struct {
	name     string
	spanName string
	events   *[]string // Shared start and close events, kept as each one happens // 所有 hook 共用的启动与关闭事件，按发生顺序记
}

func (h *mockSpanHook) Start(ctx context.Context, spanName string) {
	h.spanName = spanName
	*h.events = append(*h.events, "start "+h.name)
}

func (h *mockSpanHook) Close() {
	*h.events = append(*h.events, "close "+h.name)
}

func TestRunSpanHooks(t *testing.T) {
	var events []string
	var hooks []*mockSpanHook
	newHookFunc := func(name string) NewSpanHookFunc {
		return func() SpanHook {
			hook := &mockSpanHook{name: name, events: &events}
			hooks = append(hooks, hook)
			return hook
		}
	}

	cleanup := RunSpanHooks(context.Background(), []NewSpanHookFunc{newHookFunc("parent"), newHookFunc("child")}, "test-span")
	require.Equal(t, []string{"start parent", "start child"}, events)
	for _, hook := range hooks {
		require.Equal(t, "test-span", hook.spanName)
	}

	// The span started last ends first, the same as defer
	// 后启动的 span 先结束，与 defer 一致
	cleanup()
	require.Equal(t, []string{"start parent", "start child", "close child", "close parent"}, events)
}

func TestRunSpanHooks_Empty(t *testing.T) {
	cleanup := RunSpanHooks(context.Background(), nil, "empty-span")
	require.NotPanics(t, cleanup)
}
