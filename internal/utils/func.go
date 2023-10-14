package utils

func Some[T comparable](predicate func(T) bool, input []T) bool {

	for _, elem := range input {

		if predicate(elem) {
			return true
		}

	}

	return false
}

func Iterate(count int) []int {
	cnts := make([]int, count)
	for i := 0; i < count; i++ {
		cnts[i] = i
	}
	return cnts
}
