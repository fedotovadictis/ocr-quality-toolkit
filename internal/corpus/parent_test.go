package corpus

import "testing"

func TestValidateParentIDs(t *testing.T) {
	records := []Record{
		{
			ID:    "page-001",
			Image: "images/page-001.png",
		},
		{
			ID:       "page-001__noise-light",
			ParentID: "page-001",
			Image:    "images/page-001__noise-light.png",
		},
	}

	if err := ValidateParentIDs(records); err != nil {
		t.Fatalf("ValidateParentIDs returned error: %v", err)
	}
}

func TestValidateParentIDsMissingParent(t *testing.T) {
	records := []Record{
		{
			ID:       "page-001__noise-light",
			ParentID: "page-001",
			Image:    "images/page-001__noise-light.png",
		},
	}

	err := ValidateParentIDs(records)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestValidateParentIDsOriginalRecord(t *testing.T) {
	records := []Record{
		{
			ID:    "page-001",
			Image: "images/page-001.png",
		},
	}

	if err := ValidateParentIDs(records); err != nil {
		t.Fatalf("ValidateParentIDs returned error: %v", err)
	}
}
