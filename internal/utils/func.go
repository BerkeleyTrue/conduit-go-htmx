package utils

func some[T comparable](predicate func(T) bool, input []T) bool {

	for _, elem := range input {

		if predicate(elem) {
			return true
		}

	}

	return false
}
