package generator

import (
	"strings"
	"testing"
)

func TestMakePageIDStable(t *testing.T) {
	options := PageOptions{
		Width:      600,
		Height:     800,
		Margin:     40,
		FontPath:   `C:\fonts\regular.ttf`,
		FontSize:   24,
		LineHeight: 32,
		Seed:       42,
	}

	first := MakePageID("text-001", "Привет", options)
	second := MakePageID("text-001", "Привет", options)

	if first != second {
		t.Fatalf("expected stable ID, got %q and %q", first, second)
	}
}

func TestMakePageIDDifferentSeed(t *testing.T) {
	firstOptions := PageOptions{
		Width:      600,
		Height:     800,
		Margin:     40,
		FontPath:   "regular.ttf",
		FontSize:   24,
		LineHeight: 32,
		Seed:       1,
	}

	secondOptions := firstOptions
	secondOptions.Seed = 2

	first := MakePageID("text-001", "Привет", firstOptions)
	second := MakePageID("text-001", "Привет", secondOptions)

	if first == second {
		t.Fatalf("expected different IDs for different seeds, got %q", first)
	}
}

func TestMakePageIDIgnoresFontDirectory(t *testing.T) {
	windowsOptions := PageOptions{
		Width:      600,
		Height:     800,
		Margin:     40,
		FontPath:   `C:\fonts\regular.ttf`,
		FontSize:   24,
		LineHeight: 32,
		Seed:       42,
	}

	unixOptions := windowsOptions
	unixOptions.FontPath = "/usr/share/fonts/regular.ttf"

	windowsID := MakePageID("text-001", "Привет", windowsOptions)
	unixID := MakePageID("text-001", "Привет", unixOptions)

	if windowsID != unixID {
		t.Fatalf("expected IDs to match, got %q and %q", windowsID, unixID)
	}
}

func TestMakePageIDFormat(t *testing.T) {
	options := PageOptions{
		Width:      600,
		Height:     800,
		Margin:     40,
		FontPath:   "regular.ttf",
		FontSize:   24,
		LineHeight: 32,
		Seed:       42,
	}

	id := MakePageID("text-001", "Привет", options)

	if !strings.HasPrefix(id, "synthetic-text-001-") {
		t.Fatalf("unexpected ID format: %q", id)
	}
}
