package enumerators

// SequenceEqual determines whether two enumerators are equal by comparing their elements.
// Both enumerators are consumed during the comparison.
// Returns true if both sequences contain the same elements in the same order.
func SequenceEqual[T comparable](first, second Enumerator[T]) (bool, error) {
	defer first.Dispose()
	defer second.Dispose()

	for {
		firstHasNext := first.MoveNext()
		secondHasNext := second.MoveNext()

		// Different lengths
		if firstHasNext != secondHasNext {
			return false, nil
		}

		// Both ended
		if !firstHasNext {
			break
		}

		firstVal, err := first.Current()
		if err != nil {
			return false, err
		}

		secondVal, err := second.Current()
		if err != nil {
			return false, err
		}

		if firstVal != secondVal {
			return false, nil
		}
	}

	// Check for errors
	if err := first.Err(); err != nil {
		return false, err
	}
	if err := second.Err(); err != nil {
		return false, err
	}

	return true, nil
}
