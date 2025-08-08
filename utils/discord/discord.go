package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/models"
	"net/http"
	"strings"
)

func SendRentalNotification(webhookURL string, rental models.Rental) error {
	// Build a description string with all QualifiedUnitTypes
	var unitDetails []string
	for _, unit := range rental.QualifiedUnitTypes {
		unitDetails = append(unitDetails, fmt.Sprintf(
			"```Unit Size: %s\nAMI: %s\nMonthly Rent: $%d\nUnits Available: %d\nAsset Limit: $%d\nIncome Range: %s```",
			unit.UnitSize,
			unit.Ami,
			unit.MonthlyRent,
			unit.UnitsAvailable,
			unit.AssetLimit,
			unit.IncomeRange,
		))
	}
	description := fmt.Sprintf(
		"**Neighborhood:** %s\n**Borough**: %s\n**Units:**\n%s",
		rental.Neighborhood,
		rental.Borough,
		strings.Join(unitDetails, "\n\n"),
	)

	message := models.DiscordWebhookPayload{
		Content: nil,
		Embeds: []models.DiscordEmbed{
			{
				Title:       rental.LotteryName,
				Description: description,
				URL:         rental.LotteryURL,
				Color:       nil,
				Image: struct {
					URL string `json:"url"`
				}{
					URL: rental.PhotoURL,
				},
			},
		},
		Username:   "NYC Housing Connect",
		AvatarURL:  "https://www.premiumsvg.com/wimg_thumb/product-free-svg-house-clipart.webp",
		Attachments: []interface{}{},
	}

	// Send the message to the Discord webhook
	return sendMessage(webhookURL, message)
}

func sendMessage(webhookURL string, message models.DiscordWebhookPayload) error {
	// Create a new HTTP client
	client := &http.Client{}

	// Convert the message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	return nil
}