package protomask

import (
	"testing"

	"github.com/olexnzarov/protomask/internal/pbtest"
)

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func TestHelperAll(t *testing.T) {
	book := &pbtest.Book{
		Id: 1605,
		Price: &pbtest.Price{
			Cents: 1500,
		},
	}

	mask := All(book)
	if !mask.IsValid(book) {
		t.Fatal("invalid field mask")
		return
	}

	paths := mask.GetPaths()
	expectedPaths := []string{"id", "price"}
	if len(paths) != len(expectedPaths) {
		t.Log("field mask has unexpected length")
		t.Fail()
	}

	for _, p := range paths {
		if !contains(expectedPaths, p) {
			t.Logf("field mask contains an unexpected path: '%s'", p)
			t.Fail()
		}
	}
}

func TestFieldMaskIsValid(t *testing.T) {
	book := &pbtest.Book{}
	mask := &fieldMask{
		paths: []string{"Name"},
	}
	if mask.IsValid(book) {
		t.Fatal("IsValid failed to validate an invalid mask")
		return
	}
	if mask.IsValid(nil) {
		t.Fatal("IsValid failed to validate a mask against a nil message")
		return
	}
}

func TestMaskFieldNilReference(t *testing.T) {
	var mask *fieldMask
	if mask.IsValid(&pbtest.Book{}) {
		t.Fatal("IsValid return true on a nil reference")
		return
	}
	if mask.GetPaths() != nil {
		t.Fatal("GetPaths return an array on a nil reference")
		return
	}
}
