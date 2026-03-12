package authkratos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBooleanToNum(t *testing.T) {
	require.Equal(t, 1, BooleanToNum(true))
	require.Equal(t, 0, BooleanToNum(false))
}
