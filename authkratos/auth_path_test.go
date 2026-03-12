package authkratos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewInclude(t *testing.T) {
	scope := NewInclude("/api/v1/ping", "/api/v1/health")
	require.True(t, scope.Match("/api/v1/ping"))
	require.True(t, scope.Match("/api/v1/health"))
	require.False(t, scope.Match("/api/v1/other"))
}

func TestNewExclude(t *testing.T) {
	scope := NewExclude("/api/v1/ping")
	require.False(t, scope.Match("/api/v1/ping"))
	require.True(t, scope.Match("/api/v1/other"))
}

func TestRouteScope_Opposite(t *testing.T) {
	scope := NewInclude("/api/v1/ping")
	opposite := scope.Opposite()
	require.Equal(t, EXCLUDE, opposite.Side)
	require.False(t, opposite.Match("/api/v1/ping"))
	require.True(t, opposite.Match("/api/v1/other"))
}
