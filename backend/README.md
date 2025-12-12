# Backend

Go REST API for ETH Staking Analytics with Redis caching and PostgreSQL.

## Overview

Provides real-time LST analytics via REST endpoints:
- CoinGecko price data & on-chain TVL metrics
- APR calculations & stability ratings (1-10 scale)
- Redis caching for performance

## Architecture

```
HTTP → API Handlers → Services → Redis/PostgreSQL
```

## Tech Stack

**Go 1.21+** • **PostgreSQL** • **Redis** • **Chi Router** • **CoinGecko API**

## Project Structure

```
backend/
├── main.go                 # Entry point
├── internal/
│   ├── api/               # HTTP handlers & responses
│   ├── server/            # Server management & DI
│   ├── services/          # Business logic (tokens, valuation)
│   ├── db/                # PostgreSQL layer
│   └── cache/             # Redis layer
├── Dockerfile             # Container build
└── schema.sql            # Database schema
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

## Valuation Metrics

**APR Calculation**: 1-year average from monthly price performance
**Stability Rating**: 1-10 scale based on price volatility (10 = most stable)
**TVL**: On-chain total supply via ERC20 contracts
**Valuation Remarks**: 5-level assessment (Very Undervalued → Very Overvalued)

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

## Key Features

- **Redis Caching**: 1hr price data, 10min valuations, 1min on-chain data
- **Error Handling**: Graceful fallbacks to cached data on failures
- **Security**: Input validation, prepared statements, env-based config
- **Monitoring**: Health checks, structured logging, metrics endpoint

## Deployment Options

- **Render** (recommended for free tier)
- **Docker containers**
- **Kubernetes**
- **Cloud Run**

## Future Enhancements (considerations)

- GraphQL API support
- WebSocket real-time updates
- Multi-blockchain support
- Advanced caching strategies
