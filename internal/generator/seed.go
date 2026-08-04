package generator

import "math/rand"

// NewRandomSource создает детерминированный генератор случайных чисел.
func NewRandomSource(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}
