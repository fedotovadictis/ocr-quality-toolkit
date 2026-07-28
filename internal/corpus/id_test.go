package corpus

import (
	"strings"
	"testing"
)

func TestMakeMWSIDNormalizesSeparators(t *testing.T) {
	windowsPath := makeMWSID("17", `images\document.jpg`)
	unixPath := makeMWSID("17", "images/document.jpg")

	if windowsPath != unixPath {
		t.Fatalf("ids differ: %q and %q", windowsPath, unixPath)
	}
}
func TestMakeMWSIDFormat(t *testing.T) {
	id := makeMWSID("17", "images/document.jpg")

	if !strings.HasPrefix(id, "mws-17-") {
		t.Fatalf("unexpected id: %q", id)
	}

	if len(id) != len("mws-17-")+12 {
		t.Fatalf("unexpected id length: %d", len(id))
	}
}
