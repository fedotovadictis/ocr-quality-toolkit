package generator

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
)

// NewRandomSource создаёт детерминированный генератор случайных чисел.
func NewRandomSource(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

// DeriveSeed создаёт стабильный seed для конкретной записи и преобразования.
func DeriveSeed(
	globalSeed string,
	recordID string,
	transformName string,
) int64 {
	input := globalSeed + ":" + recordID + ":" + transformName
	sum := sha256.Sum256([]byte(input))

	return int64(binary.BigEndian.Uint64(sum[:8]))
}
