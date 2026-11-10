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

	intervalMinutes := appConfig.IntervalMinutes
	if intervalMinutes <= 0 {
		intervalMinutes = 1
	}
	log.Printf("Check interval set to %d minute(s)", intervalMinutes)

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}

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

		// Calculate sleep duration until the next interval aligned to minute boundaries (e.g., 12:01:00, 12:02:00).
		current := time.Now().In(loc)
		next := nextScheduledTime(current, intervalMinutes, loc)
		sleepDuration := next.Sub(time.Now().In(loc))
		if sleepDuration < 0 {
			sleepDuration = 0
		}

		log.Printf("Lottery check finished. Sleeping until %s (in %v)...\n", next.Format(time.RFC1123), sleepDuration)
		time.Sleep(sleepDuration)
	}
}

// nextScheduledTime calculates the next scheduled check time based on interval in minutes.
// It aligns checks to wall-clock minute boundaries (e.g., 12:01:00, 12:02:00, etc.).
func nextScheduledTime(current time.Time, intervalMinutes int, loc *time.Location) time.Time {
	if intervalMinutes <= 0 {
		intervalMinutes = 1
	}
	if loc != nil {
		current = current.In(loc)
	} else {
		loc = current.Location()
	}

	totalMinutes := current.Hour()*60 + current.Minute()
	nextTotalMinutes := (totalMinutes/intervalMinutes + 1) * intervalMinutes

	nextHour := nextTotalMinutes / 60
	nextMinute := nextTotalMinutes % 60
	next := time.Date(current.Year(), current.Month(), current.Day(), nextHour, nextMinute, 0, 0, loc)

	for !next.After(current) {
		next = next.Add(time.Duration(intervalMinutes) * time.Minute)
	}

	return next
}
