package local

import (
	"encoding/json"
	"fmt"
	"main/models"
	"os"
)

// WriteRentalsToFile serializes rentals to JSON and writes to a file
func WriteRentalsToFile(rentals []models.Rental, filename string) error {
	data, err := json.MarshalIndent(rentals, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal rentals to JSON: %w", err)
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}
	return nil
}

// ReadRentalsFromFile reads rentals from a JSON file and returns them as a slice of models.Rental
func ReadRentalsFromFile(filename string) ([]models.Rental, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	var rentals []models.Rental
	err = json.Unmarshal(data, &rentals)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return rentals, nil
}

func NewRentals(oldRentals []models.Rental, newRentals []models.Rental) []models.Rental {
	// Create a map to track existing rental IDs
	existingIDs := make(map[int]bool)
	for _, rental := range oldRentals {
		existingIDs[rental.LotteryID] = true
	}

	// Filter out new rentals that already exist in old rentals
	var uniqueNewRentals []models.Rental
	for _, rental := range newRentals {
		if !existingIDs[rental.LotteryID] {
			uniqueNewRentals = append(uniqueNewRentals, rental)
		}
	}

	return uniqueNewRentals
}