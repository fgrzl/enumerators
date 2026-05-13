package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertIntEnumeratorYields(t *testing.T, e enumerators.Enumerator[int], want []int) {
	t.Helper()
	for _, w := range want {
		require.True(t, e.MoveNext())
		got, err := e.Current()
		require.NoError(t, err)
		assert.Equal(t, w, got)
	}
	assert.False(t, e.MoveNext())
	assert.NoError(t, e.Err())
}
