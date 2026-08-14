package report

import "testing"

func TestBuildAlignment(t *testing.T) {
	items := BuildAlignment(
		"Привет мир",
		"Привет мр",
	)

	if len(items) == 0 {
		t.Fatal("expected non-empty alignment")
	}

	foundDifference := false

	for _, item := range items {
		if item.Type != "equal" {
			foundDifference = true
			break
		}
	}

	if !foundDifference {
		t.Fatal("expected alignment to contain a difference")
	}
}
