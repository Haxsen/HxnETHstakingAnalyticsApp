# HxnETHstakingAnalyticsApp

## ETH staking analytics

Lightweight Liquid Staking Token Analytics Dashboard

ETH staking analytics is a minimal full-stack Web3 analytics project that compares leading Ethereum Liquid Staking Tokens (LSTs) using 1-year price history and live on-chain activity.

The application ranks tokens from "undervalued" → "overvalued" using an average APR metric derived from monthly price performance and visualizes their 1-year price trends on a single chart.

## Features

1. **Dashboard**
   - Interactive 1-year price comparison chart (all LSTs colored separately).
   - Sortable valuation table ranked by average APR (undervalue → overvalue).
   - Displays:
     - Current price (USD & ETH)
     - 1Y Avg APR (valuation metric)
     - TVL (ETH & USD)
     - Token contract links

2. **API Data Fetcher**
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
     - `GET /api/token/:id/history`
     - `GET /api/token/:id/valuation`
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

- stETH – Lido
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

- **Phase 1 — Infra & Repo Setup (4 hours)**
  - Initialize git repo.
  - Create Terraform configuration for Render:
    - PostgreSQL (Free Tier).
    - Backend web service.
    - Frontend static service.
  - Apply infra and provision services.

- **Phase 2 — Database Schema (2 hours)**
  - Create minimal SQL schema:
    - Tables:
      - tokens (basic metadata and contract addresses)

- **Phase 3 — Backend & API Integration (10–14 hours)**
  - Initialize Go application with modules.
  - Set up Redis connection for caching.
  - Connect PostgreSQL for minimal token metadata storage.
  - Implement API data fetchers:
    - CoinGecko price API integration with Redis caching.
    - Ethereum RPC calls for on-chain data with caching.
  - Implement Redis caching layer for API responses.
  - Expose REST API endpoints with cached data.

- **Phase 4 — APR Valuation Logic (4–6 hours)**
  - Implement APR calculation using cached price data.
  - Compute 12-month averages from API responses.
  - Cache valuation results in Redis.
  - Add sorting and ranking logic for dashboard.

- **Phase 5 — Frontend (6–8 hours)**
  - Next.js + TypeScript application.
  - Build price comparison chart.
  - Build sortable valuation table.
  - Connect to backend APIs.

- **Phase 6 — CI/CD & Cron (4 hours)**
  - GitHub Actions:
    - Test backend build.
    - Auto-deploy to Render on push.
  - Cron job (daily):
    - Trigger `/api/snapshot`.
    - Store daily DB snapshots.

- **Phase 7 — Docs & Demo (2–4 hours)**
  - Final README.
  - Screenshots + 30 sec demo GIF.
  - Publish demo URL.

## Estimated Effort

| Phase | Time |
|-------|------|
| Infra + Setup | 4h |
| DB Schema | 2h |
| Backend & Indexer | 10–14h |
| Valuation Logic | 4–6h |
| Frontend | 6–8h |
| CI/CD | 4h |
| Docs & Polish | 2–4h |
| **Total** | ≈ 34 – 42 hours |

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
- **Backend:** Go + database/sql + Redis + standard library
- **Cache:** Redis (in-memory data store)
- **Indexer:** go-ethereum getLogs poller
- **Database:** PostgreSQL (Render)
- **APIs:** CoinGecko, Ethereum RPC
- **Infra:** Terraform (Render provider), Docker, GitHub Actions CI/CD

## Demo

🔗 Live Dashboard: TBD  
📦 Source Code: GitHub link here

## Author

Built as a Web3 backend + devops portfolio project showcasing:

✅ Data indexing  
✅ Blockchain RPC & contract calls  
✅ DeFi valuation analytics  
✅ Infra-as-code & CI/CD deployment
