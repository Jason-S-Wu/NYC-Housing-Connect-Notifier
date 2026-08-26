package main

import (
	"log"
	"main/utils/config"
	"main/utils/discord"
	"main/utils/fetch_housing_connect"
	"main/utils/local"
	"os"
	"time"
)

func main() {
	configPath := "config.json"
	if len(os.Args) > 1 && os.Args[1] != "" {
		configPath = os.Args[1]
	}

	appConfig, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config file %q: %v", configPath, err)
	}

	log.Printf("Loaded %d user configuration(s) from %s", len(appConfig.Users), configPath)

	for {
		log.Println("Starting lottery check...")

		for _, userConfig := range appConfig.Users {
			log.Printf("Checking lotteries for user %q...", userConfig.Name)

			rentalData, err := fetch_housing_connect.FetchAllLotteries(userConfig.HouseholdIncome, userConfig.HouseholdSize)
			if err != nil {
				log.Printf("Error fetching lotteries for user %q: %v", userConfig.Name, err)
				continue
			}

			oldRentals, err := local.ReadRentalsFromFile(userConfig.SaveFileName)
			if err != nil {
				log.Printf("Error reading old rentals for user %q from file %q: %v", userConfig.Name, userConfig.SaveFileName, err)
				continue
			}

			newRentals := local.NewRentals(oldRentals, rentalData)

			if len(newRentals) > 0 {
				log.Printf("Found %d new rental(s) for user %q.", len(newRentals), userConfig.Name)
				for _, rental := range newRentals {
					log.Println(rental.String())
					if err := discord.SendRentalNotification(userConfig.DiscordWebhook, rental); err != nil {
						log.Printf("Failed to send notification for user %q: %v", userConfig.Name, err)
					}
				}
			} else {
				log.Printf("No new rentals found for user %q.", userConfig.Name)
			}

			if err := local.WriteRentalsToFile(rentalData, userConfig.SaveFileName); err != nil {
				log.Printf("Failed to save rentals for user %q to file %q: %v", userConfig.Name, userConfig.SaveFileName, err)
			}
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
