package utils

import "unicode"

func IsAnyLowercase(input []rune) bool {
	return Some(unicode.IsLower, input)
}

func IsAnyUppercase(input []rune) bool {
  return Some(unicode.IsUpper, input)
}

func IsAnyDigit(input []rune) bool {
	return Some(unicode.IsDigit, input)
}

func IsAnySymbol(input []rune) bool {
  return Some(unicode.IsSymbol, input)
}
