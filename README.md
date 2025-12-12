# HxnETHstakingAnalyticsApp

## ETH staking analytics

### Lightweight Liquid Staking Token Analytics Dashboard

ETH staking analytics is a minimal full-stack Web3 analytics project that compares leading Ethereum Liquid Staking Tokens (LSTs) using 1-year price history data.

The application visualizes LST token performance against ETH on an interactive chart, with each token represented by a differently colored line showing its historical price trends.

## Features

1. **Interactive Price Chart**
   - 1-year historical price comparison for all LSTs vs ETH
   - Each token displayed as a differently colored line
   - Hover tooltips showing exact values and dates
   - Responsive design for desktop and mobile

2. **Token Management**
   - Dynamic token list from backend API
   - Support for multiple LST providers
   - Easy to add new tokens

3. **API Data Fetcher**
   - Fetches real-time on-chain data via Ethereum RPC calls.
   - Retrieves contract data for selected LST tokens:
     - Current balances and supply
     - Recent transaction activity
   - Data cached in Redis for performance.
   - Lightweight alternative to full blockchain indexing.

3. **Backend API**
   - Go service with Redis caching that:
     - Fetches real-time price data from CoinGecko API.
     - Retrieves on-chain data via Ethereum RPC calls.
     - Calculates APR metrics on-demand with cached results.
     - Minimal database usage for token metadata only.
     - Serves REST endpoints for frontend with Redis caching layer.
   - Endpoints:
     - `GET /api/tokens`
     - `GET /api/token/{tokenSymbol}/history`
     - `GET /api/token/{tokenSymbol}/valuation`
     - `GET /api/valuations` (sortable table data)
     - `POST /api/cache/refresh` (manual cache refresh)

4. **PostgreSQL Database**
   - Minimal storage for essential data:
     - tokens (basic metadata and contract addresses)
   - DB is hosted on Render free-tier Postgres.

5. **Infrastructure as Code**
   - Terraform creates:
     - Render backend service
     - Render frontend static site
     - Render Postgres database
     - Render Redis cache
   - Git-based continuous deploys (WIP).

## Tracked LSTs

Current set:

- wstETH – Lido
- ankrETH – Ankr
- rETH – Rocket Pool
- wBETH – Binance
- pufETH – Puffer Finance
- LSETH – Liquid Collective
- RSETH – Kelp DAO
- METH – Mantle
- CBETH – Coinbase
- TETH – Treehouse
- SFRXETH – Frax
- CDCETH – Crypto.com
- UNIETH – Universal (RockX / Bedrock)

## Key Features

- **Advanced Valuation**: APR calculations, stability rating (1-10 scale), and price-based valuation remarks
- **Real-time TVL**: On-chain total supply data for accurate circulating supply metrics
- **Comprehensive Analytics**: 1-year price history with interactive charting
- **Performance Optimized**: Redis caching for fast data retrieval
- **Production Ready**: Containerized deployment with infrastructure as code

## Architecture

```
                    +--------------------+
                    |    Frontend        |
                    |    (Next.js)       |
                    +----------+---------+
                               |
                               v
+------------+        +--------------------+       +------------------+
| CoinGecko  | -----> |     Backend        | ----> |   PostgreSQL     |
|   API      |        |  (Go API + Redis)  |       |  (Render DB)     |
+------------+        +--------------------+       +------------------+
                               ^                        ^
                               |                        |
                  +--------------------------+          |
                  |   On-chain Indexer       |          |
                  |   (RPC getLogs polling)  |          |
                  +--------------------------+          |
                               ^                     +--------+
                               |                     |  Redis |
                          Ethereum RPC               |  Cache |
                                                     +--------+
```

## Local Development

### Backend
```bash
cd backend
cp .env.example .env
go mod tidy
swag init
go run main.go
```

### Frontend
```bash
cd frontend
pnpm install
pnpm run dev
```

### Terraform (Deploy to Render)
```bash
cd infra/terraform
terraform init
terraform apply
```

## Limitations (MVP)

- Render free Postgres is limited (~1GB) and may expire after 30 days; daily SQL backups should be exported externally.
- CoinGecko API rate-limited; all price calls are cached and executed in batch once per day.
- RPC indexer is optimized only for recent activity, not full historical chain scans.
- Not an official hosted The Graph subgraph (lightweight DB indexer replacement).

## Tech Stack

- **Frontend:** Next.js + TypeScript + Apache ECharts
- **Backend:** Go + database/sql + Redis + Chi router
- **Cache:** Redis (in-memory data store)
- **Indexer:** go-ethereum getLogs poller
- **Database:** PostgreSQL (Render)
- **APIs:** CoinGecko, Ethereum RPC
- **Infra:** Terraform (Render provider), Docker, GitHub Actions CI/CD

## Demo

🔗 Live Dashboard: TBD  
📦 Source Code: GitHub link here

## License

Copyright © 2025 Haxsen (Hassan Ali). All rights reserved.

See [LICENSE](LICENSE) for details.

## Author

I (haxsen / Hassan Ali) built this as a Web3 frontend + backend + devops portfolio project showcasing:

✅ Data indexing
✅ Blockchain RPC & contract calls
✅ DeFi valuation analytics
✅ Infra-as-code & CI/CD deployment
