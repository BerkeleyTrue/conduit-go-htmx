package utils

import "unicode"

// IsAnyLowercase checks if any of the runes in the input slice are lowercase.
func IsAnyLowercase(input []rune) bool {
	return Some(unicode.IsLower, input)
}

// IsAnyUppercase checks if any of the runes in the input slice are uppercase.
func IsAnyUppercase(input []rune) bool {
  return Some(unicode.IsUpper, input)
}

// IsAnyDigit checks if any of the runes in the input slice are digits.
func IsAnyDigit(input []rune) bool {
	return Some(unicode.IsDigit, input)
}

// IsAnySymbol checks if any of the runes in the input slice are symbols.
func IsAnySymbol(input []rune) bool {
  return !All(unicode.IsLetter, input)
}
