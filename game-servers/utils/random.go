package utils

import "math/rand"

// generates random integer btw min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
