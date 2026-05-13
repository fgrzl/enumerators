package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestDisposable_Interface(t *testing.T) {
	// This test verifies that all enumerators implement Disposable

	// Arrange - create various enumerators
	slice := enumerators.Slice([]int{1, 2, 3})
	empty := enumerators.Empty[int]()
	rangeEnum := enumerators.Range(0, 3, func(i int) int { return i })

	// Act & Assert - verify they all implement Disposable
	var disposables []enumerators.Disposable
	disposables = append(disposables, slice)
	disposables = append(disposables, empty)
	disposables = append(disposables, rangeEnum)

	// Should not panic when calling Dispose
	for _, d := range disposables {
		d.Dispose()
	}
}

func TestDisposable_ChainedDisposal(t *testing.T) {
	// Arrange
	disposed1 := false
	disposed2 := false

	base1 := enumerators.Cleanup(
		enumerators.Slice([]int{1, 2}),
		func() { disposed1 = true },
	)

	base2 := enumerators.Cleanup(
		enumerators.Slice([]int{3, 4}),
		func() { disposed2 = true },
	)

	chained := enumerators.Chain(base1, base2)

	// Act
	chained.Dispose()

	// Assert
	assert.True(t, disposed1, "First enumerator should be disposed")
	assert.True(t, disposed2, "Second enumerator should be disposed")
}

func TestDisposable_FilteredDisposal(t *testing.T) {
	// Arrange
	disposed := false
	base := enumerators.Cleanup(
		enumerators.Slice([]int{1, 2, 3, 4, 5}),
		func() { disposed = true },
	)

	filtered := enumerators.Filter(base, func(x int) bool { return x%2 == 0 })

	// Act
	filtered.Dispose()

	// Assert
	assert.True(t, disposed, "Base enumerator should be disposed through filter")
}

func TestDisposable_TransformationChainDisposal(t *testing.T) {
	// Arrange
	disposed := false
	base := enumerators.Cleanup(
		enumerators.Slice([]int{1, 2, 3, 4, 5}),
		func() { disposed = true },
	)

	// Create a chain of transformations
	filtered := enumerators.Filter(base, func(x int) bool { return x > 2 })
	taken := enumerators.Take(filtered, 2)
	skipped := enumerators.SkipIf(taken, func(x int) bool { return false })

	// Act
	skipped.Dispose()

	// Assert
	assert.True(t, disposed, "Base enumerator should be disposed through transformation chain")
}

func TestDisposable_ToSliceDisposal(t *testing.T) {
	// Arrange
	disposed := false
	base := enumerators.Cleanup(
		enumerators.Slice([]int{1, 2, 3}),
		func() { disposed = true },
	)

	// Act
	result, err := enumerators.ToSlice(base)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
	assert.True(t, disposed, "ToSlice should dispose the enumerator")
}

func TestDisposable_ForEachDisposal(t *testing.T) {
	// Arrange
	disposed := false
	base := enumerators.Cleanup(
		enumerators.Slice([]int{1, 2, 3}),
		func() { disposed = true },
	)

	// Act
	err := enumerators.ForEach(base, func(x int) error { return nil })

	// Assert
	assert.NoError(t, err)
	assert.True(t, disposed, "ForEach should dispose the enumerator")
}

func TestDisposable_ConsumeDisposal(t *testing.T) {
	// Arrange
	disposed := false
	base := enumerators.Cleanup(
		enumerators.Slice([]int{1, 2, 3}),
		func() { disposed = true },
	)

	// Act
	err := enumerators.Consume(base)

	// Assert
	assert.NoError(t, err)
	assert.True(t, disposed, "Consume should dispose the enumerator")
}
