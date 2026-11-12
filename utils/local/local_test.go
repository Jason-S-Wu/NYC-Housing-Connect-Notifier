package local

import (
	"main/models"
	"os"
	"path/filepath"
	"testing"
)

func TestNewRentals(t *testing.T) {
	oldRentals := []models.Rental{
		{LotteryID: 1, LotteryName: "Apartment 1"},
		{LotteryID: 2, LotteryName: "Apartment 2"},
	}

	newRentals := []models.Rental{
		{LotteryID: 2, LotteryName: "Apartment 2 (updated)"},
		{LotteryID: 3, LotteryName: "Apartment 3"},
	}

	result := NewRentals(oldRentals, newRentals)
	if len(result) != 1 {
		t.Fatalf("expected 1 new rental, got %d", len(result))
	}
	if result[0].LotteryID != 3 {
		t.Errorf("expected rental ID 3, got %d", result[0].LotteryID)
	}
}

func TestMergeRentals(t *testing.T) {
	tests := []struct {
		name        string
		oldRentals  []models.Rental
		newRentals  []models.Rental
		expectedIDs []int
	}{
		{
			name:        "both empty",
			oldRentals:  nil,
			newRentals:  nil,
			expectedIDs: []int{},
		},
		{
			name:       "empty old rentals",
			oldRentals: nil,
			newRentals: []models.Rental{
				{LotteryID: 10, LotteryName: "A"},
				{LotteryID: 20, LotteryName: "B"},
			},
			expectedIDs: []int{10, 20},
		},
		{
			name: "empty new rentals retains old rentals",
			oldRentals: []models.Rental{
				{LotteryID: 10, LotteryName: "A"},
				{LotteryID: 20, LotteryName: "B"},
			},
			newRentals:  nil,
			expectedIDs: []int{10, 20},
		},
		{
			name: "preserves previously seen rentals missing from new response and adds new ones",
			oldRentals: []models.Rental{
				{LotteryID: 1, LotteryName: "A"},
				{LotteryID: 2, LotteryName: "B"},
			},
			newRentals: []models.Rental{
				{LotteryID: 2, LotteryName: "B updated"},
				{LotteryID: 3, LotteryName: "C"},
			},
			expectedIDs: []int{1, 2, 3},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			merged := MergeRentals(tc.oldRentals, tc.newRentals)
			if len(merged) != len(tc.expectedIDs) {
				t.Fatalf("expected %d rentals, got %d", len(tc.expectedIDs), len(merged))
			}
			for i, id := range tc.expectedIDs {
				if merged[i].LotteryID != id {
					t.Errorf("at index %d: expected LotteryID %d, got %d", i, id, merged[i].LotteryID)
				}
			}
		})
	}
}

func TestMergeRentalsUpdatesData(t *testing.T) {
	oldRentals := []models.Rental{
		{
			LotteryID:   1,
			LotteryName: "Old Name",
			QualifiedUnitTypes: []models.Unit{
				{UnitSize: "Studio", MonthlyRent: 1500},
			},
		},
	}

	newRentals := []models.Rental{
		{
			LotteryID:   1,
			LotteryName: "New Name",
			QualifiedUnitTypes: []models.Unit{
				{UnitSize: "Studio", MonthlyRent: 1600},
			},
		},
	}

	merged := MergeRentals(oldRentals, newRentals)
	if len(merged) != 1 {
		t.Fatalf("expected 1 rental, got %d", len(merged))
	}
	if merged[0].LotteryName != "New Name" {
		t.Errorf("expected updated LotteryName 'New Name', got %q", merged[0].LotteryName)
	}
	if merged[0].QualifiedUnitTypes[0].MonthlyRent != 1600 {
		t.Errorf("expected updated rent 1600, got %d", merged[0].QualifiedUnitTypes[0].MonthlyRent)
	}
}

func TestReadWriteRentalsFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "local-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "subdir", "rentals.json")

	// Read non-existent file should return empty slice and no error
	rentals, err := ReadRentalsFromFile(filePath)
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if len(rentals) != 0 {
		t.Fatalf("expected empty slice, got %d items", len(rentals))
	}

	// Write rentals
	toWrite := []models.Rental{
		{LotteryID: 100, LotteryName: "Test 100"},
		{LotteryID: 200, LotteryName: "Test 200"},
	}
	if err := WriteRentalsToFile(toWrite, filePath); err != nil {
		t.Fatalf("WriteRentalsToFile failed: %v", err)
	}

	// Read rentals back
	readBack, err := ReadRentalsFromFile(filePath)
	if err != nil {
		t.Fatalf("ReadRentalsFromFile failed: %v", err)
	}
	if len(readBack) != 2 {
		t.Fatalf("expected 2 rentals, got %d", len(readBack))
	}
	if readBack[0].LotteryID != 100 || readBack[1].LotteryID != 200 {
		t.Errorf("readBack mismatch: %+v", readBack)
	}
}
