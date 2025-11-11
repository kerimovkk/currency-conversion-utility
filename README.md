# Currency Conversion Utility

A CLI utility for converting currencies using the CoinMarketCap API.

## Requirements

- Go 1.20 or higher
- CoinMarketCap API key (get one at https://pro.coinmarketcap.com)

## Installation

### Clone the repository

```bash
git clone https://github.com/kerimovkk/currency-conversion-utility.git
cd currency-conversion-utility
```

### Install dependencies

```bash
go mod download
```

### Build the application

```bash
go build -o app cmd/app/main.go
```

## Configuration

The application uses environment variables for configuration. Copy the example file and add your API key:

```bash
cp .env.example .env
```

Edit `.env` and set your CoinMarketCap API key:

```
CMC_API_KEY=your-api-key-here
CMC_API_URL=https://pro-api.coinmarketcap.com
```

**Note:** This application requires a production CoinMarketCap API key to get real market data.

## Usage

### Basic conversion

```bash
./app 123.45 USD BTC
```

Output:
```
123.45 USD = 0.00421337 BTC
```

### Verbose mode

```bash
./app --verbose 100 BTC USD
```

Output:
```
============================================================
CURRENCY CONVERSION RESULT
============================================================
Original Amount:    100 BTC
Converted Amount:   2934500 USD
------------------------------------------------------------
Exchange Rate:      1 BTC = 29345 USD
Last Updated:       2025-11-08 12:34:56 UTC
Query Time:         2025-11-08 12:35:10 UTC
============================================================
```

### Show help

```bash
./app --help
```

### Show version

```bash
./app --version
```

## Examples

Convert USD to Bitcoin:
```bash
./app 1000 USD BTC
```

Convert Bitcoin to Ethereum:
```bash
./app 1 BTC ETH
```

Convert EUR to GBP:
```bash
./app 50 EUR GBP
```

## Development

### Running Tests

Run all tests:
```bash
go test ./...
```

Run tests with verbose output:
```bash
go test ./... -v
```

Run tests with coverage:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Testing Individual Packages

```bash
# Test domain layer
go test ./internal/domain -v

# Test use case layer
go test ./internal/usecase -v

# Test CLI adapter
go test ./internal/adapter/cli -v

# Test retry mechanism
go test ./pkg/retry -v
```

## Error Handling

The application handles various error scenarios:

- **Invalid input**: Amount must be greater than zero, currency symbols must be valid
- **API errors**:
  - 401/403: Invalid or missing API key
  - 429: Rate limit exceeded (automatic retry with backoff)
  - 5xx: Server errors (automatic retry with backoff)

