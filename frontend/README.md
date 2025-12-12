# Frontend

Next.js dashboard for ETH Staking Analytics with interactive LST price charts and valuation tables.

## Overview

Responsive web app displaying comparative analytics for Ethereum Liquid Staking Tokens:
- Interactive multi-line price charts (1-year historical data)
- Sortable valuation table (APR, Stability Rating, TVL, Remarks)
- Real-time data from backend API
- Mobile-responsive design with dark/light themes

## Architecture

```
Frontend (Next.js) → Backend API → Redis Cache / External APIs
                        ↓
                   PostgreSQL (token metadata)
```

## Tech Stack

- **Framework**: Next.js 14+ (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Charts**: Apache ECharts
- **HTTP Client**: Fetch API (built-in)
- **Package Manager**: pnpm
- **Build Tool**: Next.js built-in
- **Deployment**: Render Static Site

## Project Structure

```
frontend/
├── app/                    # Next.js App Router pages
│   ├── layout.tsx         # Root layout
│   ├── page.tsx           # Home page with chart
│   └── globals.css        # Global styles
├── components/            # Reusable UI components
│   ├── Chart.tsx          # Main price chart component
│   ├── TokenSelector.tsx  # Token selection/filtering
│   └── Header.tsx         # App header
├── lib/                   # Utility functions and API clients
│   ├── api.ts             # Backend API client
│   ├── chartConfig.ts     # ECharts configuration
│   └── types.ts           # TypeScript type definitions
├── public/                # Static assets
│   └── favicon.ico
├── package.json           # Dependencies and scripts
├── tailwind.config.js     # Tailwind CSS configuration
├── next.config.js         # Next.js configuration
├── tsconfig.json          # TypeScript configuration
├── .env.local             # Local environment variables
└── README.md              # This file
```

## Prerequisites

- Node.js 18.17 or later
- pnpm 8.0 or later
- Backend API running (for development)

## Setup

1. **Install dependencies:**
   ```bash
   pnpm install
   ```

2. **Environment variables:**
   ```bash
   cp .env.example .env.local
   # Edit .env.local with your backend API URL
   ```

3. **Run development server:**
   ```bash
   pnpm run dev
   ```

4. **Open browser:**
   Navigate to `http://localhost:3000`

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API base URL | Yes | `http://localhost:8080` |

## API Integration

The frontend connects to the backend REST API for data:

### Endpoints Used

- `GET /api/tokens` - Fetch list of tracked LST tokens
- `GET /api/token/{symbol}/history` - Get 1-year price history for a token

### API Client

Located in `lib/api.ts`, provides typed functions for API calls with error handling and caching.

## Features Implemented
- Interactive Apache ECharts price comparison chart
- 1-year historical data for all LST tokens vs ETH
- Responsive design with mobile optimization
- Real-time API integration with error handling
- Clean UI with Tailwind CSS and theme support
- Sortable table showing APR, Stability Rating (1-10), TVL, Remarks
- Color-coded valuation remarks (green/red indicators)
- Mobile-responsive table with horizontal scroll
- Real-time data updates from backend API

## Development

### Available Scripts

```bash
# Start development server
pnpm run dev

# Build for production
pnpm run build

# Start production server
pnpm run start

# Run linting
pnpm run lint

# Run type checking
pnpm run type-check
```

### Code Style

- Use TypeScript for all new code
- Follow Next.js and React best practices
- Use Tailwind CSS for styling
- Component naming: PascalCase
- File naming: kebab-case for pages, PascalCase for components

## Chart Implementation

The main chart uses Apache ECharts with:
- Multi-line series for each LST token
- ETH as baseline comparison
- Responsive design with mobile optimization
- Hover tooltips showing exact values and dates
- Legend for token identification
- Configurable time ranges

Chart configuration is centralized in `lib/chartConfig.ts`.

## Building

```bash
pnpm run build
```

This creates an optimized production build in the `.next` directory.

## Testing

```bash
# Run tests (when implemented)
pnpm run test

# Run tests with coverage
pnpm run test:coverage
```

## Deployment

The frontend is deployed as a static site on Render:

1. Build the application: `pnpm run build`
2. Static files are served from the `out` directory
3. Auto-deploys on git pushes to main branch
4. Environment variables set in Render dashboard

## Performance

- Static generation for improved loading times
- Lazy loading of chart components
- Optimized bundle splitting
- Image optimization with Next.js

## Browser Support

- Modern browsers (Chrome, Firefox, Safari, Edge)
- Mobile browsers (iOS Safari, Chrome Mobile)
- IE11 not supported

## Future Enhancements (considerations)

- Real-time WebSocket updates for price data
- Advanced chart features (zoom, annotations)
- Token comparison tables
- APR valuation display
- Dark mode theme
- PWA capabilities
- Multi-language support
