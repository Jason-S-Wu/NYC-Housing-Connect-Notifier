package models

import (
	"reflect"
	"strconv"
	"strings"
)

type SearchLotteriesResponse struct {
	Sales []struct {
		LotteryID          int    `json:"lotteryId,omitempty"`
		LotteryName        string `json:"lotteryName,omitempty"`
		LotteryDescription string `json:"lotteryDescription,omitempty"`
		EndIn              int    `json:"endIn,omitempty"`
		Prices             any    `json:"prices,omitempty"`
		Rents              any    `json:"rents,omitempty"`
		MaxIncome          int    `json:"maxIncome,omitempty"`
		MinIncome          any    `json:"minIncome,omitempty"`
		MinHouseholdSize   int    `json:"minHouseholdSize,omitempty"`
		MaxHouseholdSize   int    `json:"maxHouseholdSize,omitempty"`
		NeighborhoodIDs    string `json:"neighborhoodIDs,omitempty"`
		AmenityIDs         any    `json:"amenityIDs,omitempty"`
		UnitTypeIds        string `json:"unitTypeIds,omitempty"`
		LotteryStartDate   any    `json:"lotteryStartDate,omitempty"`
		LotteryEndDate     string `json:"lotteryEndDate,omitempty"`
		DefaultPhoto       string `json:"defaultPhoto,omitempty"`
		PropertyTypeID     int    `json:"propertyTypeId,omitempty"`
		DefaultPhotoStream string `json:"defaultPhotoStream,omitempty"`
		Trains             string `json:"trains,omitempty"`
		Markers            []struct {
			Name    string `json:"name,omitempty"`
			Lat     string `json:"lat,omitempty"`
			Lng     string `json:"lng,omitempty"`
			Address string `json:"address,omitempty"`
			City    string `json:"city,omitempty"`
			State   string `json:"state,omitempty"`
			Zip     string `json:"zip,omitempty"`
		} `json:"markers,omitempty"`
		Borough              string `json:"borough,omitempty"`
		Neighborhood         string `json:"neighborhood,omitempty"`
		Amenities            any    `json:"amenities,omitempty"`
		Studios              int    `json:"studios,omitempty"`
		OneBR                int    `json:"oneBR,omitempty"`
		TwoBR                int    `json:"twoBR,omitempty"`
		ThreeBR              int    `json:"threeBR,omitempty"`
		FourBR               int    `json:"fourBR,omitempty"`
		FiveBR               int    `json:"fiveBR,omitempty"`
		SixBR                int    `json:"sixBR,omitempty"`
		Units                int    `json:"units,omitempty"`
		IsApplied            bool   `json:"isApplied,omitempty"`
		LotteryPreferenceID  any    `json:"lotteryPreferenceId,omitempty"`
		IsMitchelLamaLottery bool   `json:"isMitchelLamaLottery,omitempty"`
	} `json:"sales,omitempty"`
	Rentals []struct {
		LotteryID          int    `json:"lotteryId,omitempty"`
		LotteryName        string `json:"lotteryName,omitempty"`
		LotteryDescription string `json:"lotteryDescription,omitempty"`
		EndIn              int    `json:"endIn,omitempty"`
		Prices             any    `json:"prices,omitempty"`
		Rents              string `json:"rents,omitempty"`
		MaxIncome          int    `json:"maxIncome,omitempty"`
		MinIncome          int    `json:"minIncome,omitempty"`
		MinHouseholdSize   int    `json:"minHouseholdSize,omitempty"`
		MaxHouseholdSize   int    `json:"maxHouseholdSize,omitempty"`
		NeighborhoodIDs    string `json:"neighborhoodIDs,omitempty"`
		AmenityIDs         string `json:"amenityIDs,omitempty"`
		UnitTypeIds        string `json:"unitTypeIds,omitempty"`
		LotteryStartDate   any    `json:"lotteryStartDate,omitempty"`
		LotteryEndDate     string `json:"lotteryEndDate,omitempty"`
		DefaultPhoto       string `json:"defaultPhoto,omitempty"`
		PropertyTypeID     int    `json:"propertyTypeId,omitempty"`
		DefaultPhotoStream string `json:"defaultPhotoStream,omitempty"`
		Trains             string `json:"trains,omitempty"`
		Markers            []struct {
			Name    string `json:"name,omitempty"`
			Lat     string `json:"lat,omitempty"`
			Lng     string `json:"lng,omitempty"`
			Address string `json:"address,omitempty"`
			City    string `json:"city,omitempty"`
			State   string `json:"state,omitempty"`
			Zip     string `json:"zip,omitempty"`
		} `json:"markers,omitempty"`
		Borough              string `json:"borough,omitempty"`
		Neighborhood         string `json:"neighborhood,omitempty"`
		Amenities            string `json:"amenities,omitempty"`
		Studios              int    `json:"studios,omitempty"`
		OneBR                int    `json:"oneBR,omitempty"`
		TwoBR                int    `json:"twoBR,omitempty"`
		ThreeBR              int    `json:"threeBR,omitempty"`
		FourBR               int    `json:"fourBR,omitempty"`
		FiveBR               int    `json:"fiveBR,omitempty"`
		SixBR                int    `json:"sixBR,omitempty"`
		Units                int    `json:"units,omitempty"`
		IsApplied            bool   `json:"isApplied,omitempty"`
		LotteryPreferenceID  any    `json:"lotteryPreferenceId,omitempty"`
		IsMitchelLamaLottery bool   `json:"isMitchelLamaLottery,omitempty"`
	} `json:"rentals,omitempty"`
}

type GetAdvertisement4RentResponse []struct {
	Ami       string `json:"ami"`
	UnitTypes []struct {
		UnitSize              string `json:"unitSize"`
		UnitLayoutTypeID      int    `json:"unitLayoutTypeId"`
		MonthlyRent           string `json:"monthlyRent"`
		UnitsAvailable        string `json:"unitsAvailable"`
		Arrow                 string `json:"arrow"`
		HouseholdSize         string `json:"householdSize"`
		AnnualHouseholdIncome string `json:"annualHouseholdIncome"`
		AssetLimit            string `json:"assetLimit"`
		IsPBV                 bool   `json:"isPBV"`
	} `json:"unitTypes"`
}

type DiscordWebhookPayload struct {
	Content any `json:"content"`
	Embeds  []DiscordEmbed `json:"embeds"`
	Username    string `json:"username"`
	AvatarURL   string `json:"avatar_url"`
	Attachments []any  `json:"attachments"`
}

type DiscordEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Color       any    `json:"color"`
	Image       struct {
		URL string `json:"url"`
	} `json:"image"`
}

type Rental struct {
	LotteryID    int    
	LotteryName  string
	Neighborhood string
	Borough      string
	LotteryURL   string
	QualifiedUnitTypes []Unit
	PhotoURL   string
}

type Unit struct {
	Ami string
	UnitSize string
	AssetLimit int
	IncomeRange string
	MonthlyRent int
	UnitsAvailable int
}


func (r Rental) String() string {
	s := "Lottery ID: \"" + strconv.Itoa(r.LotteryID) + "\", Name: \"" + r.LotteryName + "\", Neighborhood: \"" + r.Neighborhood + "\", Borough: \"" + r.Borough + "\", URL: \"" + r.LotteryURL + "\"\n"
	for _, unit := range r.QualifiedUnitTypes {
		s += "  " + unit.String() + "\n"
	}
	return s
}

func (u Unit) String() string {
	return "UnitSize: " + u.UnitSize + ", AMI: " + u.Ami + ", MonthlyRent: $" + strconv.Itoa(u.MonthlyRent) + ", UnitsAvailable: " + strconv.Itoa(u.UnitsAvailable) + ", AssetLimit: $" + strconv.Itoa(u.AssetLimit) + ", IncomeRange: " + u.IncomeRange
}

// StripWhitespace trims whitespace from all string fields in any struct (shallow, not recursive)
func StripWhitespace(s any) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String && field.CanSet() {
			field.SetString(strings.TrimSpace(field.String()))
		}
	}
}
