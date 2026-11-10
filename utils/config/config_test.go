package config

import (
	"encoding/json"
	"os"
	"testing"
)

func TestValidateConfigInterval(t *testing.T) {
	tests := []struct {
		name             string
		inputJSON        string
		expectedInterval int
		expectErr        bool
	}{
		{
			name: "default when interval is not set",
			inputJSON: `{
				"users": [
					{
						"name": "test",
						"discord_webhook": "https://discord.com/webhook",
						"household_income": 50000,
						"household_size": 1,
						"save_file_name": "test.json"
					}
				]
			}`,
			expectedInterval: 1,
			expectErr:        false,
		},
		{
			name: "interval specified as interval",
			inputJSON: `{
				"interval": 5,
				"users": [
					{
						"name": "test",
						"discord_webhook": "https://discord.com/webhook",
						"household_income": 50000,
						"household_size": 1,
						"save_file_name": "test.json"
					}
				]
			}`,
			expectedInterval: 5,
			expectErr:        false,
		},
		{
			name: "interval specified as interval_minutes",
			inputJSON: `{
				"interval_minutes": 10,
				"users": [
					{
						"name": "test",
						"discord_webhook": "https://discord.com/webhook",
						"household_income": 50000,
						"household_size": 1,
						"save_file_name": "test.json"
					}
				]
			}`,
			expectedInterval: 10,
			expectErr:        false,
		},
		{
			name: "negative interval returns error",
			inputJSON: `{
				"interval": -2,
				"users": [
					{
						"name": "test",
						"discord_webhook": "https://discord.com/webhook",
						"household_income": 50000,
						"household_size": 1,
						"save_file_name": "test.json"
					}
				]
			}`,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var cfg AppConfig
			err := json.Unmarshal([]byte(tc.inputJSON), &cfg)
			if err != nil {
				t.Fatalf("unexpected unmarshal error: %v", err)
			}

			err = validateConfig(&cfg)
			if (err != nil) != tc.expectErr {
				t.Fatalf("validateConfig() error = %v, expectErr %v", err, tc.expectErr)
			}
			if !tc.expectErr {
				if cfg.IntervalMinutes != tc.expectedInterval {
					t.Errorf("cfg.IntervalMinutes = %d; want %d", cfg.IntervalMinutes, tc.expectedInterval)
				}
				if cfg.Interval != tc.expectedInterval {
					t.Errorf("cfg.Interval = %d; want %d", cfg.Interval, tc.expectedInterval)
				}
			}
		})
	}
}

func TestLoadConfigFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := `{
		"interval": 3,
		"users": [
			{
				"name": "alice",
				"discord_webhook": "https://discord.com/webhook",
				"household_income": 80000,
				"household_size": 2,
				"save_file_name": "data/alice.json"
			}
		]
	}`
	if _, err := tempFile.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tempFile.Close()

	cfg, err := LoadConfig(tempFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.IntervalMinutes != 3 {
		t.Errorf("cfg.IntervalMinutes = %d; want 3", cfg.IntervalMinutes)
	}
}
