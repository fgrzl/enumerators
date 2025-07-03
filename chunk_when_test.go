package enumerators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChunkWhen_SplitsOnEven(t *testing.T) {
	// Arrange
	input := Slice([]int{1, 3, 4, 5, 6, 7, 2, 9})

	chunks := ChunkWhen(input, false, func(_ bool, item int) (bool, bool, error) {
		// Start a new chunk when the number is even
		split := item%2 == 0
		return split, split, nil
	})

	var result [][]int

	// Act
	for chunks.MoveNext() {
		chunk, err := chunks.Current()
		require.NoError(t, err)

		var group []int
		for chunk.MoveNext() {
			item, err := chunk.Current()
			require.NoError(t, err)
			group = append(group, item.Item) // unwrap KeyedItem
		}
		result = append(result, group)
	}

	// Assert
	assert.Equal(t, [][]int{
		{1, 3},
		{4, 5},
		{6, 7},
		{2, 9},
	}, result)
}
