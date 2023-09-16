package utils

import "fmt"

func GetAtIndex[T any](slice []T, index int) (T, error) {
  var res T
  for i, v := range slice {
    if i == index {
      return v, nil
    }
  }

  return res, fmt.Errorf("Index out of range")
}
