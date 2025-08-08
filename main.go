package main

import (
	"log"
	"os"
	"strconv"

	"main/models"
	"main/utils/discord"
	"main/utils/fetch_housing_connect"
	"main/utils/local"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	webhookURL := os.Getenv("DISCORD_WEBHOOK")
	householdIncomeStr := os.Getenv("HOUSEHOLD_INCOME")
	householdSizeStr := os.Getenv("HOUSEHOLD_SIZE")
	saveFileName := os.Getenv("SAVE_FILE_NAME")

	householdIncome, err := strconv.Atoi(householdIncomeStr)
	if err != nil {
		log.Fatalf("Invalid HOUSEHOLD_INCOME: %v", err)
	}
	householdSize, err := strconv.Atoi(householdSizeStr)
	if err != nil {
		log.Fatalf("Invalid HOUSEHOLD_SIZE: %v", err)
	}

	rentalData, err := fetch_housing_connect.FetchAllLotteries(householdIncome, householdSize)
	if err != nil {
		log.Fatalf("Error fetching lotteries: %v", err)
	}

	oldRentals, err := local.ReadRentalsFromFile(saveFileName)
	if err != nil {
		log.Printf("Error reading old rentals from file: %v", err)
		oldRentals = []models.Rental{} // If file read fails, start with an empty slice
	}

	newRentals := local.NewRentals(oldRentals, rentalData)

	println("New Rentals:")
	for _, rental := range newRentals {
		println(rental.String())
		discord.SendRentalNotification(webhookURL, rental)
	}

	local.WriteRentalsToFile(rentalData, saveFileName)
}
