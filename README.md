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

2. **Lightweight "Subgraph" (Indexer)**
   - Polls on-chain logs via Ethereum RPC.
   - Tracks recent activity for selected LST contracts:
     - Transfers
     - Deposits / Withdrawals (where available)
   - Writes events and balances into PostgreSQL.
   - This project calls it "subgraph-style indexing" rather than a full hosted The Graph deployment for simplicity and free-tier compatibility.

3. **Backend API**
   - Node.js + TypeScript service that:
     - Pulls 1-year daily OHLC price data from CoinGecko.
     - Fetches total supply from contracts via RPC.
     - Calculates monthly APR metrics.
     - Aggregates data from PostgreSQL (on-chain events).
     - Serves REST endpoints for frontend.
   - Endpoints:
     - `GET /api/tokens`
     - `GET /api/token/:id/history`
     - `GET /api/token/:id/valuation`
     - `POST /api/snapshot` (daily cron snapshot job)

4. **PostgreSQL Database**
   - Stores:
     - tokens
     - events (on-chain activity)
     - daily_snapshots (price + supply + TVL)
     - valuation_metrics
   - DB is hosted on Render free-tier Postgres.

5. **Infrastructure as Code**
   - Terraform creates:
     - Render backend service
     - Render frontend static site
     - Render Postgres database
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
                    |    (React)         |
                    +----------+---------+
                               |
                               v
+------------+        +--------------------+       +------------------+
| CoinGecko  | -----> |     Backend        | ----> |   PostgreSQL     |
|   API      |        |  (Node + TS API)   |       |  (Render DB)     |
+------------+        +--------------------+       +------------------+
                               ^
                               |
                  +--------------------------+
                  |   On-chain Indexer       |
                  |   (RPC getLogs polling)  |
                  +--------------------------+
                               ^
                               |
                          Ethereum RPC
```

## Development Phases

- **Phase 1 — Infra & Repo Setup (4 hours)**
  - Initialize git repo.
  - Create Terraform configuration for Render:
    - PostgreSQL (Free Tier).
    - Backend web service.
    - Frontend static service.
  - Apply infra and provision services.

- **Phase 2 — Database Schema (4 hours)**
  - Create Prisma schema:
    - Tables:
      - tokens
      - events
      - daily_snapshots
      - valuation_metrics

- **Phase 3 — Backend & Indexer (10–14 hours)**
  - Initialize Node + TypeScript server.
  - Connect PostgreSQL.
  - Implement Ethereum RPC indexer:
    - `provider.getLogs()` polling loop.
    - Store recent Transfer, Deposit, Withdraw events.
  - Integrate CoinGecko API:
    - Fetch daily prices (365 days).
    - Cache results.
  - Expose public API endpoints.

- **Phase 4 — APR Valuation Logic (4–6 hours)**
  - Implement monthly APR calculation.
  - Aggregate 12-month averages.
  - Store results in DB.
  - Add sorting and ranking logic.

- **Phase 5 — Frontend (6–8 hours)**
  - Vite + React SPA.
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
| DB Schema | 4h |
| Backend & Indexer | 10–14h |
| Valuation Logic | 4–6h |
| Frontend | 6–8h |
| CI/CD | 4h |
| Docs & Polish | 2–4h |
| **Total** | ≈ 36 – 44 hours |

## Local Development

### Backend
```bash
cd backend
cp .env.example .env
pnpm install
pnpm dev
```

### Frontend
```bash
cd frontend
pnpm install
pnpm dev
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

- **Frontend:** React + Vite + Chart.js/Recharts
- **Backend:** Node.js + TypeScript + Express/Fastify
- **Indexer:** ethers.js getLogs poller
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
