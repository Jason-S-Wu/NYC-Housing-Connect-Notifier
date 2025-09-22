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

### 2. Configure Environment Variables

Copy `.env.sample` to `.env` and fill in your values:

```env
DISCORD_WEBHOOK=your_discord_webhook_url
HOUSEHOLD_INCOME=your_household_income
HOUSEHOLD_SIZE=your_household_size
SAVE_FILE_NAME=rentals.json
```

**Variables:**

- `DISCORD_WEBHOOK`: Your Discord webhook URL
- `HOUSEHOLD_INCOME`: Your household's annual income (number)
- `HOUSEHOLD_SIZE`: Number of people in your household (number)
- `SAVE_FILE_NAME`: File to store seen lotteries (default: `rentals.json`)

### 3. Run Locally

```sh
go run main.go
```

### 4. Run with Docker

Build and run the container:

```sh
docker build -t nyc-housing-connect-notifier .
docker run --env-file .env -v $(pwd)/data:/app/data nyc-housing-connect-notifier
```

Or use Docker Compose:

```sh
docker-compose up --build
```

## File Structure

- `main.go` - Entry point
- `models/` - Data models
- `utils/`
  - `fetch_housing_connect/` - Fetches and filters lotteries
  - `discord/` - Sends Discord notifications
  - `local/` - Local file read/write helpers

# Screenshots
<img width="671" height="1104" alt="image" src="https://github.com/user-attachments/assets/85f05d9c-f562-4099-b7f3-d8c41480c832" />
