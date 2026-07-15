package model

import "testing"

func TestPageQueryNormalizesUnsafeInput(t *testing.T) {
	query := PageQuery{Page: -1, PageSize: 1000}

	if got := query.GetPage(); got != 1 {
		t.Fatalf("GetPage() = %d, want 1", got)
	}
	if got := query.GetPageSize(); got != 100 {
		t.Fatalf("GetPageSize() = %d, want 100", got)
	}
	if got := query.GetOffset(); got != 0 {
		t.Fatalf("GetOffset() = %d, want 0", got)
	}
}
