# Backend

Go-based REST API service for the ETH Staking Analytics application with Redis caching and minimal database usage.

## Overview

The backend provides REST API endpoints for the frontend dashboard, focusing on:
- Real-time price data from CoinGecko API
- On-chain data via Ethereum RPC calls
- Redis caching for performance
- Minimal PostgreSQL usage for token metadata

## Architecture

```
HTTP Requests → API Handlers → Services → Redis Cache / External APIs
                                      ↓
                                PostgreSQL (tokens only)
```

## Tech Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL (minimal usage)
- **Cache**: Redis (in-memory data store)
- **HTTP Router**: Chi
- **External APIs**: CoinGecko, Ethereum RPC

## Project Structure

```
backend/
├── main.go                 # Application entry point (~25 lines)
├── go.mod                  # Go modules file
├── go.sum                  # Dependency checksums
├── .env.example           # Environment variables template
├── README.md              # This file
├── schema.sql             # Database schema
├── docs/                  # Generated Swagger documentation
├── BACKEND_REFACTOR_PLAN.md # Refactor documentation
├── internal/
│   ├── api/               # HTTP transport layer
│   │   ├── handlers.go    # HTTP handlers with dependency injection
│   │   └── responses.go   # Common JSON response helpers
│   ├── server/            # Server management & DI container
│   │   └── server.go      # Server struct with clean startup/shutdown
│   ├── services/          # Business logic layer
│   │   ├── token_service.go    # Token business operations
│   │   ├── valuation_service.go # Valuation calculations & caching
│   │   ├── cache.go       # Caching wrapper functions
│   │   ├── coingecko.go   # CoinGecko API client
│   │   ├── tvl.go         # On-chain TVL fetching via ERC20 totalSupply
│   │   └── valuation.go   # APR calculations, stability rating, valuation remarks
│   ├── db/                # Database layer
│   │   ├── connection.go  # PostgreSQL connection
│   │   └── models.go      # Data models & queries
│   └── cache/             # Redis caching layer
│       └── redis.go       # Redis client setup
```

## Prerequisites

- Go 1.21 or later
- PostgreSQL database
- Redis instance
- CoinGecko API key (optional, for higher rate limits)

## Setup

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your values
   ```

3. **Database setup:**
   ```bash
   # Run schema.sql against your PostgreSQL database
   psql -d your_database < schema.sql
   ```

4. **Run the application:**
   ```bash
   go run main.go
   ```

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Yes | - |
| `REDIS_URL` | Redis connection URL | Yes | `redis://localhost:6379` |
| `COINGECKO_API_KEY` | CoinGecko API key | No | - |
| `ETHEREUM_RPC_URL` | Ethereum RPC endpoint | Yes | - |
| `PORT` | Server port | No | `8080` |
| `LOG_LEVEL` | Logging level | No | `info` |

## Database Schema

### tokens table
```sql
CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    contract_address VARCHAR(42) NOT NULL UNIQUE,
    decimals INTEGER NOT NULL DEFAULT 18,
    blockchain VARCHAR(20) NOT NULL DEFAULT 'ethereum',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/tokens` | List all tracked tokens |
| `GET` | `/api/token/{tokenSymbol}/history` | Get 1-year price history for a token (ETH denominated) |
| `GET` | `/api/token/{tokenSymbol}/valuation` | Get APR valuation metrics for a token |
| `GET` | `/api/valuations` | Get valuation metrics for all tokens (sortable table data) |
| `POST` | `/api/cache/refresh` | Manually refresh Redis cache |
| `GET` | `/health` | Health check endpoint |
| `GET` | `/swagger/*` | Interactive API documentation |

## Valuation Methodology

### Primary Metric: 1-Year Monthly Average APR

Each token's valuation score is computed from monthly price performance over 1 year:

**Step 1: Monthly Price Averages**
```
# Group 365 daily prices into 12 monthly chunks (~30 days each)
monthly_avg[m] = mean(daily_prices[chunk_m])
```

**Step 2: Monthly Returns**
```
monthly_return[1] = 0  # First month return (compared to itself)
monthly_return[m] = monthly_avg[m] - monthly_avg[m-1]  # For m = 2 to 12
# This gives 12 monthly return values
```

**Step 3: Annualized APR**
```
apr = sum(monthly_return[1..12])
```

### Stability Rating

Calculated as the coefficient of variation of daily returns, then normalized to a 1-10 rating scale:

```
stability_raw = 1 / (1 + std_dev(daily_return[1..365]) / abs(mean(daily_return[1..365])))
stability_rating = round(1 + (stability_raw - min_stability) / (max_stability - min_stability) * 9)
```

- **10/10**: Most stable token (lowest volatility)
- **1/10**: Least stable token (highest volatility)
- Rating is relative to other tokens in the current dataset

### TVL (Total Value Locked)

- Fetched via ERC20 `totalSupply()` contract calls on Ethereum mainnet
- Represents circulating supply of LST tokens
- Cached for 5 minutes due to on-chain data volatility
- Used as secondary ranking factor alongside APR and stability

### Valuation Remarks

Based on current price vs. expected price projection:

**Inputs:**
- Current Price: Today's price
- Last Month Average: Average price of most recent 30 days
- Average Monthly Return: APR ÷ 12 (average monthly price change)

**Expected Price Formula:**
```
expected_price = (average_monthly_return ÷ 2) + last_month_average
```

**Valuation Logic:**
- **Very Undervalued**: Current Price >1% below Expected Price
- **Undervalued**: Current Price 0.1%-1% below Expected Price
- **Fair Value**: Current Price ±0.1% of Expected Price
- **Overvalued**: Current Price 0.1%-1% above Expected Price
- **Very Overvalued**: Current Price >1% above Expected Price

## Phase 4 Features (Enhanced Valuation)

### Valuation Calculations
- **APR Calculation**: 1-year monthly average APR from 30-day price chunks (360 days total)
- **Monthly Processing**: Groups daily prices into 12 monthly averages, calculates differences
- **Stability Rating**: 1-10 rating scale based on coefficient of variation (10/10 = most stable)
- **Valuation Remarks**: 5-level assessment based on current price vs. expected price projection
- **Expected Price Formula**: `(average_monthly_return ÷ 2) + last_month_average`

### Enhanced Caching
- **Valuation Data**: 10-minute cache TTL for computed metrics
- **TVL Data**: 5-minute cache TTL for on-chain data
- **Price History**: 1-hour cache TTL (CoinGecko daily data)

### New Services
- `valuation.go`: Monthly APR calculations, stability scoring, 5-level valuation system
- `tvl.go`: On-chain TVL fetching via ERC20 totalSupply() contract calls

## Development

### Running locally:
```bash
go run main.go
```

### Building:
```bash
go build -o app main.go
./app
```

### Docker Build:
```bash
# Build Docker image
docker build -t eth-staking-analytics-backend .

# Run container
docker run -p 8080:8080 \
  -e DATABASE_URL="your_db_url" \
  -e REDIS_URL="your_redis_url" \
  -e COINGECKO_API_KEY="your_key" \
  -e ETHEREUM_RPC_URL="https://ethereum-rpc.publicnode.com" \
  eth-staking-backend
```

### Testing:
```bash
go test ./...
```

### Code formatting:
```bash
go fmt ./...
```

## Caching Strategy

- **Price history**: Cached for 1 hour (CoinGecko daily data doesn't change frequently)
- **On-chain data**: Cached for 1 minute
- **Valuation metrics**: Cached for 10 minutes
- **Token metadata**: Cached indefinitely (changes rarely)

## Error Handling

- API errors return appropriate HTTP status codes
- Database connection issues are logged and retried
- External API failures fall back to cached data when available
- Redis connection failures disable caching (app continues to work)

## Monitoring

- Structured logging with configurable levels
- Health check endpoint: `GET /health`
- Metrics endpoint: `GET /metrics` (future)

## Security

- Input validation on all API endpoints
- SQL injection prevention via prepared statements
- Rate limiting on external API calls
- Environment variables for sensitive configuration

## Deployment

The backend is designed to run in containers and can be deployed to:
- Render (web service)
- Docker containers
- Kubernetes
- Cloud Run

## Contributing

1. Follow Go conventions and formatting
2. Add tests for new functionality
3. Update documentation for API changes
4. Use meaningful commit messages

## Future Enhancements

- GraphQL API support
- WebSocket real-time updates
- Advanced caching strategies
- Multi-blockchain support
- API rate limiting per user
