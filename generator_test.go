package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldGenerateKeyValuePairsWhenCreatingFromMap(t *testing.T) {
	// Arrange
	m := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	// Act
	enumerator := enumerators.GenerateFromMap(m)
	var results []*enumerators.KeyValuePair[string, int]
	for enumerator.MoveNext() {
		current, err := enumerator.Current()
		require.NoError(t, err)
		results = append(results, current)
	}

	// Assert
	expectedResults := []*enumerators.KeyValuePair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
		{Key: "three", Value: 3},
	}
	assert.ElementsMatch(t, expectedResults, results)
	assert.NoError(t, enumerator.Err())
}

func TestShouldPanicWhenGenerateCalledWithNilNext(t *testing.T) {
	// Arrange & Act & Assert
	assert.Panics(t, func() {
		enumerators.Generate[int](nil)
	})
}

func FuzzGenerate(f *testing.F) {
	f.Add(1, true, false)
	f.Add(0, false, true)
	f.Fuzz(func(t *testing.T, value int, hasNext bool, hasError bool) {
		// If hasError, hasNext must be false per contract
		if hasError {
			hasNext = false
		}
		var err error
		if hasError {
			err = assert.AnError
		}
		index := 0
		gen := enumerators.Generate(func() (int, bool, error) {
			if index == 0 {
				index++
				return value, hasNext, err
			}
			return 0, false, nil
		})
		defer gen.Dispose()

		if hasNext {
			assert.True(t, gen.MoveNext())
			curr, currErr := gen.Current()
			assert.NoError(t, currErr)
			assert.Equal(t, value, curr)
		} else {
			assert.False(t, gen.MoveNext())
		}
	})
}
