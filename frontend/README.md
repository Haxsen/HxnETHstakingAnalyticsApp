This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.

----------------------------------------

# Frontend

Next.js-based frontend application for the ETH Staking Analytics dashboard, providing an interactive visualization of Liquid Staking Token performance.

## Overview

The frontend is a responsive web application that displays comparative price analytics for Ethereum Liquid Staking Tokens (LSTs). It features an interactive multi-line chart showing 1-year historical price data for all tracked LSTs against ETH, with real-time data fetched from the backend API.

Key features:
- Interactive price comparison chart using Apache ECharts
- Responsive design for desktop and mobile
- Token selection and filtering
- Real-time data updates from backend APIs
- Clean, modern UI with hover tooltips and legends

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

## Phase 3 Development Steps

This section outlines the implementation plan for Phase 3 frontend MVP (estimated 6-8 hours):

1. **Initialize Next.js Project** (30 min)
   - Create Next.js 14 app with TypeScript using `npx create-next-app`
   - Set up App Router structure
   - Configure basic project settings

2. **Install Dependencies** (30 min)
   - Install Apache ECharts and React wrapper
   - Add Tailwind CSS for styling
   - Install additional utilities (date-fns for date handling)
   - Configure pnpm as package manager

3. **Project Structure Setup** (30 min)
   - Create directory structure (app/, components/, lib/)
   - Set up Tailwind configuration
   - Configure TypeScript paths and settings
   - Create environment variable template

4. **API Integration** (1 hour)
   - Define TypeScript interfaces for API responses
   - Create API client functions in `lib/api.ts`
   - Implement error handling and loading states
   - Add basic caching for API calls

5. **Chart Component Development** (2 hours)
   - Create ECharts configuration for multi-line price chart
   - Implement token color coding and legends
   - Add hover tooltips with price and date information
   - Make chart responsive for mobile devices
   - Handle data loading and empty states

6. **Main Dashboard Page** (1 hour)
   - Build main page layout with header and chart container
   - Integrate chart component with API data
   - Add token selector/filter functionality
   - Implement loading states and error boundaries

7. **Responsive Design & Styling** (1 hour)
   - Apply Tailwind CSS for modern, clean UI
   - Ensure mobile-first responsive design
   - Add proper spacing, typography, and color scheme
   - Test across different screen sizes

8. **Testing & Optimization** (30 min)
   - Test API integration with backend
   - Verify chart renders correctly with real data
   - Optimize bundle size and loading performance
   - Add basic error handling for edge cases

9. **Build & Deployment Prep** (30 min)
   - Run production build and verify output
   - Configure static export for Render deployment
   - Update environment variables for production
   - Test build locally before deployment

## Phase 4 — Enhanced Features (6-8 hours)

### Simple Valuation Table (3 hours)
- Add a sortable table below the existing chart showing: Token, APR, Stability, TVL, Remarks
- Basic column sorting for APR, Stability, TVL (ascending/descending)
- Color-code remarks: green for "Undervalued"/"Very Undervalued", red for "Overvalued"/"Very Overvalued"
- Keep table responsive and simple with horizontal scroll on mobile

### API Integration (2 hours)
- Update `types.ts` with proper `ValuationData` interface matching backend
- Add `fetchValuations()` function to `api.ts` for `/api/valuations` endpoint
- Add loading/error states for valuation data in main page
- Handle API errors gracefully with user-friendly messages

### Layout Update (1 hour)
- Add valuation table section below existing chart on main page
- Simple section header: "Valuation Metrics"
- Maintain existing responsive design and theme support
- Keep donate button and theme toggle in header

### Basic Polish (30 min - 1 hour)
- Format numbers nicely: APR as percentage (2 decimals), TVL in millions/billions with M/B suffix
- Add "Last updated" timestamp for valuation data
- Test with real backend data and verify sorting works
- Ensure table looks good on mobile devices

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

## Contributing

1. Follow the existing code style and structure
2. Add TypeScript types for new data structures
3. Test components on multiple screen sizes
4. Update this README for any architectural changes

## Future Enhancements

- Real-time WebSocket updates for price data
- Advanced chart features (zoom, annotations)
- Token comparison tables
- APR valuation display
- Dark mode theme
- PWA capabilities
- Multi-language support
