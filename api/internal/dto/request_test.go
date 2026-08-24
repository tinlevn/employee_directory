package dto

import "testing"

func TestPaginationQueryNormalize(t *testing.T) {
	query := PaginationQuery{Page: 0, PageSize: 500}
	query.Normalize()

	if query.Page != 1 {
		t.Fatalf("expected page 1, got %d", query.Page)
	}
	if query.PageSize != 100 {
		t.Fatalf("expected page size 100, got %d", query.PageSize)
	}
	if query.Offset() != 0 {
		t.Fatalf("expected offset 0, got %d", query.Offset())
	}
}

func TestPaginationQueryOffset(t *testing.T) {
	query := PaginationQuery{Page: 3, PageSize: 25}
	if query.Offset() != 50 {
		t.Fatalf("expected offset 50, got %d", query.Offset())
	}
}
