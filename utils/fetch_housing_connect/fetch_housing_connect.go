package fetch_housing_connect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/models"
	"net/http"
	"regexp"
	"strings"
)

const SEARCH_LOTTERY_URL = "https://a806-housingconnectapi.nyc.gov/HPDPublicAPI/api/Lottery/SearchLotteries"
const SEARCH_LOTTERY_PAYLOAD = `{
  "UnitTypes": [],
  "NearbyPlaces": [],
  "NearbySubways": [],
  "Amenities": [],
  "Applied": null,
  "HPDUserId": null,
  "Boroughs": [],
  "Neighborhoods": [],
  "HouseholdSize": null,
  "Income": "",
  "HouseholdType": 1,
  "OwnerTypes": [],
  "PreferanceTypes": [],
  "LotteryTypes": [],
  "Min": null,
  "Max": null
}`
const GET_LOTTERY_INFO_URL = "https://housingconnect.nyc.gov/PublicWeb/details/"
const GET_AMI_LOTTERY_INFO_URL = "https://a806-housingconnectapi.nyc.gov/HPDPublicAPI/api/LotteryConfig/GetAdvertisement4Rent?lotteryid="
const GET_PHOTO_URL = "https://a806-housingconnectapi.nyc.gov/MailTemplates/photos/"

func FetchAllLotteries(income int, householdSize int)  ([]models.Rental, error) {
	// use an HTTP client to send a POST request to the SEARCH_LOTTERY_URL with the SEARCH_LOTTERY_PAYLOAD
	client := &http.Client{}

	// Convert the payload string to a bytes.Buffer
	payloadBytes := []byte(SEARCH_LOTTERY_PAYLOAD)
	req, err := http.NewRequest("POST", SEARCH_LOTTERY_URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var searchResponse models.SearchLotteriesResponse
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	var rentals = searchResponse.Rentals

	var allRentals []models.Rental

	for _, rental := range rentals {
		var rentalInfo models.Rental

		rentalInfo.LotteryID = rental.LotteryID
		rentalInfo.Neighborhood = rental.Neighborhood
		rentalInfo.Borough = rental.Borough
		rentalInfo.LotteryURL = GET_LOTTERY_INFO_URL + fmt.Sprint(rental.LotteryID)
		rentalInfo.LotteryName = rental.LotteryName
		rentalInfo.PhotoURL = GET_PHOTO_URL + rental.DefaultPhoto + ".png"

		// Fetch qualified units for the rental
		qualifiedUnits, err := getQualifiedUnits(rental.LotteryID, householdSize, income)
		if err != nil {
			return nil, fmt.Errorf("failed to get qualified units: %w", err)
		}
		rentalInfo.QualifiedUnitTypes = qualifiedUnits

		models.StripWhitespace(&rentalInfo)

		if len(rentalInfo.QualifiedUnitTypes) != 0 {
			allRentals = append(allRentals, rentalInfo)
		}
	}
	
	return allRentals, nil
}

// getQualifiedUnits fetches and filters units by household size and income
func getQualifiedUnits(lotteryID int, householdSize int, income int) ([]models.Unit, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", GET_AMI_LOTTERY_INFO_URL+fmt.Sprint(lotteryID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response models.GetAdvertisement4RentResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	var units []models.Unit
	if len(response) == 0 {
		return units, nil
	}

	for _, resp := range response {
		for _, unitType := range resp.UnitTypes {
			allowedSizes := parseHouseholdSizes(unitType.HouseholdSize)
			incomeRanges := parseIncomeRanges(unitType.AnnualHouseholdIncome)
			// Find the index of the household size
			idx := -1
			for i, sz := range allowedSizes {
				if sz == householdSize {
					idx = i
					break
				}
			}
			if idx == -1 || idx >= len(incomeRanges) {
				continue // not eligible for this unit
			}
			minIncome, maxIncome := incomeRanges[idx][0], incomeRanges[idx][1]
			if income < minIncome || income > maxIncome {
				continue // not eligible for this unit
			}
			assetLimit := parseDollarString(unitType.AssetLimit)
			monthlyRent := parseDollarString(unitType.MonthlyRent)
			unitsAvailable := parseIntString(unitType.UnitsAvailable)
			unit := models.Unit{
				Ami:             resp.Ami,
				UnitSize:        unitType.UnitSize,
				AssetLimit:      assetLimit,
				MonthlyRent:     monthlyRent,
				UnitsAvailable:  unitsAvailable,
				IncomeRange:     fmt.Sprintf("$%d - $%d", minIncome, maxIncome),
			}
			units = append(units, unit)
		}
	}
	
	return units, nil
}

// parseHouseholdSizes extracts household sizes from the HTML string
func parseHouseholdSizes(html string) []int {
	var sizes []int
	// Use regex to extract numbers before "people"
	re := regexp.MustCompile(`(\d+)\s*people`)
	matches := re.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) > 1 {
			var n int
			fmt.Sscanf(match[1], "%d", &n)
			if n > 0 {
				sizes = append(sizes, n)
			}
		}
	}
	return sizes
}

// parseIncomeRanges extracts min/max income for each household size from the HTML string
func parseIncomeRanges(html string) [][2]int {
	var ranges [][2]int
	// Use regex to extract income ranges like $37,852.00-$45,360.00
	re := regexp.MustCompile(`\$([\d,\.]+)\s*-\s*\$([\d,\.]+)`)
	matches := re.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) == 3 {
			min := parseDollarString(match[1])
			max := parseDollarString(match[2])
			ranges = append(ranges, [2]int{min, max})
		}
	}
	return ranges
}

// parseDollarString parses a string like "$1,234.00" or "1,234.00" to int (removes $ and commas)
func parseDollarString(s string) int {
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	var val float64
	fmt.Sscanf(s, "%f", &val)
	return int(val)
}

// parseIntString parses a string like "6" to int
func parseIntString(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	var val int
	fmt.Sscanf(s, "%d", &val)
	return val
}