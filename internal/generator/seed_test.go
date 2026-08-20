package generator

import "testing"

func TestNewRandomSourceSameSeed(t *testing.T) {
	r1 := NewRandomSource(42)
	r2 := NewRandomSource(42)

	for i := 0; i < 100; i++ {
		if r1.Intn(1_000_000) != r2.Intn(1_000_000) {
			t.Fatal("same seed must produce same sequence")
		}
	}
}

func TestNewRandomSourceDifferentSeed(t *testing.T) {
	r1 := NewRandomSource(1)
	r2 := NewRandomSource(2)

	different := false

	for i := 0; i < 20; i++ {
		if r1.Intn(1_000_000) != r2.Intn(1_000_000) {
			different = true
			break
		}
	}

	if !different {
		t.Fatal("different seeds must produce different sequences")
	}
}
func TestDeriveSeedStable(t *testing.T) {
	first := DeriveSeed(
		"internship-2026",
		"page-001",
		"noise-light",
	)

	second := DeriveSeed(
		"internship-2026",
		"page-001",
		"noise-light",
	)

	if first != second {
		t.Fatalf(
			"expected stable seed, got %d and %d",
			first,
			second,
		)
	}
}

func TestDeriveSeedDiffersByRecord(t *testing.T) {
	first := DeriveSeed(
		"internship-2026",
		"page-001",
		"noise-light",
	)

	second := DeriveSeed(
		"internship-2026",
		"page-002",
		"noise-light",
	)

	if first == second {
		t.Fatal("expected different seeds for different records")
	}
}

func TestDeriveSeedDiffersByTransform(t *testing.T) {
	first := DeriveSeed(
		"internship-2026",
		"page-001",
		"noise-light",
	)

	second := DeriveSeed(
		"internship-2026",
		"page-001",
		"grayscale",
	)

	if first == second {
		t.Fatal("expected different seeds for different transforms")
	}
}
