package utils

import "unicode"

func IsAnyLowercase(input []rune) bool {
	return some(unicode.IsLower, input)
}

func IsAnyUppercase(input []rune) bool {
  return some(unicode.IsUpper, input)
}

func IsAnyDigit(input []rune) bool {
	return some(unicode.IsDigit, input)
}

func IsAnySymbol(input []rune) bool {
  return some(unicode.IsSymbol, input)
}
