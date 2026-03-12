package authkratos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockSpanHook struct {
	started  bool
	closed   bool
	spanName string
}

func (h *mockSpanHook) Start(ctx context.Context, spanName string) {
	h.started = true
	h.spanName = spanName
}

func (h *mockSpanHook) Close() {
	h.closed = true
}

func TestRunSpanHooks(t *testing.T) {
	var hooks []*mockSpanHook

	newHookFunc := func() SpanHook {
		hook := &mockSpanHook{}
		hooks = append(hooks, hook)
		return hook
	}

	cleanup := RunSpanHooks(context.Background(), []NewSpanHookFunc{newHookFunc, newHookFunc}, "test-span")

	require.Len(t, hooks, 2)
	for _, hook := range hooks {
		require.True(t, hook.started)
		require.Equal(t, "test-span", hook.spanName)
		require.False(t, hook.closed)
	}

	cleanup()

	for _, hook := range hooks {
		require.True(t, hook.closed)
	}
}

func TestRunSpanHooksEmpty(t *testing.T) {
	cleanup := RunSpanHooks(context.Background(), nil, "empty-span")
	cleanup() // should not panic
}
