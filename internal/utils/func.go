package utils

func Some[T comparable](predicate func(T) bool, input []T) bool {

	for _, elem := range input {

		if predicate(elem) {
			return true
		}

	}

	return false
}

func All[T comparable](predicate func(T) bool, input []T) bool {

	for _, elem := range input {

		if !predicate(elem) {
			return false
		}

	}

	return true
}

func Iterate(count int) []int {
	cnts := make([]int, count)
	for i := 0; i < count; i++ {
		cnts[i] = i
	}
	return cnts
}

func Map[T, U any](mapper func(T) U, input []T) []U {
	output := make([]U, len(input))
	for idx, elem := range input {
		output[idx] = mapper(elem)
	}
	return output
}

func Difference[T comparable](input1 []T, input2 []T) []T {
  output := make([]T, 0)

  for _, elem := range input1 {
    if !Some(func(e T) bool { return e == elem }, input2) {
      output = append(output, elem)
    }
  }

  return output
}
