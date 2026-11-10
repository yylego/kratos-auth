package authkratos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewInclude(t *testing.T) {
	scope := NewInclude("/api/v1/ping", "/api/v1/health")
	require.True(t, scope.Match("/api/v1/ping"))
	require.True(t, scope.Match("/api/v1/health"))
	require.False(t, scope.Match("/api/v1/users"))
}

func TestNewExclude(t *testing.T) {
	scope := NewExclude("/api/v1/ping")
	require.False(t, scope.Match("/api/v1/ping"))
	require.True(t, scope.Match("/api/v1/users"))
}

func TestRouteScope_Match(t *testing.T) {
	require.True(t, NewInclude("/api/v1/ping").Match("/api/v1/ping"))
	require.False(t, NewInclude("/api/v1/ping").Match("/api/v1/users"))
	require.False(t, NewExclude("/api/v1/ping").Match("/api/v1/ping"))
	require.True(t, NewExclude("/api/v1/ping").Match("/api/v1/users"))

	scope := &RouteScope{Side: "OTHER", OperationSet: map[Operation]bool{}}
	require.Panics(t, func() { scope.Match("/api/v1/ping") })
}

func TestRouteScope_Opposite(t *testing.T) {
	scope := NewInclude("/api/v1/ping")
	opposite := scope.Opposite()
	require.Equal(t, EXCLUDE, opposite.Side)
	require.False(t, opposite.Match("/api/v1/ping"))
	require.True(t, opposite.Match("/api/v1/users"))

	// The two scopes hold separate sets, so changing one leaves the first one as is
	// 两个范围各有一份集合，改一个不影响另一个
	opposite.OperationSet["/api/v1/users"] = true
	require.False(t, scope.Match("/api/v1/users"))

	back := opposite.Opposite()
	require.Equal(t, INCLUDE, back.Side)
}

func TestSelectSide_Opposite(t *testing.T) {
	require.Equal(t, EXCLUDE, INCLUDE.Opposite())
	require.Equal(t, INCLUDE, EXCLUDE.Opposite())
	require.Panics(t, func() { SelectSide("OTHER").Opposite() })
}
