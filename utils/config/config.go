package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type AppConfig struct {
	Users []UserConfig `json:"users"`
}

type UserConfig struct {
	Name            string `json:"name"`
	DiscordWebhook  string `json:"discord_webhook"`
	HouseholdIncome int    `json:"household_income"`
	HouseholdSize   int    `json:"household_size"`
	SaveFileName    string `json:"save_file_name"`
}

func LoadConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validateConfig(cfg *AppConfig) error {
	if len(cfg.Users) == 0 {
		return fmt.Errorf("config must include at least one user in users")
	}

	for i := range cfg.Users {
		user := &cfg.Users[i]

		if strings.TrimSpace(user.Name) == "" {
			user.Name = fmt.Sprintf("user-%d", i+1)
		}
		if strings.TrimSpace(user.DiscordWebhook) == "" {
			return fmt.Errorf("users[%d].discord_webhook is required", i)
		}
		if user.HouseholdIncome <= 0 {
			return fmt.Errorf("users[%d].household_income must be greater than 0", i)
		}
		if user.HouseholdSize <= 0 {
			return fmt.Errorf("users[%d].household_size must be greater than 0", i)
		}
		if strings.TrimSpace(user.SaveFileName) == "" {
			return fmt.Errorf("users[%d].save_file_name is required", i)
		}
	}

	return nil
}
