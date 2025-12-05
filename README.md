# HxnETHstakingAnalyticsApp

## ETH staking analytics

Lightweight Liquid Staking Token Analytics Dashboard

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
   - Git-based continuous deploys.

## Tracked LSTs (MVP)

Initial set:

- wstETH – Lido
- ankrETH – Ankr
- rETH – Rocket Pool
- wBETH – Binance
- pufETH – Puffer Finance

## Valuation Methodology

**Metric: 12-Month Average APR**

Each token's valuation score is computed from monthly price performance:

```
monthly_return[i] = (price_end[i] / price_start[i]) - 1
monthly_apr[i] = monthly_return[i] * 12

avg_12mo_apr = mean(monthly_apr[1..12])
```

Tokens are sorted:

- LOW avg_12mo_apr → "Undervalued"
- HIGH avg_12mo_apr → "Overvalued"

This avoids simple single-period overfitting and smooths volatility across time.

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

## Development Phases

- **Phase 1 — Infrastructure Setup (4 hours)** ✅ DONE
  - Initialize git repo.
  - Create Terraform configuration for Render:
    - PostgreSQL (Free Tier).
    - Backend web service.
    - Frontend static service.
  - Apply infra and provision services.

- **Phase 2 — Core Backend APIs (8 hours)** ✅ DONE
  - Create minimal SQL schema with LST tokens
  - Initialize Go application with Chi router
  - Set up Redis connection for caching
  - Connect PostgreSQL for token metadata
  - Implement CoinGecko price API with rate limiting
  - Expose REST endpoints: `/api/tokens`, `/api/token/{symbol}/history`

- **Phase 3 — Frontend MVP (6-8 hours)** 🎯 NEXT
  - Next.js + TypeScript application
  - Interactive price comparison chart (all LSTs vs ETH)
  - Different colored lines for each token
  - Connect to backend APIs
  - Responsive design

- **Phase 4 — Enhanced Features (8-10 hours)**
  - APR valuation calculations and rankings
  - Sortable valuation table
  - TVL data integration
  - Advanced caching strategies

- **Phase 5 — Production & CI/CD (4 hours)**
  - GitHub Actions for automated testing
  - Auto-deploy to Render on push
  - Environment configuration
  - Monitoring and logging

- **Phase 6 — Polish & Demo (2-4 hours)**
  - UI/UX improvements
  - Performance optimizations
  - Documentation updates
  - Demo deployment

## Estimated Effort

| Phase | Time | Status |
|-------|------|--------|
| Phase 1: Infrastructure | 4h | ✅ Done |
| Phase 2: Core Backend | 8h | ✅ Done |
| Phase 3: Frontend MVP | 6-8h | 🎯 Next |
| Phase 4: Enhanced Features | 8-10h | Planned |
| Phase 5: Production & CI/CD | 4h | Planned |
| Phase 6: Polish & Demo | 2-4h | Planned |
| **Total** | ≈ 32 – 42 hours | |

## Local Development

### Backend
```bash
cd backend
cp .env.example .env
go mod tidy
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
export RENDER_API_KEY="your_api_key_here"
terraform init
terraform apply
```

## Limitations (MVP)

- Render free Postgres is limited (~1GB) and may expire after 30 days; daily SQL backups should be exported externally.
- CoinGecko API rate-limited; all price calls are cached and executed in batch once per day.
- RPC indexer is optimized only for recent activity, not full historical chain scans.
- Not an official hosted The Graph subgraph (lightweight DB indexer replacement).

## Planned Enhancements

- Migrate indexer to The Graph decentralized network subgraph.
- DEX liquidity depth comparisons via Uniswap V3 subgraph.
- Wallet analytics (LST holdings across users).
- Peg deviation alerts.
- Automatic ETL to object storage (S3/GCS) for backups.

## Tech Stack

- **Frontend:** Next.js + TypeScript + Chart.js/Recharts
- **Backend:** Go + database/sql + Redis + Chi router
- **Cache:** Redis (in-memory data store)
- **Indexer:** go-ethereum getLogs poller
- **Database:** PostgreSQL (Render)
- **APIs:** CoinGecko, Ethereum RPC
- **Infra:** Terraform (Render provider), Docker, GitHub Actions CI/CD

## Demo

🔗 Live Dashboard: TBD  
📦 Source Code: GitHub link here

## Author

I (haxsen) built this as a Web3 frontend + backend + devops portfolio project showcasing:

✅ Data indexing  
✅ Blockchain RPC & contract calls  
✅ DeFi valuation analytics  
✅ Infra-as-code & CI/CD deployment
