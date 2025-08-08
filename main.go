package main

import (
	"log"
	"main/models"
	"main/utils/discord"
	"main/utils/fetch_housing_connect"
	"main/utils/local"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	webhookURL := os.Getenv("DISCORD_WEBHOOK")
	householdIncomeStr := os.Getenv("HOUSEHOLD_INCOME")
	householdSizeStr := os.Getenv("HOUSEHOLD_SIZE")
	saveFileName := os.Getenv("SAVE_FILE_NAME")
	sleepDurationHrsStr := os.Getenv("SLEEP_DURATION_HOURS")

	householdIncome, err := strconv.Atoi(householdIncomeStr)
	if err != nil {
		log.Fatalf("Invalid HOUSEHOLD_INCOME: %v", err)
	}
	householdSize, err := strconv.Atoi(householdSizeStr)
	if err != nil {
		log.Fatalf("Invalid HOUSEHOLD_SIZE: %v", err)
	}

	sleepDurationHrs, err := strconv.Atoi(sleepDurationHrsStr)
	if err != nil {
		sleepDurationHrs = 6 // Default to 6 hours if parsing fails
	}

	for {
		rentalData, err := fetch_housing_connect.FetchAllLotteries(householdIncome, householdSize)
		if err != nil {
			log.Printf("Error fetching lotteries: %v", err)
		} else {
			oldRentals, err := local.ReadRentalsFromFile(saveFileName)
			if err != nil {
				log.Printf("Error reading old rentals from file: %v", err)
				oldRentals = []models.Rental{}
			}

			newRentals := local.NewRentals(oldRentals, rentalData)

			println("New Rentals:")
			for _, rental := range newRentals {
				println(rental.String())
				discord.SendRentalNotification(webhookURL, rental)
			}

			local.WriteRentalsToFile(rentalData, saveFileName)
		}

		current := time.Now()
		now := time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location())
		next := now.Truncate(time.Hour).Add(time.Duration(sleepDurationHrs) * time.Hour)
		if !next.After(now) {
			next = next.Add(time.Duration(sleepDurationHrs) * time.Hour)
		}
		sleepDuration := next.Sub(current)
		log.Printf("Sleeping until next interval at %s (in %v)...\n", next.Format(time.RFC1123), sleepDuration)
		time.Sleep(sleepDuration)
	}
}
