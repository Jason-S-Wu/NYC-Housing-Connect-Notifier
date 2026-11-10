# NYC-Housing-Connect-Notifier

Automatically notifies you via Discord when new NYC Housing Connect rental lotteries become available that match your household size and income.

## Features

- Fetches current NYC Housing Connect rental lotteries
- Filters by household income and size
- Sends notifications to a Discord channel via webhook
- Remembers previously seen lotteries to avoid duplicate notifications

## How It Works

1. Fetches all available rental lotteries from the NYC Housing Connect API
2. Filters for units you qualify for based on your income and household size
3. Sends a detailed notification to your Discord channel for each new lottery
4. Saves the list of seen lotteries to avoid duplicate alerts

## Requirements

- Go 1.20+
- Docker (optional, for containerized deployment)

## Setup

### 1. Clone the repository

```sh
git clone https://github.com/Jason-S-Wu/NYC-Housing-Connect-Notifier.git
cd NYC-Housing-Connect-Notifier
```

### 2. Configure Users

Copy `config.sample.json` to `config.json` and fill in your values:

```json
{
  "interval": 1,
  "users": [
    {
      "name": "alice",
      "discord_webhook": "https://discord.com/api/webhooks/...",
      "household_income": 90000,
      "household_size": 2,
      "save_file_name": "data/alice-rentals.json"
    },
    {
      "name": "bob",
      "discord_webhook": "https://discord.com/api/webhooks/...",
      "household_income": 120000,
      "household_size": 3,
      "save_file_name": "data/bob-rentals.json"
    }
  ]
}
```

**Configuration options:**

- `interval`: Interval in minutes between lottery checks (optional integer, default: `1`). Aligns runs to clock minute boundaries (e.g., 12:01:00, 12:02:00). Can also be specified as `interval_minutes`.

**Fields per user:**

- `name`: A label shown in logs
- `discord_webhook`: Discord webhook URL
- `household_income`: Annual household income (number)
- `household_size`: Number of people in household (number)
- `save_file_name`: File path where seen lotteries are stored

### 3. Run Locally

```sh
go run main.go
```

Optional: pass a custom config file path:

```sh
go run main.go path/to/config.json
```

### 4. Run with Docker

Build and run the container:

```sh
docker build -t nyc-housing-connect-notifier .
docker run -v $(pwd)/config.json:/app/config.json:ro -v $(pwd)/data:/app/data nyc-housing-connect-notifier
```

Or use Docker Compose:

```sh
docker-compose up --build
```

## File Structure

- `main.go` - Entry point
- `config.sample.json` - Example configuration for one or more users
- `models/` - Data models
- `utils/`
  - `config/` - Config file loading and validation
  - `fetch_housing_connect/` - Fetches and filters lotteries
  - `discord/` - Sends Discord notifications
  - `local/` - Local file read/write helpers

# Screenshots

<img width="671" height="1104" alt="image" src="https://github.com/user-attachments/assets/85f05d9c-f562-4099-b7f3-d8c41480c832" />
