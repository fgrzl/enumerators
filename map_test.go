package enumerators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToMap(t *testing.T) {
	// Arrange
	input := Slice([]string{"a", "bb", "ccc"})

	// Act
	result, err := ToMap(
		input,
		func(s string) int { return len(s) }, // key: length of string
		func(s string) string { return s + "-x" }, // value: original + "-x"
	)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "a-x", result[1])
	assert.Equal(t, "bb-x", result[2])
	assert.Equal(t, "ccc-x", result[3])
}
