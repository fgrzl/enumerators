package enumerators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChunkByKey_GroupsByPrefix(t *testing.T) {
	// Arrange
	input := Slice([]string{
		"apple", "apricot",
		"banana", "blueberry",
		"cherry",
		"date", "dragonfruit",
	})

	chunks := ChunkByKey(input, func(s string) string {
		return string(s[0]) // group by first letter
	})

	var result [][]string

	// Act
	for chunks.MoveNext() {
		chunk, err := chunks.Current()
		require.NoError(t, err)

		var group []string
		for chunk.Chunk.MoveNext() {
			item, err := chunk.Chunk.Current()
			require.NoError(t, err)
			group = append(group, item)
		}
		result = append(result, group)
	}

	// Assert
	assert.Equal(t, [][]string{
		{"apple", "apricot"},
		{"banana", "blueberry"},
		{"cherry"},
		{"date", "dragonfruit"},
	}, result)
}
