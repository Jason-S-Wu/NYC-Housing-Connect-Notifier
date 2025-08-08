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

	householdIncome, err := strconv.Atoi(householdIncomeStr)
	if err != nil {
		log.Fatalf("Invalid HOUSEHOLD_INCOME: %v", err)
	}
	householdSize, err := strconv.Atoi(householdSizeStr)
	if err != nil {
		log.Fatalf("Invalid HOUSEHOLD_SIZE: %v", err)
	}

	for {
		log.Println("Starting lottery check...")

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

			if len(newRentals) > 0 {
				log.Println("Found new rentals:")
				for _, rental := range newRentals {
					log.Println(rental.String())
					discord.SendRentalNotification(webhookURL, rental)
				}
			} else {
				log.Println("No new rentals found.")
			}

			local.WriteRentalsToFile(rentalData, saveFileName)
		}
		
		// Calculate sleep duration until the next 6-hour interval (e.g., 12 AM, 6 AM, 12 PM, 6 PM) in Eastern US time.
		loc, err := time.LoadLocation("America/New_York")
		if err != nil {
			log.Fatalf("Failed to load timezone: %v", err)
		}
		current := time.Now().In(loc)
		// Get the current hour of the day.
		hour := current.Hour()
		// Find the next 6-hour multiple.
		nextHour := (hour/6)*6 + 6
		// If the next hour is today, but has already passed, set it for the next day.
		if nextHour >= 24 {
			nextHour -= 24
			current = current.Add(24 * time.Hour)
		}
		
		// Create the next scheduled time.
		next := time.Date(current.Year(), current.Month(), current.Day(), nextHour, 0, 0, 0, loc)
		sleepDuration := next.Sub(time.Now().In(loc))

		log.Printf("Lottery check finished. Sleeping until %s (in %v)...\n", next.Format(time.RFC1123), sleepDuration)
		time.Sleep(sleepDuration)
	}
}