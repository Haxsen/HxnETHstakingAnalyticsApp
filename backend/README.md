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
├── main.go                 # Application entry point with routes & handlers
├── go.mod                  # Go modules file
├── go.sum                  # Dependency checksums
├── .env.example           # Environment variables template
├── README.md              # This file
├── schema.sql             # Database schema
├── docs/                  # Generated Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── db/                # Database layer
│   │   ├── connection.go  # PostgreSQL connection
│   │   └── models.go      # Data models & queries
│   ├── cache/             # Redis caching layer
│   │   └── redis.go       # Redis client setup
│   └── services/          # Business logic & external APIs
│       ├── coingecko.go   # CoinGecko API client
│       └── cache.go       # Caching wrapper functions
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
| `POST` | `/api/cache/refresh` | Manually refresh Redis cache |
| `GET` | `/health` | Health check endpoint |
| `GET` | `/swagger/*` | Interactive API documentation |

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
